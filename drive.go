package gdutils

import (
	"google.golang.org/api/drive/v3"
	"io"
)

type Drive struct {
	s  *ServerClient
	id string
}

// GetFiles 获取文件与文件夹列表
func (d *Drive) GetFiles(FolderId string) (folders []*Folder, files []*File, err error) {
	if FolderId == "" {
		FolderId = d.id
	}
	FileList, err := d.s.Files.List().
		Corpora("drive").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		DriveId(d.id).
		Q("'" + FolderId + "' in parents").
		Do()
	if err != nil {
		return nil, nil, err
	}

	for i, v := range FileList.Files {
		if v.MimeType == folderType {
			folders = append(folders, &Folder{
				File: FileList.Files[i],
				d:    d,
			})
		} else {
			files = append(files, &File{
				File: FileList.Files[i],
				s:    d.s,
			})
		}
	}
	return
}

func (d *Drive) Upload(FileName string, FolderId string, Reader io.Reader) (*File, error) {
	if FolderId == "" {
		FolderId = d.id
	}
	return d.s.Upload(FileName, FolderId, Reader)
}

func (d *Drive) CreateFolder(FolderId string, FolderName string) (*Folder, error) {
	if FolderId == "" {
		FolderId = d.id
	}
	f, e := d.s.Files.Create(&drive.File{
		Name:    FolderName,
		DriveId: d.id, Parents: []string{FolderId},
		MimeType: folderType,
	}).SupportsAllDrives(true).
		SupportsTeamDrives(true).Do()
	return &Folder{
		File: f,
		d:    d,
	}, e
}
