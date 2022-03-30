package monitor

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type C_sns struct {
	cfg aws.Config

	s_region  string
	s_message string
	s_mobile  string
	s_title   string
}

// 문자 발송 함수
func Send_sns(_s_message, _s_mobile string) error {
	t := C_sns{}
	err := t.Init_config()
	if err != nil {
		return err
	}

	err = t.Send(_s_message, _s_mobile)
	if err != nil {
		return err
	}

	return nil
}

func (t *C_sns) Init_config() error {
	_cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		fmt.Print(err)
		return err
	}
	t.cfg = _cfg
	return nil
}

func (t *C_sns) Send(_s_message, _s_mobile string) error {

	client := sns.NewFromConfig(t.cfg)

	_, err := client.Publish(context.TODO(), &sns.PublishInput{
		Subject:     aws.String("Server Err alert"),
		Message:     &_s_message,
		PhoneNumber: &_s_mobile,
	})

	if err != nil {
		fmt.Print(err)
		return err
	}
	return nil
}
