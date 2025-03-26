# youtube-subscription-browser
dashboard with information about your youtube subscriptions

## Prerequisites
You will need to activate the YouTube data API and set up OAuth credentials.
I used the instructions from the Go quickstart: https://developers.google.com/youtube/v3/quickstart/go
to get the `client_secret.json` file.

The `client_secret.json` file will need to be in the same directory as the final executables.

While developing you will also need to add "Test Users" from the "Audience" page after setting up
the OAuth consent. See: https://www.youtube.com/watch?v=bkZns_VOB6Io

## Running

### Initialize the database
You will need to create the `youtube.sqlite` file and set up the database schema. The `init-db.go`
script will run through the schema setup for you.
```bash
go run scripts/init-db/main.go
```

### Use the YouTube Data API to Populate the Database
There seems to be a quirk that the YouTube data api will only return 1000 subscriptions. If you have more than that
it is recommended you use Google Takeout to get your list of subscriptions as a CSV.

```bash
go run scripts/fill-db/main.go
```

### Use a CSV file to Populate the Database
Any CSV should work as long as it has the following:
* a header row
* the first column is the channel ID

```bash
go run main.go channels.csv
```

#### Example Takeout File
```csv
Channel Id,Channel Url,Channel Title
UCK8sQmJBp8GCxrOtXWBpyEA,http://www.youtube.com/channel/UCK8sQmJBp8GCxrOtXWBpyEA,Google
UCtXKDgv1AVoG88PLl8nGXmw,http://www.youtube.com/channel/UCtXKDgv1AVoG88PLl8nGXmw,Google TechTalks
UCJS9pqu9BzkAMNTmzNMNhvg,http://www.youtube.com/channel/UCJS9pqu9BzkAMNTmzNMNhvg,Google Cloud Tech
UCbmNph6atAoGfqLoCL_duAg,http://www.youtube.com/channel/UCbmNph6atAoGfqLoCL_duAg,Talks at Google
UC_x5XG1OV2P6uZZ5FSM9Ttw,http://www.youtube.com/channel/UC_x5XG1OV2P6uZZ5FSM9Ttw,Google for Developers

```

### Authorizing the YouTube Data API

After running, open the given URL and give access to the API. After authorizing, I was redirected to a localhost URL
that gave me a not found error. I just had to copy the `code` parameter from the URL and copy/paste it into the 
terminal and press enter. I might have used the wrong client type, but "other" wasn't available like the quickstart
said to use. This is needed even if using a CSV file in order to get keyword, topic, and other channel data.
