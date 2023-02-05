package handlers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/jakenichols2719/simpleblog/pkg/events"
	"github.com/jakenichols2719/simpleblog/pkg/listings"
	"github.com/jakenichols2719/simpleblog/pkg/posts"
)

func HandleCreateEvent(ctx context.Context, event *events.Event) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}
	// Create a new post
	p := posts.NewPost(event.WithContent)
	err = posts.Create(p, posts.Connect(cfg))
	if err != nil {
		return err
	}
	// Create a new listing using the event listing information
	l := event.WithListing
	l.PostID = p.PostID
	l.LiveVersionID = p.VersionID
	l.UploadTimestamp = event.Timestamp
	// Create a new listing
	err = listings.Create(l, listings.Connect(cfg))
	if err != nil {
		return err
	}
	return nil
}

func HandleUpdateEvent(ctx context.Context, event *events.Event) error {
	if event.Target == "" {
		return fmt.Errorf("attempted to invoke a UpdateEvent without a target")
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}

	// Read the old listing for the target
	listingsClient := listings.Connect(cfg)
	l, err := listings.ReadOne(event.Target, listingsClient)
	if err != nil {
		return err
	}
	// If there's new content, upload it and update the VersionID
	if event.WithContent != "" {
		p := posts.NewVersion(event.Target, event.WithContent)
		posts.Create(p, posts.Connect(cfg))
		l.UpdateVersionID(p.VersionID)
	}
	// If there's new listing data, update the old listing.
	if event.WithListing != nil {
		l.UpdatePostInfo(event.WithListing)
	}
	// Send an update to the listing database
	err = listings.Update(l, listingsClient)
	if err != nil {
		return err
	}
	return err
}
