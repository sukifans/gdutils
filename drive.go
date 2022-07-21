package gdutils

import (
	"google.golang.org/api/drive/v3"
	"io"
)

type Drive struct {
	*drive.Drive
	s *ServerClient
}

func (d *Drive) Refresh() error {
	t, e := d.s.GetDrive(d.Id)
	if e != nil {
		return e
	}
	d.Drive = t.Drive
	return nil
}

func (d *Drive) GetFolder(FolderId string) (*Folder, error) {
	f, e := d.s.GetFile(FolderId)
	if e != nil {
		return nil, e
	}
	if f.MimeType != folderType {
		return nil, ErrNotFolder
	}
	return &Folder{
		File: f,
		d:    d,
	}, nil
}

// ListFiles 获取文件与文件夹列表
func (d *Drive) ListFiles(FolderId string) (folders []*Folder, files []*File, err error) {
	if FolderId == "" {
		FolderId = d.Id
	}
	FileList, err := d.s.Files.List().
		Corpora("drive").
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		DriveId(d.Id).
		Q("'" + FolderId + "' in parents").
		Do()
	if err != nil {
		return nil, nil, err
	}

	for i, v := range FileList.Files {
		f := File{
			File: FileList.Files[i],
			s:    d.s,
		}
		if v.MimeType == folderType {
			folders = append(folders, &Folder{
				File: &f,
				d:    d,
			})
		} else {
			files = append(files, &f)
		}
	}
	return
}

func (d *Drive) Upload(FileName string, FolderId string, Reader io.Reader) (*File, error) {
	if FolderId == "" {
		FolderId = d.Id
	}
	return d.s.Upload(FileName, FolderId, Reader)
}

func (d *Drive) CreateFolder(FolderId string, FolderName string) (*Folder, error) {
	if FolderId == "" {
		FolderId = d.Id
	}
	f, e := d.s.Files.Create(&drive.File{
		Name:    FolderName,
		DriveId: d.Id, Parents: []string{FolderId},
		MimeType: folderType,
	}).SupportsAllDrives(true).
		SupportsTeamDrives(true).Do()
	return &Folder{
		File: &File{
			File: f,
			s:    d.s,
		},
		d: d,
	}, e
}
