package gdutils

import (
	"io"
)

type Folder struct {
	*File
	d *Drive
}

func (f *Folder) Upload(FileName string, Reader io.Reader) (*File, error) {
	return f.d.Upload(FileName, f.Id, Reader)
}

func (f *Folder) CreateFolder(Name string) (*Folder, error) {
	return f.d.CreateFolder(f.Id, Name)
}
