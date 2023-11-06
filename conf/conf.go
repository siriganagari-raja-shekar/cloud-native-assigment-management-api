package conf

import (
	"csye6225-mainproject/log"
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
	logger := log.GetLoggerInstance()
	err := godotenv.Load(getEnvDir())
	if err != nil {
		logger.Warn(fmt.Sprintf("%v", err.Error()))
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
