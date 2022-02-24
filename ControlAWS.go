package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	domaintypes "github.com/aws/aws-sdk-go-v2/service/route53domains/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// 구조체 정의
// -----------------------------------------------------------------------------//

type AWS_Ctrl struct {
	TD__S3      C__S3
	TD__RDS     C__RDS
	TD__Route53 C__Route53
}

// RDS DB 구조체
type C__RDS struct {
	s__db_name string
	s__id      string
	s__pw      string
}

// S3 구조체
type C__S3 struct {
	s__bucket__name string
	s__bucket__acl  string
	s__object__name string
	s__object__acl  string
}

// Route53 구조체
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


// 실행 함수 모음
// ------------------------------------------------------------------------------ //

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

// 도메인 등록 함수
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

// 도메인 삭제 함수
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

// DB 생성 함수
func (t *C__RDS) DB__create(_s__db__name, _s__id, _s__pw string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	t.s__db_name = _s__db__name
	t.s__id = _s__id
	t.s__pw = _s__pw

	client := rds.NewFromConfig(cfg)

	// DB 생성 INPUT 값 설정
	db_make_input := &rds.CreateDBClusterInput{
		DBClusterIdentifier: aws.String(t.s__db_name),
		Engine:              aws.String("aurora"),
		EngineMode:          aws.String("serverless"),
		MasterUsername:      aws.String(t.s__id),
		MasterUserPassword:  aws.String(t.s__pw),
		ScalingConfiguration: &rdstypes.ScalingConfiguration{
			AutoPause:             aws.Bool(true),
			MinCapacity:           aws.Int32(1),
			MaxCapacity:           aws.Int32(32),
			SecondsBeforeTimeout:  aws.Int32(300),
			SecondsUntilAutoPause: aws.Int32(300),
			TimeoutAction:         aws.String("ForceApplyCapacityChange"),
		},
	}

	//DB 생성 (makeDBCluster 함수 실행)
	_, err = makeDBCluster(context.TODO(), client, db_make_input)
	if err != nil {
		return err
	}
	return nil
}

// DB 삭제 함수
func (t *C__RDS) DB__remove(_s__db__name string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	t.s__db_name = _s__db__name

	client := rds.NewFromConfig(cfg)

	// DB 삭제 INPUT 값 설정
	db_remove_input := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(t.s__db_name),
		SkipFinalSnapshot:   *aws.Bool(true),
	}

	// DB 삭제 (removeDBCluster 함수 실행)
	_, err = removeDBCluster(context.TODO(), client, db_remove_input)
	if err != nil {
		return err
	}
	return nil
}

// DB 조회 함수
func (t *C__RDS) DB__getinfo(_s__db__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := rds.NewFromConfig(cfg)

	t.s__db_name = _s__db__name

	// DB 조회 INPUT 값 설정
	db__find_input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(t.s__db_name),
	}

	// DB 조회
	result, err := getDBCluster(context.TODO(), client, db__find_input)
	if err != nil {
		return err
	}
	fmt.Println("all done")
	printinfo := result.DBClusters

	fmt.Printf("%v", printinfo)
	return nil
}

// Bucket 생성
func (t *C__S3) Bucket__create(_s__bucket__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s__bucket__name = _s__bucket__name

	bucket__make_input := &s3.CreateBucketInput{
		Bucket: aws.String(t.s__bucket__name),
	}

	_, err = makeBucket(context.TODO(), client, bucket__make_input)

	if err != nil {
		return err
	}

	return nil
}

// Bucket 삭제
func (t *C__S3) Bucket__remove(_s__bukcket__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s__bucket__name = _s__bukcket__name

	bucket__remove_input := &s3.DeleteBucketInput{
		Bucket: aws.String(t.s__bucket__name),
	}

	_, err = removeBucket(context.TODO(), client, bucket__remove_input)
	if err != nil {
		return err
	}
	return nil
}

