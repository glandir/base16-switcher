package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type UpdateCmd struct {}

func (self *UpdateCmd) Run() error {
	config := LoadConfig()

	UpdateTemplates(&config)

	DownloadSchemeLists(&config)
	PullSchemes(&config)

	return nil
}

func UpdateTemplates(config *Config) {
	templatesDir := filepath.Join(config.xdgDirs.CacheHome(), "templates")

	err := os.MkdirAll(templatesDir, os.ModePerm)
	assert(err, "Could not create directory ", templatesDir)

	fmt.Print("Pulling templates\n-----------------\n")

	for name, app := range config.Applications {
		templateDir := filepath.Join(templatesDir, name)
		PullOrClone(name, app.Url, templateDir)
	}
}

func DownloadSchemeLists(config *Config) {
	schemeListPath := filepath.Join(config.xdgDirs.CacheHome(), "schemes-list.yaml")
	out, err := os.Create(schemeListPath)
	assert(err, "Error creating file ", schemeListPath)
	defer out.Close()

	for _, url := range config.SchemeSources {
		resp, err := http.Get(url)
		assert(err, "Could not downlad ", url)
		defer resp.Body.Close()

		// ATM, just ignore unavailable ones.
		if resp.StatusCode == http.StatusOK {
			_, err = io.Copy(out, resp.Body)
			assert(err, "Error writing to file ", schemeListPath)
		}
	}
}

func PullSchemes(config *Config) {
	fmt.Print("\n\nPulling schemes\n---------------\n")
	schemes := ReadYamlFile(config.xdgDirs.QueryCache("schemes-list.yaml"))

	schemesDir := filepath.Join(config.xdgDirs.CacheHome(), "schemes")
	err := os.MkdirAll(schemesDir, os.ModePerm)
	assert(err, "Could not create directory ", schemesDir)

	for name, location := range schemes {
		schemeDir := filepath.Join(schemesDir, name)
		PullOrClone(name, location, schemeDir)
	}
}
