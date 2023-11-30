package gdown

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/schollz/progressbar/v3"
	"google.golang.org/api/drive/v3"
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

/* Get File Details */
func GetFile(driveClient *drive.Service, fileId string) (*File, error) {
	/* Get File details as response */
	fileRes, err := driveClient.Files.Get(fileId).Fields("id", "name", "mimeType", "size").Do()
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
func GetFolder(driveClient *drive.Service, folderId string) (*Folder, error) {

	/* Get Folder details as response */
	folderRes, err := driveClient.Files.Get(folderId).Fields("id", "name", "mimeType").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve folder: %v", err)
	}

	/* Fetching Files List */
	filesListRes, err := driveClient.Files.List().Q(fmt.Sprintf("\"%s\" in parents", folderId)).Fields("files(id, name, mimeType, size)").OrderBy("name").Do()
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
			folder, err := GetFolder(driveClient, file.Id)
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
func DownloadFile(driveClient *drive.Service, file *File, folderPath string) error {
	/* Send Download File Request */
	fileRes, err := driveClient.Files.Get(file.Id).Download()
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
func DownloadFilesInFolder(driveClient *drive.Service, folder *Folder) {
	/* Derive Folder path */
	folderPath := fmt.Sprintf("%s/%s", DownloadFolderPath, folder.Name)
	/* Download all the files one at a time */
	for _, file := range folder.Files {
		err := DownloadFile(driveClient, &file, folderPath)
		if err != nil {
			log.Printf("Error Downloading file: %s, %v", file.Name, err)
		}
	}
}
