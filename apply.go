package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type ApplyCmd struct {
	Name string `arg:"" optional:"" name:"scheme-name"`
}

func (self *ApplyCmd) Run() error {
	config := LoadConfig()

	if len(self.Name) == 0 {
		if len(config.DefaultColorscheme) == 0 {
			panic("")
		}
		self.Name = config.DefaultColorscheme
	}

	availableSchemes := AvailableSchemes(&config)
	schemePath, ok := availableSchemes[self.Name]
	if !ok {
		assert(fmt.Errorf("Scheme %s does not exist.", self.Name))
	}

	fmt.Printf("Applying scheme %s\n\n", self.Name)

	scheme := ReadYamlFile(schemePath)
	scheme = ConvertScheme(scheme)

	for appName, application := range config.Applications {
		for template, destination := range application.Files {
			destination := os.ExpandEnv(destination)

			templateContent := ReadTemplate(appName, template, &config)
			templateContent = Apply(scheme, templateContent)

			WriteFile(destination, templateContent)
		}
	}

	for _, application := range config.Applications {
		for _, hook := range application.Hooks {
			err := exec.Command("sh", "-c", hook).Run()
			if err != nil {
				fmt.Println("Hook execution failed:", err.Error())
			}
		}
	}

	return nil
}

func Apply(scheme map[string]string, template string) string {
	for key, value := range scheme {
		template = strings.ReplaceAll(template, "{{"+key+"}}", value)
	}
	return template
}

func ConvertScheme(scheme map[string]string) map[string]string {
	result := make(map[string]string)

	result["scheme-author"] = scheme["author"]
	result["scheme-name"] = scheme["name"]

	for _, n := range []string{
		"00", "01", "02", "03",
		"04", "05", "06", "07",
		"08", "09", "0A", "0B",
		"0C", "0D", "0E", "0F",
	} {
		name := "base" + n

		baseHex := scheme[name]

		r := baseHex[0:2]
		g := baseHex[2:4]
		b := baseHex[4:6]

		result[name+"-hex-r"] = r
		result[name+"-hex-g"] = g
		result[name+"-hex-b"] = b
		result[name+"-rgb-r"] = toRgb(r)
		result[name+"-rgb-g"] = toRgb(g)
		result[name+"-rgb-b"] = toRgb(b)
		result[name+"-dec-r"] = toDec(r)
		result[name+"-dec-g"] = toDec(g)
		result[name+"-dec-b"] = toDec(b)

		result[name+"-hex"] = baseHex
		result[name+"-hex-bgr"] = b + g + r
	}

	return result
}

func toRgb(hex string) string {
	n, err := strconv.ParseInt(hex, 16, 32)
	assert(err, "Could not convert hex value ", hex, " to rgb component")
	return strconv.FormatInt(n, 10)
}

func toDec(hex string) string {
	n, err := strconv.ParseInt(hex, 16, 32)
	assert(err, "Could not convert hex value ", hex, " to decimal rgb component")
	f := float64(n) / 256.0
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func ReadTemplate(appName string, template string, config *Config) string {
	relativeTemplPath := "templates/" +
		appName + "/templates/" + template + ".mustache"

	templatePath := config.xdgDirs.QueryCache(relativeTemplPath)
	if len(templatePath) == 0 {
		assert(fmt.Errorf("Could not find template %s in cache ", relativeTemplPath))
	}
	templateBytes, err := ioutil.ReadFile(templatePath)
	assert(err, "Could not read template file ", templatePath)
	templateContent := string(templateBytes)
	return templateContent
}

func WriteFile(path string, content string) {
	dirpath := filepath.Dir(path)
	err := os.MkdirAll(dirpath, os.ModeDir)
	assert(err, "Could not create directory ", dirpath)

	outfile, err := os.Create(path)
	assert(err, "Could not create file ", path)
	defer outfile.Close()

	_, err = outfile.WriteString(content)
	assert(err, "Could not write to file ", content)
}
