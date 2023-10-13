package archive

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// Extract tarball to dir
func UnTar(tarball string, dir string) (err error) {
	reader, err := os.Open(tarball)
	if err != nil {
		return
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)

	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		path := filepath.Join(dir, hdr.Name)
		mode := hdr.FileInfo().Mode()
		switch hdr.Typeflag {
		case tar.TypeReg:
			file, err := os.Create(path)
			if err != nil {
				return err
			}
			defer file.Close()

			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
		case tar.TypeDir:
			err = os.MkdirAll(path, mode)
			if err != nil {
				return err
			}
		case tar.TypeLink:
			err = os.Link(filepath.Join(dir, filepath.Clean(hdr.Linkname)), path)
			if err != nil {
				return err
			}
		case tar.TypeSymlink:
			err = os.Symlink(filepath.Clean(hdr.Linkname), path)
			if err != nil {
				return err
			}
		case tar.TypeXGlobalHeader, tar.TypeXHeader:
			continue

		}

	}

	return
}
