package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"git.ewintr.nl/gte/internal/task"
	_ "modernc.org/sqlite"
)

type sqliteMigration string

var sqliteMigrations = []sqliteMigration{
	`CREATE TABLE task ("id" TEXT UNIQUE, "version" INTEGER, "folder" TEXT, "action" TEXT, "project" TEXT, "due" TEXT, "recur" TEXT)`,
	`CREATE TABLE system ("latest_sync" INTEGER)`,
	`INSERT INTO system (latest_sync) VALUES (0)`,
	`CREATE TABLE local_id ("id" TEXT UNIQUE, "local_id" INTEGER UNIQUE)`,
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

func (s *Sqlite) LatestSync() (time.Time, error) {
	rows, err := s.db.Query(`SELECT latest_sync FROM system`)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	defer rows.Close()

	rows.Next()
	var latest int64
	if err := rows.Scan(&latest); err != nil {
		return time.Time{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return time.Unix(latest, 0), nil
}

func (s *Sqlite) SetTasks(tasks []*task.Task) error {
	// set tasks
	if _, err := s.db.Exec(`DELETE FROM task`); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	localIdMap := map[string]int{}
	for _, t := range tasks {
		var recurStr string
		if t.Recur != nil {
			recurStr = t.Recur.String()
		}
		_, err := s.db.Exec(`
INSERT INTO task
(id, version, folder, action, project, due, recur)
VALUES
(?, ?, ?, ?, ?, ?, ?)`,
			t.Id, t.Version, t.Folder, t.Action, t.Project, t.Due.String(), recurStr)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}

		localIdMap[t.Id] = 0
	}

	// set local_ids
	rows, err := s.db.Query(`SELECT id, local_id FROM local_id`)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var localId int
		if err := rows.Scan(&id, &localId); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		if _, ok := localIdMap[id]; ok {
			localIdMap[id] = localId
		}
	}

	if _, err := s.db.Exec(`DELETE FROM local_id`); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	var used []int
	for _, localId := range localIdMap {
		if localId != 0 {
			used = append(used, localId)
		}
	}

	for id, localId := range localIdMap {
		if localId == 0 {
			newLocalId := NextLocalId(used)
			localIdMap[id] = newLocalId
			used = append(used, newLocalId)
		}
	}

	for id, localId := range localIdMap {
		if _, err := s.db.Exec(`
INSERT INTO local_id
(id, local_id)
VALUES
(?, ?)`, id, localId); err != nil {
			return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
	}

	// update system
	if _, err := s.db.Exec(`
UPDATE system
SET latest_sync = ?`,
		time.Now().Unix()); err != nil {
		return fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return nil
}

func (s *Sqlite) FindAllInFolder(folder string) ([]*task.LocalTask, error) {
	rows, err := s.db.Query(`
SELECT id, version, folder, action, project, due, recur
FROM task
WHERE folder = ?`, folder)
	if err != nil {
		return []*task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return tasksFromRows(rows)
}

func (s *Sqlite) FindAllInProject(project string) ([]*task.LocalTask, error) {
	rows, err := s.db.Query(`
SELECT id, version, folder, action, project, due, recur
FROM task
WHERE project = ?`, project)
	if err != nil {
		return []*task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	return tasksFromRows(rows)
}

func (s *Sqlite) FindById(id string) (*task.LocalTask, error) {
	var folder, action, project, due, recur string
	var version int
	row := s.db.QueryRow(`
SELECT version, folder, action, project, due, recur
FROM task
WHERE id = ?
LIMIT 1`, id)
	if err := row.Scan(&version, &folder, &action, &project, &due, &recur); err != nil {
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
		}}, nil
}

func (s *Sqlite) FindByLocalId(localId int) (*task.LocalTask, error) {
	var id string
	row := s.db.QueryRow(`SELECT id FROM local_id WHERE local_id = ?`, localId)
	if err := row.Scan(&id); err != nil {
		return &task.LocalTask{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	t, err := s.FindById(id)
	if err != nil {
		return &task.LocalTask{}, nil
	}

	return t, nil
}

func (s *Sqlite) LocalIds() (map[string]int, error) {
	rows, err := s.db.Query(`SELECT id, local_id FROM local_id`)
	if err != nil {
		return map[string]int{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
	}

	idMap := map[string]int{}
	defer rows.Close()
	for rows.Next() {
		var id string
		var local_id int
		if err := rows.Scan(&id, &local_id); err != nil {
			return map[string]int{}, fmt.Errorf("%w: %v", ErrSqliteFailure, err)
		}
		idMap[id] = local_id
	}

	return idMap, nil
}

func tasksFromRows(rows *sql.Rows) ([]*task.LocalTask, error) {
	tasks := []*task.LocalTask{}

	defer rows.Close()
	for rows.Next() {
		var id, folder, action, project, due, recur string
		var version int
		if err := rows.Scan(&id, &version, &folder, &action, &project, &due, &recur); err != nil {
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
			}})
	}

	return tasks, nil
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
