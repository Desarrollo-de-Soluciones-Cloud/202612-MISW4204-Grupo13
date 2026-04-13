package domain

type WorkspaceType string

const (
	CourseType  WorkspaceType = "course"
	ProjectType WorkspaceType = "project"
)

func IsValidWorkspaceType(state WorkspaceType) bool {
	switch state {
	case CourseType, ProjectType:
		return true
	default:
		return false
	}
}

func ValidWorkspaceTypesString() string {
	return string(CourseType) + ", " + string(ProjectType)
}

type WorkspaceState string

const (
	ActiveState WorkspaceState = "active"
	ClosedState WorkspaceState = "closed"
)

func IsValidWorkspaceState(state WorkspaceState) bool {
	switch state {
	case ActiveState, ClosedState:
		return true
	default:
		return false
	}
}

func ValidWorkspaceStatesString() string {
	return string(ActiveState) + ", " + string(ClosedState)
}