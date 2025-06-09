package state

import (
	"context"
	"io"
	"time"

	"github.com/sohWenMing/gator/internal/config"
	"github.com/sohWenMing/gator/internal/database"
)

type State struct {
	config       *config.Config
	writer       io.Writer
	queries      *database.Queries
	stateContext StateContext
}

type StateContext struct {
	Context    context.Context
	CancelFunc context.CancelFunc
}

func InitState(w io.Writer) *State {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	return &State{
		writer: w,
		stateContext: StateContext{
			ctx, cancelFunc,
		},
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
func (s *State) GetQueries() *database.Queries {
	return s.queries
}
func (s *State) GetStateContext() *StateContext {
	return &s.stateContext
}
