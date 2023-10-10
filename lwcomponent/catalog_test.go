package lwcomponent_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/lacework/go-sdk/api"
	"github.com/lacework/go-sdk/internal/file"
	"github.com/lacework/go-sdk/internal/lacework"
	"github.com/lacework/go-sdk/lwcomponent"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestCatalogNewCatalog(t *testing.T) {
	var (
		prefix            = "testNewCatalog"
		apiComponentCount = 4
	)

	t.Run("ok", func(t *testing.T) {
		fakeServer := lacework.MockServer()
		defer fakeServer.Close()

		fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
			fmt.Fprint(w, generateComponentsResponse(prefix, apiComponentCount))
		})

		client, _ := api.NewClient("catalog_test",
			api.WithToken("TOKEN"),
			api.WithURL(fakeServer.URL()),
		)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)
	})

	t.Run("installed", func(t *testing.T) {
		fakeServer := lacework.MockServer()
		defer fakeServer.Close()

		_, home := FakeHome()
		defer ResetHome(home)

		fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
			fmt.Fprint(w, generateComponentsResponse(prefix, apiComponentCount))
		})

		client, _ := api.NewClient("catalog_test",
			api.WithToken("TOKEN"),
			api.WithURL(fakeServer.URL()),
		)

		name := fmt.Sprintf("%s-%d", prefix, 1)
		version := fmt.Sprintf("%d.0.0", 1)

		CreateLocalComponent(name, version, false)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)
		assert.Equal(t, lwcomponent.Installed, component.Status)
	})

	t.Run("invalid api semantic version", func(t *testing.T) {
		fakeServer := lacework.MockServer()
		defer fakeServer.Close()

		_, home := FakeHome()
		defer ResetHome(home)

		fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
			fmt.Fprint(w, generateInvalidComponentsResponse())
		})

		client, _ := api.NewClient("catalog_test",
			api.WithToken("TOKEN"),
			api.WithURL(fakeServer.URL()),
		)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.Nil(t, catalog)
		assert.NotNil(t, err)
	})

	t.Run("invalid local semantic version", func(t *testing.T) {
		fakeServer := lacework.MockServer()
		defer fakeServer.Close()

		_, home := FakeHome()
		defer ResetHome(home)

		fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
			fmt.Fprint(w, generateComponentsResponse(prefix, 1))
		})

		client, _ := api.NewClient("catalog_test",
			api.WithToken("TOKEN"),
			api.WithURL(fakeServer.URL()),
		)

		name := fmt.Sprintf("%s-%d", prefix, 1)

		CreateLocalComponent(name, "invalid-version", false)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		assert.Equal(t, lwcomponent.UnknownStatus, component.Status)

	})

}

func TestCatalogComponentCount(t *testing.T) {
	var (
		prefix           = "testCount"
		apiCount         = 5
		deprecatedCount  = 1
		developmentCount = 3
	)

	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
		fmt.Fprint(w, generateComponentsResponse(prefix, apiCount))
	})

	client, _ := api.NewClient("catalog_test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)

	t.Run("no api components installed", func(t *testing.T) {
		_, home := FakeHome()
		defer ResetHome(home)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)
		assert.Equal(t, apiCount, catalog.ComponentCount())
	})

	t.Run("api components installed", func(t *testing.T) {
		_, home := FakeHome()
		defer ResetHome(home)

		for i := 0; i < apiCount; i++ {
			CreateLocalComponent(fmt.Sprintf("%s-%d", prefix, i), fmt.Sprintf("%d.0.0", i), false)
		}

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)
		assert.Equal(t, apiCount, catalog.ComponentCount())
	})

	t.Run("deprecated components", func(t *testing.T) {
		_, home := FakeHome()
		defer ResetHome(home)

		for i := 0; i < deprecatedCount; i++ {
			CreateLocalComponent(fmt.Sprintf("deprecated-%d", i), fmt.Sprintf("%d.0.0", i), false)
		}

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)
		assert.Equal(t, apiCount+deprecatedCount, catalog.ComponentCount())
	})

	t.Run("development components", func(t *testing.T) {
		_, home := FakeHome()
		defer ResetHome(home)

		for i := 0; i < developmentCount; i++ {
			CreateLocalComponent(fmt.Sprintf("dev-%d", i), fmt.Sprintf("%d.0.0", i), true)
		}

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)
		assert.Equal(t, apiCount+developmentCount, catalog.ComponentCount())
	})

	t.Run("all components", func(t *testing.T) {
		_, home := FakeHome()
		defer ResetHome(home)

		for i := 0; i < deprecatedCount; i++ {
			CreateLocalComponent(fmt.Sprintf("all-deprecated-%d", i), fmt.Sprintf("%d.0.0", i), false)
		}

		for i := 0; i < developmentCount; i++ {
			CreateLocalComponent(fmt.Sprintf("all-dev-%d", i), fmt.Sprintf("%d.0.0", i), true)
		}

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)
		assert.Equal(t, apiCount+deprecatedCount+developmentCount, catalog.ComponentCount())
	})
}

