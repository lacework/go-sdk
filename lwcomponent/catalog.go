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

	path := filepath.Join(cacheDir, componentCacheDir)

	if _, err = os.Stat(path); err != nil {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return "", err
		}
	}

	return path, nil
}

type Catalog struct {
	client *api.Client

	Components       map[string]CDKComponent
	stageConstructor StageConstructor
}

func (c *Catalog) ComponentCount() int {
	return len(c.Components)
}

func (c *Catalog) Persist() {
	// @jon-stewart: TODO: store catalog on disk
}

// Return a CDKComponent that is present on the host.
func (c *Catalog) GetComponent(name string) (*CDKComponent, error) {
	component, exists := c.Components[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("component %s not found", name))
	}

	return &component, nil
}

func (c *Catalog) ListComponentVersions(component *CDKComponent) (versions []*semver.Version) {
	if component.apiInfo == nil {
		return
	}

	versions = component.apiInfo.AllVersions()

	return
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

	if len(response.Data) == 0 {
		err = errors.New("Invalid API response")
		return
	}

	data := response.Data[0]

	component.InstallMessage = data.InstallMessage
	component.UpdateMessage = data.UpdateMessage

	stage, err := c.stageConstructor(component.Name, data.ArtifactUrl)
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

	var rawComponents []api.LatestComponentVersion

	if len(response.Data) > 0 {
		rawComponents = response.Data[0].Components
	}

	cdkComponents := make(map[string]CDKComponent, len(rawComponents)+len(localComponents))

	for _, c := range rawComponents {
		ver, err := semver.NewVersion(c.Version)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("component '%s' version '%s'", c.Name, c.Version))
		}

		allVersions, err := listComponentVersions(client, c.Id)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("unable to fetch component '%s' versions", c.Name))
		}

		api := NewAPIInfo(c.Id, c.Name, ver, allVersions, c.Description, c.Size, c.Deprecated, Type(c.ComponentType))

		host, found := localComponents[c.Name]
		if found {
			delete(localComponents, c.Name)
		}

		component := NewCDKComponent(c.Name, c.Description, Type(c.ComponentType), api, host)

		cdkComponents[c.Name] = component
	}

	for _, localHost := range localComponents {
		if localHost.Development() {
			devInfo, err := NewDevInfo(localHost.Dir())
			if err != nil {
				return nil, err
			}

			cdkComponents[localHost.Name()] = NewCDKComponent(
				localHost.Name(),
				devInfo.Desc,
				devInfo.ComponentType,
				nil,
				localHost)
		} else {
			// @jon-stewart: persisted API info
			cdkComponents[localHost.Name()] = NewCDKComponent(localHost.Name(), "", BinaryType, nil, localHost)
		}
	}

	return &Catalog{client, cdkComponents, stageConstructor}, nil
}

func LocalComponents() (components []CDKComponent, err error) {
	var localHost map[string]HostInfo

	localHost, err = loadLocalComponents()

	for _, l := range localHost {
		if l.Development() {
			devInfo, err := NewDevInfo(l.Dir())
			if err != nil {
				return nil, err
			}

			components = append(components, NewCDKComponent(l.Name(), devInfo.Desc, devInfo.ComponentType, nil, l))
		} else {
			// @jon-stewart: persisted API info
			components = append(components, NewCDKComponent(l.Name(), "", BinaryType, nil, l))
		}
	}

	return
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

func listComponentVersions(client *api.Client, componentId int32) (versions []*semver.Version, err error) {
	response, err := client.V2.Components.ListComponentVersions(componentId, operatingSystem, architecture)
	if err != nil {
		return nil, err
	}

	var rawVersions []string

	if len(response.Data) > 0 {
		rawVersions = response.Data[0].Versions
	}

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

// Returns the directory that the component executable and configuration is stored in.
func componentDirectory(componentName string) (string, error) {
	dir, err := CatalogCacheDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(dir, componentName), nil
}
