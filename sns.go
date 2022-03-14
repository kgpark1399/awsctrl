package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type C_sns struct {
	s_title   string
	s_region  string
	s_acid    string
	s_ackey   string
	s_session string
	s_number  string
	cfg       aws.Config
}

func Send_Alram(_s_acid, _s_ackey, _s_session, _s_region, _s_message, _s_number string) (string, error) {
	cSMS := New_C_monitor()
	cSMS.Init(_s_acid, _s_ackey, _s_session, _s_region)
	return cSMS.Send(_s_message, _s_number)
}

func New_C_monitor() *C_sns {
	c := &C_sns{}
	return c
}

func (t *C_sns) Init(_s_acid, _s_ackey, _s_session, _s_region string) {

	t.cfg = aws.Config{
		Region:      _s_region,
		Credentials: credentials.NewStaticCredentialsProvider(_s_acid, _s_ackey, _s_session),
	}
}

func (t *C_sns) Send(_s_message, _s_number string) (string, error) {

	client := sns.NewFromConfig(t.cfg)

	result, err := client.Publish(context.TODO(), &sns.PublishInput{
		Subject:     aws.String("Server Err"),
		Message:     &_s_message,
		PhoneNumber: &_s_number,
	})

	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}
