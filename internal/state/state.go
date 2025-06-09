package state

import (
	"io"

	"github.com/sohWenMing/gator/internal/config"
)

type State struct {
	config *config.Config
	writer io.Writer
}

func InitState(w io.Writer) *State {
	return &State{
		writer: w,
	}
}

//writer passed into InitState so that state can also be flexibly used for testing, writing to places other that os.Stdout

func (s *State) SetConfig(c *config.Config) {
	s.config = c
}
func (s *State) GetConfig() *config.Config {
	return s.config
}
func (s *State) GetWriter() io.Writer {
	return s.writer
}
