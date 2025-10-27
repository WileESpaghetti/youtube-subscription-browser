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

type CacheError struct{}

func (e *CacheError) Error() string {
	return "cache error"
}

type Cache interface {
	Put(key string, item any) error
	Get(key string, item any) error
	Has(key string) bool
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
	parts := strings.Split(key, CacheKeySeparator)         // TODO maybe turn into directories?
	fileName := strings.Join(append(parts), "-") + ".json" // FIXME might want to do some extra processing like putting stuff in directory

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

func (fc *fileCache) Get(key string, item any) error {
	fileName := cacheKeyToFileName(key)
	cacheFile := filepath.Join(fc.root, fileName)

	f, err := os.Open(cacheFile)
	if err != nil {
		return fmt.Errorf("%s not found: %w", key, fmt.Errorf("error opening cache file: file = %s : %s\n", cacheFile, err))
	}

	jd := json.NewDecoder(f)
	err = jd.Decode(item)
	if err != nil {
		return fmt.Errorf("%s not found: %w", key, fmt.Errorf("error decoding cache file: file = %s : %s\n", cacheFile, err))
	}

	return nil
}

// Has checks if the given key exists in the cache.
// For FileCache, this means that the cache file exists, but does not necessarily mean you have permissions to read
// it, or that it can be successfully deserialized.
func (fc *fileCache) Has(key string) bool {
	fileName := cacheKeyToFileName(key)
	cacheFile := filepath.Join(fc.root, fileName)
	_, err := os.Stat(cacheFile)
	return err == nil
}

func cacheKeyToFileName(key string) string {
	// channel = channel-$ID.json
	return strings.Join(strings.Split(key, CacheKeySeparator), "-") + ".json"
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

func (nc *nullCache) Get(key string, item any) error {
	return fmt.Errorf("%s not found", key)
}

func (nc *nullCache) Has(key string) bool {
	return false
}

// refreshCache

type refreshCache struct {
	cache Cache
}

// NewRefreshCache creates a cache that allows overwriting an existing cache, but never returns anything.
// This is useful if you want to force a cache refresh.
// Prefer using NewCache() instead.
func NewRefreshCache(cache ...Cache) *refreshCache {
	rc := &refreshCache{}

	if len(cache) == 0 {
		rc.cache = NewCache()
		return rc
	}

	if cache[0] == nil { // FIXME might need to jump through some hoops to find out if this is REALLY nil
		rc.cache = NewCache()
	} else {
		rc.cache = cache[0]
	}

	return rc
}

func (rc *refreshCache) init() error {
	// FIXME should this clear the cache?
	return nil
}

func (rc *refreshCache) Put(key string, item any) error {
	return rc.cache.Put(key, item)
}

func (rc *refreshCache) Get(key string, item any) error {
	return fmt.Errorf("%s not found", key)
}

func (rc *refreshCache) Has(key string) bool {
	// wrapped cache might have it, but we want to pretend we don't, to encourage overwriting what exists
	return false
}
