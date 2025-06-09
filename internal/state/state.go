package state

import (
	"io"

	"github.com/sohWenMing/gator/internal/config"
	"github.com/sohWenMing/gator/internal/database"
)

type State struct {
	config  *config.Config
	writer  io.Writer
	queries *database.Queries
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
func (s *State) SetQueries(q *database.Queries) {
	s.queries = q
}
