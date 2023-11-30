package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type File struct {
	Id       string
	Name     string
	MimeType string
	Size     int64
}

const folderMimeType = "application/vnd.google-apps.folder"

type Folder struct {
	Id      string
	Name    string
	Folders []Folder
	Files   []File
}

func Prettify(data interface{}) string {
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("❌ Could not convert data to json")
		return ""
	}
	return string(d)
}

/* Get File Details */
func GetFile(driveService *drive.Service, fileId string) (*File, error) {
	/* Get File details as response */
	fileRes, err := driveService.Files.Get(fileId).Fields("id", "name", "mimeType", "size").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file: %v", err)
	}

	// fmt.Printf("fileId: %+v\n", fileRes.Id)
	// fmt.Printf("fileName: %+v\n", fileRes.Name)
	// fmt.Printf("fileMimeType: %+v\n", fileRes.MimeType)
	// fmt.Printf("fileSize: %+v\n", fileRes.Size)

	return &File{
		Id:       fileRes.Id,
		Name:     fileRes.Name,
		MimeType: fileRes.MimeType,
		Size:     fileRes.Size,
	}, nil
}

/* Get Folder Details */
func GetFolder(driveService *drive.Service, folderId string) (*Folder, error) {

	/* Get Folder details as response */
	folderRes, err := driveService.Files.Get(folderId).Fields("id", "name", "mimeType").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve folder: %v", err)
	}

	/* Fetching Files List */
	filesListRes, err := driveService.Files.List().Q(fmt.Sprintf("\"%s\" in parents", folderId)).Fields("files(id, name, mimeType, size)").OrderBy("name").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file list: %v", err)
	}

	/* Constructing File and Folder lists */
	filesList := []File{}
	foldersList := []Folder{}
	for _, file := range filesListRes.Files {
		if file.MimeType != folderMimeType {
			filesList = append(filesList, File{
				Id:       file.Id,
				Name:     file.Name,
				MimeType: file.MimeType,
				Size:     file.Size,
			})
		} else {
			folder, err := GetFolder(driveService, file.Id)
			if err != nil {
				continue
			}
			foldersList = append(foldersList, *folder)
		}
	}

	/* Creating Parent Folder Struct Object */
	return &Folder{
		Id:      folderRes.Id,
		Name:    folderRes.Name,
		Folders: foldersList,
		Files:   filesList,
	}, nil
}

const DownloadFolderPath = "downloads"

/* Download a Single file */
func DownloadFile(driveService *drive.Service, file *File, folderPath string) error {
	/* Send Download File Request */
	fileRes, err := driveService.Files.Get(file.Id).Download()
	if err != nil {
		return fmt.Errorf("could not download file, %v", err)
	}
	defer fileRes.Body.Close()

	/* Create Folder */
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create download directory: %s, %v", DownloadFolderPath, err)
	}

	/* Open File Handle */
	filePath := fmt.Sprintf("%s/%s", folderPath, file.Name)
	fileHandle, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("could not open and/or create file: %s, %v", filePath, err)
	}
	defer fileHandle.Close()

	fmt.Printf("> %s\n", filePath)
	/* setup progress bar */
	bar := progressbar.DefaultBytes(
		fileRes.ContentLength,
		"downloading",
	)

	/* Download into opened file */
	io.Copy(io.MultiWriter(fileHandle, bar), fileRes.Body)

	return nil
}

/* Download all Files in a Folder */
func DownloadFilesInFolder(driveService *drive.Service, folder *Folder) {
	/* Derive Folder path */
	folderPath := fmt.Sprintf("%s/%s", DownloadFolderPath, folder.Name)
	/* Download all the files one at a time */
	for _, file := range folder.Files {
		err := DownloadFile(driveService, &file, folderPath)
		if err != nil {
			log.Printf("Error Downloading file: %s, %v", file.Name, err)
		}
	}
}

func main() {
	fmt.Println("GDOWN CLI")

	ctx := context.Background()

	/* interface for reading credentials.json file */
	type Credential struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}

	/* Reading credentials.json file */
	content, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalln("Could not read file credentials.json")
	}

	/* Converting json to Credentials */
	var cred Credential
	err = json.Unmarshal(content, &cred)
	if err != nil {
		log.Fatalln("Could not convert json to Credentials interface")
	}

	/* Creating JWT Config object */
	config := &jwt.Config{
		Email:      cred.ClientEmail,
		PrivateKey: []byte(cred.PrivateKey),
		Scopes: []string{
			// See information about your Google Drive files
			"https://www.googleapis.com/auth/drive.metadata.readonly",
			// See and download all your Google Drive Files
			"https://www.googleapis.com/auth/drive.readonly",
		},
		TokenURL: google.JWTTokenURL,
	}

	/* Getting client from the jwt config */
	client := config.Client(ctx)

	/* Creating drive service for the client */
	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	/* Creating drive service using the credentials file and scopes */
	// driveService, err := drive.NewService(context.Background(), option.WithCredentialsFile("credentials.json"), option.WithScopes(
	// 	"https://www.googleapis.com/auth/drive.metadata.readonly",
	// 	"https://www.googleapis.com/auth/drive.readonly",
	// ))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	/* Fetching file details */
	// file1.txt
	// fileId := "1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ"
	// file, err := GetFile(driveService, fileId)
	// if err != nil {
	// 	log.Fatalf("Error getting File.\n%v", err)
	// }

	// fmt.Printf("file: %s\n", Prettify(file))

	// err = DownloadFile(driveService, file, downloadFolderPath)
	// if err != nil {
	// 	log.Fatalf("Error downloading file.\n%v", err)
	// }

	/* Fetching folder details */
	// gdown folder
	folderId := "1SVHxav6Y5LoYbdgfx2MSsdYlT74RTjej"
	folder, err := GetFolder(driveService, folderId)
	if err != nil {
		log.Fatalf("Error getting Folder.\n%v", err)
	}

	fmt.Printf("folder: %s\n", Prettify(folder))

	/* Download all files in the folder */
	DownloadFilesInFolder(driveService, folder)
}
