package gdutils

import (
	"google.golang.org/api/drive/v3"
	"io"
)

type ServerClient struct {
	Server   *drive.Service
	FoldType string
}

//获取为共享drive的id
//不包含自己的drive

func (s *ServerClient) GetDriveList(PageSize int64) (*drive.DriveList, error) {
	return s.Server.Drives.List().PageSize(PageSize).Do()
}

// GetFiles 获取文件列表
func (s *ServerClient) GetFiles(DriveId string) ([]*drive.File, error) {
	FileList, err := s.Server.Files.List().Corpora("drive").IncludeItemsFromAllDrives(true).SupportsAllDrives(true).DriveId(DriveId).Do()
	if err != nil {

		return nil, err

	}

	return FileList.Files, err
}

// GetFolders 获取文件夹列表
func (s *ServerClient) GetFolders(DriveId string) ([]*drive.File, error) {
	FileList, err := s.Server.Files.List().
		Corpora("drive").Q("mimeType='application/vnd.google-apps.folder'").
		OrderBy("createdTime desc").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		DriveId(DriveId).Do()
	if err != nil {
		return nil, err
	}

	return FileList.Files, err
}

// Upload 上传文件
func (s *ServerClient) Upload(FileName string, FolderId string, Reader io.Reader) (*drive.File, error) {
	return s.Server.Files.Create(&drive.File{
		Name:    FileName,
		Parents: []string{FolderId}},
	).Media(Reader).SupportsAllDrives(true).Do()
}

func (s ServerClient) CreateFolder(FolderName string, DriveId string) (*drive.File, error) {
	return s.Server.Files.Create(&drive.File{
		Name:    FolderName,
		DriveId: DriveId, Parents: []string{DriveId},
		MimeType: s.FoldType,
	}).SupportsAllDrives(true).SupportsTeamDrives(true).Do()
}
