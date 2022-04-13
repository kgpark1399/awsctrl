package monitor

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"gopkg.in/ini.v1"
)

type C_aws_sms struct {
	cfg aws.Config

	// AWS 접속 인증
	s_access__id  string
	s_access__key string

	s_region  string
	s_message string
	s_mobile  string
	s_title   string
}

// AWS 환경설정 ini 파일 호출 및 변수 저장
func (t *C_aws_sms) Init() error {

	read, err := ini.Load("config.ini")
	if err != nil {
		log.Println("[ERROR] Fail to read config.ini file : ", err)
		return err
	}

	t.s_access__id = read.Section("AWS").Key("S_access__id").String()
	t.s_access__key = read.Section("AWS").Key("S_access__key").String()
	t.s_region = read.Section("AWS").Key("S_region").String()

	t.cfg = aws.Config{
		Region:      *aws.String(t.s_region),
		Credentials: credentials.NewStaticCredentialsProvider(t.s_access__id, t.s_access__key, ""),
	}

	return nil
}

// 경고 알림 문자 발송
func (t *C_aws_sms) Send(_s_message, _s_mobile string) error {

	err := t.Init()
	if err != nil {
		return err
	}

	client := sns.NewFromConfig(t.cfg)

	_, err = client.Publish(context.TODO(), &sns.PublishInput{
		Subject:     aws.String("Server Err alert"),
		Message:     &_s_message,
		PhoneNumber: &_s_mobile,
	})

	if err != nil {
		log.Println("[ERROR] Failed to aws : ", err)
		return err
	}
	return nil
}
