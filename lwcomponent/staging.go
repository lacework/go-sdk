package lwcomponent

import (
	"archive/tar"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/internal/file"
	dircopy "github.com/otiai10/copy"
	"github.com/pkg/errors"
)

type StageConstructor func(name, artifactUrl string, size int64) (stage Stager, err error)

type Stager interface {
	Close()

	Commit(string) error

	Directory() string

	Download(progressClosure func(filepath string, sizeB int64)) error

	Filename() string

	Signature() (sig []byte, err error)

	Unpack() error

	Validate() error
}

type stageTarGz struct {
	artifactUrl *url.URL
	dir         string
	name        string
	size        int64
}

func NewStageTarGz(name, artifactUrl string, size int64) (stage Stager, err error) {
	dir, err := os.MkdirTemp("", "cdk-component-stage-tar-gz-")
	if err != nil {
		return
	}

	_url, err := url.Parse(artifactUrl)
	if err != nil {
		os.RemoveAll(dir)
		return
	}

	stage = &stageTarGz{artifactUrl: _url, dir: dir, name: name, size: size}

	return
}

func (s *stageTarGz) Close() {
	os.RemoveAll(s.dir)
}

func (s *stageTarGz) Commit(targetDir string) (err error) {
	_, err = os.Stat(s.dir)
	if os.IsNotExist(err) {
		err = errors.New("component not staged")
		return
	}

	_, err = os.Stat(targetDir)
	if os.IsNotExist(err) {
		err = errors.New("target install directory doesn't exist")
		return
	}

	if err = dircopy.Copy(s.dir, targetDir); err != nil {
		return
	}

	return
}

func (s *stageTarGz) Directory() string {
	return s.dir
}

func (s *stageTarGz) Filename() string {
	return filepath.Base(s.artifactUrl.Path)
}

func (s *stageTarGz) Download(progressClosure func(filepath string, sizeB int64)) error {
	fileName := filepath.Base(s.artifactUrl.Path)

	path := filepath.Join(s.dir, fileName)

	if _, err := os.Create(path); err != nil {
		return err
	}

	go progressClosure(path, s.size*1024)

	return DownloadFile(path, s.artifactUrl.String())
}

func (s *stageTarGz) Signature() ([]byte, error) {
	_, err := os.Stat(s.dir)
	if os.IsNotExist(err) {
		return nil, errors.New("component not staged")
	}

	path := filepath.Join(s.dir, SignatureFile)
	if !file.FileExists(path) {
		return nil, errors.New("missing .signature file")
	}

	sig, err := os.ReadFile(path)
	if err != nil {
		return sig, err
	}

	// Artifact signature may or may not be b64encoded
	decoded_sig, err := base64.StdEncoding.DecodeString(string(sig))
	if err == nil {
		return decoded_sig, nil
	}
	return sig, nil
}

func (s *stageTarGz) Unpack() (err error) {
	fileName := filepath.Base(s.artifactUrl.Path)

	gzFile := filepath.Join(s.dir, fileName)
	tarball := filepath.Join(s.dir, strings.TrimRight(fileName, ".gz"))

	err = gunzip(gzFile, tarball)
	if err != nil {
		return
	}

	err = unTar(tarball, s.dir)
	if err != nil {
		return
	}

	os.Remove(gzFile)
	os.Remove(tarball)

	return nil
}

func (s *stageTarGz) Validate() error {
	data, err := os.ReadFile(filepath.Join(s.dir, VersionFile))
	if err != nil {
		return err
	}

	version := string(data)

	_, err = semver.NewVersion(strings.TrimSpace(version))
	if err != nil {
		return errors.Errorf("invalid staged semantic version '%s' for component '%s'", version, s.name)
	}

	if !file.FileExists(filepath.Join(s.dir, SignatureFile)) {
		return errors.Errorf("missing file '%s'", s.name)
	}

	path := filepath.Join(s.dir, s.name)

	if operatingSystem == "windows" {
		path = fmt.Sprintf("%s.exe", path)
	}

	if !file.FileExists(path) {
		return errors.Errorf("missing file '%s'", path)
	}

	return nil
}

// Inflate GZip file.
//
// Writes decompressed data to target path.
func gunzip(source string, target string) (err error) {
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

func unTar(tarball string, dir string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()

	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		path := filepath.Join(dir, header.Name)

		info := header.FileInfo()
		if info.IsDir() {
			if err := os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}

	return nil
}
