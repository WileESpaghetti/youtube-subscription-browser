// Sample Go code for user authorization

package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"slices"
	"strings"
)

const missingClientSecretsMessage = `
Please configure OAuth 2.0
`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func channelsListByUsername(service *youtube.Service, part string, forUsername string) {
	call := service.Channels.List([]string{part})
	call = call.ForUsername(forUsername)
	response, err := call.Do()
	handleError(err, "")
	fmt.Println(fmt.Sprintf("This channel's ID is %s. Its title is '%s', "+
		"and it has %d views.",
		response.Items[0].Id,
		response.Items[0].Snippet.Title,
		response.Items[0].Statistics.ViewCount))
}

func listSubscriptions(y *youtube.Service, db *sql.DB) {

	call := y.Subscriptions.List([]string{"snippet", "contentDetails", "id"}).
		Mine(true).
		MaxResults(25)

	var allSubscriptions []*youtube.Subscription

	err := call.Pages(context.TODO(), func(page *youtube.SubscriptionListResponse) error {
		allSubscriptions = append(allSubscriptions, page.Items...)
		return nil
	})

	if err != nil {
		log.Fatalf("Error listing subscriptions: %v", err)
	}

	if len(allSubscriptions) == 0 {
		fmt.Println("No subscriptions found.")
		return
	}

	channelIds := make([]string, 0, len(allSubscriptions))
	for _, s := range allSubscriptions {
		channelIds = append(channelIds, s.Snippet.ResourceId.ChannelId)
	}

	// break here
	populateDatabase(y, db, channelIds)
}

func populateDatabase(y *youtube.Service, db *sql.DB, channelIDs []string) {
	channels := make([]*youtube.Channel, 0, len(channelIDs))

	apiIdLimit := 50 // FIXME not sure if this is documented somewhere, but I found it on a stack overflow
	channelIdCount := 0
	for page := range slices.Chunk(channelIDs, apiIdLimit) {
		// FIXME add random sleep between pages
		// FIXME how to detect if a channel has been removed? API returns no data in that case, which is confusing when counts don't align
		call2 := y.Channels.List([]string{"snippet", "brandingSettings", "id", "statistics", "topicDetails"}).Id(page...)
		err := call2.Pages(context.TODO(), func(page *youtube.ChannelListResponse) error {
			channelIdCount = channelIdCount + len(page.Items)
			channels = append(channels, page.Items...)
			return nil
		})
		if err != nil {
			log.Fatalf("Error listing channels: %v", err)
		}
	}

	insertCount := 0
	fmt.Println("Your Subscriptions:")
	for _, subscription := range channels {
		insertCount += 1
		fmt.Printf("- %d, %s (Channel ID: %s)\n", insertCount, subscription.Snippet.Title, subscription.Id)

		_, err := db.Exec("INSERT INTO channels(youtube_id, title, description, custom_url, branding_title, branding_description, subscriber_count, video_count) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
			subscription.Id,
			subscription.Snippet.Title,
			subscription.Snippet.Description,
			subscription.Snippet.CustomUrl,
			subscription.BrandingSettings.Channel.Title,
			subscription.BrandingSettings.Channel.Description,
			subscription.Statistics.ViewCount,
			subscription.Statistics.VideoCount,
		)
		if err != nil {
			fmt.Printf("...unable to save: %s\n", err)
		}

		// TODO get channelID
		channelID, err := getChannelID(db, subscription.Id)
		if err != nil {
			fmt.Printf("...unable to get channel row ID for subscription: %s\n", err)
			continue // can not save topic/keyword associations if we do not have a channel
		}

		if subscription.TopicDetails != nil { // work around nil pointer panic on some stuff
			fmt.Printf("- %d, %s ...(Topics ID: %s)\n", insertCount, subscription.TopicDetails.TopicIds, subscription.Id)
			topicIDs, err := getTopicIDs(db, subscription.TopicDetails.TopicIds)
			if err != nil {
				fmt.Printf("...unable to get topic ids: %s\n", err)
			}

			err = saveTopicAssociations(db, channelID, topicIDs)
			if err != nil {
				fmt.Printf("...unable to save topic ids: %s\n", err)
			}
		} else {
			fmt.Printf("...topic information not found: %#v\n", subscription)
		}

		keywords, err := splitKeywords(subscription.BrandingSettings.Channel.Keywords)
		if err != nil {
			fmt.Printf("...unable to split keywords: %s\n", err)
		}

		fmt.Printf("- %d, %s ...(Keywords: %#v)\n", insertCount, keywords, subscription.Id)

		err = saveKeywords(db, keywords)
		if err != nil {
			fmt.Printf("...unable to save keywords: %s\n", err)
		}

		keywordIDs, err := getKeywordIDs(db, keywords)
		if err != nil {
			fmt.Printf("...unable to get keyword ids: %s\n", err)
		}

		err = saveKeywordAssociations(db, channelID, keywordIDs)
		if err != nil {
			fmt.Printf("...unable to save keyword ids: %s\n", err)
		}
	}
}

