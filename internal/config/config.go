package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Config struct {
	LoggedInUserName string `json:"current_username"`
	jsonPath         string
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

func (c *Config) UpdateLoggedInUserName(username string) {
	c.LoggedInUserName = username
}
func (c *Config) SetJsonPath(jsonPath string) {
	c.jsonPath = jsonPath
}
func (c *Config) GetJsonPath() string {
	return c.jsonPath
}

func (c *Config) String() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintln("{"))
	b.WriteString(fmt.Sprintln("CurrentUserName: ", c.LoggedInUserName))
	b.WriteString(fmt.Sprintln("jsonPath: ", c.jsonPath))
	b.WriteString(fmt.Sprintln("}"))
	return b.String()
}
