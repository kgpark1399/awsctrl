package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	domaintypes "github.com/aws/aws-sdk-go-v2/service/route53domains/types"
)

type AWS_ctrl struct {
	C_route53 C_Route53
}

type C_Route53 struct {
	s_record__name        string
	s_record__type        string
	s_record__value       string
	s_record__action      string
	s_hostzone__id        string
	s_domain__name        string
	s_domain__admin_email string
	b_domain__auto_renew  bool
	cfg                   aws.Config
}

func (t *C_Route53) aws_config() error {
	_cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return err
	}
	t.cfg = _cfg
	return nil
}

func init() {
	var err error
	load := C_Route53{}
	load.aws_config()
	if err != nil {
		fmt.Printf("err : %v", err)
		return
	}
}

// 레코드 변경 함수(추가, 삭제)
// (sample) s__record__name = "test.devtoolstest2.com"
func (t *C_Route53) Record_change(_s_hostzone__id, _s_record_name, _s_record_type, _s_record__value, _s_record__action string) error {

	client := route53.NewFromConfig(t.cfg)

	t.s_hostzone__id = _s_hostzone__id
	t.s_record__name = _s_record_name
	t.s_record__type = _s_record_type
	t.s_record__value = _s_record__value
	t.s_record__action = _s_record__action

	record_create_input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: &t.s_hostzone__id,
		ChangeBatch: &route53types.ChangeBatch{
			Changes: []route53types.Change{
				{
					Action: route53types.ChangeAction(t.s_record__action),
					ResourceRecordSet: &route53types.ResourceRecordSet{
						Name: &t.s_record__name,
						Type: route53types.RRType(t.s_record__type),
						TTL:  aws.Int64(300),
						ResourceRecords: []route53types.ResourceRecord{
							{
								Value: &t.s_record__value,
							},
						},
					},
				},
			},
		},
	}
	// func ChangeRecord 함수 실행(레코드 추가)
	_, err := changeRecord(context.TODO(), client, record_create_input)
	if err != nil {
		return err
	}
	return nil
}

func (t *C_Route53) Domain_register(_s_domain__name, _s_domain__admin_email string, _b_domain__auto_renew bool) error {
	client := route53domains.NewFromConfig(t.cfg)

	t.s_domain__name = _s_domain__name
	t.s_domain__admin_email = _s_domain__admin_email
	t.b_domain__auto_renew = _b_domain__auto_renew

	domain_register_input := &route53domains.RegisterDomainInput{
		AdminContact: &domaintypes.ContactDetail{
			AddressLine1:     aws.String("1 Main Street"),
			AddressLine2:     aws.String("2 Main Street"),
			City:             aws.String("Suwon"),
			ContactType:      domaintypes.ContactTypePerson,
			CountryCode:      domaintypes.CountryCodeKr,
			Email:            &t.s_domain__admin_email,
			FirstName:        aws.String("devtools"),
			LastName:         aws.String("ltd"),
			OrganizationName: aws.String("devtoolstest"),
			PhoneNumber:      aws.String("+82.0212341234"),
			ZipCode:          aws.String("12345"),
		},
		DomainName:                      &t.s_domain__name,
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
			Email:            &t.s_domain__admin_email,
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
			Email:            &t.s_domain__admin_email,
			FirstName:        aws.String("devtools"),
			LastName:         aws.String("ltd"),
			OrganizationName: aws.String("devtoolstest"),
			PhoneNumber:      aws.String("+82.0212341234"),
			ZipCode:          aws.String("12345"),
		},
		AutoRenew: &t.b_domain__auto_renew,
	}

	_, err := createDomain(context.TODO(), client, domain_register_input)
	if err != nil {
		return err
	}
	return nil
}

func (t *C_Route53) Domain_remove(_s_domain__name string) error {
	client := route53domains.NewFromConfig(t.cfg)
	t.s_domain__name = _s_domain__name

	domain_remove_input := &route53domains.DeleteDomainInput{
		DomainName: &t.s_domain__name,
	}

	_, err := removeDomain(context.TODO(), client, domain_remove_input)
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
