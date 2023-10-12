package archive

import (
	"compress/gzip"
	"io"
	"os"
)

// Inflate GZip file.
//
// Writes decompressed data to target path.
func Gunzip(source string, target string) (err error) {
	reader, err := os.Open(source)
	if err != nil {
		return
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return
	}
	defer archive.Close()

	writer, err := os.Create(target)
	if err != nil {
		return
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)

	return
}
