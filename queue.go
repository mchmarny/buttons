package main

import (
	"context"
	"errors"

	"cloud.google.com/go/pubsub"
)

// queue pushes events to pubsub topic
type queue struct {
	client *pubsub.Client
	topic  *pubsub.Topic
}

// newQueue is invoked once per Storable life cycle to configure the store
func newQueue(ctx context.Context, projectID, topicName string) (q *queue, err error) {
	logger.Print("Init Queue...")

	if projectID == "" {
		return nil, errors.New("projectID not set")
	}

	if topicName == "" {
		return nil, errors.New("topicName not set")
	}

	if ctx == nil {
		return nil, errors.New("context not set")
	}

	c, e := pubsub.NewClient(ctx, projectID)
	if e != nil {
		return nil, e
	}

	o := &queue{
		client: c,
		topic:  c.Topic(topicName),
	}

	return o, nil
}

// push persist the content
func (q *queue) push(ctx context.Context, data []byte) error {
	msg := &pubsub.Message{Data: data}
	result := q.topic.Publish(ctx, msg)
	_, err := result.Get(ctx)
	return err
}