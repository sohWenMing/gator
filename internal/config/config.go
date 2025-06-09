package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	CurrentUserName string `json:"current_username"`
}

func Read(path string) (c *Config, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, _ := io.ReadAll(file)
	var configVar Config
	err = json.Unmarshal(bytes, &configVar)
	if err != nil {
		return nil, err
	}
	return &configVar, nil
}

func (c *Config) UpdateCurrentUserName(username string) {
	c.CurrentUserName = username
}

func (c *Config) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintln("{"))
	b.WriteString(fmt.Sprintln("CurrentUserName: ", c.CurrentUserName))
	b.WriteString(fmt.Sprintln("}"))
	return b.String()
}

func WriteConfigToFile(c Config, filePath string) error {
	marshalledConfig, err := json.Marshal(c)
	if err != nil {
		return err
	}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("error occured here")
		return err
	}
	defer file.Close()
	return os.WriteFile(filePath, marshalledConfig, 0644)
}
