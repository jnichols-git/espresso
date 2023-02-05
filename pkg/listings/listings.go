package listings

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

/*
- post_id (string)
- upload_date (timestamp)
- live_version_id (string)
- archive_version_ids ([]string)
- title (string)
- description (string)
- tags ([]string)
*/

const (
	partitionKey = "post_id"
	sortKey      = "upload_timestamp"
)

type Listing struct {
	PostID            string   `json:"post_id" dynamodbav:"post_id"`
	UploadTimestamp   int64    `json:"upload_timestamp" dynamodbav:"upload_timestamp"`
	UpdateTimestamp   int64    `json:"update_timestamp" dynamodbav:"update_timestamp"`
	LiveVersionID     string   `json:"live_version_id" dynamodbav:"live_version_id"`
	ArchiveVersionIDs []string `json:"archive_version_ids" dynamodbav:"archive_version_ids"`
	Title             string   `json:"title" dynamodbav:"title"`
	Description       string   `json:"description" dynamodbav:"description"`
	Tags              []string `json:"tags" dynamodbav:"tags"`
}

func New(pID string, uTS int64, vID string, avIDs []string, title, desc string, tags []string) *Listing {
	return &Listing{
		PostID:            pID,
		UploadTimestamp:   uTS,
		UpdateTimestamp:   uTS,
		LiveVersionID:     vID,
		ArchiveVersionIDs: avIDs,
		Title:             title,
		Description:       desc,
		Tags:              tags,
	}
}

func (l *Listing) UpdateVersionID(newID string) {
	if l.ArchiveVersionIDs == nil {
		l.ArchiveVersionIDs = []string{l.LiveVersionID}
	} else {
		l.ArchiveVersionIDs = append(l.ArchiveVersionIDs, l.LiveVersionID)
	}
	l.LiveVersionID = newID
	l.UpdateTimestamp = time.Now().UTC().Unix()
}

// Update a listing with a new one.
// This updates the Title, Description, and Tags of the calling object with those of the second.
func (l *Listing) UpdatePostInfo(other *Listing) {
	if other == nil {
		return
	}
	if other.Title != "" {
		l.Title = other.Title
	}
	if other.Description != "" {
		l.Description = other.Description
	}
	if other.Tags != nil {
		l.Tags = other.Tags
	}
}

func (l *Listing) DynamoDBKey() map[string]types.AttributeValue {
	pk, _ := attributevalue.Marshal(l.PostID)
	sk, _ := attributevalue.Marshal(l.UploadTimestamp)
	return map[string]types.AttributeValue{
		partitionKey: pk,
		sortKey:      sk,
	}
}

func (l *Listing) UnmarshalDynamoDBAV(doc map[string]types.AttributeValue) (err error) {
	return attributevalue.UnmarshalMap(doc, l)
}

func (l *Listing) MarshalDynamoDBAV() (doc map[string]types.AttributeValue, err error) {
	return attributevalue.MarshalMap(l)
}
