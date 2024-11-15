package configparser

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Exclusions Exclusions `yaml:"exclusions"`
	Changes    Changes    `yaml:"changes"`
	FolderPath string     `yaml:"folder_path"`
  CountFunctions int `yaml:"count_functions"`
  MaxActionsPerFunc int `yaml:"max_count_actions_per_func"`
	Insertions []string   `yaml:"insertions"`
}

type Exclusions struct {
	Extensions []string `yaml:"extensions"`
	Files      []string `yaml:"files"`
	Folders    []string `yaml:"folders"`
}

type Changes struct {
	Classes    []string `yaml:"classes"`
	Funcs      []string `yaml:"funcs"`
	Enums      []string `yaml:"enums"`
	Structs    []string `yaml:"structs"`
	Extensions []string `yaml:"extensions"`
}

func MustLoadByPath(path string) Config {
	stream, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var cfg Config
	err = yaml.Unmarshal(stream, &cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
