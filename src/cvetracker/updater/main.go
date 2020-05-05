package updater

import (
	"cvetracker/util"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"go.uber.org/zap"
	"gopkg.in/yaml.v2"
)

// Main ...
func Main() {

	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)

	var fileName string
	flag.StringVar(&fileName, "c", "", "config file")
	flag.Parse()

	fileNames := []string{"cveupdater.config", "cveupdater.yaml", "cveupdater.json"}

	if fileName != "" {
		if !util.FileExists(fileName) {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("File %s not found", fileName))
			os.Exit(2)
		}
	} else {

		for _, tmpFileName := range fileNames {
			if util.FileExists(tmpFileName) {
				fileName = tmpFileName
				break
			}
		}
	}

	if fileName == "" {
		fmt.Fprintln(os.Stderr, "Use arg -c to specify the config filename")
		os.Exit(2)
	}

	config, err := getConfigFromFile(fileName)
	if err != nil {
		panic(err)
	}

	config.AssertValid()

	updater := &updater{config: config}
	updater.run()

}

func getConfigFromFile(fileName string) (*Config, error) {

	rawFile, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config *Config

	yamlErr := yaml.Unmarshal(rawFile, &config)
	if yamlErr == nil {
		zap.L().Debug("Config file is YAML")
	} else {

		jsonErr := json.Unmarshal(rawFile, &config)
		if jsonErr == nil {
			zap.L().Debug("Config file is JSON")
		} else {
			return nil, errors.New("Unable to parse config file as json or yaml: " + yamlErr.Error() + " : " + jsonErr.Error())
		}

	}

	return config, nil
}

// WriteExampleConfig ...
func WriteExampleConfig(fileName string) {

	if util.FileExists(fileName) {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("File %s already exist!", fileName))
		return
	}

	err := ioutil.WriteFile(fileName, []byte(ExampleConfig().YamlString()), 0400)
	if err != nil {
		panic(err)
	}
}
