package gdutils

import (
	"google.golang.org/api/drive/v3"
	"net/http"
)

type File struct {
	*drive.File
	s *ServerClient
}

func (f *File) Refresh() error {
	t, e := f.s.GetFile(f.Id)
	if e != nil {
		return e
	}
	f.File = t.File
	return nil
}

func (f *File) Rename(name string) error {
	n, e := f.s.Files.Update(f.Id, &drive.File{Name: name}).
		SupportsAllDrives(true).Do()
	if e != nil {
		return e
	}
	f.File = n
	return nil
}

func (f *File) Download(opt *DownloadOpt) (*http.Response, error) {
	return f.s.Download(f.Id, opt)
}

func (f *File) Delete() error {
	return f.s.Delete(f.Id)
}
