package infrastructure

import "testing"

func TestPubSubReportJobPublisherCloseWithNilFields(t *testing.T) {
	publisher := &PubSubReportJobPublisher{}
	if err := publisher.Close(); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
