package gdutils

import (
	"google.golang.org/api/drive/v3"
	"net/http"
)

type File struct {
	*drive.File
	s *ServerClient
}

func (f *File) Rename(name string) error {
	name, f.Name = f.Name, name
	n, e := f.s.Files.Update(f.Id, f.File).Do()
	if e != nil {
		f.Name = name
		return e
	}
	f.File = n
	return nil
}

func (f *File) Download() (*http.Response, error) {
	return f.s.Download(f.Id)
}

func (f *File) Delete() error {
	return f.s.Delete(f.Id)
}
