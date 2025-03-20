// Sample Go code for user authorization

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"slices"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
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
		Mine(true)

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

	channels := make([]*youtube.Channel, 0, len(allSubscriptions))

	apiIdLimit := 50 // FIXME not sure if this is documented somewhere, but I found it on a stack overflow
	for page := range slices.Chunk(channelIds, apiIdLimit) {
		call2 := y.Channels.List([]string{"snippet", "brandingSettings", "id", "statistics", "topicDetails"}).Id(page...)
		err = call2.Pages(context.TODO(), func(page *youtube.ChannelListResponse) error {
			channels = append(channels, page.Items...)
			return nil
		})
		if err != nil {
			log.Fatalf("Error listing channels: %v", err)
		}
	}

	fmt.Println("Your Subscriptions:")
	for _, subscription := range channels {
		fmt.Printf("- %s (Channel ID: %s)\n", subscription.Snippet.Title, subscription.Id)

		_, err = db.Exec("INSERT INTO channels(youtube_id, title, description, custom_url, branding_title, branding_description, subscriber_count, video_count) VALUES(?, ?, ?, ?, ?, ?, ?, ?)",
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
	}

}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
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
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	//channelsListByUsername(service, "snippet,contentDetails,statistics", "GoogleDevelopers")
	listSubscriptions(service, db)
}
