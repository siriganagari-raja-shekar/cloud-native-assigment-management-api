package conf

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
)

type Configuration struct {
}

func (c *Configuration) Set() {
	setupEnv()
}

func setupEnv() {
	err := godotenv.Load(getEnvDir())
	if err != nil {
		fmt.Printf("%v\n", err.Error())
	}
}

func getEnvDir() string {
	cwd, _ := os.Getwd()

	for {
		goModPath := filepath.Join(cwd, "go.mod")
		_, err := os.Stat(goModPath)

		if err == nil {
			return filepath.Join(cwd, ".env")
		} else {
			parent := filepath.Dir(cwd)
			if parent == cwd {
				break
			}
			cwd = parent
		}
	}
	return ".env"
}
