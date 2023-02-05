package posts

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func Connect(cfg aws.Config) *s3.Client {
	return s3.NewFromConfig(cfg)
}

func Create(p *PostVersion, client *s3.Client) (err error) {
	fmt.Printf("Creating post %s/%s\n", p.PostID, p.VersionID)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("posts/%s/%s.md", p.PostID, p.VersionID)),
		Body:   strings.NewReader(p.Content),
	})
	if err != nil {
		fmt.Printf("Failed to create post %s/%s: %s\n", p.PostID, p.VersionID, err.Error())
		return err
	}
	return nil
}

func Read(pID, vID string, client *s3.Client) (p *PostVersion, err error) {
	fmt.Printf("Reading post %s/%s\n", pID, vID)
	res, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET")),
		Key:    aws.String(fmt.Sprintf("posts/%s/%s.md", pID, vID)),
	})
	if err != nil {
		fmt.Printf("Failed to read post %s/%s: %s\n", pID, vID, err.Error())
		return nil, err
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Failed to read post %s/%s: %s\n", pID, vID, err.Error())
		return nil, err
	}
	p = &PostVersion{
		PostID:    pID,
		VersionID: vID,
	}
	p.Content = string(body)
	return p, nil
}
