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
	cdkCacheName      string = "cdk_cache"
	featureFlag       string = "PUBLIC.cdk.v3"
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

// Return a CDKComponent that is present on the host.
func (c *Catalog) GetComponent(name string) (*CDKComponent, error) {
	component, exists := c.Components[name]
	if !exists {
		return nil, errors.New(fmt.Sprintf("component %s not found", name))
	}

	return &component, nil
}

func (c *Catalog) ListComponentVersions(component *CDKComponent) (versions []*semver.Version, err error) {
	if component.ApiInfo == nil {
		err = errors.Errorf("component '%s' api info not available", component.Name)
		return
	}

	versions = component.ApiInfo.AllVersions
	if versions != nil {
		return
	}

	return listComponentVersions(c.client, component.ApiInfo.Id)
}

func (c *Catalog) PrintComponents() [][]string {
	result := [][]string{}

	for _, component := range c.Components {
		result = append(result, component.PrintSummary())
	}

	return result
}

func (c *Catalog) Stage(
	component *CDKComponent,
	version string,
	progressClosure func(filepath string, sizeB int64)) (stageClose func(), err error) {
	var (
		semv *semver.Version
	)

	stageClose = func() {}

	if version == "" {
		semv = component.ApiInfo.Version
	} else {
		semv, err = semver.NewVersion(version)
		if err != nil {
			return
		}
	}

	if component.HostInfo != nil {
		var installedVersion *semver.Version

		installedVersion, err = component.HostInfo.Version()
		if err != nil {
			return
		}

		if installedVersion.Equal(semv) {
			err = errors.Errorf("version '%s' already installed", semv.String())
			return
		}
	}

	response, err := c.client.V2.Components.FetchComponentArtifact(
		component.ApiInfo.Id,
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

	stage, err := c.stageConstructor(component.Name, data.ArtifactUrl, data.Size)
	if err != nil {
		return
	}

	if err = stage.Download(progressClosure); err != nil {
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

	component.HostInfo, err = CreateHostInfo(componentDir, component.Description, component.Type)
	if err != nil {
		return
	}

	if component.ApiInfo != nil &&
		(component.ApiInfo.ComponentType == BinaryType || component.ApiInfo.ComponentType == CommandType) {
		if err := os.Chmod(filepath.Join(componentDir, component.Name), 0744); err != nil {
			return errors.Wrap(err, "unable to make component executable")
		}
	}

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

func NewCatalog(
	client *api.Client,
	stageConstructor StageConstructor,
) (*Catalog, error) {
	if stageConstructor == nil {
		return nil, errors.New("StageConstructor is not specified to create new catalog")
	}

	response, err := client.V2.Components.ListComponents(operatingSystem, architecture)
	if err != nil {
		return nil, err
	}

	var rawComponents []api.LatestComponentVersion

	if len(response.Data) > 0 {
		rawComponents = response.Data[0].Components
	}

	cdkComponents := make(map[string]CDKComponent, len(rawComponents))

	for _, c := range rawComponents {
		ver, err := semver.NewVersion(c.Version)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("component '%s' version '%s'", c.Name, c.Version))
		}

		var allVersions []*semver.Version

		apiInfo := NewAPIInfo(c.Id, c.Name, ver, allVersions, c.Description, c.Size, c.Deprecated, Type(c.ComponentType))
		cdkComponents[c.Name] = NewCDKComponent(c.Name, c.Description, Type(c.ComponentType), apiInfo, nil)
	}

	components, err := mergeComponents(cdkComponents)
	if err != nil {
		return nil, err
	}

	return &Catalog{client, components, stageConstructor}, nil
}

func NewCachedCatalog(
	client *api.Client,
	stageConstructor StageConstructor,
	cachedComponentsApiInfo map[string]*ApiInfo,
) (*Catalog, error) {
	if stageConstructor == nil {
		return nil, errors.New("StageConstructor is not specified to create new catalog")
	}

	cachedComponents := make(map[string]CDKComponent, len(cachedComponentsApiInfo))

	for _, a := range cachedComponentsApiInfo {
		cachedComponents[a.Name] = NewCDKComponent(a.Name, a.Desc, a.ComponentType, a, nil)
	}

	components, err := mergeComponents(cachedComponents)
	if err != nil {
		return nil, err
	}

	return &Catalog{client, components, stageConstructor}, nil
}

// mergeComponents combines the passed in components with the local components
func mergeComponents(components map[string]CDKComponent) (allComponents map[string]CDKComponent, err error) {
	localComponents, err := LoadLocalComponents()
	if err != nil {
		return
	}

	allComponents = make(map[string]CDKComponent, len(localComponents)+len(components))

	for _, c := range components {
		var hostInfo *HostInfo
		component, ok := localComponents[c.Name]
		if ok {
			hostInfo = component.HostInfo
			delete(localComponents, c.Name)
		}
		allComponents[c.Name] = NewCDKComponent(c.Name, c.Description, c.Type, c.ApiInfo, hostInfo)
	}

	for _, c := range localComponents {
		allComponents[c.Name] = c
	}

	return
}

func LoadLocalComponents() (components map[string]CDKComponent, err error) {
	cacheDir, err := CatalogCacheDir()
	if err != nil {
		return
	}

	subDir, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}

	components = make(map[string]CDKComponent, len(subDir))

	// Prototype backwards compatibility
	prototypeState, err := LocalState()
	if err != nil {
		prototypeState = new(State)
		err = nil
	}
	prototypeComponents := make(map[string]Component, len(prototypeState.Components))
	for _, component := range prototypeState.Components {
		prototypeComponents[component.Name] = component
	}

	for _, file := range subDir {
		if !file.IsDir() {
			continue
		}

		hostInfo, _ := NewHostInfo(filepath.Join(cacheDir, file.Name()))
		if hostInfo == nil {

			component, found := prototypeComponents[file.Name()]
			if !found {
				continue
			}

			hostInfo, err = CreateHostInfo(filepath.Join(cacheDir, file.Name()), component.Description, component.Type)
			if err != nil {
				return nil, err
			}
		}

		if hostInfo.Development() {
			devInfo, err := newDevInfo(hostInfo.Dir)
			if err != nil {
				return nil, err
			}
			components[hostInfo.Name()] = NewCDKComponent(hostInfo.Name(), devInfo.Desc, devInfo.ComponentType, nil, hostInfo)
		} else {
			components[hostInfo.Name()] = NewCDKComponent(
				hostInfo.Name(),
				hostInfo.Description,
				hostInfo.ComponentType,
				nil,
				hostInfo)
		}
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
