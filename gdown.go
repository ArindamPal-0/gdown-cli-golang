package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

type Folder struct {
	Id      string
	Name    string
	Folders []Folder
	Files   []File
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
	fileId := "1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ"
	/* Get File details as response */
	fileRes, err := driveService.Files.Get(fileId).Fields("id", "name", "mimeType", "size").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve file: %v", err)
	}

	// fmt.Printf("fileId: %+v\n", fileRes.Id)
	// fmt.Printf("fileName: %+v\n", fileRes.Name)
	// fmt.Printf("fileMimeType: %+v\n", fileRes.MimeType)
	// fmt.Printf("fileSize: %+v\n", fileRes.Size)

	file := File{
		Id:       fileRes.Id,
		Name:     fileRes.Name,
		MimeType: fileRes.MimeType,
		Size:     fileRes.Size,
	}

	fmt.Printf("%#v\n", file)

	/* Fetching folder details */
	// gdown folder
	folderId := "1SVHxav6Y5LoYbdgfx2MSsdYlT74RTjej"
	/* Get Folder details as response */
	folderRes, err := driveService.Files.Get(folderId).Fields("id", "name", "mimeType", "size").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve folder: %v", err)
	}

	folder := Folder{
		Id:      folderRes.Id,
		Name:    folderRes.Name,
		Folders: []Folder{},
		Files:   []File{},
	}

	fmt.Printf("%#v\n", folder)
}
