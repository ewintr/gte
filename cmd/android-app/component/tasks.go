package component

import (
	"fmt"
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

	latestFetch, latestDisp, err := t.local.LatestSyncs()
	if err != nil {
		return 0, 0, err
	}
	// use unix timestamp for time comparison, because time.Before and
	// time.After depend on a monotonic clock and on my phone the
	// monotonic clock stops ticking when it goes to suspended sleep
	if latestFetch.Add(15*time.Minute).Unix() > time.Now().Unix() || latestDisp.Add(2*time.Minute).Unix() > time.Now().Unix() {
		return countDisp, 0, nil
	}

	res, err := process.NewFetch(t.remote, t.local, task.FOLDER_PLANNED).Process()
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

func (t *Tasks) Add(fields map[string]string) error {
	update := &task.LocalUpdate{
		Fields: []string{},
	}
	if len(fields["action"]) != 0 {
		update.Action = fields["action"]
		update.Fields = append(update.Fields, task.FIELD_ACTION)
	}
	if len(fields["project"]) != 0 {
		update.Project = fields["project"]
		update.Fields = append(update.Fields, task.FIELD_PROJECT)
	}
	due := task.NewDateFromString(fields["due"])
	if !due.IsZero() {
		update.Due = due
		update.Fields = append(update.Fields, task.FIELD_DUE)
	}
	recur := task.NewRecurrer(fields["recur"])
	if recur != nil {
		update.Recur = recur
		update.Fields = append(update.Fields, task.FIELD_RECUR)
	}
	if len(update.Fields) == 0 {
		return fmt.Errorf("no fields in new task")
	}

	if err := process.NewNew(t.local, update).Process(); err != nil {
		return err
	}

	return nil
}

func (t *Tasks) Update(id, newDue string) error {
	due := task.NewDateFromString(newDue)
	localTask, err := t.local.FindById(id)
	if err != nil {
		return err
	}
	update := &task.LocalUpdate{
		ForVersion: localTask.Version,
		Due:        due,
		Fields:     []string{task.FIELD_DUE},
	}

	return process.NewUpdate(t.local, localTask.Id, update).Process()
}
