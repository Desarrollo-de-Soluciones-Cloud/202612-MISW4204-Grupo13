package infrastructure_test

import (
	reportsApplication "backend/internal/reports/application"
	reportsInfrastructure "backend/internal/reports/infrastructure"
	"context"
	"testing"
	"time"
)

func TestNewPubSubReportJobPublisherWithEmulatorAndPublishFailsCleanly(t *testing.T) {
	t.Setenv("PUBSUB_EMULATOR_HOST", "127.0.0.1:8681")

	publisher, err := reportsInfrastructure.NewPubSubReportJobPublisher(context.Background(), "test-project", "weekly-reports")
	if err != nil {
		t.Fatalf("expected publisher creation with emulator host, got %v", err)
	}
	defer func() { _ = publisher.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err = publisher.PublishWeeklyReportJob(ctx, reportsApplication.WeeklyReportJobMessage{
		WorkspaceID: 1,
		WeekID:      2,
		UserID:      3,
	})
	if err == nil {
		t.Fatalf("expected publish to fail without emulator server")
	}
}
