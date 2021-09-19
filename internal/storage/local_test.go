package storage_test

import (
	"sort"
	"testing"

	"ewintr.nl/go-kit/test"
	"ewintr.nl/gte/internal/storage"
	"ewintr.nl/gte/internal/task"
)

func TestNextLocalId(t *testing.T) {
	for _, tc := range []struct {
		name string
		used []int
		exp  int
	}{
		{
			name: "empty",
			used: []int{},
			exp:  1,
		},
		{
			name: "not empty",
			used: []int{5},
			exp:  6,
		},
		{
			name: "multiple",
			used: []int{2, 3, 4},
			exp:  5,
		},
		{
			name: "holes",
			used: []int{1, 5, 8},
			exp:  9,
		},
		{
			name: "expand limit",
			used: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			exp:  11,
		},
		{
			name: "wrap if possible",
			used: []int{8, 9},
			exp:  1,
		},
		{
			name: "find hole",
			used: []int{1, 2, 3, 4, 5, 7, 8, 9},
			exp:  6,
		},
		{
			name: "dont wrap if expanded before",
			used: []int{15, 16},
			exp:  17,
		},
		{
			name: "do wrap if expanded limit is reached",
			used: []int{99},
			exp:  1,
		},
		{
			name: "sync bug",
			used: []int{151, 956, 955, 150, 154, 155, 145, 144,
				136, 152, 148, 146, 934, 149, 937, 135, 140, 139,
				143, 137, 153, 939, 138, 953, 147, 141, 938, 142,
			},
			exp: 957,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			test.Equals(t, tc.exp, storage.NextLocalId(tc.used))
		})
	}
}

func TestMergeNewTaskSet(t *testing.T) {
	task1 := &task.Task{Id: "id-1", Version: 1, Action: "action-1"}
	task1v2 := &task.Task{Id: "id-1", Version: 2, Action: "action-1v2"}
	task2 := &task.Task{Id: "id-2", Version: 2, Action: "action-2"}
	emptyUpdate := &task.LocalUpdate{}

	t.Run("local ids are added", func(t *testing.T) {
		act1 := storage.MergeNewTaskSet([]*task.LocalTask{}, []*task.Task{task1})
		test.Assert(t, len(act1) == 1, "length was not 1")
		test.Equals(t, 1, act1[0].LocalId)

		act2 := storage.MergeNewTaskSet(act1, []*task.Task{task1, task2})
		var actIds []int
		for _, t := range act2 {
			actIds = append(actIds, t.LocalId)
		}
		sort.Ints(actIds)
		test.Equals(t, []int{1, 2}, actIds)
	})

	for _, tc := range []struct {
		name     string
		oldTasks []*task.LocalTask
		newTasks []*task.Task
		exp      []*task.LocalTask
	}{
		{
			name:     "add tasks and find local ids",
			oldTasks: []*task.LocalTask{},
			newTasks: []*task.Task{task1, task2},
			exp: []*task.LocalTask{
				{Task: *task1, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED},
				{Task: *task2, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED},
			},
		},
		{
			name: "update existing task",
			oldTasks: []*task.LocalTask{
				{Task: *task1, LocalUpdate: emptyUpdate},
				{Task: *task2, LocalId: 2, LocalUpdate: emptyUpdate},
			},
			newTasks: []*task.Task{task1v2, task2},
			exp: []*task.LocalTask{
				{Task: *task1v2, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED},
				{Task: *task2, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED},
			},
		},
		{
			name: "remove deleted task",
			oldTasks: []*task.LocalTask{
				{Task: *task1, LocalUpdate: emptyUpdate},
				{Task: *task2, LocalUpdate: emptyUpdate},
			},
			newTasks: []*task.Task{task2},
			exp: []*task.LocalTask{
				{Task: *task2, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED},
			},
		},
		{
			name: "remove only outdated updates",
			oldTasks: []*task.LocalTask{
				{
					Task: *task1,
					LocalUpdate: &task.LocalUpdate{
						ForVersion: 1,
						Project:    "project-v2",
					},
				},
				{
					Task: *task2,
					LocalUpdate: &task.LocalUpdate{
						ForVersion: 2,
						Project:    "project-v3",
					},
				},
			},
			newTasks: []*task.Task{task1v2, task2},
			exp: []*task.LocalTask{
				{Task: *task1v2, LocalUpdate: emptyUpdate, LocalStatus: task.STATUS_FETCHED},
				{
					Task: *task2,
					LocalUpdate: &task.LocalUpdate{
						ForVersion: 2,
						Project:    "project-v3",
					},
					LocalStatus: task.STATUS_UPDATED,
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			sExp := task.ById(tc.exp)
			sAct := task.ById(storage.MergeNewTaskSet(tc.oldTasks, tc.newTasks))
			for i := range sAct {
				sAct[i].LocalId = 0
			}
			sort.Sort(sExp)
			sort.Sort(sAct)
			test.Equals(t, sExp, sAct)
		})
	}
}
