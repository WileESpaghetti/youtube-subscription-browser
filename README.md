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
go run init-db.go
```

```bash
go run main.go
```

Open the given URL and give access to the API. After authorizing, I was redirected to a localhost URL that gave me a
not found error. I just had to copy the `code` parameter from the URL and copy/paste it into the terminal and press
enter. I might have used the wrong client type, but "other" wasn't available like the quickstart said to use.
