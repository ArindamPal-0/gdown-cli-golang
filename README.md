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

> During the authentication process, paste the URL in browser and it'll generate redirect url with `authorization key`, paste the authorization key in the terminal to authenticate.
>
> redirect url: `http://localhost:8000/auth/google/callback?state=state-token&code=4%2F0AfJohXnfe9bYEgbx2xhRGm35swCYmPr3yWnNQ7Qwswv7l6ro8R2s2mPv1p5TSlF0vcfMdA&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.metadata.readonly+https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fdrive.readonly`
> 
> `code=4%2F0AfJohXnfe9bYEgbx2xhRGm35swCYmPr3yWnNQ7Qwswv7l6ro8R2s2mPv1p5TSlF0vcfMdA` This part contains the authorization key but it is browser encoded.
> 
> So the authorization key will be `4/0AfJohXnfe9bYEgbx2xhRGm35swCYmPr3yWnNQ7Qwswv7l6ro8R2s2mPv1p5TSlF0vcfMdA`, i.e. `%2F` translates to `/`

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
