package env

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	configJsonPath string
}

func ReadEnv(path string) (e *EnvVars, err error) {
	err = godotenv.Load(path)
	if err != nil {
		return nil, err
	}
	return &EnvVars{
		configJsonPath: os.Getenv("CONFIG_JSON_PATH"),
	}, nil
}

func (e *EnvVars) GetConfigJsonPath() string {
	return e.configJsonPath
}
