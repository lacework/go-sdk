package lwcomponent

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/internal/file"
	"github.com/pkg/errors"
)

type StageConstructor func(name, artifactUrl string, size int64) (stage Stager, err error)

type Stager interface {
	Close()

	Commit(string) error

	Directory() string

	Download() error

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

	dirEntries, err := os.ReadDir(s.dir)
	if err != nil {
		return
	}

	for _, entry := range dirEntries {
		src := filepath.Join(s.dir, entry.Name())
		dst := filepath.Join(targetDir, entry.Name())

		err = os.Rename(src, dst)
		if err != nil {
			return err
		}
	}

	return
}

func (s *stageTarGz) Directory() string {
	return s.dir
}

func (s *stageTarGz) Download() (err error) {
	fileName := filepath.Base(s.artifactUrl.Path)

	path := filepath.Join(s.dir, fileName)

	err = DownloadFile(path, s.artifactUrl.String(), s.size, 0)
	if err != nil {
		return
	}

	return
}

func (s *stageTarGz) Signature() (sig []byte, err error) {
	_, err = os.Stat(s.dir)
	if os.IsNotExist(err) {
		err = errors.New("component not staged")
		return
	}

	path := filepath.Join(s.dir, SignatureFile)
	if !file.FileExists(path) {
		err = errors.New("missing .signature file")
		return
	}

	sig, err = os.ReadFile(path)
	if err != nil {
		return
	}

	return
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

func (s *stageTarGz) Validate() (err error) {
	data, err := os.ReadFile(filepath.Join(s.dir, VersionFile))
	if err != nil {
		return
	}

	version := string(data)

	_, err = semver.NewVersion(strings.TrimSpace(version))
	if err != nil {
		return errors.Errorf("invalid staged semantic version '%s' for component '%s'", version, s.name)
	}

	if !file.FileExists(filepath.Join(s.dir, SignatureFile)) {
		return errors.Errorf("missing file '%s'", s.name)
	}

	if !file.FileExists(filepath.Join(s.dir, s.name)) {
		return errors.Errorf("missing file '%s'", s.name)
	}

	return
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

func unTar(tarball string, dir string) (err error) {
	reader, err := os.Open(tarball)
	if err != nil {
		return
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
			if err = os.MkdirAll(path, info.Mode()); err != nil {
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

	return
}
