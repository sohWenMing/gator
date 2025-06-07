package main

import (
	"fmt"
	"log"

	"github.com/sohWenMing/gator/internal/config"
	"github.com/sohWenMing/gator/internal/env"
)

func main() {
	readEnvVars, err := env.ReadEnv("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	jsonFilename := readEnvVars.GetConfigJsonPath()
	jsonPath := fmt.Sprintf("../../%s", jsonFilename)
	readConfig, err := config.Read(jsonPath)
	if err != nil {
		log.Fatal(err)
	}

	readConfig.UpdateCurrentUserName("nindgabeet")
	err = config.WriteConfigToFile(*readConfig, jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	readConfig, err = config.Read(jsonPath)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(readConfig.String())

}

// Read the config file.
// Set the current user to "lane" (actually, you should use your name instead) and update the config file on disk.
// Read the config file again and print the contents of the config struct to the terminal.
