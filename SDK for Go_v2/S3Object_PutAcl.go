package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3PutObjectAclAPI interface {
	PutObjectAcl(ctx context.Context,
		params *s3.PutObjectAclInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectAclOutput, error)
}

func ChangeObjectAcl(c context.Context, api S3PutObjectAclAPI, input *s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error) {
	return api.PutObjectAcl(c, input)
}

func S3PutObjectAcl() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := s3.NewFromConfig(cfg)

	input := &s3.PutObjectAclInput{
		Bucket: aws.String("devtoolstestbucket1231"),
		Key:    aws.String("asdasdasd.txt"),
		ACL:    types.ObjectCannedACLPublicReadWrite,
	}

	_, err = ChangeObjectAcl(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Could not change acls ")
	}
}