// Buekct 권한 변경
func (t *C__S3) Bucket__change(_s__bucket__name, _s__bucket__acl string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s__bucket__name = _s__bucket__name
	t.s__bucket__acl = _s__bucket__acl

	bucket__change_input := &s3.PutBucketAclInput{
		Bucket: &t.s__bucket__name,
		ACL:    s3types.BucketCannedACL(t.s__bucket__acl),
	}

	_, err = ChangeBucketAcl(context.TODO(), client, bucket__change_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 업로드
// (sample) _s_ob_name = "C:\\Temp\\test1.txt"
func (t *C__S3) Object__uplaod(_s__bucket__name, _s__object__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s__bucket__name = _s__bucket__name
	t.s__object__name = _s__object__name

	file, err := os.Open(t.s__object__name)
	if err != nil {
		return err
	}

	defer file.Close()

	object__upload_input := &s3.PutObjectInput{
		Bucket: &t.s__bucket__name,
		Key:    &t.s__object__name,
		Body:   file,
	}

	_, err = putFile(context.TODO(), client, object__upload_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 삭제
func (t *C__S3) Object__remove(_s__bucket__name, _s__object__name string) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil
	}

	t.s__bucket__name = _s__bucket__name
	t.s__object__name = _s__object__name

	client := s3.NewFromConfig(cfg)

	object__remove_input := &s3.DeleteObjectInput{
		Bucket: &t.s__bucket__name,
		Key:    &t.s__object__name,
	}

	_, err = DeleteItem(context.TODO(), client, object__remove_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 권한 변경
// (sample) _s_ob_acl = "private"
func (t *C__S3) Object__change(_s__bucket__name, _s__object__name, _s__object__acl string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s__bucket__name = _s__bucket__name
	t.s__object__name = _s__object__name
	t.s__object__acl = _s__object__acl

	object__change_input := &s3.PutObjectAclInput{
		Bucket: &t.s__bucket__name,
		Key:    &t.s__object__name,
		ACL:    s3types.ObjectCannedACL(t.s__object__acl),
	}
	_, err = changeObjectAcl(context.TODO(), client, object__change_input)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}

// ------------------------------------------------------------------------------ //

// AWS API 인터페이스 모음
// <------------------------------------------------------------------------->

// DB 생성 API 정의
type RDSCreateDBClusterAPI interface {
	CreateDBCluster(ctx context.Context,
		params *rds.CreateDBClusterInput,
		optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
}

// DB 삭제 API 정의
type RDSDeleteDBClusterAPI interface {
	DeleteDBCluster(ctx context.Context,
		params *rds.DeleteDBClusterInput,
		optFns ...func(*rds.Options)) (*rds.DeleteDBClusterOutput, error)
}

// DB 조회 API 정의
type RDSDescribeDBClustersAPI interface {
	DescribeDBClusters(ctx context.Context,
		params *rds.DescribeDBClustersInput,
		optFns ...func(*rds.Options)) (*rds.DescribeDBClustersOutput, error)
}

// 버캣 생성 API 정의
type S3CreateBucketAPI interface {
	CreateBucket(ctx context.Context,
		params *s3.CreateBucketInput,
		optFns ...func(*s3.Options)) (*s3.CreateBucketOutput, error)
}

// 버캣 삭제 API 정의
type S3DeleteBucketAPI interface {
	DeleteBucket(ctx context.Context,
		params *s3.DeleteBucketInput,
		optFns ...func(*s3.Options)) (*s3.DeleteBucketOutput, error)
}

// 버캣 권한 변경 함수
type S3PutBucketAclAPI interface {
	PutBucketAcl(ctx context.Context,
		params *s3.PutBucketAclInput,
		optFns ...func(*s3.Options)) (*s3.PutBucketAclOutput, error)
}

// Object 업로드 API 정의
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

// Object 삭제 API 정의
type S3DeleteObjectAPI interface {
	DeleteObject(ctx context.Context,
		params *s3.DeleteObjectInput,
		optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// Object 권한변경 API 정의
type S3PutObjectAclAPI interface {
	PutObjectAcl(ctx context.Context,
		params *s3.PutObjectAclInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectAclOutput, error)
}

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

// <------------------------------------------------------------------------->

// AWS API 함수 모음
// <------------------------------------------------------------------------->

// API DB 생성 함수
func makeDBCluster(c context.Context, api RDSCreateDBClusterAPI, input *rds.CreateDBClusterInput) (*rds.CreateDBClusterOutput, error) {
	return api.CreateDBCluster(c, input)
}

// API DB 삭제 함수
func removeDBCluster(c context.Context, api RDSDeleteDBClusterAPI, input *rds.DeleteDBClusterInput) (*rds.DeleteDBClusterOutput, error) {
	return api.DeleteDBCluster(c, input)
}

// API DB 조회 함수
func getDBCluster(c context.Context, api RDSDescribeDBClustersAPI, input *rds.DescribeDBClustersInput) (*rds.DescribeDBClustersOutput, error) {
	return api.DescribeDBClusters(c, input)
}

// API Bucket 생성 함수
func makeBucket(c context.Context, api S3CreateBucketAPI, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return api.CreateBucket(c, input)
}

// API Bucket 삭제 함수
func removeBucket(c context.Context, api S3DeleteBucketAPI, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return api.DeleteBucket(c, input)
}

// API Bucket 권한 변경 함수
func ChangeBucketAcl(c context.Context, api S3PutBucketAclAPI, input *s3.PutBucketAclInput) (*s3.PutBucketAclOutput, error) {
	return api.PutBucketAcl(c, input)
}

// API Object 업로드 함수
func putFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

// API Object 삭제 함수
func DeleteItem(c context.Context, api S3DeleteObjectAPI, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
	return api.DeleteObject(c, input)
}

// API Object 권한 변경 함수
func changeObjectAcl(c context.Context, api S3PutObjectAclAPI, input *s3.PutObjectAclInput) (*s3.PutObjectAclOutput, error) {
	return api.PutObjectAcl(c, input)
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

// <------------------------------------------------------------------------->

func Test_tmp2(_t *testing.T) {
	ctrlroute53 := C__Route53{}
	// ctrlroute53.Reocrd__change("Z07036302DA6YT5R06WYS", "test22.devtoolstest2.com", "A", "1.1.1.1", "DELETE")
	// ctrlroute53.Domain__register("devtoolstest1231.com.", "kgpark@devtools.kr", false)
	ctrlroute53.Domain__remove("devtoolstest1231.com")
}
