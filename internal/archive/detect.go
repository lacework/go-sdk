package archive

import (
	"os"
	"github.com/gabriel-vasile/mimetype"
	"path/filepath"
	"strings"
)

func detectFileType(file string) (mimeType string, err error) {

	fDescriptor, err := os.Open(file)
	if err != nil {
		return
	}
	defer fDescriptor.Close()
	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)

	_, err = fDescriptor.Read(head)
	if err != nil {
		return
	}

	mtype := mimetype.Detect(head)
	mimeType = mtype.String()

	return
}

func FileIsGZ(file string) (isGZ bool, err error) {
	mtype, err := detectFileType(file)
	if err != nil {
		return
	}
	isGZ = mtype == "application/gzip"
	return
}

func FileIsTar(file string) (isTar bool, err error) {
	mtype, err := detectFileType(file)
	if err != nil {
		return
	}
	isTar = mtype == "application/x-tar"
	return
}

func DetectTGZAndUnpack(filePath string, targetDir string) (err error) {
	// detect if file is a tar gz and extract to targetDir
	fileName := filepath.Base(filePath)
	baseFileName := strings.ReplaceAll(fileName, ".tar.gz", "")
	baseFileName = strings.ReplaceAll(baseFileName, ".tgz", "")
	unpackDir, err := os.MkdirTemp("", "temp-cdk-unpack-archive-")
	if err != nil {
		return
	}
	defer os.RemoveAll(unpackDir)

	tarFile := filepath.Join(unpackDir, baseFileName + ".tar")
	isGZ, err := FileIsGZ(filePath)
	if err != nil {
		return
	}

	if isGZ {
		if err := Gunzip(filePath, tarFile); err != nil {
			return err
		}
	} else {

		return
	}
	
	isTar, err := FileIsTar(tarFile)
	if err != nil {
		return
	} 
	
	if isTar {
		if err = UnTar(tarFile, targetDir); err != nil {
			return
		}
	} else {
		return
	}
	return
}