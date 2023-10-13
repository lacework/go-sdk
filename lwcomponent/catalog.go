package lwcomponent

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"aead.dev/minisign"
	"github.com/Masterminds/semver"
	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/cache"
	"github.com/pkg/errors"
)

const (
	componentCacheDir string = "components"
	featureFlag       string = "PUBLIC.cdk.v1"
	operatingSystem   string = runtime.GOOS
	architecture      string = runtime.GOARCH
)

func CatalogV1Enabled(client *api.Client) bool {
	response, err := client.V2.FeatureFlags.GetFeatureFlagsMatchingPrefix(featureFlag)
	if err != nil {
		return false
	}

	return len(response.Data.Flags) >= 1
}

// Returns the local directory that Components will be stored in.
func CatalogCacheDir() (string, error) {
	cacheDir, err := cache.CacheDir()
	if err != nil {
		return "", errors.Wrap(err, "unable to locate components directory")
	}

	return filepath.Join(cacheDir, componentCacheDir), nil
}

type Catalog struct {
	client *api.Client

	Components       map[string]CDKComponent
	stageConstructor StageConstructor
}

func (c *Catalog) ComponentCount() int {
	return len(c.Components)
}

func (c *Catalog) Cache() {
	// @jon-stewart: TODO: cache catalog
}

// Return a CDKComponent that is present on the host.
func (c *Catalog) GetComponent(name string) (*CDKComponent, error) {
	component, exists := c.Components[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("component %s not found", name))
	}

	return &component, nil
}

func (c *Catalog) ListComponentVersions(component *CDKComponent) (versions []*semver.Version, err error) {
	if component.apiInfo == nil {
		return
	}

	response, err := c.client.V2.Components.ListComponentVersions(component.apiInfo.Id(), operatingSystem, architecture)
	if err != nil {
		return nil, err
	}

	rawVersions := response.Data[0].Versions

	versions = make([]*semver.Version, len(rawVersions))
	for idx, v := range rawVersions {
		ver, err := semver.NewVersion(v)
		if err != nil {
			return nil, err
		}

		versions[idx] = ver
	}

	return versions, nil
}

func (c *Catalog) PrintComponents() [][]string {
	result := [][]string{}

	for _, component := range c.Components {
		result = append(result, component.PrintSummary())
	}

	return result
}

func (c *Catalog) Stage(component *CDKComponent, version string) (stageClose func(), err error) {
	var (
		semv *semver.Version
	)

	stageClose = func() {}

	if version == "" {
		semv = component.apiInfo.LatestVersion()
	} else {
		semv, err = semver.NewVersion(version)
		if err != nil {
			return
		}
	}

	if component.hostInfo != nil {
		var installedVersion *semver.Version

		installedVersion, err = component.hostInfo.Version()
		if err != nil {
			return
		}

		if installedVersion.Equal(semv) {
			err = errors.Errorf("version '%s' already installed", semv.String())
			return
		}
	}

	response, err := c.client.V2.Components.FetchComponentArtifact(
		component.apiInfo.Id(),
		operatingSystem,
		architecture,
		semv.String())
	if err != nil {
		return
	}

	stage, err := c.stageConstructor(component.Name, response.Data[0].ArtifactUrl)
	if err != nil {
		return
	}

	if err = stage.Download(); err != nil {
		stage.Close()
		return
	}

	if err = stage.Unpack(); err != nil {
		stage.Close()
		return
	}

	if err = stage.Validate(); err != nil {
		stage.Close()
		return
	}

	component.stage = stage
	stageClose = stage.Close

	return
}

func (c *Catalog) Verify(component *CDKComponent) (err error) {
	data, err := os.ReadFile(filepath.Join(component.stage.Directory(), component.Name))
	if err != nil {
		return
	}

	sig, err := component.stage.Signature()
	if err != nil {
		return
	}

	rootPublicKey := minisign.PublicKey{}
	if err := rootPublicKey.UnmarshalText([]byte(publicKey)); err != nil {
		return errors.Wrap(err, "unable to load root public key")
	}

	return verifySignature(rootPublicKey, data, sig)
}

func (c *Catalog) Install(component *CDKComponent) (err error) {
	if component.stage == nil {
		return errors.Errorf("component '%s' not staged", component.Name)
	}

	componentDir, err := componentDirectory(component.Name)
	if err != nil {
		return
	}

	err = os.MkdirAll(componentDir, os.ModePerm)
	if err != nil {
		return
	}

	err = component.stage.Commit(componentDir)
	if err != nil {
		return
	}

	component.hostInfo = NewHostInfo(componentDir)

	return
}

// Delete a CDKComponent
//
// Remove the Component install directory and all sub-directory.  This function will not return an
// error if the Component is not installed.
func (c *Catalog) Delete(component *CDKComponent) (err error) {
	componentDir, err := componentDirectory(component.Name)
	if err != nil {
		return
	}

	_, err = os.Stat(componentDir)
	if err != nil {
		return errors.Errorf("component not installed. Try running 'lacework component install %s'", component.Name)
	}

	os.RemoveAll(componentDir)

	return
}

func NewCatalog(client *api.Client, stageConstructor StageConstructor) (*Catalog, error) {
	if stageConstructor == nil {
		return nil, errors.New("nil Catalog StageConstructor")
	}

	localComponents, err := loadLocalComponents()
	if err != nil {
		return nil, err
	}

	response, err := client.V2.Components.ListComponents(operatingSystem, architecture)
	if err != nil {
		return nil, err
	}

	rawComponents := response.Data[0].Components

	cdkComponents := make(map[string]CDKComponent, len(rawComponents)+len(localComponents))

	for _, c := range rawComponents {
		ver, err := semver.NewVersion(c.Version)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("component '%s' version '%s'", c.Name, c.Version))
		}

		api := NewAPIInfo(c.Id, c.Name, ver, c.Description, c.Size)

		host, found := localComponents[c.Name]
		if found {
			delete(localComponents, c.Name)
		}

		component := NewCDKComponent(c.Name, Type(c.ComponentType), api, host)

		cdkComponents[c.Name] = component
	}

	for _, localHost := range localComponents {
		// @jon-stewart: TODO: local specification
		cdkComponents[localHost.Name()] = NewCDKComponent(localHost.Name(), BinaryType, nil, localHost)
	}

	return &Catalog{client, cdkComponents, stageConstructor}, nil
}

func loadLocalComponents() (local map[string]HostInfo, err error) {
	cacheDir, err := CatalogCacheDir()
	if err != nil {
		return
	}

	subDir, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}

	local = make(map[string]HostInfo, len(subDir))

	for _, file := range subDir {
		if !file.IsDir() {
			continue
		}

		local[file.Name()] = NewHostInfo(filepath.Join(cacheDir, file.Name()))
	}

	return
}

// func isDevelopmentComponent(path string, name string) bool {
// 	return file.FileExists(filepath.Join(path, name, DevelopmentFile))
// }

// Returns the directory that the component executable and configuration is stored in.
func componentDirectory(componentName string) (string, error) {
	dir, err := CatalogCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, componentName), nil
}
