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
	// Get all workspaces for the professor
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

	// Extract workspace IDs
	workspaceIDs := make([]uint, len(workspaces))
	for i, ws := range workspaces {
		workspaceIDs[i] = ws.ID
	}

	// Get assignments for those workspaces with monitor and assistant roles
	assignments, err := uc.assignmentRepository.FindByWorkspaceIDsAndRoles(
		workspaceIDs,
		[]assignmentDomain.AssignmentRole{assignmentDomain.RoleMonitor, assignmentDomain.RoleAssistant},
	)
	if err != nil {
		return nil, err
	}

	// Process monitors
	monitorsInfo := []MonitorAssistantInfo{}
	seenMonitorUserIDs := make(map[uint]bool)

	for _, assignment := range assignments {
		if assignment.Role != assignmentDomain.RoleMonitor {
			continue
		}

		if seenMonitorUserIDs[assignment.UserID] {
			continue
		}
		seenMonitorUserIDs[assignment.UserID] = true

		user, err := uc.userRepository.FindByID(assignment.UserID)
		if err != nil {
			continue
		}

		monitorsInfo = append(monitorsInfo, MonitorAssistantInfo{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Role:        string(assignmentDomain.RoleMonitor),
			WeeklyHours: assignment.WeeklyHours,
		})
	}

	// Process assistants
	assistantsInfo := []MonitorAssistantInfo{}
	seenAssistantUserIDs := make(map[uint]bool)

	for _, assignment := range assignments {
		if assignment.Role != assignmentDomain.RoleAssistant {
			continue
		}

		if seenAssistantUserIDs[assignment.UserID] {
			continue
		}
		seenAssistantUserIDs[assignment.UserID] = true

		user, err := uc.userRepository.FindByID(assignment.UserID)
		if err != nil {
			continue
		}

		assistantsInfo = append(assistantsInfo, MonitorAssistantInfo{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Role:        string(assignmentDomain.RoleAssistant),
			WeeklyHours: assignment.WeeklyHours,
		})
	}

	return &ListWorkspaceMonitorsAndAssistantsOutput{
		Monitors:   monitorsInfo,
		Assistants: assistantsInfo,
	}, nil
}
