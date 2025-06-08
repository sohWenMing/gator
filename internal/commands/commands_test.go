package commands

import (
	"reflect"
	"testing"
)

func TestParseCommand(t *testing.T) {
	commandMap := InitCommandMap()

	type test struct {
		name          string
		expectedCmd   Command
		expectedArgs  []string
		userInputArgs []string
		isErrExpected bool
	}

	tests := []test{
		{
			"test basic get login command",
			LoginCommand,
			[]string{"nindgabeet"},
			[]string{"gator", "login", "nindgabeet"},
			false,
		},
		{
			"test empty input",
			Command{},
			[]string{},
			[]string{},
			true,
		},
		{
			"test default to help, no arg after exectuable",
			HelpCommand,
			[]string{},
			[]string{"gator"},
			false,
		},
		{
			"test default to help, cmd not found ",
			HelpCommand,
			[]string{"arg1", "arg2"},
			[]string{"gator", "failCmd", "arg1", "arg2"},
			false,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			cmd, args, err := commandMap.ParseCommand(test.userInputArgs)
			switch test.isErrExpected {
			case true:
				if err == nil {
					t.Errorf("expected error, didn't get one")
					return
				}
			default:
				if err != nil {
					t.Errorf("didn't expect error, got %v", err)
					return
				}
			}
			if cmd.name != test.expectedCmd.name || cmd.description != test.expectedCmd.description {
				t.Errorf("\ngot %v\nwant %v", cmd, test.expectedCmd)
				return
			}
			if !reflect.DeepEqual(args, test.expectedArgs) {
				t.Errorf("\ngot %v\nwant %v", args, test.expectedArgs)
				return
			}
		})
	}
}
