package screen

type SaveConfigRequest struct {
	Fields map[string]string
}

type SaveNewTaskRequest struct {
	Fields map[string]string
}

type SyncTasksRequest struct{}

type MarkTaskDoneRequest struct {
	ID string
}
