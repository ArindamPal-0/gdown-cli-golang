package gdown

import (
	"fmt"
	"io"
	"os"
	"path"

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

func ListItemDetails(id string) error {
	/* Get Drive Client */
	driveClient, err := NewDriveClientUsingServiceAccount()
	if err != nil {
		return fmt.Errorf("could not get Drive Client.\n%v", err)
	}

	/* Get File */
	file, err := GetFile(driveClient, id)
	if err != nil {
		return fmt.Errorf("could not get File.\n%v", err)
	}

	/* List File Details */
	if file.MimeType != folderMimeType {
		ListFileDetails(file)
		return nil
	}

	/* Get Folder Details */
	folder, err := GetFolder(driveClient, id)
	if err != nil {
		return fmt.Errorf("could not get Folder.\n%v", err)
	}

	/* List Folder Details */
	ListFolderDetails(folder)

	return nil
}

func DownloadItem(id string, alsoListDetails bool) error {
	/* Get Drive Client */
	driveClient, err := NewDriveClientUsingServiceAccount()
	if err != nil {
		return fmt.Errorf("could not get Drive Client.\n%v", err)
	}

	/* Get File */
	file, err := GetFile(driveClient, id)
	if err != nil {
		return fmt.Errorf("could not get File.\n%v", err)
	}

	/* Download File (item is a file) */
	if file.MimeType != folderMimeType {
		err := downloadFile(driveClient, file, downloadFolderPath, alsoListDetails)
		if err != nil {
			return fmt.Errorf("could not download item.\n%v", err)
		}

		return nil
	}

	/* Get Folder */
	folder, err := GetFolder(driveClient, id)
	if err != nil {
		return fmt.Errorf("could not get Folder.\n%v", err)
	}

	/* Download Folder (item is a folder) */
	err = downloadFilesInFolder(driveClient, folder, downloadFolderPath, alsoListDetails)
	if err != nil {
		return fmt.Errorf("could not download item.\n%v", err)
	}

	return nil
}

/* Get File Details */
func GetFile(driveClient *drive.Service, fileId string) (*File, error) {
	/* Get File details as response */
	fileRes, err := driveClient.Files.Get(fileId).Fields("id", "name", "mimeType", "size").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file.\n%v", err)
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
		return nil, fmt.Errorf("unable to retrieve folder.\n%v", err)
	}

	/* Fetching Files List */
	filesListRes, err := driveClient.Files.List().Q(fmt.Sprintf("\"%s\" in parents", folderId)).Fields("files(id, name, mimeType, size)").OrderBy("name").Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve file list.\n%v", err)
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

/* List File Details */
func ListFileDetails(file *File) {
	fmt.Printf("Id: %s\n", file.Id)
	fmt.Printf("Name: %s\n", file.Name)
	fmt.Printf("MimeType: %s\n", file.MimeType)
	fmt.Printf("Size: %d\n", file.Size)
}

/* List Folder Details */
func ListFolderDetails(folder *Folder) {
	fmt.Printf("Id: 📁 %s\n", folder.Id)
	fmt.Printf("Name: %s\n", folder.Name)

	for _, file := range folder.Files {
		ListFileDetails(&file)
	}
}

var downloadFolderPath string = "downloads"

/* Set download folder path */
func SetDownloadFolderPath(path string) {
	downloadFolderPath = path
}

/* Download File to default download folder path */
func DownloadFile(driveClient *drive.Service, file *File) error {
	return downloadFile(driveClient, file, downloadFolderPath, false)
}

/* List details and Download File to default download folder path */
func ListDetailsAndDownloadFile(driveClient *drive.Service, file *File) error {
	return downloadFile(driveClient, file, downloadFolderPath, true)
}

/* Download folder to default download path */
func DownloadFilesInFolder(driveClient *drive.Service, folder *Folder) error {
	return downloadFilesInFolder(driveClient, folder, downloadFolderPath, false)
}

/* List folder details and Download Folder to default download path */
func ListDetailsAndDownloadfilesInFolder(driveClient *drive.Service, folder *Folder) error {
	return downloadFilesInFolder(driveClient, folder, downloadFolderPath, true)
}

/* Download a Single file to folderPath */
func downloadFile(driveClient *drive.Service, file *File, folderPath string, alsoListDetails bool) error {
	/* Send Download File Request */
	fileRes, err := driveClient.Files.Get(file.Id).Download()
	if err != nil {
		return fmt.Errorf("could not download file, %v", err)
	}
	defer fileRes.Body.Close()

	/* Create Folder */
	err = os.MkdirAll(folderPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create download directory: %s.\n%v", downloadFolderPath, err)
	}

	filePath := path.Join(folderPath, file.Name)

	/* Open File Handle */
	fileHandle, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("could not open and/or create file: %s.\n%v", filePath, err)
	}
	defer fileHandle.Close()

	/* List file details */
	if alsoListDetails {
		ListFileDetails(file)
	}
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
func downloadFilesInFolder(driveClient *drive.Service, folder *Folder, downloadPath string, alsoListDetails bool) error {
	/* Derive Folder path */
	folderPath := path.Join(downloadPath, folder.Name)

	/* List folder details */
	if alsoListDetails {
		fmt.Printf("Id: 📁 %s\n", folder.Id)
		fmt.Printf("Name: %s\n", folder.Name)
	}

	/* Download all the files one at a time */
	for _, file := range folder.Files {
		err := downloadFile(driveClient, &file, folderPath, alsoListDetails)
		if err != nil {
			return fmt.Errorf("error Downloading file: %s.\n%v", file.Name, err)
		}
	}

	return nil
}
