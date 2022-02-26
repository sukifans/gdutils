package gdutils

import (
	"google.golang.org/api/drive/v3"
	"io"
)

type Folder struct {
	*drive.File
	d *Drive
}

func (f *Folder) Upload(FileName string, Reader io.Reader) (*File, error) {
	return f.d.Upload(FileName, f.Id, Reader)
}

func (f *Folder) CreateFolder(Name string) (*Folder, error) {
	return f.d.CreateFolder(f.Id, Name)
}
