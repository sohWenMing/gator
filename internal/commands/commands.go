package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sohWenMing/gator/internal/database"
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
	RegisterCommand = Command{
		"register",
		"registers a new user to the application",
		RegisterCallBack,
	}
	ResetCommand = Command{
		"reset",
		"resets the database to an empty state",
		ResetCallBack,
	}
	UsersCommand = Command{
		"users",
		"lists all the users that are available in the system",
		UsersCallBack,
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
	returnedMap["register"] = RegisterCommand
	returnedMap["reset"] = ResetCommand
	returnedMap["users"] = UsersCommand
	return returnedMap
}

func LoginCallBack(s *state.State, args []string) error {
	if len(args) != 1 {
		return errors.New("number of args passed into login command should only be 1, being the user to login")
	}
	nameForLogin := args[0]
	user, err := s.GetQueries().GetUser(s.GetStateContext().Context, nameForLogin)
	if err != nil {
		errMsg := err.Error()
		if isNotFound := strings.Contains(errMsg, "no rows in result set"); isNotFound {
			msg := fmt.Sprintf("user with name %s cannot be found", nameForLogin)
			utils.WriteLine(s.GetWriter(), fmt.Sprintln(msg))
			utils.WriteLine(s.GetWriter(), fmt.Sprintln("current user logged in:", s.GetConfig().LoggedInUserName))
			return nil
		} else {
			return err
		}
	}
	userNameToUpdate := user.Name
	s.GetConfig().UpdateLoggedInUserName(userNameToUpdate)
	err = WriteConfigToFile(s)
	if err != nil {
		return err
	}
	utils.WriteLine(s.GetWriter(), fmt.Sprintln("current user logged in:", s.GetConfig().LoggedInUserName))
	return nil
}
func ResetCallBack(s *state.State, args []string) error {
	if len(args) != 0 {
		return errors.New("reset cannot be called with arguments")
	}
	err := s.GetQueries().DeleteAllUsers(s.GetStateContext().Context)
	if err != nil {
		return err
	}
	s.GetConfig().LoggedInUserName = ""
	err = WriteConfigToFile(s)
	if err != nil {
		return err
	}
	utils.WriteLine(s.GetWriter(), "database has been reset to blank slate")
	return nil
}
func UsersCallBack(s *state.State, args []string) error {
	if len(args) != 0 {
		return errors.New("users cannot be called with arguments")
	}
	users, err := s.GetQueries().GetAllUsers(s.GetStateContext().Context)
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.GetConfig().LoggedInUserName {
			utils.WriteLine(s.GetWriter(), fmt.Sprintf("* %s (current)", user.Name))
		} else {
			utils.WriteLine(s.GetWriter(), fmt.Sprintf("* %s", user.Name))
		}
	}
	return nil
}
func HelpCallBack(s *state.State, args []string) error {
	return nil

}
func RegisterCallBack(s *state.State, args []string) error {
	if len(args) != 1 {
		return errors.New("number of args passed into register command should only be 1, being the user to register")
	}

	user, err := s.GetQueries().CreateUser(s.GetStateContext().Context, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      args[0],
	})
	if err != nil {
		errMsg := err.Error()
		if isDup := strings.Contains(errMsg, "duplicate key value violates unique constraint"); isDup == true {
			message := fmt.Sprintln("user has already been created with name: ", args[0])
			utils.WriteLine(s.GetWriter(), message)
			return nil
		} else {
			return err
		}
	} else {
		message := fmt.Sprintln("user has been created with name: ", user.Name)
		utils.WriteLine(s.GetWriter(), message)
		return nil
	}
}

func WriteConfigToFile(s *state.State) error {
	marshalledConfig, err := json.Marshal(s.GetConfig())
	if err != nil {
		return err
	}
	file, err := os.Open(s.GetConfig().GetJsonPath())
	if err != nil {
		fmt.Println("error occured here")
		return err
	}
	defer file.Close()
	return os.WriteFile(s.GetConfig().GetJsonPath(), marshalledConfig, 0644)
}
