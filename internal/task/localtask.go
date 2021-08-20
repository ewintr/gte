package task

type LocalTask struct {
	Task
	LocalId int
}

type ByDue []*LocalTask

func (lt ByDue) Len() int           { return len(lt) }
func (lt ByDue) Swap(i, j int)      { lt[i], lt[j] = lt[j], lt[i] }
func (lt ByDue) Less(i, j int) bool { return lt[j].Due.After(lt[i].Due) }
