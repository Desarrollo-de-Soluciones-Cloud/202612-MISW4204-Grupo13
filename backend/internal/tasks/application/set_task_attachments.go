package application

import "backend/internal/tasks/domain"

type SetTaskAttachmentsInput struct {
	ID          uint
	Attachments []domain.TaskAttachment
}

type SetTaskAttachments struct {
	repository domain.TaskRepository
}

func NewSetTaskAttachments(repo domain.TaskRepository) *SetTaskAttachments {
	return &SetTaskAttachments{repository: repo}
}

func (uc *SetTaskAttachments) Execute(input SetTaskAttachmentsInput) (*TaskOutput, error) {
	task, err := uc.repository.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	task.Attachments = input.Attachments
	if err := uc.repository.UpdateAttachments(task.ID, task.Attachments); err != nil {
		return nil, err
	}

	task.Attachments = domain.NormalizeTaskAttachmentsForPersistence(task.Attachments)
	return newTaskOutput(task), nil
}