func TestCatalogGetComponent(t *testing.T) {
	var (
		fakeServer = lacework.MockServer()
		prefix     = "testGet"
		count      = 3
	)

	defer fakeServer.Close()

	fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
		fmt.Fprint(w, generateComponentsResponse(prefix, count))
	})

	client, _ := api.NewClient("catalog_test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)

	t.Run("found", func(t *testing.T) {
		var (
			name    = fmt.Sprintf("%s-%d", prefix, 1)
			version = "1.0.0"
		)

		_, home := FakeHome()
		defer ResetHome(home)

		CreateLocalComponent(name, version, true)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)
		assert.Equal(t, lwcomponent.Installed, component.Status)
	})

	t.Run("not found", func(t *testing.T) {
		_, home := FakeHome()
		defer ResetHome(home)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.Nil(t, err)

		component, err := catalog.GetComponent("component-example")
		assert.Nil(t, component)
		assert.NotNil(t, err)
	})

	t.Run("development", func(t *testing.T) {
		var (
			name    = "development"
			version = "1.1.0"
		)

		_, home := FakeHome()
		defer ResetHome(home)

		CreateLocalComponent(name, version, true)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)
		assert.Equal(t, lwcomponent.Development, component.Status)

		ver := component.InstalledVersion()
		assert.Equal(t, version, ver.String())
	})

	t.Run("deprecated", func(t *testing.T) {
		var (
			name    = "deprecated"
			version = "1.1.0"
		)

		_, home := FakeHome()
		defer ResetHome(home)

		CreateLocalComponent(name, version, false)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)
		assert.Equal(t, lwcomponent.Deprecated, component.Status)

		ver := component.InstalledVersion()
		assert.Equal(t, version, ver.String())
	})
}

func TestCatalogListComponentVersions(t *testing.T) {
	var (
		prefix = "testCatalogListComponentVersions"
		count  = 4
	)

	_, home := FakeHome()
	defer ResetHome(home)

	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
		fmt.Fprint(w, generateComponentsResponse(prefix, count))
	})

	client, _ := api.NewClient("catalog_test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)

	t.Run("ok", func(t *testing.T) {
		name := fmt.Sprintf("%s-%d", prefix, 1)
		versions := []string{"0.1.0", "1.1.1", "3.0.1", "5.4.3"}

		fakeServer.MockAPI("Components/1", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
			fmt.Fprint(w, generateComponentVersionsResponse(name, versions))
		})

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		vers, err := catalog.ListComponentVersions(component)
		assert.Nil(t, err)

		for idx, v := range versions {
			assert.Equal(t, v, vers[idx].String())
		}
	})

	t.Run("invalid semantic version", func(t *testing.T) {
		name := fmt.Sprintf("%s-%d", prefix, 2)
		versions := []string{"0.1.0", "1.invalid.1", "3.0.1", "5.4.3"}

		fakeServer.MockAPI("Components/2", func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
			fmt.Fprint(w, generateComponentVersionsResponse(name, versions))
		})

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		vers, err := catalog.ListComponentVersions(component)
		assert.Nil(t, vers)
		assert.NotNil(t, err)
	})
}

func TestCatalogStage(t *testing.T) {
	var (
		apiComponentCount int    = 4
		prefix            string = "staging"
		version           string = "1.0.0"
		name              string = fmt.Sprintf("%s-1", prefix)
	)

	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	_, home := FakeHome()
	defer ResetHome(home)

	fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
		fmt.Fprint(w, generateComponentsResponse(prefix, apiComponentCount))
	})

	fakeServer.MockAPI("Components/Artifact/1", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")

		if r.URL.Query().Get("version") != version {
			http.Error(w, "component version not found", http.StatusNotFound)
		} else {
			fmt.Fprint(w, generateFetchResponse(1, name, version, ""))
		}
	})

	client, _ := api.NewClient("catalog_test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()),
	)

	catalog, err := lwcomponent.NewCatalog(client, newTestStage)
	assert.NotNil(t, catalog)
	assert.Nil(t, err)

	t.Run("ok", func(t *testing.T) {
		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		stageClose, err := catalog.Stage(component, version)
		assert.Nil(t, err)
		defer stageClose()
	})

	t.Run("already installed", func(t *testing.T) {
		CreateLocalComponent(name, version, false)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		stageClose, err := catalog.Stage(component, version)
		assert.NotNil(t, err)
		defer stageClose()
	})

	t.Run("invalid semantic version", func(t *testing.T) {
		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		stageClose, err := catalog.Stage(component, "invalid-version")
		assert.NotNil(t, err)
		defer stageClose()
	})

	t.Run("version not found", func(t *testing.T) {
		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		stageClose, err := catalog.Stage(component, "99.99.99")
		assert.NotNil(t, err)
		defer stageClose()
	})

}

