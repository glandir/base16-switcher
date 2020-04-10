package main

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type ListCmd struct {}

func (self *ListCmd) Run() error {
	config := LoadConfig()
	availableSchemes := AvailableSchemes(&config)
	keys := make([]string, 0, len(availableSchemes))
	for k := range availableSchemes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Println(k)
	}
	return nil
}

// Returns a map of {schemename: path}
func AvailableSchemes(config *Config) map[string]string {
	configDir := config.xdgDirs.CacheHome()
	glob := configDir + "/schemes/*/*.yaml"
	availableSchemes, err := filepath.Glob(glob)
	assert(err, "Failed to resolve glob ", glob)

	result := make(map[string]string)
	for _, e := range availableSchemes {
		filename := strings.TrimSuffix(filepath.Base(e), ".yaml")
		result[filename] = e
	}

	return result
}
