package main

import (
	"fmt"
	"log"

	"github.com/arindampal-0/gdown-cli-golang/internals/gdown"
)

func main() {
	fmt.Println("GDOWN CLI")

	/* Get Drive client using service account */
	driveClient, err := gdown.NewDriveClientUsingServiceAccount()
	if err != nil {
		log.Fatalf("Error creating Drive Service.\n%v", err)
	}

	/* Fetching file details */
	// file1.txt
	fileId := "1NuuL9qNo5BJYnfNqN_lxBOUN0P-AociQ"
	file, err := gdown.GetFile(driveClient, fileId)
	if err != nil {
		log.Fatalf("Error getting File.\n%v", err)
	}

	fmt.Printf("file: %s\n", gdown.Prettify(file))

	/* Download a single file */
	err = gdown.DownloadFile(driveClient, file, gdown.DownloadFolderPath)
	if err != nil {
		log.Fatalf("Error downloading file.\n%v", err)
	}

	/* Fetching folder details */
	// gdown folder
	folderId := "1SVHxav6Y5LoYbdgfx2MSsdYlT74RTjej"
	folder, err := gdown.GetFolder(driveClient, folderId)
	if err != nil {
		log.Fatalf("Error getting Folder.\n%v", err)
	}

	fmt.Printf("folder: %s\n", gdown.Prettify(folder))

	/* Download all files in the folder */
	gdown.DownloadFilesInFolder(driveClient, folder)
} //
