package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go-mod.ewintr.nl/gte/internal/task"
	"github.com/google/uuid"
	_ "modernc.org/sqlite"
)

type sqliteMigration string

var sqliteMigrations = []sqliteMigration{
	`CREATE TABLE task ("id" TEXT UNIQUE, "version" INTEGER, "folder" TEXT, "action" TEXT, "project" TEXT, "due" TEXT, "recur" TEXT)`,
	`CREATE TABLE system ("latest_sync" INTEGER)`,
	`INSERT INTO system (latest_sync) VALUES (0)`,
	`CREATE TABLE local_id ("id" TEXT UNIQUE, "local_id" INTEGER UNIQUE)`,
	`ALTER TABLE local_id RENAME TO local_task`,
	`ALTER TABLE local_task ADD COLUMN local_update TEXT`,
	`ALTER TABLE task ADD COLUMN local_id INTEGER`,
	`ALTER TABLE task ADD COLUMN local_update TEXT`,
	`UPDATE task SET local_id = (SELECT local_id FROM local_task WHERE local_task.id=task.id)`,
	`UPDATE task SET local_update = (SELECT local_update FROM local_task WHERE local_task.id=task.id)`,
	`DROP TABLE local_task`,
	`ALTER TABLE task ADD COLUMN local_status TEXT`,
	`UPDATE task SET local_status = "fetched"`,
	`DROP TABLE system`,
	`CREATE TABLE system ("latest_fetch" INTEGER, "latest_dispatch" INTEGER)`,
	`INSERT INTO system (latest_fetch, latest_dispatch) VALUES (0, 0)`,
}

var (
	ErrInvalidConfiguration     = errors.New("invalid configuration")
	ErrIncompatibleSQLMigration = errors.New("incompatible migration")
	ErrNotEnoughSQLMigrations   = errors.New("already more migrations than wanted")
	ErrSqliteFailure            = errors.New("sqlite returned an error")
)

type SqliteConfig struct {
	DBPath string
}

// Sqlite is an sqlite implementation of LocalRepository
type Sqlite struct {
	db *sql.DB
}

func NewSqlite(conf *SqliteConfig) (*Sqlite, error) {
	db, err := sql.Open("sqlite", conf.DBPath)
	if err != nil {
		return &Sqlite{}, fmt.Errorf("%w: %v", ErrInvalidConfiguration, err)
	}

	s := &Sqlite{
		db: db,
	}

	if err := s.migrate(sqliteMigrations); err != nil {
		return &Sqlite{}, err
	}

	return s, nil
}

func (s *Sqlite) LatestSyncs() (time.Time, time.Time, error) {
	rows, err := s.db.Query(`SELECT strftime('%s', latest_fetch), strftime('%s', latest_dispatch) FROM system`)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	defer rows.Close()

	rows.Next()
	var latest_fetch, latest_dispatch int64
	if err := rows.Scan(&latest_fetch, &latest_dispatch); err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return time.Unix(latest_fetch, 0), time.Unix(latest_dispatch, 0), nil
}

