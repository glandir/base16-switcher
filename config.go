package main

import (
	"io/ioutil"
	"os"

	"github.com/OpenPeeDeeP/xdg"
	"gopkg.in/yaml.v2"
)

type Application struct {
	Url   string
	Files map[string]string
	Hooks []string
}

type Config struct {
	DefaultColorscheme string   `yaml:"default-colorscheme"`
	SchemeSources      []string `yaml:"scheme-sources"`
	Applications       map[string]Application
	xdgDirs            *xdg.XDG `yaml:",ignore"`
}

func LoadConfig() Config {
	var config Config
	config.xdgDirs = xdg.New("base16-switcher", "")

	configDir := config.xdgDirs.ConfigHome()
	err := os.MkdirAll(configDir, os.ModePerm)
	assert(err, "Could not create directory ", configDir)

	cacheDir := config.xdgDirs.CacheHome()
	err = os.MkdirAll(cacheDir, os.ModePerm)
	assert(err, "Could not create directory ", configDir)

	configPath := config.xdgDirs.QueryConfig("config.yaml")
	file, err := ioutil.ReadFile(configPath)
	assert(err, "Could not read configuration file ", configPath)

	err = yaml.Unmarshal(file, &config)
	assert(err, "Could not parse ", configPath)

	return config
}
