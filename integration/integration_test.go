package integration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sohWenMing/gator/internal/database"
	"github.com/sohWenMing/gator/internal/env"
	"github.com/sohWenMing/gator/internal/state"
)

var (
	testState   *state.State
	testEnvVars *env.EnvVars
	buf         *bytes.Buffer = &bytes.Buffer{}
	workingDir  string        = "/home/nindgabeet/workspace/github.com/sohWenMing/gator"
)

func TestMain(m *testing.M) {
	startCmd := exec.Command("make", "run-pg-dev")
	startCmd.Dir = workingDir
	if err := startCmd.Run(); err != nil {
		fmt.Println("postgres container failed to launch: ", err)
		os.Exit(1)
	}
	testState = LoadState(buf)
	testEnvVars, err := GetEnvVars("../.env")
	if err != nil {
		fmt.Println("env vars failed to load: ", err)
		os.Exit(1)
	}
	queries, err := ConnectToDB(testEnvVars.GetDBConnectionString())
	if err != nil {
		fmt.Println("error loading queries: ", err)
		os.Exit(1)
	}
	testState.SetQueries(queries)

	err = PingDB(testState)
	if err != nil {
		fmt.Println("Ping db failed: ", err)
		os.Exit(1)
	} else {
		fmt.Println("All setup passed for integration test")
	}
	code := m.Run()

	err = testState.GetQueries().DeleteAllUsers(testState.GetStateContext().Context)
	if err != nil {
		fmt.Println("error on cleanup of DB: ", err)
	}
	stopCmd := exec.Command("make", "stop-pg-dev")
	stopCmd.Dir = workingDir
	if err := stopCmd.Run(); err != nil {
		fmt.Println("cleanup refused to run", err)
		os.Exit(1)
	}
	os.Exit(code)
}

func TestCreateUser(t *testing.T) {
	type test struct {
		testName      string
		inputName     string
		isErrExpected bool
	}

	tests := []test{
		{
			"initial test with empty database",
			"nindgabeet",
			false,
		},
		{
			"duplicate should fail",
			"nindgabeet",
			true,
		}, {
			"null should fail",
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			createUserParams := database.CreateUserParams{
				ID:        uuid.New(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Name:      test.inputName,
			}
			user, err := testState.GetQueries().CreateUser(
				testState.GetStateContext().Context, createUserParams,
			)
			switch test.isErrExpected {
			case false:
				if err != nil {
					t.Errorf("didn't expect err, got %v", err)
				}
				if user.Name != test.inputName {
					t.Errorf("\ngot: %s\nwant: %s", user.Name, test.inputName)
				}
			case true:
				if err == nil {
					fmt.Println("user name: ", user.Name)
					t.Errorf("expected err, didn't get one")
				}
			}
		})
	}
}
