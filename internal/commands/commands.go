package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sohWenMing/gator/internal/database"
	"github.com/sohWenMing/gator/internal/helper"
	rssparser "github.com/sohWenMing/gator/internal/rss_parser"
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
	AggCommand = Command{
		"agg",
		"gets the information from the rss aggreagator",
		AggCallBack,
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
	returnedMap["agg"] = AggCommand
	return returnedMap
}

func AggCallBack(s *state.State, args []string) error {
	fmt.Println("agg callback was called")
	var wg sync.WaitGroup
	urls := s.GetAggregatorUrls()

	type ParsedRssFeedToError struct {
		rssFeed rssparser.ParsedRssFeed
		err     error
	}

	rssFeedToErrChan := make(chan ParsedRssFeedToError)
	for _, url := range urls {
		fmt.Println("checking url: ", url)
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelFunc()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
			req.Header.Set("User-Agent", "gator")
			if err != nil {
				rssFeedToErrChan <- ParsedRssFeedToError{
					rssparser.ParsedRssFeed{}, err,
				}
				return
			}
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				rssFeedToErrChan <- ParsedRssFeedToError{
					rssparser.ParsedRssFeed{}, err,
				}
				return
			}
			resBytes, err := io.ReadAll(res.Body)
			if err != nil {
				rssFeedToErrChan <- ParsedRssFeedToError{
					rssparser.ParsedRssFeed{}, err,
				}
				return
			}
			parsedRSS, err := rssparser.ParseRSSXML(resBytes)
			if err != nil {
				rssFeedToErrChan <- ParsedRssFeedToError{
					rssparser.ParsedRssFeed{}, err,
				}
				return
			}
			rssFeedToErrChan <- ParsedRssFeedToError{
				parsedRSS, nil,
			}
			return
		}(url)
	}
	go func() {
		wg.Wait()
		close(rssFeedToErrChan)
	}()
	for result := range rssFeedToErrChan {
		if result.err == nil {
			utils.WriteLine(s.GetWriter(), result.rssFeed.String())
		}
	}
	return nil
}

func LoginCallBack(s *state.State, args []string) error {
	if len(args) != 1 {
		return errors.New("number of args passed into login command should only be 1, being the user to login")
	}
	nameForLogin := args[0]
	contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
	user, err := s.GetQueries().GetUser(contextStruct.Context, nameForLogin)
	defer contextStruct.CancelFunc()

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
	contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
	err := s.GetQueries().DeleteAllUsers(contextStruct.Context)
	defer contextStruct.CancelFunc()
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
	contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
	users, err := s.GetQueries().GetAllUsers(contextStruct.Context)
	defer contextStruct.CancelFunc()
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

	contextStruct := helper.SpawnTimeOutContext(context.Background(), 10*time.Second)
	user, err := s.GetQueries().CreateUser(contextStruct.Context, database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      args[0],
	})
	defer contextStruct.CancelFunc()
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
