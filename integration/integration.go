package integration

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/sohWenMing/gator/internal/database"
	"github.com/sohWenMing/gator/internal/env"
	"github.com/sohWenMing/gator/internal/state"
	"github.com/sohWenMing/gator/internal/utils"
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
		for i := range 60 {
			_, err := s.GetQueries().Ping(s.GetStateContext().Context)
			if err != nil {
				pingFailLine := fmt.Sprintln("attempt to ping db failed: ", i)
				utils.WriteLine(s.GetWriter(), pingFailLine)
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
