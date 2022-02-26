package gdutils

import "google.golang.org/api/drive/v3"

type File drive.File

// Remove
// Deprecated: todo
func (f *File) Remove() error {
	return nil
}
