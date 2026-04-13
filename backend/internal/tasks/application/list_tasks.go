package application

import "backend/internal/tasks/domain"

type ListTasksOutput struct {
	Tasks []TaskOutput `json:"tasks"`
}

type ListTasks struct {
	repository domain.TaskRepository
}

func NewListTasks(repo domain.TaskRepository) *ListTasks {
	return &ListTasks{repository: repo}
}

func (uc *ListTasks) Execute() (*ListTasksOutput, error) {
	tasks, err := uc.repository.FindAll()
	if err != nil {
		return nil, err
	}

	output := make([]TaskOutput, len(tasks))
	for i := range tasks {
		output[i] = *newTaskOutput(&tasks[i])
	}

	return &ListTasksOutput{Tasks: output}, nil
}
