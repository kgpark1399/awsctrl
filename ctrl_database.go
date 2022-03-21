package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type C_sns struct {
	sTitle   string
	sRegion  string
	sAcid    string
	sAckey   string
	sSession string
	sTopic   string
	cfg      aws.Config
}

// 문자 발송 함수
func Send_sns(_sMessage string) (string, error) {
	cSMS := New_C_sns()
	cSMS.Init("id", "key", "session", "us-east-1")
	return cSMS.Send(_sMessage, "arn:")
}

// 문자 발송 함수
// func Send_sns(_sAcid, _sAckey, _sSession, _sRegion, _sMessage, _sTopic string) (string, error) {
// 	cSMS := New_C_sns()
// 	cSMS.Init(_sAcid, _sAckey, _sSession, _sRegion)
// 	return cSMS.Send(_sMessage, _sTopic)
// }

// --------------------------------------------------------------------------------------
func New_C_sns() *C_sns {
	c := &C_sns{}
	return c
}

func (t *C_sns) Init(_sAcid, _sAckey, _sSession, _sRegion string) {

	// aws config(id,passwd)
	t.cfg = aws.Config{
		Region:      _sRegion,
		Credentials: credentials.NewStaticCredentialsProvider(_sAcid, _sAckey, _sSession),
	}
}

func (t *C_sns) Send(_sMessage, _sTopic string) (string, error) {

	client := sns.NewFromConfig(t.cfg)

	result, err := client.Publish(context.TODO(), &sns.PublishInput{
		Subject:  aws.String("Server Err"),
		Message:  &_sMessage,
		TopicArn: &_sTopic,
	})

	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}
