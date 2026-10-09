package image

import "io"

type Upload struct {
	Size int64
	File io.ReadSeeker
}
