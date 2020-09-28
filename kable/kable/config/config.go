package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/redradrat/kable/pkg/kable/errors"

	"github.com/mitchellh/go-homedir"
)

const (
	ConfigFileName = "kableconfig.json"
)

var currentConfig *Config
var ConfigPath string
var RootDir string
var UserDir string
var CurDir string
var RepoDir string
var ConceptDir string
var cfgHierarchy []string

// Config represents the Kable currentConfig object
type Config struct {
	APIVersion          int  `json:"apiVersion"`
	UseKeychainProvider bool `json:"useKeychainProvider"`
}

func init() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	RootDir = "/etc/kable/"
	UserDir = filepath.Join(home + "/.kable/")
	RepoDir = filepath.Join(UserDir, "/repos/")
	ConceptDir = filepath.Join(UserDir, "/concepts/")
	CurDir = "./"
	cfgHierarchy = []string{
		filepath.Join(CurDir, ConfigFileName),
		filepath.Join(UserDir, ConfigFileName),
		filepath.Join(RootDir, ConfigFileName),
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
		ConfigPath = cfgFile
		cfgFileObj, err = os.Open(cfgFile)
		if err != nil {
			return "", fmt.Errorf("Cannot open config file: %s \n", err)
		}
	} else {
		// Use default cfg hierarchy.
		for _, path := range cfgHierarchy {
			ConfigPath = path
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
			return ConfigPath, fmt.Errorf("Cannot write config file: %s \n", err)
		}
	} else {
		var repoConf Config
		content, err := ioutil.ReadFile(ConfigPath)
		if err != nil {
			return ConfigPath, err
		}
		if err = json.Unmarshal(content, &repoConf); err != nil {
			return ConfigPath, err
		}
		if err := setCurrentConfig(&repoConf); err != nil {
			return ConfigPath, fmt.Errorf("Cannot write config file: %s \n", err)
		}
	}

	return ConfigPath, nil
}

// initConfig initializes a fresh config file at ~/.config/kable/kableconfig.json, if not yet initialized
func initConfig() error {
	if configSet() {
		return nil
	}
	if err := setCurrentConfig(&Config{
		APIVersion:          1,
		UseKeychainProvider: true,
	}); err != nil {
		return err
	}
	ConfigPath = filepath.Join(UserDir, ConfigFileName)
	if err := writeConfig(ConfigPath); err != nil {
		return err
	}
	return nil
}

func writeConfig(path string) error {
	if !configSet() {
		return errors.ConfigNotInitializedError
	}
	filepath.Dir(ConfigPath)
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
	return nil, errors.ConfigNotInitializedError
}

func setCurrentConfig(config *Config) error {
	if configSet() {
		return errors.ConfigAlreadyInitializedError
	}
	currentConfig = config
	return nil
}

func configSet() bool {
	if currentConfig != nil && ConfigPath != "" {
		return true
	}
	return false
}