func getChannelID(db *sql.DB, youtubeID string) (int, error) {
	var id int

	err := db.QueryRow("SELECT id FROM channels WHERE youtube_id = ?", youtubeID).
		Scan(&id)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return id, err
	case err != nil:
		return id, err
	default:
		return id, nil
	}
}

func saveKeywordAssociations(db *sql.DB, channelID int, keywordIDs []int) error {
	// FIXME since keyword should be unique we could probably do this with an insert with select/join to automatically
	//  fetch the ids
	if len(keywordIDs) == 0 {
		return nil
	}

	for _, keywordID := range keywordIDs {
		result, err := db.Exec("INSERT INTO channels_keywords(channel_id, keyword_id) VALUES(?, ?)", channelID, keywordID)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func saveKeywords(db *sql.DB, keywords []string) error {
	// TODO probably want to remove '#' from hashtag keywords at some point
	// TODO seems like some people also use multiple hashtags as a single keyword. may want to additionally split that
	if len(keywords) == 0 {
		return nil
	}

	for _, k := range keywords {
		result, err := db.Exec("INSERT INTO keywords(keyword) VALUES(?)", k)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func getKeywordIDs(db *sql.DB, keywords []string) ([]int, error) {
	if len(keywords) == 0 {
		return nil, nil
	}

	// handle IN clause placeholders
	keywordPlaceholders := strings.Repeat("?,", len(keywords))
	keywordPlaceholders = keywordPlaceholders[:len(keywordPlaceholders)-1] // strip off the trailing ,
	args := make([]interface{}, 0, len(keywords))
	for _, id := range keywords {
		args = append(args, id)
	}

	queryTopicIDs := fmt.Sprintf("SELECT id FROM keywords WHERE keyword in (%s)", keywordPlaceholders)
	rows, err := db.Query(queryTopicIDs, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0, len(keywords))
	for rows.Next() {
		var keywordID int
		if err := rows.Scan(&keywordID); err != nil {
			return nil, err
		}

		ids = append(ids, keywordID)
	}

	if rerr := rows.Close(); rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

// splitKeywords splits the list of keywords provided by the YouTube data API.
// The keywords are separated by a space, but if a keyword should contain
// multiple words then those words will be quoted. This format allows us to
// treat the keyword list as a space separated CSV record.
func splitKeywords(s string) ([]string, error) {
	if len(s) == 0 {
		return nil, nil
	}

	keywords := make([]string, 0)

	buf := strings.NewReader(s)

	splitter := csv.NewReader(buf)
	splitter.Comma = ' '

	records, err := splitter.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for _, r := range records {
		// there should only be one record, but we'll assume that's not the case
		for _, k := range r {
			keywords = append(keywords, strings.ToLower(k))
		}
	}

	return keywords, nil
}

func saveTopicAssociations(db *sql.DB, channelID int, topicIDs []int) error {
	if len(topicIDs) == 0 {
		return nil
	}

	for _, topicID := range topicIDs {
		result, err := db.Exec("INSERT INTO channels_topics(channel_id, topic_id) VALUES(?, ?)", channelID, topicID)
		if err != nil {
			return err
		}

		_, err = result.RowsAffected()
		if err != nil {
			return err
		}
	}

	return nil
}

func getTopicIDs(db *sql.DB, topicIDs []string) ([]int, error) {
	if len(topicIDs) == 0 {
		return nil, nil
	}

	// handle IN clause placeholders
	topicPlaceholders := strings.Repeat("?,", len(topicIDs))
	topicPlaceholders = topicPlaceholders[:len(topicPlaceholders)-1] // strip off the trailing ,
	args := make([]interface{}, 0, len(topicIDs))
	for _, id := range topicIDs {
		args = append(args, id)
	}

	queryTopicIDs := fmt.Sprintf("SELECT id FROM topics WHERE topic_id in (%s)", topicPlaceholders)
	rows, err := db.Query(queryTopicIDs, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := make([]int, 0, len(topicIDs))
	for rows.Next() {
		var topicID int
		if err := rows.Scan(&topicID); err != nil {
			return nil, err
		}

		ids = append(ids, topicID)
	}

	if rerr := rows.Close(); rerr != nil {
		return nil, err
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func getChannelsFromTakeout(file string) ([]string, error) {
	if len(file) == 0 {
		// no file given
		return nil, errors.New("no input file given")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	splitter := csv.NewReader(f)

	records, err := splitter.ReadAll()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	channels := make([]string, 0, len(records))

	for _, r := range records[1:] {
		channels = append(channels, r[0])
	}

	return channels, nil
}

func main() {
	ctx := context.Background()

	b, err := os.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	dbFile := "youtube.sqlite"
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFile))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if len(os.Args) > 1 {
		fmt.Printf("Using CSV file: %#v\n", os.Args)
		channels, err := getChannelsFromTakeout(os.Args[1])
		fmt.Printf("channel cound: %d\n", len(channels))
		if err != nil {
			fmt.Printf("can not read input file: '%s':, %s\n", os.Args[1], err)
			os.Exit(1)
		}
		populateDatabase(service, db, channels)
	} else {
		//channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
		listSubscriptions(service, db)
	}
}
