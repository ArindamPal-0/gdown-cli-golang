package gdown

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var configFolderPath string = (func() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Could not get User Config Directory.\n%v", err)
	}

	return path.Join(configDir, "gdown")
})()

/* Drive Client using Service Account */
func NewDriveClientUsingServiceAccount() (*drive.Service, error) {
	ctx := context.Background()

	/* interface for reading credentials.json file */
	type Credential struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
	}

	/* Reading credentials.json file */
	var credentialsPath string = path.Join(configFolderPath, "service-account", "credentials.json")
	content, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %v.\n%v", credentialsPath, err)
	}

	/* Converting json to Credentials */
	var cred Credential
	err = json.Unmarshal(content, &cred)
	if err != nil {
		return nil, fmt.Errorf("could not convert json to Credentials interface.\n%v", err)
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
	driveClient, err := drive.NewService(ctx, option.WithHTTPClient(client))

	/* Creating drive service using the credentials file and scopes */
	// driveService, err := drive.NewService(context.Background(), option.WithCredentialsFile("credentials.json"), option.WithScopes(
	// 	"https://www.googleapis.com/auth/drive.metadata.readonly",
	// 	"https://www.googleapis.com/auth/drive.readonly",
	// ))

	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client.\n%v", err)
	}

	return driveClient, nil
}

/* Drive Client using Oauth2 */
func NewDriveClientUsingOauth2() {}
