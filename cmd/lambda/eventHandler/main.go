package main

import (
	"context"
	"encoding/json"
	"fmt"

	aws_events "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jakenichols2719/simpleblog/pkg/events"
	"github.com/jakenichols2719/simpleblog/pkg/handlers"
)

func HandleEvent(ctx context.Context, awsEvent aws_events.CloudWatchEvent) (err error) {
	event := &events.Event{}
	err = json.Unmarshal(awsEvent.Detail, event)
	if err != nil {
		return err
	}
	if event.Kind == events.CreateEvent {
		return handlers.HandleCreateEvent(ctx, event)
	}
	if event.Kind == events.UpdateEvent {
		return handlers.HandleUpdateEvent(ctx, event)
	}
	return fmt.Errorf("invalid event kind %s", event.Kind)
}

func main() {
	lambda.Start(HandleEvent)
}
