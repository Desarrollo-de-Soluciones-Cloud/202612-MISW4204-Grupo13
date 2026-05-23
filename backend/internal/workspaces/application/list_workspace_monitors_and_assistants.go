package application

import (
	assignmentDomain "backend/internal/assignments/domain"
	usersDomain "backend/internal/users/domain"
	workspacesDomain "backend/internal/workspaces/domain"
)

type ListWorkspaceMonitorsAndAssistantsInput struct {
	ProfessorID uint
}

type MonitorAssistantInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	WeeklyHours int    `json:"weekly_hours"`
}

type ListWorkspaceMonitorsAndAssistantsOutput struct {
	Monitors   []MonitorAssistantInfo `json:"monitors"`
	Assistants []MonitorAssistantInfo `json:"assistants"`
}

type ListWorkspaceMonitorsAndAssistants struct {
	workspaceRepository workspacesDomain.WorkspaceRepository
	assignmentRepository assignmentDomain.AssignmentRepository
	userRepository      usersDomain.UserRepository
}

func NewListWorkspaceMonitorsAndAssistants(
	workspaceRepository workspacesDomain.WorkspaceRepository,
	assignmentRepository assignmentDomain.AssignmentRepository,
	userRepository usersDomain.UserRepository,
) *ListWorkspaceMonitorsAndAssistants {
	return &ListWorkspaceMonitorsAndAssistants{
		workspaceRepository:  workspaceRepository,
		assignmentRepository: assignmentRepository,
		userRepository:       userRepository,
	}
}

func (uc *ListWorkspaceMonitorsAndAssistants) Execute(input ListWorkspaceMonitorsAndAssistantsInput) (*ListWorkspaceMonitorsAndAssistantsOutput, error) {
	workspaces, err := uc.workspaceRepository.FindByUserID(input.ProfessorID)
	if err != nil {
		return nil, err
	}

	if len(workspaces) == 0 {
		return &ListWorkspaceMonitorsAndAssistantsOutput{
			Monitors:   []MonitorAssistantInfo{},
			Assistants: []MonitorAssistantInfo{},
		}, nil
	}

	assignments, err := uc.assignmentRepository.FindByWorkspaceIDsAndRoles(
		extractWorkspaceIDs(workspaces),
		[]assignmentDomain.AssignmentRole{assignmentDomain.RoleMonitor, assignmentDomain.RoleAssistant},
	)
	if err != nil {
		return nil, err
	}

	return &ListWorkspaceMonitorsAndAssistantsOutput{
		Monitors:   uc.buildMonitorAssistantInfo(assignments, assignmentDomain.RoleMonitor),
		Assistants: uc.buildMonitorAssistantInfo(assignments, assignmentDomain.RoleAssistant),
	}, nil
}

func extractWorkspaceIDs(workspaces []workspacesDomain.Workspace) []uint {
	workspaceIDs := make([]uint, len(workspaces))
	for i, ws := range workspaces {
		workspaceIDs[i] = ws.ID
	}
	return workspaceIDs
}

func (uc *ListWorkspaceMonitorsAndAssistants) buildMonitorAssistantInfo(assignments []assignmentDomain.Assignment, role assignmentDomain.AssignmentRole) []MonitorAssistantInfo {
	result := make([]MonitorAssistantInfo, 0)
	seenUserIDs := make(map[uint]bool)

	for _, assignment := range assignments {
		if assignment.Role != role || seenUserIDs[assignment.UserID] {
			continue
		}

		seenUserIDs[assignment.UserID] = true
		user, err := uc.userRepository.FindByID(assignment.UserID)
		if err != nil {
			continue
		}

		result = append(result, MonitorAssistantInfo{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Role:        string(role),
			WeeklyHours: assignment.WeeklyHours,
		})
	}

	return result
}
