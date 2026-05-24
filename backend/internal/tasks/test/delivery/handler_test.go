package delivery_test

import (
	tasksDomain "backend/internal/tasks/domain"
	usersDomain "backend/internal/users/domain"
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	testTasksPath          = "/tasks"
	testTaskByIDPath       = "/tasks/1"
	testTaskDownloadPath   = "/tasks/1/attachments/att_1/download"
	testHeaderContentType  = "Content-Type"
	testApplicationJSON    = "application/json"
	testTaskTitle          = "Task title"
	testTaskDescription    = "Task description"
	testTaskWeekStartDate  = "2026-04-06"
	testAttachmentFileName = "evidence.pdf"
	testAttachmentPDF      = "application/pdf"
	testAttachmentFilePath = "attachments/task_1/file_1_evidence.pdf"
)

func TestCreateTaskUnauthorized(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _ := newTaskHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, testTasksPath, nil)

	handler.CreateTask(c)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCreateTaskBadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _ := newTaskHandlerForTest(t)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, testTasksPath, bytes.NewBufferString(`{"assignment_id":"bad"}`))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authenticatedUser(1, usersDomain.RoleAdmin))

	handler.CreateTask(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateTaskSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, _ := newTaskHandlerForTest(t)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]any{
		"assignment_id":   1,
		"title":           testTaskTitle,
		"description":     testTaskDescription,
		"status":          "abierto",
		"spent_hours":     2,
		"observations":    "",
		"week_start_date": testTaskWeekStartDate,
	})
	req := httptest.NewRequest(http.MethodPost, testTasksPath, bytes.NewBuffer(body))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authenticatedUser(99, usersDomain.RoleAdmin))

	handler.CreateTask(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if len(repo.tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(repo.tasks))
	}
}

func TestListTasksFiltersByAccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, _ := newTaskHandlerForTest(t)
	seedTask(t, repo, 10, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), nil)
	seedTask(t, repo, 99, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testTasksPath, nil)
	c.Set("current_user", authenticatedUser(10, usersDomain.RoleMonitor))

	handler.ListTasks(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var response struct {
		Tasks []any `json:"tasks"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if len(response.Tasks) != 1 {
		t.Fatalf("expected 1 visible task, got %d", len(response.Tasks))
	}
}

func TestGetTaskByIDForbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, _ := newTaskHandlerForTest(t)
	task := seedTask(t, repo, 10, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testTaskByIDPath, nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(200, usersDomain.RoleAssistant))

	handler.GetTaskByID(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d for task %d", w.Code, task.ID)
	}
}

func TestUpdateTaskBadID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _ := newTaskHandlerForTest(t)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPut, "/tasks/bad", nil)
	c.Params = gin.Params{{Key: "id", Value: "bad"}}
	c.Set("current_user", authenticatedUser(1, usersDomain.RoleAdmin))

	handler.UpdateTask(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteTaskSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, storage := newTaskHandlerForTest(t)
	task := seedTask(t, repo, 10, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), []tasksDomain.TaskAttachment{
		{ID: "att_1", Name: testAttachmentFileName, FilePath: testAttachmentFilePath, ContentType: testAttachmentPDF, Size: 4},
	})
	storage.uploaded[testAttachmentFilePath] = []byte("test")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodDelete, testTaskByIDPath, nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(99, usersDomain.RoleAdmin))

	handler.DeleteTask(c)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if _, ok := repo.tasks[task.ID]; ok {
		t.Fatal("expected task deleted")
	}
	if len(storage.deleted) != 1 {
		t.Fatalf("expected 1 deleted attachment, got %d", len(storage.deleted))
	}
}

func TestDownloadAttachmentSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, storage := newTaskHandlerForTest(t)
	task := seedTask(t, repo, 10, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), []tasksDomain.TaskAttachment{
		{ID: "att_1", Name: testAttachmentFileName, FilePath: testAttachmentFilePath, ContentType: testAttachmentPDF, Size: 4},
	})
	storage.uploaded[testAttachmentFilePath] = []byte("test")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testTaskDownloadPath, nil)
	c.Params = gin.Params{
		{Key: "id", Value: "1"},
		{Key: "attachmentId", Value: "att_1"},
	}
	c.Set("current_user", authenticatedUser(99, usersDomain.RoleAdmin))

	handler.DownloadAttachment(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d for task %d", w.Code, task.ID)
	}
	if body := w.Body.String(); body != "test" {
		t.Fatalf("expected attachment body, got %q", body)
	}
}

func TestCreateTaskMultipartWithAttachment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, storage := newTaskHandlerForTest(t)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	_ = writer.WriteField("assignment_id", "1")
	_ = writer.WriteField("title", testTaskTitle)
	_ = writer.WriteField("description", testTaskDescription)
	_ = writer.WriteField("status", "abierto")
	_ = writer.WriteField("spent_hours", "2")
	_ = writer.WriteField("observations", "")
	_ = writer.WriteField("week_start_date", testTaskWeekStartDate)
	fileWriter, err := writer.CreateFormFile("attachments", "evidence.txt")
	if err != nil {
		t.Fatalf("expected form file, got %v", err)
	}
	if _, err := fileWriter.Write([]byte("hello")); err != nil {
		t.Fatalf("expected file write, got %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("expected writer close, got %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, testTasksPath, &body)
	req.Header.Set(testHeaderContentType, writer.FormDataContentType())
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authenticatedUser(99, usersDomain.RoleAdmin))

	handler.CreateTask(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if len(repo.tasks) != 1 {
		t.Fatalf("expected 1 task, got %d", len(repo.tasks))
	}
	if len(storage.uploaded) != 1 {
		t.Fatalf("expected 1 uploaded file, got %d", len(storage.uploaded))
	}
}

func TestCreateTaskForbiddenForForeignOperationalUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, _, _ := newTaskHandlerForTest(t)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]any{
		"assignment_id":   1,
		"title":           testTaskTitle,
		"description":     testTaskDescription,
		"status":          "abierto",
		"spent_hours":     2,
		"observations":    "",
		"week_start_date": testTaskWeekStartDate,
	})
	req := httptest.NewRequest(http.MethodPost, testTasksPath, bytes.NewBuffer(body))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("current_user", authenticatedUser(999, usersDomain.RoleMonitor))

	handler.CreateTask(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestDownloadAttachmentForbiddenForForeignUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, storage := newTaskHandlerForTest(t)
	seedTask(t, repo, 10, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), []tasksDomain.TaskAttachment{
		{ID: "att_1", Name: testAttachmentFileName, FilePath: testAttachmentFilePath, ContentType: testAttachmentPDF, Size: 4},
	})
	storage.uploaded[testAttachmentFilePath] = []byte("test")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, testTaskDownloadPath, nil)
	c.Params = gin.Params{{Key: "id", Value: "1"}, {Key: "attachmentId", Value: "att_1"}}
	c.Set("current_user", authenticatedUser(200, usersDomain.RoleAssistant))

	handler.DownloadAttachment(c)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestUpdateTaskSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler, repo, _ := newTaskHandlerForTest(t)
	seedTask(t, repo, 10, 1, time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC), nil)

	w := httptest.NewRecorder()
	body, _ := json.Marshal(map[string]any{
		"assignment_id":   1,
		"title":           "Updated title",
		"description":     "Updated description",
		"status":          "finalizado",
		"spent_hours":     2,
		"observations":    "ok",
		"week_start_date": testTaskWeekStartDate,
	})
	req := httptest.NewRequest(http.MethodPut, testTaskByIDPath, bytes.NewBuffer(body))
	req.Header.Set(testHeaderContentType, testApplicationJSON)
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: "1"}}
	c.Set("current_user", authenticatedUser(99, usersDomain.RoleAdmin))

	handler.UpdateTask(c)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}