type testStage struct {
	dir string
}

// Close implements lwcomponent.Stager.
func (t *testStage) Close() {
	os.RemoveAll(t.dir)
}

// Commit implements lwcomponent.Stager.
func (*testStage) Commit(string) error {
	return nil
}

// Directory implements lwcomponent.Stager.
func (t *testStage) Directory() string {
	return t.dir
}

// Download implements lwcomponent.Stager.
func (*testStage) Download() error {
	return nil
}

// Signature implements lwcomponent.Stager.
func (t *testStage) Signature() (sig []byte, err error) {
	path := filepath.Join(t.dir, lwcomponent.SignatureFile)
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

// Unpack implements lwcomponent.Stager.
func (*testStage) Unpack() error {
	return nil
}

// Validate implements lwcomponent.Stager.
func (*testStage) Validate() error {
	return nil
}

func newTestStage(name, artifactUrl string) (stage lwcomponent.Stager, err error) {
	stage = &testStage{}

	return
}

func TestCatalogVerify(t *testing.T) {

	t.Run("", func(t *testing.T) {})
	t.Run("", func(t *testing.T) {})
	t.Run("", func(t *testing.T) {})
	t.Run("", func(t *testing.T) {})
	t.Run("", func(t *testing.T) {})
}

func TestCatalogInstall(t *testing.T) {
	var (
		apiComponentCount int    = 4
		prefix            string = "staging"
		version           string = "1.0.0"
	)

	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	_, home := FakeHome()
	defer ResetHome(home)

	for i := 0; i < apiComponentCount; i++ {
		name := fmt.Sprintf("%s-%d", prefix, i)

		fakeServer.MockAPI(fmt.Sprintf("Components/Artifact/%d", i), func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, generateFetchResponse(int32(i), name, version, ""))
		})
	}

	fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
		fmt.Fprint(w, generateComponentsResponse(prefix, apiComponentCount))
	})

	client, _ := api.NewClient("catalog_test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()))

	catalog, err := lwcomponent.NewCatalog(client, newTestStage)
	assert.NotNil(t, catalog)
	assert.Nil(t, err)

	t.Run("ok", func(t *testing.T) {
		name := fmt.Sprintf("%s-1", prefix)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		stageClose, err := catalog.Stage(component, version)
		assert.Nil(t, err)
		defer stageClose()

		dir, _ := lwcomponent.CatalogCacheDir()
		os.MkdirAll(filepath.Join(dir, name), os.ModePerm)

		err = catalog.Install(component)
		assert.Nil(t, err)
	})

	t.Run("not staged", func(t *testing.T) {
		name := fmt.Sprintf("%s-2", prefix)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		dir, _ := lwcomponent.CatalogCacheDir()
		os.MkdirAll(filepath.Join(dir, name), os.ModePerm)

		err = catalog.Install(component)
		assert.NotNil(t, err)
	})
}

