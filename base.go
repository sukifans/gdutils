package gdutils

import "errors"

const folderType = "application/vnd.google-apps.folder"

var (
	ErrNotFolder           = errors.New("file is not a folder")
	ErrOperationNotSupport = errors.New("operation not support")
)
