package monitor

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type C_aws_sms struct {
	cfg aws.Config

	s_region  string
	s_message string
	s_mobile  string
	s_title   string
}

// AWS 환경설정 호출 및 접속
func (t *C_aws_sms) Init() error {
	_cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Println("[ERROR] Failed to connect aws configure : ", err)
		return err
	}
	t.cfg = _cfg
	return nil
}

// 경고 알림 문자 발송
func (t *C_aws_sms) Send__sms(_s_message, _s_mobile string) error {

	client := sns.NewFromConfig(t.cfg)

	_, err := client.Publish(context.TODO(), &sns.PublishInput{
		Subject:     aws.String("Server Err alert"),
		Message:     &_s_message,
		PhoneNumber: &_s_mobile,
	})

	if err != nil {
		log.Println("[ERROR] Failed to aws publish input : ", err)
		return err
	}
	return nil
}

// Send__sms 구조체 없이 실행
func Send__alert_sms(_s_message, _s_mobile string) error {
	t := C_aws_sms{}
	err := t.Init()
	if err != nil {
		return err
	}

	err = t.Send__sms(_s_message, _s_mobile)
	if err != nil {
		return err
	}
	return nil
}
