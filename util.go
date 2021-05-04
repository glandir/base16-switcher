package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v2"
)

func PullOrClone(name string, location string, targetDir string) {
	repo, err := git.PlainOpen(targetDir)

	if err != nil {
		fmt.Println("Cloning", name, "from", location)

		_, err := git.PlainClone(targetDir, false, &git.CloneOptions{
			URL:               location,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		assert(err, "Error pulling from ", location)
	} else {
		fmt.Println("Pulling", name, "from", location)

		w, err := repo.Worktree()
		assert(err, "Error retrieving work tree")
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
	}
}

func ReadYamlFile(path string) (map[string]string) {
	result := make(map[string]string)

	file, err := ioutil.ReadFile(path)
	assert(err, "Error reading file ", path)

	err = yaml.Unmarshal(file, result)
	assert(err, "Could not parse ", path)

	return result
}

func assert(err error, messages ...interface{}) {
	if err != nil {
		fmt.Print(messages...)
		if len(messages) != 0 {
			fmt.Print(" : ")
		}
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
