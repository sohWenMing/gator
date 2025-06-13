package integration

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sohWenMing/gator/internal/commands"
	"github.com/sohWenMing/gator/internal/database"
	"github.com/sohWenMing/gator/internal/env"
	"github.com/sohWenMing/gator/internal/helper"
	"github.com/sohWenMing/gator/internal/state"
)

var (
	testState   *state.State
	testEnvVars *env.EnvVars
	buf         *bytes.Buffer = &bytes.Buffer{}
	workingDir  string        = "/home/nindgabeet/workspace/github.com/sohWenMing/gator"
	commandMap  commands.CommandMap
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

	cfg, err := LoadConfigAndSetJson("../", testEnvVars)
	if err != nil {
		fmt.Println("Load config failed", err)
		os.Exit(1)
	}

	testState.SetConfig(cfg)
	commandMap = commands.InitCommandMap()

	code := m.Run()

	contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
	defer contextStruct.CancelFunc()
	err = testState.GetQueries().DeleteAllUsers(contextStruct.Context)
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
			contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
			defer contextStruct.CancelFunc()
			user, err := testState.GetQueries().CreateUser(
				contextStruct.Context, createUserParams,
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

	contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
	defer contextStruct.CancelFunc()
	err := testState.GetQueries().DeleteAllUsers(contextStruct.Context)
	if err != nil {
		t.Errorf("clearing of all users failed")
		return
	}
}

func TestListAllUsers(t *testing.T) {
	usernames := []string{"test1", "test2", "test3"}
	for _, username := range usernames {
		userToCreate := database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      username,
		}
		contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
		_, err := testState.GetQueries().CreateUser(contextStruct.Context, userToCreate)
		contextStruct.CancelFunc()
		if err != nil {
			t.Errorf("error on creating user: %v", err)
			return
		}
	}

	loginArgs := []string{"gator", "login", "test2"}
	loginCmd, args, err := commandMap.ParseCommand(loginArgs)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	loginCallBack := commandMap[loginCmd.GetName()].CallBack
	err = loginCallBack(testState, args)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	want := "current user logged in: test2"
	got := strings.ReplaceAll(buf.String(), "\n", "")
	if got != want {
		t.Errorf("\ngot: %q\nwant: %q", got, want)
	}
	buf.Reset()
	userArgs := []string{"gator", "users"}
	userCmd, args, err := commandMap.ParseCommand(userArgs)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	usersCallBack := commandMap[userCmd.GetName()].CallBack
	err = usersCallBack(testState, args)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	got = buf.String()
	checkStrings := []string{"test1", "test2 (current)", "test3"}
	for _, checkString := range checkStrings {
		if !strings.Contains(got, checkString) {
			t.Errorf("%s could not be found in output of users callback", checkString)
		}
	}
}
