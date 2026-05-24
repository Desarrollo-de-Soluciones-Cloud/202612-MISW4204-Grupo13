package infrastructure

import (
	"context"
	"encoding/json"

	reportsApplication "backend/internal/reports/application"

	"cloud.google.com/go/pubsub"
)

type PubSubReportJobPublisher struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

func NewPubSubReportJobPublisher(
	ctx context.Context,
	projectID string,
	topicName string,
) (*PubSubReportJobPublisher, error) {
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}

	return &PubSubReportJobPublisher{
		client: client,
		topic:  client.Topic(topicName),
	}, nil
}

func (p *PubSubReportJobPublisher) PublishWeeklyReportJob(
	ctx context.Context,
	job reportsApplication.WeeklyReportJobMessage,
) error {
	payload, err := json.Marshal(job)
	if err != nil {
		return err
	}

	result := p.topic.Publish(ctx, &pubsub.Message{Data: payload})
	_, err = result.Get(ctx)
	return err
}

func (p *PubSubReportJobPublisher) Close() error {
	if p.topic != nil {
		p.topic.Stop()
	}

	if p.client != nil {
		return p.client.Close()
	}

	return nil
}
