package lwcomponent

import (
	"bytes"
	"encoding/json"
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
	featureFlag       string = "PUBLIC.cdk.v1"
	operatingSystem   string = runtime.GOOS
	architecture      string = runtime.GOARCH
)

func CatalogV1Enabled(client *api.Client) bool {
	return true
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

func (c *Catalog) Persist() error {
	cacheDir, err := CatalogCacheDir()
	if err != nil {
		return err
	}

	cdkCacheFile := filepath.Join(cacheDir, cdkCacheName)
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(c.Components); err != nil {
		return err
	}
	if err := os.WriteFile(cdkCacheFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (c *Catalog) PersistComponent(component *CDKComponent) error {
	cdkCache, err := loadCdkCache()
	if err != nil {
		return err
	}

	cacheDir, err := CatalogCacheDir()
	if err != nil {
		return err
	}

	if component.ApiInfo != nil && len(component.ApiInfo.AllVersions) == 0 {
		allVersions, err := listComponentVersions(c.client, component.ApiInfo.Id)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("unable to fetch component '%s' versions", component.Name))
		}
		component.ApiInfo.AllVersions = allVersions
	}
	cdkCache[component.Name] = *component

	cdkCacheFile := filepath.Join(cacheDir, cdkCacheName)
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(cdkCache); err != nil {
		return err
	}
	if err := os.WriteFile(cdkCacheFile, buf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
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
		err = errors.Errorf("component '%s' api info  already installed", component.Name)
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
		semv = &component.ApiInfo.Version
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

	if component.ApiInfo != nil &&
		(component.ApiInfo.ComponentType == BinaryType || component.ApiInfo.ComponentType == CommandType) {
		if err := os.Chmod(filepath.Join(componentDir, component.Name), 0744); err != nil {
			return errors.Wrap(err, "unable to make component executable")
		}
	}

	component.HostInfo = NewHostInfo(componentDir)

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
	includeComponentVersions bool,
) (*Catalog, error) {
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

		var allVersions []*semver.Version
		if includeComponentVersions {
			allVersions, err = listComponentVersions(client, c.Id)
			if err != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("unable to fetch component '%s' versions", c.Name))
			}
		}

		apiInfo := NewAPIInfo(c.Id, c.Name, ver, allVersions, c.Description, c.Size, c.Deprecated, Type(c.ComponentType))

		hostInfo, found := localComponents[c.Name]
		if found {
			delete(localComponents, c.Name)
		}

		component := NewCDKComponent(c.Name, c.Description, Type(c.ComponentType), apiInfo, hostInfo)

		cdkComponents[c.Name] = component
	}

	for _, localHost := range localComponents {
		if localHost.Development() {
			devInfo, err := NewDevInfo(localHost.Dir)
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
			// Unknown components
			cdkComponents[localHost.Name()] = NewCDKComponent(localHost.Name(), "", EmptyType, nil, localHost)
		}
	}

	return &Catalog{client, cdkComponents, stageConstructor}, nil
}

func LocalComponents() (components []CDKComponent, err error) {
	cdkCache, err := loadCdkCache()
	if err != nil {
		return
	}

	localComponents, err := loadLocalComponents()
	if err != nil {
		return
	}

	// Load components stored in cdk_cache
	for _, c := range cdkCache {
		hostInfo, ok := localComponents[c.Name]
		if ok {
			delete(localComponents, c.Name)
		}
		components = append(components, NewCDKComponent(c.Name, c.Description, c.Type, c.ApiInfo, hostInfo))
	}

	// Load components in the components folder but not stored in cdk_cache
	for _, hostInfo := range localComponents {
		if hostInfo.Development() {
			devInfo, err := NewDevInfo(hostInfo.Dir)
			if err != nil {
				return nil, err
			}
			components = append(components, NewCDKComponent(hostInfo.Name(), devInfo.Desc, devInfo.ComponentType, nil, hostInfo))
		} else {
			// Unknown components
			components = append(components, NewCDKComponent(hostInfo.Name(), "", EmptyType, nil, hostInfo))
		}
	}

	return
}

func loadCdkCache() (map[string]CDKComponent, error) {
	var cdkCache map[string]CDKComponent

	cacheDir, err := CatalogCacheDir()
	if err != nil {
		return nil, err
	}

	cacheBytes, err := os.ReadFile(filepath.Join(cacheDir, cdkCacheName))
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(cacheBytes, &cdkCache); err != nil {
		return nil, err
	}

	return cdkCache, nil
}

func loadLocalComponents() (local map[string]*HostInfo, err error) {
	cacheDir, err := CatalogCacheDir()
	if err != nil {
		return
	}

	subDir, err := os.ReadDir(cacheDir)
	if err != nil {
		return
	}

	local = make(map[string]*HostInfo, len(subDir))

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
