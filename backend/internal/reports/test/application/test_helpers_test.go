package application_test

import (
	"context"
	applicationpkg "backend/internal/reports/application"
	assignmentsDomain "backend/internal/assignments/domain"
	reportsDomain "backend/internal/reports/domain"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"io"
	"os"
	"path/filepath"
)

type mockReportRepository struct {
	reports  map[uint]*reportsDomain.Report
	nextID   uint
	saveErr  error
	findErr  error
	listErr  error
}

type mockWorkspaceReader struct {
	workspace *workspacesDomain.Workspace
	err       error
}

type mockWeekReader struct {
	week *weeksDomain.Week
	err  error
}

type mockAssignmentReader struct {
	assignments []assignmentsDomain.Assignment
	err         error
}

type mockTaskReader struct {
	tasks []tasksDomain.Task
	err   error
}

type mockUserReader struct {
	user *usersDomain.User
	err  error
}

type mockPDFGenerator struct {
	err error
}

type mockAIReportGenerator struct {
	text string
	err  error
}

type mockReportFileStorage struct {
	objectName  string
	contentType string
	uploaded    []byte
	err         error
}

func newMockReportRepository() *mockReportRepository {
	return &mockReportRepository{
		reports: make(map[uint]*reportsDomain.Report),
		nextID:  1,
	}
}

func (m *mockReportRepository) Save(report *reportsDomain.Report) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	report.ID = m.nextID
	m.nextID++
	copyReport := *report
	m.reports[report.ID] = &copyReport
	return nil
}

func (m *mockReportRepository) FindByID(id uint) (*reportsDomain.Report, error) {
	if m.findErr != nil {
		return nil, m.findErr
	}
	report, ok := m.reports[id]
	if !ok {
		return nil, reportsDomain.ErrReportNotFound
	}
	copyReport := *report
	return &copyReport, nil
}

func (m *mockReportRepository) FindAll(workspaceID uint, weekID *uint, userID *uint) ([]reportsDomain.Report, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}

	result := make([]reportsDomain.Report, 0)
	for _, report := range m.reports {
		if report.WorkspaceID != workspaceID {
			continue
		}
		if weekID != nil && report.WeekID != *weekID {
			continue
		}
		if userID != nil && report.UserID != *userID {
			continue
		}
		result = append(result, *report)
	}
	return result, nil
}

func (m *mockReportRepository) AutoMigrate() error { return nil }

func (m *mockWorkspaceReader) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.workspace, nil
}

func (m *mockWeekReader) FindByID(id uint) (*weeksDomain.Week, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.week, nil
}

func (m *mockAssignmentReader) FindAllByWorkspaceID(workspaceID uint) ([]assignmentsDomain.Assignment, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.assignments, nil
}

func (m *mockTaskReader) FindAllByWorkspaceAndWeek(workspaceID uint, weekID uint, weekInitialDate string) ([]tasksDomain.Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.tasks, nil
}

func (m *mockUserReader) FindByID(id uint) (*usersDomain.User, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.user, nil
}

func (m *mockPDFGenerator) Generate(filePath string, title string, lines []string) error {
	if m.err != nil {
		return m.err
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(title), 0o644)
}

func (m *mockAIReportGenerator) GenerateWeeklyReport(input applicationpkg.AIWeeklyReportInput) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	return m.text, nil
}

func (m *mockReportFileStorage) Upload(ctx context.Context, objectName string, reader io.Reader, contentType string) error {
	if m.err != nil {
		return m.err
	}
	payload, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	m.objectName = objectName
	m.contentType = contentType
	m.uploaded = payload
	return nil
}
