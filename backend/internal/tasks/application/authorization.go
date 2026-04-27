package application

import usersDomain "backend/internal/users/domain"

func canReadAllTasks(role usersDomain.UserRole) bool {
	return role == usersDomain.RoleProfessor || role == usersDomain.RoleAdmin
}

func isOperationalRole(role usersDomain.UserRole) bool {
	return role == usersDomain.RoleMonitor || role == usersDomain.RoleAssistant
}
