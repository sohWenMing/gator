package database

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
)

var workingDir string = "/home/nindgabeet/workspace/github.com/sohWenMing/gator"

func TestMain(m *testing.M) {
	startCmd := exec.Command("make", "run-pg")
	startCmd.Dir = workingDir
	if err := startCmd.Run(); err != nil {
		fmt.Println("postgres container failed to launch", err)
		os.Exit(1)
	}
	code := m.Run()
	stopCmd := exec.Command("make", "stop-pg")
	stopCmd.Dir = workingDir
	if err := stopCmd.Run(); err != nil {
		fmt.Println("cleanup refused to run", err)
		os.Exit(1)
	}
	os.Exit(code)

}

func TestDBConnection(t *testing.T) {
	user := "postgres"
	password := "postgres"
	dbname := "testdb"
	host := "localhost"
	port := 5432
	connString := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		user, password, dbname, host, port)

	for i := 0; i < 20; i++ {
		_, err := ConnectToDB(connString)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		} else {
			return
		}
	}
	t.Errorf("did not connect in time ")
}
