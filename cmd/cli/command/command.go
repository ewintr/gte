package command

import (
	"errors"
	"fmt"

	"git.ewintr.nl/gte/internal/configuration"
	"git.ewintr.nl/gte/internal/process"
	"git.ewintr.nl/gte/internal/storage"
	"git.ewintr.nl/gte/internal/task"
	"git.ewintr.nl/gte/pkg/mstore"
)

var (
	ErrInitCommand   = errors.New("could not initialize command")
	ErrFailedCommand = errors.New("could not execute command")
)

type Result struct {
	Message string
}

type Command interface {
	Do() (Result, error)
}

func Parse(args []string, conf *configuration.Configuration) (Command, error) {
	if len(args) == 0 {
		return NewEmpty()
	}

	cmd, _ := args[0], args[1:]
	switch cmd {
	case "sync":
		return NewSync(conf)
	case "today":
		return NewToday(conf)
	default:
		return NewEmpty()
	}
}

type Empty struct{}

func NewEmpty() (*Empty, error) {
	return &Empty{}, nil
}

func (cmd *Empty) Do() (Result, error) {
	return Result{
		Message: "did nothing\n",
	}, nil
}

type Sync struct {
	syncer *process.Sync
}

func NewSync(conf *configuration.Configuration) (*Sync, error) {
	msgStore := mstore.NewIMAP(conf.IMAP())
	remote := storage.NewRemoteRepository(msgStore)
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Sync{}, fmt.Errorf("%w: %v", ErrInitCommand, err)
	}
	syncer := process.NewSync(remote, local)

	return &Sync{
		syncer: syncer,
	}, nil
}

func (s *Sync) Do() (Result, error) {
	result, err := s.syncer.Process()
	if err != nil {
		return Result{}, fmt.Errorf("%w: %v", ErrFailedCommand, err)
	}

	return Result{
		Message: fmt.Sprintf("synced %d tasks\n", result.Count),
	}, nil
}

type Today struct {
	local storage.LocalRepository
}

func NewToday(conf *configuration.Configuration) (*Today, error) {
	local, err := storage.NewSqlite(conf.Sqlite())
	if err != nil {
		return &Today{}, fmt.Errorf("%w: %v", ErrInitCommand, err)
	}

	return &Today{
		local: local,
	}, nil
}

func (t *Today) Do() (Result, error) {
	tasks, err := t.local.FindAllInFolder(task.FOLDER_PLANNED)
	if err != nil {
		return Result{}, fmt.Errorf("%w: %v", ErrFailedCommand, err)
	}

	todayTasks := []*task.Task{}
	for _, t := range tasks {
		if t.Due == task.Today || task.Today.After(t.Due) {
			todayTasks = append(todayTasks, t)
		}
	}

	if len(todayTasks) == 0 {
		return Result{
			Message: "nothing left",
		}, nil
	}

	var msg string
	for _, t := range todayTasks {
		msg += fmt.Sprintf("%s - %s\n", t.Project, t.Action)
	}

	return Result{
		Message: msg,
	}, nil
}
