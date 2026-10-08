package image

import "io"

type Upload struct {
	Filename string
	Size     int64
	File     io.ReadSeeker
}
