package importer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const CacheKeySeparator = ":"
const defaultCacheDir = ".cache"

type Cache interface {
	Put(key string, item any) error
	Get(key string) (any, error) // FIXME how to get back to original type?
	//Clear() error // TODO
}

// NewCache creates a new cache using the default cache provider
func NewCache() Cache {
	return NewFileCache(defaultCacheDir)
}

// FIleCache

type fileCache struct {
	root string
}

// NewFileCache creates a new cache that stores data as JSON files in the given directory.
// Prefer using NewCache() instead.
func NewFileCache(dir string) *fileCache {
	fc := &fileCache{root: dir}
	_ = fc.init() // FIXME how should I handle error? New*() functions don't return errors typically

	return fc
}

func (fc *fileCache) Put(key string, item any) error {
	parts := strings.Split(key, CacheKeySeparator)        // TODO maybe turn into directories?
	fileName := strings.Join(append(parts, ".json"), "-") // FIXME might want to do some extra processing like putting stuff in directory

	f, err := os.Create(filepath.Join(fc.root, fileName)) // existing file will be overwritten
	if err != nil {
		return fmt.Errorf("could not create cache file: file = \"%s\": %w", fileName, err)
	}
	defer f.Close()

	encoder := json.NewEncoder(f)

	err = encoder.Encode(item)
	if err != nil {
		return fmt.Errorf("could not serialize item: %w", err) // FIXME maybe use some to string method?
	}

	return nil
}

func (fc *fileCache) Get(key string) (any, error) {
	return nil, fmt.Errorf("not implemented")
}

// init creates the cache directory if it doesn't exist
func (fc *fileCache) init() error {
	err := os.MkdirAll(fc.root, 0755)
	if err != nil {
		path, fErr := filepath.Abs(fc.root)
		if fErr != nil {
			path = fc.root
		}

		return fmt.Errorf("could not create cache directory: dir = \"%s\": %w", path, err)
	}

	return nil
}

// nullCache

type nullCache struct{}

// NewNullCache creates a cache that does nothing. It is used if you want to disable caching.
// Prefer using NewCache() instead.
func NewNullCache() *nullCache {
	return &nullCache{}
}

func (nc *nullCache) Put(key string, item any) error {
	return nil
}

func (nc *nullCache) Get(key string) (any, error) {
	return nil, nil
}

// refreshCache

type refreshCache struct {
	cache Cache
}

// NewRefreshCache creates a cache that allows overwriting an existing cache, but never returns anything.
// This is useful if you want to force a cache refresh.
// Prefer using NewCache() instead.
func NewRefreshCache(cache Cache) *refreshCache {
	return &refreshCache{cache: cache}
}

func (rc *refreshCache) init() error {
	// FIXME should this clear the cache?
	return nil
}

func (rc *refreshCache) Put(key string, item any) error {
	return rc.cache.Put(key, item)
}

func (rc *refreshCache) Get(key string) (any, error) {
	return nil, nil
}
