package main

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	domaintypes "github.com/aws/aws-sdk-go-v2/service/route53domains/types"
)


type AWS_ctrl struct {
	TD__Route53 C__Route53
}

type C__Route53 struct {
	s__record__name        string
	s__record__type        string
	s__record__value       string
	s__record__action      string
	s__hostzone__id        string
	s__domain__name        string
	s__domain__admin_email string
	b__domain__auto_renew  bool
}

// <--------------------------------------------------------------------------------->

// 레코드 변경 함수(추가, 삭제)
// (sample) s__record__name = "test.devtoolstest2.com"
func (t *C__Route53) Reocrd__change(_s__hostzone__id, _s_record_name, _s_record_type, _s__record__value, _s__record__action string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := route53.NewFromConfig(cfg)

	t.s__hostzone__id = _s__hostzone__id
	t.s__record__name = _s_record_name
	t.s__record__type = _s_record_type
	t.s__record__value = _s__record__value
	t.s__record__action = _s__record__action

	record__create_input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &t.s__hostzone__id,
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action: route53types.ChangeAction(t.s__record__action),
					ResourceRecordSet: &route53types.ResourceRecordSet{
						Name: &t.s__record__name,
						Type: route53types.RRType(t.s__record__type),
						TTL:  aws.Int64(300),
						ResourceRecords: []route53types.ResourceRecord{
							{
								Value: &t.s__record__value,
							},
						},
					},
				},
			},
		},
	}

	// func ChangeRecord 함수 실행(레코드 추가)
	_, err = changeRecord(context.TODO(), client, record__create_input)
	if err != nil {
		return err
	}
	return nil
}

func (t *C__Route53) Domain__register(_s__domain__name, _s__domain__admin_email string, _b__domain__auto_renew bool) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	client := route53domains.NewFromConfig(cfg)

	t.s__domain__name = _s__domain__name
	t.s__domain__admin_email = _s__domain__admin_email
	t.b__domain__auto_renew = _b__domain__auto_renew

	domain__register_input := &route53domains.RegisterDomainInput{
		AdminContact: &domaintypes.ContactDetail{
			AddressLine1:     aws.String("1 Main Street"),
			AddressLine2:     aws.String("2 Main Street"),
			City:             aws.String("Suwon"),
			ContactType:      domaintypes.ContactTypePerson,
			CountryCode:      domaintypes.CountryCodeKr,
			Email:            &t.s__domain__admin_email,
			FirstName:        aws.String("devtools"),
			LastName:         aws.String("ltd"),
			OrganizationName: aws.String("devtoolstest"),
			PhoneNumber:      aws.String("+82.0212341234"),
			ZipCode:          aws.String("12345"),
		},
		DomainName:                      &t.s__domain__name,
		DurationInYears:                 aws.Int32(1),
		PrivacyProtectAdminContact:      aws.Bool(false),
		PrivacyProtectRegistrantContact: aws.Bool(false),
		PrivacyProtectTechContact:       aws.Bool(false),
		RegistrantContact: &domaintypes.ContactDetail{
			AddressLine1:     aws.String("1 Main Street"),
			AddressLine2:     aws.String("2 Main Street"),
			City:             aws.String("Suwon"),
			ContactType:      domaintypes.ContactTypePerson,
			CountryCode:      domaintypes.CountryCodeKr,
			Email:            &t.s__domain__admin_email,
			FirstName:        aws.String("devtools"),
			LastName:         aws.String("ltd"),
			OrganizationName: aws.String("devtoolstest"),
			PhoneNumber:      aws.String("+82.0212341234"),
			ZipCode:          aws.String("12345"),
		},
		TechContact: &domaintypes.ContactDetail{
			AddressLine1:     aws.String("1 Main Street"),
			AddressLine2:     aws.String("2 Main Street"),
			City:             aws.String("Suwon"),
			ContactType:      domaintypes.ContactTypePerson,
			CountryCode:      domaintypes.CountryCodeKr,
			Email:            &t.s__domain__admin_email,
			FirstName:        aws.String("devtools"),
			LastName:         aws.String("ltd"),
			OrganizationName: aws.String("devtoolstest"),
			PhoneNumber:      aws.String("+82.0212341234"),
			ZipCode:          aws.String("12345"),
		},
		AutoRenew: &t.b__domain__auto_renew,
	}

	_, err = createDomain(context.TODO(), client, domain__register_input)
	if err != nil {
		return err
	}
	return nil
}

func (t *C__Route53) Domain__remove(_s__domain__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	client := route53domains.NewFromConfig(cfg)
	t.s__domain__name = _s__domain__name

	domain__remove_input := &route53domains.DeleteDomainInput{
		DomainName: &t.s__domain__name,
	}

	_, err = removeDomain(context.TODO(), client, domain__remove_input)
	if err != nil {
		return err
	}
	return nil
}

// <------------------------------------------------------------------------------------>

//레코드 변경(추가,삭제) API 인터페이스
type ROUTE53ChangeRecordAPI interface {
	ChangeResourceRecordSets(ctx context.Context,
		params *route53.ChangeResourceRecordSetsInput,
		optFns ...func(*route53.Options)) (*route53.ChangeResourceRecordSetsOutput, error)
}

//도메인 등록 API 인터페이스
type ROUTE53RegisterDomainAPI interface {
	RegisterDomain(ctx context.Context,
		params *route53domains.RegisterDomainInput,
		optFns ...func(*route53domains.Options)) (*route53domains.RegisterDomainOutput, error)
}

// 도메인 삭제 API 인터페이스
type ROUTE53DeleteDomainAPI interface {
	DeleteDomain(ctx context.Context,
		params *route53domains.DeleteDomainInput,
		optFns ...func(*route53domains.Options)) (*route53domains.DeleteDomainOutput, error)
}

//레코드 변경(추가,삭제) 함수
func changeRecord(c context.Context, api ROUTE53ChangeRecordAPI, input *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
	return api.ChangeResourceRecordSets(c, input)
}

//도메인 등록 함수
func createDomain(c context.Context, api ROUTE53RegisterDomainAPI, input *route53domains.RegisterDomainInput) (*route53domains.RegisterDomainOutput, error) {
	return api.RegisterDomain(c, input)
}

// 도메인 삭제 함수
func removeDomain(c context.Context, api ROUTE53DeleteDomainAPI, input *route53domains.DeleteDomainInput) (*route53domains.DeleteDomainOutput, error) {
	return api.DeleteDomain(c, input)
}

// <------------------------------------------------------------------------------------>

// 함수 동작 테스트
func Test_route53(_t *testing.T) {
	ctrlroute53 := C__Route53{}
	//	ctrlroute53.Reocrd__change("Z07036302DA6YT5R06WYS", "test22.devtoolstest2.com", "A", "1.1.1.1", "DELETE")
	// ctrlroute53.Domain__register("devtoolstest1231.com.", "kgpark@devtools.kr", false)
	ctrlroute53.Domain__remove("devtoolstest1231.com")
}