func TestCatalogDelete(t *testing.T) {
	var (
		apiComponentCount int    = 4
		prefix            string = "staging"
		version           string = "1.0.0"
	)

	fakeServer := lacework.MockServer()
	defer fakeServer.Close()

	_, home := FakeHome()
	defer ResetHome(home)

	fakeServer.MockAPI("Components", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method, "Components API only accepts HTTP GET")
		fmt.Fprint(w, generateComponentsResponse(prefix, apiComponentCount))
	})

	client, _ := api.NewClient("catalog_test",
		api.WithToken("TOKEN"),
		api.WithURL(fakeServer.URL()))

	t.Run("delete-installed", func(t *testing.T) {
		name := fmt.Sprintf("%s-1", prefix)

		CreateLocalComponent(name, version, false)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		err = catalog.Delete(component)
		assert.Nil(t, err)

		dir, _ := lwcomponent.CatalogCacheDir()
		dir = filepath.Join(dir, name)

		_, err = os.Stat(dir)
		assert.NotNil(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("delete-development", func(t *testing.T) {
		name := "delete-dev"

		CreateLocalComponent(name, version, true)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		err = catalog.Delete(component)
		assert.Nil(t, err)

		dir, _ := lwcomponent.CatalogCacheDir()
		dir = filepath.Join(dir, name)

		_, err = os.Stat(dir)
		assert.NotNil(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("delete-not-installed", func(t *testing.T) {
		name := fmt.Sprintf("%s-1", prefix)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		err = catalog.Delete(component)
		assert.Nil(t, err)

		dir, _ := lwcomponent.CatalogCacheDir()
		dir = filepath.Join(dir, name)

		_, err = os.Stat(dir)
		assert.NotNil(t, err)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("delete-twice", func(t *testing.T) {
		name := fmt.Sprintf("%s-2", prefix)

		CreateLocalComponent(name, version, false)

		catalog, err := lwcomponent.NewCatalog(client, newTestStage)
		assert.NotNil(t, catalog)
		assert.Nil(t, err)

		component, err := catalog.GetComponent(name)
		assert.NotNil(t, component)
		assert.Nil(t, err)

		err = catalog.Delete(component)
		assert.Nil(t, err)

		err = catalog.Delete(component)
		assert.Nil(t, err)

		dir, _ := lwcomponent.CatalogCacheDir()
		dir = filepath.Join(dir, name)

		_, err = os.Stat(dir)
		assert.NotNil(t, err)
		assert.True(t, os.IsNotExist(err))
	})

}

func generateComponentsResponse(prefix string, count int) string {
	var (
		components = []api.LatestComponentVersion{}
		idx        int32
	)

	for idx = 0; idx < int32(count); idx++ {
		component := api.LatestComponentVersion{
			Id:      idx,
			Name:    fmt.Sprintf("%s-%d", prefix, idx),
			Version: fmt.Sprintf("%d.0.0", idx),
			Size:    512,
		}

		components = append(components, component)
	}

	response := api.ListComponentsResponse{
		Data: []api.LatestComponent{{Components: components}},
	}

	result, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

func generateInvalidComponentsResponse() string {
	response := api.ListComponentsResponse{
		Data: []api.LatestComponent{{Components: []api.LatestComponentVersion{{Id: 1, Name: "invalidVersion", Version: "invalidVersion"}}}},
	}

	result, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return string(result)
}

func generateComponentVersionsResponse(name string, versions []string) string {
	response := api.ListComponentVersionsResponse{
		Data: []api.ComponentVersions{{
			Id:       1,
			Name:     name,
			Versions: versions,
		}},
	}

	result, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)

}

func generateFetchResponse(id int32, name string, version string, url string) string {
	response := api.FetchComponentResponse{
		Data: []api.Artifact{{
			Id:             id,
			Name:           name,
			Version:        version,
			Size:           0,
			InstallMessage: "install message",
			UpdateMessage:  "update message",
			ArtifactUrl:    url,
		}},
	}

	result, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
	}

	return string(result)
}

func CreateLocalComponent(componentName string, version string, development bool) {
	dir, err := lwcomponent.CatalogCacheDir()
	if err != nil {
		panic(err)
	}

	path := filepath.Join(dir, componentName)
	if err := os.Mkdir(path, os.ModePerm); err != nil {
		panic(err)
	}

	if development {
		if _, err := os.Create(filepath.Join(path, lwcomponent.DevelopmentFile)); err != nil {
			panic(err)
		}
	}

	if err := os.WriteFile(filepath.Join(path, lwcomponent.VersionFile), []byte(version), 0666); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(path, lwcomponent.SignatureFile), []byte(version), 0666); err != nil {
		panic(err)
	}

	if err := os.WriteFile(filepath.Join(path, componentName), []byte("#!/bin/sh\necho 'hi'"), 0766); err != nil {
		panic(err)
	}
}

func FakeHome() (fake string, home string) {
	fake, err := os.MkdirTemp("", "catalog_test")
	if err != nil {
		panic(err)
	}

	home = os.Getenv("HOME")

	os.Setenv("HOME", fake)
	homedir.DisableCache = true

	cacheDir, err := lwcomponent.CatalogCacheDir()
	if err != nil {
		panic(err)
	}

	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		panic(err)
	}

	return
}

func ResetHome(dir string) {
	os.Setenv("HOME", dir)
	homedir.DisableCache = false
}
