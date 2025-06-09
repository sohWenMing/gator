package main

import (
	"fmt"
	"log"
	"os"
	"time"

	integration "github.com/sohWenMing/gator/integration"
	"github.com/sohWenMing/gator/internal/commands"
	"github.com/sohWenMing/gator/internal/utils"
)

func main() {
	envPath := os.Getenv("ENVPATH")
	if envPath == "" {
		fmt.Println("ENVPATH not found in environment, defaulting to .env")
	} else {
		fmt.Println("Found envPath: ", envPath)
	}

	readEnvVars, err := integration.GetEnvVars(envPath)
	if err != nil {
		log.Fatal(err)
	}
	// first read .env file
	cfg, err := integration.LoadConfigAndSetJson("../../", readEnvVars)
	if err != nil {
		log.Fatal(err)
	}
	/*
		get the config from the json file - as calculated from the .env file,
		relative to where this file is located within the project
	*/

	state := integration.LoadState(os.Stdout)
	// run the initial state, attaching os.Stdout as the writer

	state.SetConfig(cfg)
	//attach the config initialised to the state

	queries, err := integration.ConnectToDB(readEnvVars.GetDBConnectionString())

	if err != nil {
		log.Fatal(err)
	} else {
		utils.WriteLine(state.GetWriter(), "Connection to database established")
	}

	state.SetQueries(queries)
	// attach the queries from the dbConnection to state, to be accessed by functions

	startPingTime := time.Now()
	err = integration.PingDB(state)
	if err != nil {
		log.Fatal(err)
	} else {
		timetaken := time.Since(startPingTime)
		timeTakenLine := fmt.Sprintln("Ping succeeded. Time taken: ", timetaken)
		utils.WriteLine(state.GetWriter(), timeTakenLine)
	}
	// initiate ping test - ping is to make sure that database is ready to accept connectoins before moving on with anything else

	commandMap := commands.InitCommandMap()
	parsedCommand, args, err := commandMap.ParseCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	err = parsedCommand.CallBack(state, args)
	if err != nil {
		fmt.Println(err)
	}

}

// Read the config file.
// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
// Read the config file again and print the contents of the config struct to the terminal.
