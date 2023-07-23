package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const BucketName = "todoz"

type Store struct {
	client  *s3.Client
	context context.Context
}

func (store *Store) setupAws() (err error) {
	store.context = context.Background()

	cfg, err := config.LoadDefaultConfig(store.context)

	if err != nil {
		return err
	}

	store.client = s3.NewFromConfig(cfg)

	output, err := store.client.ListObjectsV2(store.context, &s3.ListObjectsV2Input{
		Bucket: aws.String("todoz"),
	})

	for _, object := range output.Contents {
		log.Printf("key=%s size=%d", aws.ToString(object.Key), object.Size)
	}

	if err != nil {
		log.Fatal(err)
	}

	return
}

func (store *Store) putTodoList(list TodoList) (err error) {
	if len(list.Id) < 1 {
		return errors.New("invalid id")
	}

	contentInBytes, err := json.Marshal(list)
	if err != nil {
		log.Println(err)
		return err
	}

	obj := s3.PutObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(list.Id),
		Body:   bytes.NewReader(contentInBytes),
	}
	output, err := store.client.PutObject(context.TODO(), &obj)

	if err != nil {
		return err
	}

	log.Println(output.ETag)
	return nil
}

func (store *Store) getTodoList(id string) (list TodoList, err error) {
	obj := s3.GetObjectInput{
		Bucket: aws.String(BucketName),
		Key:    aws.String(id),
	}

	response, err := store.client.GetObject(context.TODO(), &obj)
	if err != nil {
		return
	}

	list = TodoList{}
	err = json.NewDecoder(response.Body).Decode(&list)

	if err != nil {
		return
	}

	log.Println(response.ETag)
	return
}
