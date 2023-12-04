# GDOWN CLI Golang

Google DOWNloader

This utility can be used to download all the files automatically from a google drive folder one at a time (alternative to manually clicking download on each file on google drive web UI).

## Usage

### Prerequisite

Before running the `gdown` built executable, make sure to copy the Google Service Account `credentials.json` file to same folder as the `gdown` built executable.

### Show help

```shell
$  gdown
Download files and folders from Google Drive

Usage:
  gdown [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  download    Download the item
  help        Help about any command
  list        List details if the item

Flags:
  -h, --help   help for gdown

Use "gdown [command] --help" for more information about a command.
```

Help for `list` subcommand

```shell
$  gdown list --help
List details if the item

Usage:
  gdown list <id> [flags]

Flags:
  -h, --help   help for list
```

Help for `download` subcommand

```shell
$ gdown download --help
Download the item

Usage:
  gdown download <id> [flags]

Flags:
  -h, --help   help for download
  -l, --list   also list details of the item
```

### List file/folder Details

```shell
gdown list <id>
```

E.g.

```shell
$  gdown list 1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ
Id: 1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ
Name: file1.txt
MimeType: text/plain
Size: 16
```

### Download File/Folder

```shell
gdown download <id>
```

E.g.

```shell
$  gdown download 1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ
> downloads/file1.txt
downloading 100% |█████████████████████| (16/16 B, 8.5 kB/s)
```

### List and Download File/Folder

```shell
gdown download --list <id>
```

E.g.

```shell
$  gdown download --list 1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ
Id: 1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ
Name: file1.txt
MimeType: text/plain
Size: 16
> downloads/file1.txt
downloading 100% |██████████████████████| (16/16 B, 18 kB/s)
```

## Setup

### Using Libraries

- [google drive api](https://developers.google.com/drive/api/quickstart/go)
- [joho/godotenv](https://github.com/joho/godotenv)
- [schollz/progressbar](https://github.com/schollz/progressbar)
- [spf13/cobra](https://github.com/spf13/cobra)

### Google Cloud Credentials

Create a google cloud project and enable *google drive api*. After that create a new *Service Account* under *Create Credentials*. From the created Service Account, generate a new *Key* and select *Key Type* as JSON. After that credentials will be downloaded in a json file. Change it's name to `credentials.json` and put it in project directory or directory from where `gdown` will be run.

### Build and Run Instruction

#### Build

> Build (Linux)

```shell
go build -o build/gdown
```

> Build (Windows)

```shell
go build -o build/gdown.exe
```

#### Run

Help message

```shell
./build/gdown
```

List Details of File/Folder

```shell
./build/gdown list <id>
```

### Dev Setup

Running the application, show help message

```shell
go run gdown.go
```

List Details of File/Folder

```shell
go run gdown.go list <id>
```

Download File/Folder

```shell
go run gdown.go download <id>
```

List Details and Download File/Folder

```shell
go run gdown.go download --list <id>
```

### Scratch setup

```shell
go mod init github.com/arindampal-0/gdown-cli-golang
go get google.golang.org/api/drive/v3
go get golang.org/x/oauth3/google
go get github.com/joho/godotenv
go get -u github.com/schollz/progressbar/v3
go get -u github.com/spf13/cobra@latestgo get 
```

## TODO

- [ ] Google client authentication
- [x] Google service account authentication
- [x] Fetch file details
- [x] Fetch folder details and list of files
- [x] Download a file
- [x] Download all the files in a folder
- [x] Make it a cli application taking cli args
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
