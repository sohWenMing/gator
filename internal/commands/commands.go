package commands

import (
	"errors"

	"github.com/sohWenMing/gator/internal/state"
	"github.com/sohWenMing/gator/internal/utils"
)

type Command struct {
	name        string
	description string
	CallBack    func(s *state.State, args []string) error
}

func (c *Command) GetName() string {
	return c.name
}

type CommandMap map[string]Command

var (
	LoginCommand = Command{
		"login",
		"sets the logged in user to be set to what the user enetered",
		LoginCallBack,
	}
	HelpCommand = Command{
		"help",
		"prints the help and usage messages",
		HelpCallBack,
	}
)

func (cm CommandMap) ParseCommand(osArgs []string) (parsedCommand Command, args []string, err error) {
	if len(osArgs) < 1 {
		return Command{}, []string{}, errors.New("osArgs could not be parsed from user input")
	}
	// first arg should be just the name of the executable
	if len(osArgs) == 1 {
		return HelpCommand, []string{}, nil
		// if program ran and there was only one arg, that's name of executable
	}
	cmdString := utils.TrimToLower(osArgs[1])
	args = osArgs[2:]
	if cmd, ok := cm[cmdString]; ok {
		return cmd, args, nil
	} else {
		return HelpCommand, args, nil
	}
}

/*
	ParseCommand should only have the responsibility of parsing the user input and processing the command and the arguments, delegating all other
	functionality to later functions
*/

func InitCommandMap() CommandMap {
	returnedMap := CommandMap(make(map[string]Command))
	returnedMap["login"] = LoginCommand
	return returnedMap
}

func LoginCallBack(s *state.State, args []string) error {
	if len(args) != 1 {
		return errors.New("only on argument being the user to be logged in should be entered.")
	}
	userNameToUpdate := args[0]
	s.GetConfig().UpdateCurrentUserName(userNameToUpdate)
	return nil

}
func HelpCallBack(s *state.State, args []string) error {
	return nil

}
