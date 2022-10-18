package component

import (
	"sort"

	"ewintr.nl/gte/internal/process"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
	"ewintr.nl/gte/pkg/msend"
	"ewintr.nl/gte/pkg/mstore"
)

type Tasks struct {
	local  storage.LocalRepository
	remote *storage.RemoteRepository
	disp   *storage.Dispatcher
}

func NewTasks(conf *Configuration) (*Tasks, error) {
	local := storage.NewMemory()
	remote := storage.NewRemoteRepository(mstore.NewIMAP(conf.IMAP()))
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	return &Tasks{
		local:  local,
		remote: remote,
		disp:   disp,
	}, nil
}

func (t *Tasks) Today() ([]string, error) {
	reqs := process.ListReqs{
		Due:           task.Today(),
		IncludeBefore: true,
		ApplyUpdates:  true,
	}
	res, err := process.NewList(t.local, reqs).Process()
	if err != nil {
		return []string{}, err
	}
	sort.Sort(task.ByDefault(res.Tasks))

	tasks := []string{}
	for _, t := range res.Tasks {
		tasks = append(tasks, t.Action)
	}

	return tasks, nil
}

func (t *Tasks) Sync() (int, error) {
	res, err := process.NewFetch(t.remote, t.local).Process()
	if err != nil {
		return 0, err
	}
	return res.Count, nil
}
