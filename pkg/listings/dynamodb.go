package listings

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"golang.org/x/exp/slices"
)

func Connect(cfg aws.Config) *dynamodb.Client {
	return dynamodb.NewFromConfig(cfg)
}

// Create a listing in the blog table.
func Create(l *Listing, client *dynamodb.Client) (err error) {
	item, err := l.MarshalDynamoDBAV()
	if err != nil {
		return err
	}
	_, err = client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DB_TABLE")),
		Item:      item,
	})
	if err != nil {
		return err
	}
	return nil
}

// Read a listing with the given key.
func ReadOne(key string, client *dynamodb.Client) (l *Listing, err error) {
	statement := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s" = '%s'`, os.Getenv("DB_TABLE"), partitionKey, key)
	res, err := client.ExecuteStatement(context.TODO(), &dynamodb.ExecuteStatementInput{
		Statement: aws.String(statement),
	})
	if err != nil {
		return nil, err
	}
	l = &Listing{}
	err = attributevalue.UnmarshalMap(res.Items[0], l)
	if err != nil {
		return nil, err
	}
	return l, nil
}

// Read up to pc listings, ordered by UploadTimestamp, starting at the last listing
func ReadMany(pc, pn int, filter func(*Listing) bool, client *dynamodb.Client) (ls []*Listing, err error) {
	qin := &dynamodb.ScanInput{
		TableName: aws.String(os.Getenv("DB_TABLE")),
	}
	res, err := client.Scan(context.TODO(), qin)
	if err != nil {
		return nil, err
	}
	ls = make([]*Listing, 0)
	err = attributevalue.UnmarshalListOfMaps(res.Items, &ls)
	if err != nil {
		return nil, err
	}
	slices.SortFunc(ls, func(a, b *Listing) bool {
		return a.UploadTimestamp > b.UpdateTimestamp
	})
	filtered := make([]*Listing, 0)
	for _, listing := range ls {
		if filter(listing) {
			filtered = append(filtered, listing)
		}
	}
	s := pc * pn
	e := pc * (pn + 1)
	if e > len(filtered) {
		e = len(filtered)
	}
	return filtered[s:e], nil
}

// Update a listing in the blog table.
func Update(l *Listing, client *dynamodb.Client) (err error) {
	item, err := l.MarshalDynamoDBAV()
	if err != nil {
		return err
	}
	builder := expression.UpdateBuilder{}
	// Iterate through the item.
	// For anything that isn't the partition or sort key, add them to the UpdateBuilder.
	for k, v := range item {
		if k == partitionKey || k == sortKey {
			continue
		}
		builder = builder.Set(expression.Name(k), expression.Value(v))
	}
	// Create the update expression
	expr, err := expression.NewBuilder().WithUpdate(builder).Build()
	if err != nil {
		panic(err)
	}
	_, err = client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName:                 aws.String(os.Getenv("DB_TABLE")),
		Key:                       l.DynamoDBKey(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
	})
	if err != nil {
		return err
	}
	return nil
}
