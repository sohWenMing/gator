package integrationtesting

import (
	"fmt"
	"testing"

	"github.com/sohWenMing/gator/internal/config"
	"github.com/sohWenMing/gator/internal/env"
)

func TestLoadConfigFromEnv(t *testing.T) {
	envVar, err := env.ReadEnv("../.env")
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
	}
	/*
		Because it's running from within test folder and running from
		the main file would be different, we need to concat one "." to
		make the relative path correct
	*/
	jsonPath := envVar.GetConfigJsonPath()
	pathToRead := fmt.Sprintf("../%s", jsonPath)

	readConfig, err := config.Read(pathToRead)
	if err != nil {
		t.Errorf("didn't expect error, got %v", err)
		return
	}
	expectedDBUrl := "postgres://example"
	if readConfig.DbUrl != expectedDBUrl {
		t.Errorf("\ngot: %s\nwant: %s", readConfig.DbUrl, expectedDBUrl)
	}

}
