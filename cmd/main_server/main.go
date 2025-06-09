package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sohWenMing/gator/internal/commands"
	"github.com/sohWenMing/gator/internal/config"
	"github.com/sohWenMing/gator/internal/database"
	"github.com/sohWenMing/gator/internal/env"
	"github.com/sohWenMing/gator/internal/state"
	"github.com/sohWenMing/gator/internal/utils"
)

func main() {
	envPath := os.Getenv("ENVPATH")
	if envPath == "" {
		fmt.Println("ENVPATH not found in environment, defaulting to .dev")
	} else {
		fmt.Println("Found envPath: ", envPath)
	}

	readEnvVars, err := env.ReadEnv(envPath)
	if err != nil {
		log.Fatal(err)
	}
	// first read .env file

	jsonFilename := readEnvVars.GetConfigJsonPath()
	jsonPath := fmt.Sprintf("../../%s", jsonFilename)
	cfg, err := config.Read(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	/*
		get the config from the json file - as calculated from the .env file,
		relative to where this file is located within the project
	*/

	state := state.InitState(os.Stdout)
	state.SetConfig(cfg)

	queries, err := database.ConnectToDB(readEnvVars.GetDBConnectionString())
	if err != nil {
		log.Fatal(err)
	} else {
		utils.WriteLine(state.GetWriter(), "Connection to database established")
	}

	state.SetQueries(queries)

	commandMap := commands.InitCommandMap()
	parsedCommand, args, err := commandMap.ParseCommand(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	err = parsedCommand.CallBack(state, args)
	if err != nil {
		log.Fatal(err)
	}

	if parsedCommand.GetName() == "login" {
		err = config.WriteConfigToFile(*cfg, jsonPath)
		if err != nil {
			log.Fatal(err)
		} else {
			updatedUsername := args[0]
			prompt := fmt.Sprintf("logged in user has been updated to %s", updatedUsername)
			utils.WriteLine(state.GetWriter(), prompt)
		}
	}

	cfg, err = config.Read(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg.String())

}

// Read the config file.
// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
// Read the config file again and print the contents of the config struct to the terminal.
