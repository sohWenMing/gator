package integration

import (
	"context"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/sohWenMing/gator/internal/config"
	"github.com/sohWenMing/gator/internal/database"
	"github.com/sohWenMing/gator/internal/env"
	"github.com/sohWenMing/gator/internal/helper"
	"github.com/sohWenMing/gator/internal/state"
)

func LoadState(writer io.Writer) (returnedState *state.State) {
	returnedState = state.InitState(writer)
	return returnedState
}

func GetEnvVars(envPath string) (*env.EnvVars, error) {
	readEnvVars, err := env.ReadEnv(envPath)
	if err != nil {
		return nil, err
	}
	return readEnvVars, nil
}

func ConnectToDB(connectionString string) (*database.Queries, error) {
	queries, err := database.ConnectToDB(connectionString)
	if err != nil {
		return nil, err
	}
	return queries, nil
}

func PingDB(s *state.State) error {
	proceedChan := make(chan struct{})
	exitChan := make(chan error)

	go func(proceedChan chan<- struct{}, exitChan chan<- error) {
		for i := 0; i < 60; i++ {
			contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
			_, err := s.GetQueries().Ping(contextStruct.Context)
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			} else {
				proceedChan <- struct{}{}
				return
			}
		}
		exitChan <- errors.New("ping was not successful in time: program requires termination")
		return
	}(proceedChan, exitChan)

	select {
	case err := <-exitChan:
		return err
	case <-proceedChan:
		return nil
	}

}

func LoadConfigAndSetJson(relDir string, envVars *env.EnvVars) (*config.Config, error) {

	jsonFilename := envVars.GetConfigJsonPath()
	jsonPath := fmt.Sprintf("%s%s", relDir, jsonFilename)
	absJsonPath, err := filepath.Abs(jsonPath)
	if err != nil {
		return nil, err
	}
	cfg, err := config.Read(absJsonPath)
	if err != nil {
		return nil, err
	}
	cfg.SetJsonPath(absJsonPath)
	return cfg, nil
}
