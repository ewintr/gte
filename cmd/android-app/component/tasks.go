package component

import (
	"sort"
	"time"

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

func NewTasks(conf *Configuration, tasks []*task.Task) *Tasks {
	local := storage.NewMemory(tasks...)
	remote := storage.NewRemoteRepository(mstore.NewIMAP(conf.IMAP()))
	disp := storage.NewDispatcher(msend.NewSSLSMTP(conf.SMTP()))

	return &Tasks{
		local:  local,
		remote: remote,
		disp:   disp,
	}
}

func (t *Tasks) All() ([]*task.Task, error) {
	lts, err := t.local.FindAll()
	if err != nil {
		return []*task.Task{}, err
	}
	for _, lt := range lts {
		lt.ApplyUpdate()
	}
	ts := []*task.Task{}
	for _, lt := range lts {
		ts = append(ts, &lt.Task)
	}

	return ts, nil
}

func (t *Tasks) Today() (map[string]string, error) {
	reqs := process.ListReqs{
		Due:           task.Today(),
		IncludeBefore: true,
		ApplyUpdates:  true,
	}
	res, err := process.NewList(t.local, reqs).Process()
	if err != nil {
		return map[string]string{}, err
	}
	sort.Sort(task.ByDefault(res.Tasks))

	tasks := map[string]string{}
	for _, t := range res.Tasks {
		tasks[t.Id] = t.Action
	}

	return tasks, nil
}

func (t *Tasks) Sync() (int, int, error) {
	countDisp, err := process.NewSend(t.local, t.disp).Process()
	if err != nil {
		return 0, 0, err
	}

	latestFetch, err := t.local.LatestSync()
	if err != nil {
		return 0, 0, err
	}
	// use unix timestamp for time comparison, because time.Before and
	// time.After depend on a monotonic clock and in Android the
	// monotonic clock stops ticking if the phone is in suspended sleep
	if latestFetch.Add(15*time.Minute).Unix() > time.Now().Unix() {
		return countDisp, 0, nil
	}

	res, err := process.NewFetch(t.remote, t.local).Process()
	if err != nil {
		return countDisp, 0, err
	}
	return countDisp, res.Count, nil
}

func (t *Tasks) MarkDone(id string) error {
	localTask, err := t.local.FindById(id)
	if err != nil {
		return err
	}

	update := &task.LocalUpdate{
		ForVersion: localTask.Version,
		Fields:     []string{task.FIELD_DONE},
		Done:       true,
	}

	return process.NewUpdate(t.local, id, update).Process()
}
