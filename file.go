package gdutils

import (
	"google.golang.org/api/drive/v3"
	"net/http"
)

type File struct {
	*drive.File
	s *ServerClient
}

func (f *File) Download() (*http.Response, error) {
	return f.s.Download(f.Id)
}

func (f *File) Delete() error {
	return f.s.Delete(f.Id)
}
