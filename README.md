# GDOWN CLI Golang

Google DOWNloader

This utility can be used to download all the files automatically from a google drive folder one at a time (alternative to manually clicking download on each file on google drive web UI).

## Setup

Build
```shell
go build -o build/gdown
```

Run
```shell
./build/gdown
```
Make sure to run the executable from the directory containing .env file with the required env variables `GOOGLE_OAUTH_CLIENT_ID` and `GOOGLE_OAUTH_CLIENT_SECRET`. 

Or provide them while running the application like this
```shell
GOOGLE_OAUTH_CLIENT_ID={CLIENT_ID} GOOGLE_OAUTH_CLIENT_SECRET={CLIENT_SECRET} ./build/gdown
```

### Dev Setup

Setup environment variables first
```.env
GOOGLE_OAUTH_CLIENT_ID = 
GOOGLE_OAUTH_CLIENT_SECRET = 
```

```shell
go run gdown.go
```

### Scratch setup

```shell
go mod init github.com/arindampal-0/gdown-cli-golang
go get google.golang.org/api/drive/v3
go get golang.org/x/oauth3/google
go get github.com/joho/godotenv
```

## TODO

- [ ] Google client authentication
- [ ] Fetch file details
- [ ] Fetch folder details and list of files
- [ ] Download a file
- [ ] Make it a cli application taking cli args
- [ ] Download all the files in a folder
- [ ] Download a folder recursively
