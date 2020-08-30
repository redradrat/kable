package kable

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/mitchellh/go-homedir"
)

const (
	ConfigFileName                = "kableconfig.json"
	ConfigNotInitializedError     = "currentConfig is not yet initialized"
	ConfigAlreadyInitializedError = "currentConfig is already initialized"
	RepositoryAlreadyExistsError  = "repository is already configured"
	RepositoryNotExistsError      = "repository is not configured"
	RepositoryNotInitializedError = "repository is not yet initialized"
)

var currentConfig *Config
var configPath string
var rootDir string
var userDir string
var curDir string
var repoDir string
var cfgHierarchy []string

// Config represents the Kable currentConfig object
type Config struct {
	APIVersion          APIVersion   `json:"apiVersion"`
	Repositories        Repositories `json:"repositories"`
	UseKeychainProvider bool         `json:"useKeychainProvider"`
}

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rootDir = "/etc/kable/"
	userDir = filepath.Join(home + "/.kable/")
	repoDir = filepath.Join(userDir, "/repos/")
	curDir = "./"
	cfgHierarchy = []string{
		filepath.Join(curDir, ConfigFileName),
		filepath.Join(userDir, ConfigFileName),
		filepath.Join(rootDir, ConfigFileName),
	}
}

// ReadConfig tries to read Config from either a file at the path that has been given as string argument, or, if cfgFile
// is an empty string, tries default locations.
//
// Default: /etc/kable/kableconfig.json, ~/.kable/kableconfig.json, ./kableconfig.json
func ReadConfig(cfgFile string) (string, error) {
	var cfgFileObj *os.File
	var err error

	if cfgFile != "" {
		// Use config file from the flag.
		configPath = cfgFile
		cfgFileObj, err = os.Open(cfgFile)
		if err != nil {
			return "", fmt.Errorf("Cannot open config file: %s \n", err)
		}
	} else {
		// Use default cfg hierarchy.
		for _, path := range cfgHierarchy {
			configPath = path
			cfgFileObj, err = os.Open(path)
			if err == nil {
				break
			}
			if !os.IsNotExist(err) {
				return "", fmt.Errorf("Cannot open config file: %s \n", err)
			}
		}
	}

	if cfgFileObj == nil {
		if err := initConfig(); err != nil {
			return configPath, fmt.Errorf("Cannot write config file: %s \n", err)
		}
	} else {
		var repoConf Config
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			return configPath, err
		}
		if err = json.Unmarshal(content, &repoConf); err != nil {
			return configPath, err
		}
		if err := setCurrentConfig(&repoConf); err != nil {
			return configPath, fmt.Errorf("Cannot write config file: %s \n", err)
		}
	}

	return configPath, nil
}

// initConfig initializes a fresh config file at ~/.config/kable/kableconfig.json, if not yet initialized
func initConfig() error {
	if configSet() {
		return nil
	}
	if err := setCurrentConfig(&Config{
		APIVersion:          "1",
		Repositories:        map[uuid.UUID]Repository{},
		UseKeychainProvider: true,
	}); err != nil {
		return err
	}
	configPath = filepath.Join(userDir + ConfigFileName)
	if err := writeConfig(configPath); err != nil {
		return err
	}
	return nil
}

func writeConfig(path string) error {
	if !configSet() {
		return fmt.Errorf(ConfigNotInitializedError)
	}
	filepath.Dir(configPath)
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	out, err := json.MarshalIndent(currentConfig, "", "	")
	if err != nil {
		return err
	}
	_, err = file.Write(out)
	if err != nil {
		return err
	}
	return nil
}

// GetConfig returns the current Config Object. One first needs to Read or Init Config.
func GetConfig() (*Config, error) {
	if configSet() {
		return currentConfig, nil
	}
	return nil, fmt.Errorf(ConfigNotInitializedError)
}

func setCurrentConfig(config *Config) error {
	if configSet() {
		return fmt.Errorf(ConfigAlreadyInitializedError)
	}
	currentConfig = config
	return nil
}

func configSet() bool {
	if currentConfig != nil && configPath != "" {
		return true
	}
	return false
}
