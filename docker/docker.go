package docker

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/docker/docker/container"
)

const DefaultDockerRoot string = "/var/lib/docker"
const DefaultDockerVersion string = "1.10.0"

type Container struct {
	container.CommonContainer
}

// GetDockerRoot gets docker root path in Docker environment
func GetDockerRoot() string {
	dockerRoot := os.Getenv("DOCKER_ROOT")
	if dockerRoot == "" {
		dockerRoot = DefaultDockerRoot
	}
	return dockerRoot
}

// GetDockerVersion gets docker version from env
func GetDockerVersion() string {
	dockerVersion := os.Getenv("DOCKER_VERSION")
	if dockerVersion == "" {
		return DefaultDockerVersion
	}
	// FIXME: validate dockerVersion to be "version.major.minor" format
	return dockerVersion
}

// Containers returns an array of docker containers unmarshaled from config.json
// in docker root path.
func Containers() ([]Container, error) {
	containersPath := filepath.Join(GetDockerRoot(), "containers")

	containerEntries, err := ioutil.ReadDir(containersPath)
	if err != nil {
		return nil, err
	}

	containers := []Container{}

	for _, entry := range containerEntries {
		entryName := entry.Name()
		if len(entryName) != len("ffb082df6289394f4d285ef2ea31051deed699f6b352cf4109fb7e97fd15237a") {
			continue
		}

		// If docker version differs, json file's name differs, too.
		// config.json in 1.10.0-, while config.v2.json in 1.10.0+
		var configFilename string
		match, err := util.CompareDockerVersion(GetDockerVersion(), "1.10.0")
		if err != nil {
		}

		// If match is true, it means current docker version is newer or at least equal.
		if match {
			configFilename = "config.v2.json"
		} else {
			configFilename = "config.json"
		}

		containerJsonPath := filepath.Join(containersPath, entryName, configFilename)

		con, err := containerFromJson(containerJsonPath)
		if err != nil {
			continue
		}

		containers = append(containers, con)
	}
	return containers, nil
}

func containerFromJson(file string) (Container, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return Container{}, err
	}

	var con Container
	if err := json.Unmarshal(data, &con); err != nil {
		return Container{}, err
	}

	return con, nil
}