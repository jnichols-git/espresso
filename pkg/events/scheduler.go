package events

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	etypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	stypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/google/uuid"
)

type Client struct {
	p *eventbridge.Client
	s *scheduler.Client
}

func Connect(cfg aws.Config) *Client {
	return &Client{
		p: eventbridge.NewFromConfig(cfg),
		s: scheduler.NewFromConfig(cfg),
	}
}

func Schedule(e *Event, client *Client) (err error) {
	raw, err := json.Marshal(e)
	if err != nil {
		return err
	}
	// Parse timestamp from e.
	// If it's after the current time, schedule the event with client.s; otherwise push with client.p
	t := time.Unix(e.Timestamp, 0)
	schedule := false
	if t.After(time.Now()) {
		schedule = true
	}
	if schedule {
		// We need to format the event into a CloudWatchEvent.
		cwe := &events.CloudWatchEvent{
			Source:     os.Getenv("EVENT_NAME"),
			DetailType: "event",
			Detail:     json.RawMessage(raw),
		}
		raw, _ := json.Marshal(cwe)
		detail := string(raw)
		expr := t.Format(time.RFC3339)
		input := &scheduler.CreateScheduleInput{
			FlexibleTimeWindow: &stypes.FlexibleTimeWindow{
				Mode: stypes.FlexibleTimeWindowModeOff,
			},
			Name:               aws.String(uuid.NewString()),
			ScheduleExpression: aws.String(fmt.Sprintf("at(%s)", expr)),
			Target: &stypes.Target{
				Arn:     aws.String(os.Getenv("EVENT_HANDLER_ARN")),
				RoleArn: aws.String(os.Getenv("EVENT_HANDLER_ROLE_ARN")),
				Input:   aws.String(detail),
			},
			ScheduleExpressionTimezone: aws.String("Etc/UTC"),
		}
		_, err = client.s.CreateSchedule(context.TODO(), input)
	} else {
		detail := string(raw)
		input := &eventbridge.PutEventsInput{
			Entries: []etypes.PutEventsRequestEntry{
				{
					Source:       aws.String(os.Getenv("EVENT_NAME")),
					DetailType:   aws.String("event"),
					Detail:       aws.String(detail),
					EventBusName: aws.String("EVENT_BUS"),
				},
			},
			EndpointId: aws.String("t831eg9jtu.veo"),
		}
		_, err = client.p.PutEvents(context.TODO(), input)
	}
	return err
}
