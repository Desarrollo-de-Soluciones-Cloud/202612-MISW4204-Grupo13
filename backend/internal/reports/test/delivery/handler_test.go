package delivery_test

import (
	assignmentsDomain "backend/internal/assignments/domain"
	reportsApplication "backend/internal/reports/application"
	reportsDomain "backend/internal/reports/domain"
	authDomain "backend/internal/auth/domain"
	deliverypkg "backend/internal/reports/delivery"
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	weeksDomain "backend/internal/weeks/domain"
	workspacesDomain "backend/internal/workspaces/domain"
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

const expectedBadRequestFormat = "expected 400, got %d"

func newReportHandlerForTest() *deliverypkg.ReportHandler {
	return deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{})
}

type reportWorkspaceReaderStub struct {
	workspace *workspacesDomain.Workspace
	err       error
}

type reportDownloadStorageStub struct {
	payload []byte
	err     error
}

type reportRepositoryStub struct {
	reports map[uint]*reportsDomain.Report
}

type weekReaderStub struct{}
type userReaderStub struct{}
type reportAssignmentReaderStub struct{}
type reportTaskReaderStub struct {
	tasks []tasksDomain.Task
}
type reportPDFGeneratorStub struct{}
type reportAIStub struct{}
type jobPublisherStub struct{}

func (r *reportWorkspaceReaderStub) FindByID(id uint) (*workspacesDomain.Workspace, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.workspace == nil || r.workspace.ID != id {
		return nil, workspacesDomain.ErrWorkspaceNotFound
	}
	return r.workspace, nil
}

func (s *reportDownloadStorageStub) Download(ctx context.Context, objectName string) (io.ReadCloser, error) {
	if s.err != nil {
		return nil, s.err
	}
	return io.NopCloser(strings.NewReader(string(s.payload))), nil
}

func (r *reportRepositoryStub) Save(report *reportsDomain.Report) error { return nil }
func (r *reportRepositoryStub) FindByID(id uint) (*reportsDomain.Report, error) {
	if report, ok := r.reports[id]; ok {
		return report, nil
	}
	return nil, reportsDomain.ErrReportNotFound
}
func (r *reportRepositoryStub) FindAll(workspaceID uint, weekID *uint, userID *uint) ([]reportsDomain.Report, error) {
	result := make([]reportsDomain.Report, 0)
	for _, report := range r.reports {
		if report.WorkspaceID == workspaceID {
			result = append(result, *report)
		}
	}
	return result, nil
}
func (r *reportRepositoryStub) AutoMigrate() error { return nil }

func (w *weekReaderStub) FindByID(id uint) (*weeksDomain.Week, error) {
	return &weeksDomain.Week{ID: id, Number: 1, InitialDate: "2026-04-07", FinalDate: "2026-04-13"}, nil
}

func (u *userReaderStub) FindByID(id uint) (*usersDomain.User, error) {
	return &usersDomain.User{ID: id, Name: "Ana", Email: "ana@example.com", GlobalRole: usersDomain.RoleAssistant}, nil
}

func (a *reportAssignmentReaderStub) FindAllByWorkspaceID(workspaceID uint) ([]assignmentsDomain.Assignment, error) {
	return []assignmentsDomain.Assignment{{ID: 10, UserID: 2, WorkspaceID: workspaceID, Role: assignmentsDomain.RoleAssistant, WeeklyHours: 4}}, nil
}

func (t *reportTaskReaderStub) FindAllByWorkspaceAndWeek(workspaceID uint, weekID uint, weekInitialDate string) ([]tasksDomain.Task, error) {
	return t.tasks, nil
}

func (p *reportPDFGeneratorStub) Generate(filePath string, title string, lines []string) error {
	return nil
}

func (a *reportAIStub) GenerateWeeklyReport(input reportsApplication.AIWeeklyReportInput) (string, error) {
	return "summary", nil
}
func (j *jobPublisherStub) PublishWeeklyReportJob(ctx context.Context, job reportsApplication.WeeklyReportJobMessage) error {
	return nil
}

func TestGenerateWeeklyReportsUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/reports/weekly", nil)

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGenerateWeeklyReportsBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly", bytes.NewBufferString(`{"workspace_id":"bad"}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(expectedBadRequestFormat, w.Code)
	}
}

func TestListReportsBadWorkspaceFilter(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reports", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.ListReports(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(expectedBadRequestFormat, w.Code)
	}
}

func TestDownloadReportBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = gin.Params{{Key: "id", Value: "bad"}}
	c.Set("current_user", authDomain.AuthenticatedUser{
		ID:         1,
		GlobalRole: usersDomain.RoleAdmin,
	})

	handler.DownloadReport(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf(expectedBadRequestFormat, w.Code)
	}
}

func TestProcessWeeklyReportJobUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := newReportHandlerForTest()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/reports/weekly/process", nil)

	handler.ProcessWeeklyReportJob(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestProcessWeeklyReportJobBadPayload(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		ProcessWeeklyReport: reportsApplication.NewProcessWeeklyReportJob(
			reportsApplication.ProcessWeeklyReportJobDependencies{
				ReportRepo:      &reportRepositoryStub{reports: map[uint]*reportsDomain.Report{}},
				WorkspaceReader: &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, Name: "WS"}},
				WeekReader:      &weekReaderStub{},
				AssignmentReader: &reportAssignmentReaderStub{},
				TaskReader:        &reportTaskReaderStub{},
				UserReader:        &userReaderStub{},
				PDFGenerator:      &reportPDFGeneratorStub{},
				AIReportGenerator: &reportAIStub{},
			},
			nil,
		),
		PubSubPushAuthToken: "secret",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly/process?token=secret", bytes.NewBufferString(`{"message":{"data":"%%%bad%%%"}}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ProcessWeeklyReportJob(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListReportsSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	workspaceReader := &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 10, Name: "WS"}}
	reportRepo := &reportRepositoryStub{reports: map[uint]*reportsDomain.Report{
		1: {ID: 1, WorkspaceID: 1, WeekID: 2, AssignmentID: 10, UserID: 3, FilePath: "reports/file.pdf"},
	}}
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		ListReports:     reportsApplication.NewListReports(reportRepo, workspaceReader, &weekReaderStub{}, &userReaderStub{}),
		WorkspaceReader: workspaceReader,
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/reports?workspace_id=1", nil)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor})

	handler.ListReports(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestDownloadReportSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	workspaceReader := &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 10, Name: "WS"}}
	reportRepo := &reportRepositoryStub{reports: map[uint]*reportsDomain.Report{
		55: {ID: 55, WorkspaceID: 1, WeekID: 2, AssignmentID: 10, UserID: 3, FilePath: "reports/test.pdf"},
	}}
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		GetReportByID:     reportsApplication.NewGetReportByID(reportRepo),
		WorkspaceReader:   workspaceReader,
		ReportFileStorage: &reportDownloadStorageStub{payload: []byte("pdf")},
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/reports/55/download", nil)
	c.Params = gin.Params{{Key: "id", Value: "55"}}
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor})

	handler.DownloadReport(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGenerateWeeklyReportsWorkspaceNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		WorkspaceReader: &reportWorkspaceReaderStub{},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly", bytes.NewBufferString(`{"workspace_id":1,"week_id":2}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 1, GlobalRole: usersDomain.RoleAdmin})

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGenerateWeeklyReportsForbiddenWorkspace(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		WorkspaceReader: &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 99}},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly", bytes.NewBufferString(`{"workspace_id":1,"week_id":2}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor})

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestGenerateWeeklyReportsQueued(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		QueueWeeklyReports: reportsApplication.NewQueueWeeklyReports(
			&reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 10, Name: "WS"}},
			&weekReaderStub{},
			&reportAssignmentReaderStub{},
			&reportTaskReaderStub{tasks: []tasksDomain.Task{{ID: 1, AssignmentID: 10}}},
			&jobPublisherStub{},
		),
		WorkspaceReader: &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 10, Name: "WS"}},
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly", bytes.NewBufferString(`{"workspace_id":1,"week_id":2}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor})

	handler.GenerateWeeklyReports(c)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
}

func TestProcessWeeklyReportJobNotImplemented(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		PubSubPushAuthToken: "secret",
	})
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/reports/weekly/process?token=secret", bytes.NewBufferString(`{"message":{"data":"e30="}}`))
	req.Header.Set("Content-Type", "application/json")
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.ProcessWeeklyReportJob(c)

	if w.Code != http.StatusNotImplemented {
		t.Fatalf("expected 501, got %d", w.Code)
	}
}

func TestDownloadReportForbiddenForProfessor(t *testing.T) {
	gin.SetMode(gin.TestMode)
	workspaceReader := &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 99, Name: "WS"}}
	reportRepo := &reportRepositoryStub{reports: map[uint]*reportsDomain.Report{
		55: {ID: 55, WorkspaceID: 1, WeekID: 2, AssignmentID: 10, UserID: 3, FilePath: "reports/test.pdf"},
	}}
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		GetReportByID:   reportsApplication.NewGetReportByID(reportRepo),
		WorkspaceReader: workspaceReader,
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/reports/55/download", nil)
	c.Params = gin.Params{{Key: "id", Value: "55"}}
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor})

	handler.DownloadReport(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestDownloadReportWithoutStorage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	workspaceReader := &reportWorkspaceReaderStub{workspace: &workspacesDomain.Workspace{ID: 1, UserID: 10, Name: "WS"}}
	reportRepo := &reportRepositoryStub{reports: map[uint]*reportsDomain.Report{
		55: {ID: 55, WorkspaceID: 1, WeekID: 2, AssignmentID: 10, UserID: 3, FilePath: "reports/test.pdf"},
	}}
	handler := deliverypkg.NewReportHandler(deliverypkg.ReportHandlerDependencies{
		GetReportByID:   reportsApplication.NewGetReportByID(reportRepo),
		WorkspaceReader: workspaceReader,
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/reports/55/download", nil)
	c.Params = gin.Params{{Key: "id", Value: "55"}}
	c.Set("current_user", authDomain.AuthenticatedUser{ID: 10, GlobalRole: usersDomain.RoleProfessor})

	handler.DownloadReport(c)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
