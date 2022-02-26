package gdutils

import (
	"google.golang.org/api/drive/v3"
	"io"
)

type ServerClient drive.Service

func (s *ServerClient) GetDrive(DriveId string) *Drive {
	return &Drive{
		s:  s,
		id: DriveId,
	}
}

// GetFile
// Deprecated: todo
func (s *ServerClient) GetFile(FileId string) *File {
	return nil
}

// GetDriveList 获取为共享drive的id
//不包含自己的drive
func (s *ServerClient) GetDriveList(PageSize int64) (*drive.DriveList, error) {
	return s.Drives.List().PageSize(PageSize).Do()
}

// Upload 上传文件
func (s *ServerClient) Upload(FileName string, FolderId string, Reader io.Reader) (*File, error) {
	f, e := s.Files.Create(&drive.File{
		Name:    FileName,
		Parents: []string{FolderId}},
	).Media(Reader).SupportsAllDrives(true).Do()
	return (*File)(f), e
}
