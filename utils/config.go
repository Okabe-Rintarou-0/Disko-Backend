package utils

import (
	"gopkg.in/yaml.v2"
	"os"
)

func ReadConfig(path string, conf any) {
	file, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err = yaml.Unmarshal(file, conf); err != nil {
		panic(err)
	}
}