func (s *Sqlite) SetTasks(tasks []*task.Task) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	defer tx.Rollback()

	oldTasks, err := s.FindAll()
	if err != nil {
		return err
	}
	newTasks := MergeNewTaskSet(oldTasks, tasks)

	if _, err := s.db.Exec(`DELETE FROM task`); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	for _, t := range newTasks {
		var recurStr string
		if t.Recur != nil {
			recurStr = t.Recur.String()
		}

		_, err := s.db.Exec(`
INSERT INTO task
(id, local_id, version, folder, action, project, due, recur, local_update, local_status)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			t.Id, t.LocalId, t.Version, t.Folder, t.Action, t.Project, t.Due.String(), recurStr, t.LocalUpdate, t.LocalStatus)

		if err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
	}

	if _, err := s.db.Exec(`UPDATE system SET latest_fetch=DATETIME('now')`); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return nil
}

func (s *Sqlite) FindAll() ([]*task.LocalTask, error) {
	rows, err := s.db.Query(`
SELECT id, local_id, version, folder, action, project, due, recur, local_update, local_status
FROM task`)
	if err != nil {
		return []*task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	tasks := []*task.LocalTask{}
	defer rows.Close()
	for rows.Next() {
		var id, folder, action, project, due, recur, localStatus string
		var localId, version int
		var localUpdate task.LocalUpdate
		if err := rows.Scan(&id, &localId, &version, &folder, &action, &project, &due, &recur, &localUpdate, &localStatus); err != nil {
			return []*task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		tasks = append(tasks, &task.LocalTask{
			Task: task.Task{
				Id:      id,
				Version: version,
				Folder:  folder,
				Action:  action,
				Project: project,
				Due:     task.NewDateFromString(due),
				Recur:   task.NewRecurrer(recur),
			},
			LocalId:     localId,
			LocalUpdate: &localUpdate,
			LocalStatus: localStatus,
		})
	}

	return tasks, nil
}

func (s *Sqlite) FindById(id string) (*task.LocalTask, error) {
	var folder, action, project, due, recur, localStatus string
	var localId, version int
	var localUpdate task.LocalUpdate
	row := s.db.QueryRow(`
SELECT local_id, version, folder, action, project, due, recur, local_update, local_status
FROM task
WHERE task.id = ?
LIMIT 1`, id)
	if err := row.Scan(&localId, &version, &folder, &action, &project, &due, &recur, &localUpdate, &localStatus); err != nil {
		return &task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return &task.LocalTask{
		Task: task.Task{
			Id:      id,
			Version: version,
			Folder:  folder,
			Action:  action,
			Project: project,
			Due:     task.NewDateFromString(due),
			Recur:   task.NewRecurrer(recur),
		},
		LocalId:     localId,
		LocalUpdate: &localUpdate,
		LocalStatus: localStatus,
	}, nil
}

func (s *Sqlite) FindByLocalId(localId int) (*task.LocalTask, error) {
	var id string
	row := s.db.QueryRow(`SELECT id FROM task WHERE local_id = ?`, localId)
	if err := row.Scan(&id); err != nil {
		return &task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	t, err := s.FindById(id)
	if err != nil {
		return &task.LocalTask{}, nil
	}

	return t, nil
}

func (s *Sqlite) SetLocalUpdate(id string, update *task.LocalUpdate) error {
	if _, err := s.db.Exec(`
UPDATE task
SET local_update = ?, local_status = ?
WHERE id = ?`, update, task.STATUS_UPDATED, id); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return nil
}

func (s *Sqlite) MarkDispatched(localId int) error {
	if _, err := s.db.Exec(`
UPDATE task
SET local_status = ?
WHERE local_id = ?`, task.STATUS_DISPATCHED, localId); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	if _, err := s.db.Exec(`
UPDATE system
SET latest_dispatch=DATETIME('now')`); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return nil
}

func (s *Sqlite) Add(update *task.LocalUpdate) (*task.LocalTask, error) {
	rows, err := s.db.Query(`SELECT local_id FROM task`)
	if err != nil {
		return &task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	var used []int
	for rows.Next() {
		var localId int
		if err := rows.Scan(&localId); err != nil {
			return &task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		used = append(used, localId)
	}
	rows.Close()

	tsk := &task.LocalTask{
		Task: task.Task{
			Id:      uuid.New().String(),
			Version: 0,
			Folder:  task.FOLDER_NEW,
		},
		LocalId:     NextLocalId(used),
		LocalStatus: task.STATUS_UPDATED,
		LocalUpdate: update,
	}

	var recurStr string
	if tsk.LocalUpdate.Recur != nil {
		recurStr = tsk.LocalUpdate.Recur.String()
	}
	if _, err := s.db.Exec(`
INSERT INTO task
(id, local_id, version, folder, action, project, due, recur, local_status, local_update)
VALUES
(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		tsk.Id, tsk.LocalId, tsk.Version, tsk.Folder, tsk.Action, tsk.Project,
		tsk.Due.String(), recurStr, tsk.LocalStatus, tsk.LocalUpdate); err != nil {
		return &task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return tsk, nil
}

func (s *Sqlite) migrate(wanted []sqliteMigration) error {
	// admin table
	if _, err := s.db.Exec(`
CREATE TABLE IF NOT EXISTS migration
("id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, "query" TEXT)
`); err != nil {
		return err
	}

	// find existing
	rows, err := s.db.Query(`SELECT query FROM migration ORDER BY id`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	existing := []sqliteMigration{}
	for rows.Next() {
		var query string
		if err := rows.Scan(&query); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		existing = append(existing, sqliteMigration(query))
	}
	rows.Close()

	// compare
	missing, err := compareMigrations(wanted, existing)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	// execute missing
	for _, query := range missing {
		if _, err := s.db.Exec(string(query)); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}

		// register
		if _, err := s.db.Exec(`
INSERT INTO migration
(query) VALUES (?)
`, query); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
	}

	return nil
}

func compareMigrations(wanted, existing []sqliteMigration) ([]sqliteMigration, error) {
	needed := []sqliteMigration{}
	if len(wanted) < len(existing) {
		return []sqliteMigration{}, ErrNotEnoughSQLMigrations
	}

	for i, want := range wanted {
		switch {
		case i >= len(existing):
			needed = append(needed, want)
		case want == existing[i]:
			// do nothing
		case want != existing[i]:
			return []sqliteMigration{}, fmt.Errorf("%w: %v", ErrIncompatibleSQLMigration, want)
		}
	}

	return needed, nil
}
