package domain

import (
	domainpkg "backend/internal/assignments/domain"
	"errors"
	"testing"
)

func TestBuildWorkloadWithAssignmentAddsAssistantHours(t *testing.T) {
	current := domainpkg.UserAssignmentWorkload{AssistantWeeklyHours: 10, MonitorWeeklyHours: 4, MonitorAssignments: 1}

	next, err := domainpkg.BuildWorkloadWithAssignment(current, domainpkg.RoleAssistant, 6)
	if err != nil {
		t.Fatalf("expected no error building assistant workload, got %v", err)
	}

	if next.AssistantWeeklyHours != 16 {
		t.Fatalf("expected assistant hours 16, got %d", next.AssistantWeeklyHours)
	}
	if next.MonitorWeeklyHours != 4 {
		t.Fatalf("expected monitor hours 4, got %d", next.MonitorWeeklyHours)
	}
	if next.MonitorAssignments != 1 {
		t.Fatalf("expected monitor assignments 1, got %d", next.MonitorAssignments)
	}
}

func TestBuildWorkloadWithAssignmentAddsMonitorHoursAndCount(t *testing.T) {
	current := domainpkg.UserAssignmentWorkload{AssistantWeeklyHours: 10, MonitorWeeklyHours: 4, MonitorAssignments: 1}

	next, err := domainpkg.BuildWorkloadWithAssignment(current, domainpkg.RoleMonitor, 3)
	if err != nil {
		t.Fatalf("expected no error building monitor workload, got %v", err)
	}

	if next.AssistantWeeklyHours != 10 {
		t.Fatalf("expected assistant hours 10, got %d", next.AssistantWeeklyHours)
	}
	if next.MonitorWeeklyHours != 7 {
		t.Fatalf("expected monitor hours 7, got %d", next.MonitorWeeklyHours)
	}
	if next.MonitorAssignments != 2 {
		t.Fatalf("expected monitor assignments 2, got %d", next.MonitorAssignments)
	}
}

func TestValidateRF05WorkloadAssistantHoursLimit(t *testing.T) {
	err := domainpkg.ValidateRF05Workload(domainpkg.UserAssignmentWorkload{
		AssistantWeeklyHours: 23,
		MonitorWeeklyHours:   0,
		MonitorAssignments:   0,
	})
	if !errors.Is(err, domainpkg.ErrAssignmentAssistantHoursLimitExceeded) {
		t.Fatalf("expected ErrAssignmentAssistantHoursLimitExceeded, got %v", err)
	}
}

func TestValidateRF05WorkloadMonitorCountLimit(t *testing.T) {
	err := domainpkg.ValidateRF05Workload(domainpkg.UserAssignmentWorkload{
		AssistantWeeklyHours: 0,
		MonitorWeeklyHours:   10,
		MonitorAssignments:   4,
	})
	if !errors.Is(err, domainpkg.ErrAssignmentMonitorCountLimitExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorCountLimitExceeded, got %v", err)
	}
}

func TestValidateRF05WorkloadMonitorHoursLimit(t *testing.T) {
	err := domainpkg.ValidateRF05Workload(domainpkg.UserAssignmentWorkload{
		AssistantWeeklyHours: 0,
		MonitorWeeklyHours:   13,
		MonitorAssignments:   2,
	})
	if !errors.Is(err, domainpkg.ErrAssignmentMonitorHoursLimitExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorHoursLimitExceeded, got %v", err)
	}
}

func TestValidateRF05WorkloadFortyPercentLimit(t *testing.T) {
	err := domainpkg.ValidateRF05Workload(domainpkg.UserAssignmentWorkload{
		AssistantWeeklyHours: 10,
		MonitorWeeklyHours:   5,
		MonitorAssignments:   1,
	})
	if !errors.Is(err, domainpkg.ErrAssignmentMonitorFortyPercentExceeded) {
		t.Fatalf("expected ErrAssignmentMonitorFortyPercentExceeded, got %v", err)
	}
}

func TestCalculateMaxMonitorHoursFromAssistantRoundsUp(t *testing.T) {
	maxMonitorHours := domainpkg.CalculateMaxMonitorHoursFromAssistant(11)
	if maxMonitorHours != 5 {
		t.Fatalf("expected rounded-up monitor limit 5, got %d", maxMonitorHours)
	}
}

func TestValidateRF05WorkloadValidCase(t *testing.T) {
	err := domainpkg.ValidateRF05Workload(domainpkg.UserAssignmentWorkload{
		AssistantWeeklyHours: 20,
		MonitorWeeklyHours:   8,
		MonitorAssignments:   3,
	})
	if err != nil {
		t.Fatalf("expected no error for valid RF05 workload, got %v", err)
	}
}
