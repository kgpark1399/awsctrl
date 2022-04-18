package monitor

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"gopkg.in/ini.v1"
)

type C_notice__sms struct {
	cfg aws.Config

	// AWS 접속 인증
	s_aws__access_key string
	s_aws__secret_key string
	s_aws__region     string

	// 메시지 발송 설정
	s_msg__title string
	s_msg__body  string
	s_mobile     string
}

// SMS 발송
func (t *C_notice__sms) Send(_s_msg__title, _s_msg__body, _s_mobile string) error {

	err := t.Init()
	if err != nil {
		return err
	}

	client := sns.NewFromConfig(t.cfg)

	_, err = client.Publish(context.TODO(), &sns.PublishInput{
		Subject:     &_s_msg__title,
		Message:     &_s_msg__body,
		PhoneNumber: &_s_mobile,
	})

	if err != nil {
		log.Println("[ERROR] Failed to aws : ", err)
		return err
	}
	return nil
}

// AWS 환경설정 ini 파일 호출 및 변수 저장
func (t *C_notice__sms) Init() error {

	read, err := ini.Load("config.ini")
	if err != nil {
		log.Println("[ERROR] Not found config.ini file : ", err)
		return err
	}

	title := "aws_configure"
	t.s_aws__access_key = read.Section(title).Key("S_aws__access_key").String()
	t.s_aws__secret_key = read.Section(title).Key("S_aws__secret_key").String()
	t.s_aws__region = read.Section(title).Key("S_aws__region").String()

	t.cfg = aws.Config{
		Region:      *aws.String(t.s_aws__region),
		Credentials: credentials.NewStaticCredentialsProvider(t.s_aws__access_key, t.s_aws__secret_key, ""),
	}

	return nil
}
