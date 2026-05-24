package infrastructure_test

import (
	reportsInfrastructure "backend/internal/reports/infrastructure"
	"testing"
)

func TestPubSubReportJobPublisherCloseWithNilFields(t *testing.T) {
	publisher := &reportsInfrastructure.PubSubReportJobPublisher{}
	if err := publisher.Close(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
