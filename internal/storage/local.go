package storage

import (
	"errors"
	"sort"
	"time"

	"ewintr.nl/gte/internal/task"
)

var (
	ErrTaskNotFound = errors.New("task was not found")
)

type LocalRepository interface {
	LatestSyncs() (time.Time, time.Time, error) // last fetch, last dispatch, err
	SetTasks(tasks []*task.Task) error
	FindAll() ([]*task.LocalTask, error)
	FindById(id string) (*task.LocalTask, error)
	FindByLocalId(id int) (*task.LocalTask, error)
	SetLocalUpdate(id string, update *task.LocalUpdate) error
	MarkDispatched(id int) error
	Add(update *task.LocalUpdate) (*task.LocalTask, error)
}

// NextLocalId finds a new local id by incrememting to a variable limit.
//
// When tasks are edited, some get removed because they are done or deleted.
// It is very confusing if existing tasks get renumbered, or if a new one
// immediatly gets the id of an removed one. So it is better to just
// increment. However, local id's also benefit from being short, so we
// don't want to keep incrementing forever.
//
// This function takes a list if id's that are in use and sets the limit
// to the nearest power of ten depening on the current highest id used.
// The new id is an incremented one from that max. However, if the limit
// is reached, it first tries to find "holes" in the current sequence,
// starting from the bottom. If there are no holes, the limit is increased.
func NextLocalId(used []int) int {
	if len(used) == 0 {
		return 1
	}

	sort.Ints(used)
	usedMax := 1
	for _, u := range used {
		if u > usedMax {
			usedMax = u
		}
	}

	var limit int
	for limit = 1; limit <= len(used) || limit < usedMax; limit *= 10 {
	}

	newId := used[len(used)-1] + 1
	if newId < limit {
		return newId
	}

	usedMap := map[int]bool{}
	for _, u := range used {
		usedMap[u] = true
	}

	for i := 1; i < limit; i++ {
		if _, ok := usedMap[i]; !ok {
			return i
		}
	}

	return limit
}

// MergeNewTaskSet updates a local set of tasks with a remote one
//
// The new set is leading and tasks that are not in there get dismissed. Tasks that
// were created locally and got dispatched might temporarily dissappear if the
// remote inbox has a delay in processing.
func MergeNewTaskSet(oldTasks []*task.LocalTask, newTasks []*task.Task) []*task.LocalTask {

	// create lookups
	resultMap := map[string]*task.LocalTask{}
	for _, nt := range newTasks {
		resultMap[nt.Id] = &task.LocalTask{
			Task:        *nt,
			LocalId:     0,
			LocalUpdate: &task.LocalUpdate{},
			LocalStatus: task.STATUS_FETCHED,
		}
	}
	oldMap := map[string]*task.LocalTask{}
	for _, ot := range oldTasks {
		oldMap[ot.Id] = ot
	}

	// apply local id rules:
	// - keep id's that were present in the old set
	// - find new id's for new tasks
	// - assignment of local id's is non deterministic
	var used []int
	for _, ot := range oldTasks {
		if _, ok := resultMap[ot.Id]; ok {
			resultMap[ot.Id].LocalId = ot.LocalId
			used = append(used, ot.LocalId)
		}
	}
	for id, nt := range resultMap {
		if nt.LocalId == 0 {
			newLocalId := NextLocalId(used)
			resultMap[id].LocalId = newLocalId
			used = append(used, newLocalId)
		}
	}

	// apply local update rules:
	// - only keep local updates if the new task hasn't moved to a new version yet
	for _, ot := range oldTasks {
		if nt, ok := resultMap[ot.Id]; ok {
			if ot.LocalUpdate.ForVersion >= nt.Version {
				resultMap[ot.Id].LocalUpdate = ot.LocalUpdate
				resultMap[ot.Id].LocalStatus = task.STATUS_UPDATED
			}
		}
	}

	var result []*task.LocalTask
	for _, nt := range resultMap {
		result = append(result, nt)
	}
	return result
}
