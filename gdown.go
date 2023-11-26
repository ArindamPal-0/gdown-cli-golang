package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	token, err := getTokenFromFile(tokenFile)

	if !token.Valid() {
		log.Println("Token Expired.")
	}

	if err != nil {
		log.Println("Could not read token from file")
	}

	// generate new token if token file does not exists
	// or token got expired
	if err != nil || !token.Valid() {
		log.Println("Generating new token...")
		token = getTokenFromWeb(config)

		log.Println("Saving generated token to file")
		saveToken(tokenFile, token)
	}

	return config.Client(context.Background(), token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser and then input the authorization code here: %v\n", authURL)

	fmt.Printf("auth code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	return token
}

// Retrieves toke from a local file
func getTokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

// Saves a token to a file path
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credentials file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func main() {
	fmt.Println("GDOWN CLI")

	ctx := context.Background()

	_ = godotenv.Load()
	// if err != nil {
	// 	fmt.Printf("Could not load .env file: %v\n", err)
	// }

	CLIENT_ID := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	CLIENT_SECRET := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	// Check if Client ID and Secret are added
	if CLIENT_ID == "" {
		log.Fatal("GOOGLE_OAUTH_CLIENT_ID environment variable not present")
	}
	if CLIENT_SECRET == "" {
		log.Fatal("GOOGLE_OAUTH_CLIENT_SECRET environment variable not present")
	}

	// fmt.Printf("CLIENT_ID: %v\n", CLIENT_ID)
	// fmt.Printf("CLIENT_SECRET: %v\n", CLIENT_SECRET)

	config := &oauth2.Config{
		RedirectURL:  "http://localhost:8000/auth/google/callback",
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
		Scopes: []string{
			// See information about your Google Drive files
			"https://www.googleapis.com/auth/drive.metadata.readonly",
			// See and download all your Google Drive Files
			"https://www.googleapis.com/auth/drive.readonly",
		},
		Endpoint: google.Endpoint,
	}

	client := getClient(config)

	driveService, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// file1.txt
	fileId := "1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ"
	file, err := driveService.Files.Get(fileId).Fields("id", "name", "mimeType").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve file: %v", err)
	}

	fmt.Printf("fileId: %+v\n", file.Id)
	fmt.Printf("fileName: %+v\n", file.Name)
	fmt.Printf("fileMimeType: %+v\n", file.MimeType)
}
