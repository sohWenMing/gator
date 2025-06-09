package env

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type EnvVars struct {
	dbConnectionString string
	configJsonPath     string
}

func ReadEnv(path string) (e *EnvVars, err error) {
	err = godotenv.Load(path)
	if err != nil {
		return nil, err
	}
	configJsonPath := os.Getenv("CONFIG_JSON_PATH")
	dbConnectionString := os.Getenv("CONNSTRING")

	if configJsonPath == "" || dbConnectionString == "" {
		var b strings.Builder
		b.WriteString(fmt.Sprintln("configurations could not be loaded"))
		b.WriteString(fmt.Sprintln("configJsonPath: ", configJsonPath))
		b.WriteString(fmt.Sprintln("dbConnectionString", dbConnectionString))
		return nil, errors.New(b.String())
	}
	// first: try to load the envfile passed in
	return &EnvVars{
		configJsonPath:     configJsonPath,
		dbConnectionString: dbConnectionString,
	}, nil
}

func (e *EnvVars) GetConfigJsonPath() string {
	return e.configJsonPath
}
func (e *EnvVars) GetDBConnectionString() string {
	return e.dbConnectionString
}
