# GDOWN CLI Golang

Google DOWNloader

This utility can be used to download all the files automatically from a google drive folder one at a time (alternative to manually clicking download on each file on google drive web UI).

## Setup

### Google Cloud Credentials

Create a google cloud project and enable *google drive api*. After that create a new *Service Account* under *Create Credentials*. From the created Service Account, generate a new *Key* and select *Key Type* as JSON. After that credentials will be downloaded in a json file. Change it's name to `credentials.json` and put it in project directory or directory from where `gdown` will be run.

### Build and Run Instruction

Build
```shell
go build -o build/gdown
```

Run
```shell
./build/gdown
```

### Dev Setup

Running the application
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
- [x] Google service account authentication
- [x] Fetch file details
- [x] Fetch folder details and list of files
- [ ] Download a file
- [ ] Make it a cli application taking cli args
- [ ] Download all the files in a folder
- [ ] Download a folder recursively

## Common Issues

### WSL problem opening URL from terminal (used in google oauth2 login)

[No method available for opening url - wsl github issue](https://github.com/microsoft/WSL/issues/8892)

Installing `xdg-utils` and `wslu` fixes the issue.
```shell
sudo apt install xdg-utils
sudo add-apt-repository ppa:wslutilities/wslu
sudo apt update
sudo apt install wslu
```

Opening of URL from terminal is handled by [browser package](https://github.com/pkg/browser)
