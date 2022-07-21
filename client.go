package gdutils

import (
	"google.golang.org/api/drive/v3"
	"io"
	"net/http"
)

type ServerClient drive.Service

func (s *ServerClient) GetDrive(DriveId string) (*Drive, error) {
	d, e := s.Drives.Get(DriveId).Do()
	if e != nil {
		return nil, e
	}
	return &Drive{
		Drive: d,
		s:     s,
	}, nil
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
	return &File{
		File: f,
		s:    s,
	}, e
}

// Download 下载文件
func (s *ServerClient) Download(FileId string, opt *DownloadOpt) (*http.Response, error) {
	req := s.Files.Get(FileId).
		SupportsAllDrives(true)
	if opt != nil {
		header := req.Header()
		if opt.Range != "" {
			header.Set("Range", opt.Range)
		}
	}
	return req.Download()
}

// Delete 删除文件
func (s *ServerClient) Delete(FileId string) error {
	return s.Files.Delete(FileId).
		SupportsAllDrives(true).
		Do()
}
