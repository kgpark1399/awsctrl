package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	rdstypes "github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/route53domains"
	domaintypes "github.com/aws/aws-sdk-go-v2/service/route53domains/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type AWS_Ctrl struct {
	C_s3      C_S3
	C_rds     C_RDS
	C_route53 C_Route53
	C_ec2     C_EC2
	C_network C_Network
}

// EC2 구조체
type C_EC2 struct {
	s_ins__id    string
	s_ins__imgid string
	s_ins__type  string
	i_ins__count int32
	s_sg_name    string
	s_keyname    string
	s_tag__name  string
	s_tag__value string
}

// 네트워크 및 보안 설정
type C_Network struct {
	s_key_name       string
	s_sg__name       string
	s_sg__id         string
	s_sg__comment    string
	s_ins__id        string
	s_sg__cirdip     string
	s_sg__start_port int32
	s_sg__end_port   int32
	s_sg__protocal   string
}

// RDS 구조체
type C_RDS struct {
	s_db_name string
	s_id      string
	s_pw      string
}

// S3 구조체
type C_S3 struct {
	s_bucket__name string
	s_bucket__acl  string
	s_object__name string
	s_object__acl  string
}

// Route53 구조체
type C_Route53 struct {
	s_record__name        string
	s_record__type        string
	s_record__value       string
	s_record__action      string
	s_hostzone__id        string
	s_hostzone__name      string
	s_hostzone_time_stamp string
	s_domain__name        string
	s_domain__admin_email string
	b_domain__auto_renew  bool
}

// 1. EC2

// 1.1 Instance 생성
//---------------------------------------------------------------------------------//
func (t *C_EC2) Instance__create(_s_ins__imgid, _s_ins__type, _s__sg_name, _s_keyname, _s_tag__name, _s_tag__value string, _i_ins__count int32) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}
	client := ec2.NewFromConfig(cfg)
	t.s_ins__imgid = _s_ins__imgid
	t.s_ins__type = _s_ins__type
	t.s_sg_name = _s__sg_name
	t.s_keyname = _s_keyname
	t.s_tag__name = _s_tag__name
	t.s_tag__value = _s_tag__value
	t.i_ins__count = _i_ins__count

	insatnce__create_input := &ec2.RunInstancesInput{
		ImageId:        &t.s_ins__imgid,
		InstanceType:   ec2types.InstanceType(t.s_ins__type),
		MinCount:       aws.Int32(1),
		MaxCount:       &t.i_ins__count,
		SecurityGroups: []string{t.s_sg_name},
		KeyName:        &t.s_keyname,
	}

	result, err := makeIns(context.TODO(), client, insatnce__create_input)
	if err != nil {
		return err
	}

	// (sample) tag__name : Name , tag__value : webserver
	instance__tag_input := &ec2.CreateTagsInput{
		Resources: []string{*result.Instances[0].InstanceId},
		Tags: []ec2types.Tag{
			{
				Key:   &t.s_tag__name,
				Value: &t.s_tag__value,
			},
		},
	}

	_, err = makeTags(context.TODO(), client, instance__tag_input)
	if err != nil {
		return err
	}

	return nil
}

//---------------------------------------------------------------------------------//

// 1.2 EC2 Instance 삭제
//---------------------------------------------------------------------------------//
func (t *C_EC2) Instance__remove(_s_ins__id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)

	t.s_ins__id = _s_ins__id

	insatnce__remove_input := &ec2.TerminateInstancesInput{
		InstanceIds: []string{t.s_ins__id},
	}

	_, err = removeIns(context.TODO(), client, insatnce__remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 1.3 EC2 Instance 시작
//---------------------------------------------------------------------------------//
func (t *C_EC2) Instance__run(_s_ins_id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_ins__id = _s_ins_id

	insatnce__run_input := &ec2.StartInstancesInput{
		InstanceIds: []string{t.s_ins__id},
	}

	_, err = startIns(context.TODO(), client, insatnce__run_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 1.4 EC2 Instance 중지
//---------------------------------------------------------------------------------//
func (t *C_EC2) Instance__stop(_s_ins_id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_ins__id = _s_ins_id

	insatnce__stop_input := &ec2.StopInstancesInput{
		InstanceIds: []string{t.s_ins__id},
	}

	_, err = stopIns(context.TODO(), client, insatnce__stop_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 2. EC2 Netowkr & Security

// 2.1 EC2 Key-pair 생성
//---------------------------------------------------------------------------------//
func (t *C_Network) Keypair__create(_s_key_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_key_name = _s_key_name

	keypair__create_input := &ec2.CreateKeyPairInput{
		KeyName: &t.s_key_name,
		KeyType: ec2types.KeyTypeRsa,
	}

	result, err := makeKeypair(context.TODO(), client, keypair__create_input)
	if err != nil {
		return err
	}

	output_keyinfo := result.KeyMaterial
	fmt.Printf("key info : %s", *output_keyinfo)
	return nil
}

//---------------------------------------------------------------------------------//

// 2.2 EC2 Key-pair 삭제
//---------------------------------------------------------------------------------//
func (t *C_Network) Keypair__remove(_s_key_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_key_name = _s_key_name

	keypair__remove_input := &ec2.DeleteKeyPairInput{
		KeyName: &t.s_key_name,
	}

	_, err = removeKeypair(context.TODO(), client, keypair__remove_input)
	if err != nil {
		return err
	}

	return nil

}

//---------------------------------------------------------------------------------//

// 2.3 EC2 Elastic IP 할당 및 부여
//---------------------------------------------------------------------------------//
func (t *C_Network) Instance__associate_ip(_s_ins__id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_ins__id = _s_ins__id

	input := &ec2.AllocateAddressInput{
		Domain: ec2types.DomainTypeStandard,
	}

	result, err := allocateIP(context.TODO(), client, input)
	if err != nil {
		return err
	}

	ipaddress := result.PublicIp

	ipinput := &ec2.AssociateAddressInput{
		InstanceId: &t.s_ins__id,
		PublicIp:   aws.String(*ipaddress),
	}
	_, err = associateIP(context.TODO(), client, ipinput)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 2.4 Security Group 생성
//---------------------------------------------------------------------------------//
func (t *C_Network) Securitygroup__create(_s_sg__name, _s_sg__comment string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_sg__name = _s_sg__name
	t.s_sg__comment = _s_sg__comment

	sg_create_input := &ec2.CreateSecurityGroupInput{
		Description: &t.s_sg__comment,
		GroupName:   &t.s_sg__name,
	}
	_, err = createsg(context.TODO(), client, sg_create_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 2.5 _Security Group 삭제
//---------------------------------------------------------------------------------//
func (t *C_Network) Securitygroup_remove(_s_sg__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := ec2.NewFromConfig(cfg)
	t.s_sg__name = _s_sg__name

	sg_remove_input := &ec2.DeleteSecurityGroupInput{
		GroupName: &t.s_sg__name,
	}
	_, err = removesg(context.TODO(), client, sg_remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 2.6 _Security Group Inbound 규칙 추가
func (t *C_Network) Securitygroup_Inbound(_s_sg__id, _s_sg__cirdip, _s_sg__protocal string, _s_sg__start_port, _s_sg__end_port int32) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		fmt.Println(err)
	}

	client := ec2.NewFromConfig(cfg)

	t.s_sg__id = _s_sg__id
	t.s_sg__cirdip = _s_sg__cirdip
	t.s_sg__protocal = _s_sg__protocal
	t.s_sg__start_port = _s_sg__start_port
	t.s_sg__end_port = _s_sg__end_port

	sg_inbound_input := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:    &t.s_sg__id,
		CidrIp:     &t.s_sg__cirdip,
		IpProtocol: &t.s_sg__protocal,
		FromPort:   &t.s_sg__start_port,
		ToPort:     &t.s_sg__end_port,
	}

	_, err = change_sg__insbound(context.TODO(), client, sg_inbound_input)
	if err != nil {
		fmt.Println(err)
	}
}

//---------------------------------------------------------------------------------//

//---------------------------------------------------------------------------------//

// 3. RDS

// 3.1 도메인 등록
//---------------------------------------------------------------------------------//
func (t *C_Route53) Domain__register(_s_domain__name, _s_domain__admin_email string, _b_domain__auto_renew bool) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	client := route53domains.NewFromConfig(cfg)

	t.s_domain__name = _s_domain__name
	t.s_domain__admin_email = _s_domain__admin_email
	t.b_domain__auto_renew = _b_domain__auto_renew

	domain__register_input := &route53domains.RegisterDomainInput{
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

	_, err = createDomain(context.TODO(), client, domain__register_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 3.2 도메인 삭제 함수
//---------------------------------------------------------------------------------//
func (t *C_Route53) Domain__remove(_s_domain__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	client := route53domains.NewFromConfig(cfg)
	t.s_domain__name = _s_domain__name

	domain__remove_input := &route53domains.DeleteDomainInput{
		DomainName: &t.s_domain__name,
	}

	_, err = removeDomain(context.TODO(), client, domain__remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 3.3 레코드 변경 함수(추가, 삭제)
//---------------------------------------------------------------------------------//
// (sample) s__record__name = "test.devtoolstest2.com"
func (t *C_Route53) Reocrd__change(_s_hostzone__id, _s_record_name, _s_record_type, _s_record__value, _s_record__action string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return err
	}

	client := route53.NewFromConfig(cfg)

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
	_, err = changeRecord(context.TODO(), client, record_create_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 3.4 호스트존 생성
//---------------------------------------------------------------------------------//
func (t *C_Route53) Hostzone__Create(_s_hostzone__name, _s_hostzone__time_stamp string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	t.s_hostzone__name = _s_hostzone__name
	t.s_hostzone_time_stamp = _s_hostzone__time_stamp

	client := route53.NewFromConfig(cfg)

	hostzone_create_input := &route53.CreateHostedZoneInput{
		// timestamp "2021-01-25-15:38"
		CallerReference: &t.s_hostzone_time_stamp,
		Name:            &t.s_hostzone__name,
		HostedZoneConfig: &route53types.HostedZoneConfig{
			PrivateZone: *aws.Bool(false),
		},
	}

	_, err = makeHostzone(context.TODO(), client, hostzone_create_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 3.5 호스트존 삭제
//---------------------------------------------------------------------------------//
func (t *C_Route53) Hostzone__remove(_s_hostzone__id string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return err
	}

	t.s_hostzone__id = _s_hostzone__id

	client := route53.NewFromConfig(cfg)

	hostzone_remove_input := &route53.DeleteHostedZoneInput{
		Id: &t.s_hostzone__id,
	}

	_, err = removeHostzone(context.TODO(), client, hostzone_remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 4. RDS

// 4.1 DB 생성
//---------------------------------------------------------------------------------//
func (t *C_RDS) DB__create(_s_db__name, _s_id, _s_pw string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	t.s_db_name = _s_db__name
	t.s_id = _s_id
	t.s_pw = _s_pw

	client := rds.NewFromConfig(cfg)

	// DB 생성 INPUT 값 설정
	db_make_input := &rds.CreateDBClusterInput{
		DBClusterIdentifier: aws.String(t.s_db_name),
		Engine:              aws.String("aurora"),
		EngineMode:          aws.String("serverless"),
		MasterUsername:      aws.String(t.s_id),
		MasterUserPassword:  aws.String(t.s_pw),
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

//---------------------------------------------------------------------------------//

// 4.2 DB 삭제
//---------------------------------------------------------------------------------//
func (t *C_RDS) DB__remove(_s_db__name string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	t.s_db_name = _s_db__name

	client := rds.NewFromConfig(cfg)

	// DB 삭제 INPUT 값 설정
	db_remove_input := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(t.s_db_name),
		SkipFinalSnapshot:   *aws.Bool(true),
	}

	// DB 삭제 (removeDBCluster 함수 실행)
	_, err = removeDBCluster(context.TODO(), client, db_remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 4.3 DB 조회 함수
//---------------------------------------------------------------------------------//
func (t *C_RDS) DB__getinfo(_s_db__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := rds.NewFromConfig(cfg)

	t.s_db_name = _s_db__name

	// DB 조회 INPUT 값 설정
	db__find_input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(t.s_db_name),
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

//---------------------------------------------------------------------------------//

// 5. S3

// 5.1 Bucket 생성
//---------------------------------------------------------------------------------//
func (t *C_S3) Bucket__create(_s_bucket__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket__name = _s_bucket__name

	bucket_make_input := &s3.CreateBucketInput{
		Bucket: aws.String(t.s_bucket__name),
	}

	_, err = makeBucket(context.TODO(), client, bucket_make_input)

	if err != nil {
		return err
	}

	return nil
}

//---------------------------------------------------------------------------------//

// 5.2 Bucket 삭제
//---------------------------------------------------------------------------------//
func (t *C_S3) Bucket__remove(_s_bukcket__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket__name = _s_bukcket__name

	bucket_remove_input := &s3.DeleteBucketInput{
		Bucket: aws.String(t.s_bucket__name),
	}

	_, err = removeBucket(context.TODO(), client, bucket_remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 5.3 Buekct 권한 변경
//---------------------------------------------------------------------------------//
func (t *C_S3) Bucket__change(_s_bucket__name, _s_bucket__acl string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket__name = _s_bucket__name
	t.s_bucket__acl = _s_bucket__acl

	bucket__change_input := &s3.PutBucketAclInput{
		Bucket: &t.s_bucket__name,
		ACL:    s3types.BucketCannedACL(t.s_bucket__acl),
	}

	_, err = changeBucketAcl(context.TODO(), client, bucket__change_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 5.4 Object 업로드
//---------------------------------------------------------------------------------//
// (sample) _s_ob_name = "C:\\Temp\\test1.txt"
func (t *C_S3) Object__uplaod(_s_bucket__name, _s_object__name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket__name = _s_bucket__name
	t.s_object__name = _s_object__name

	file, err := os.Open(t.s_object__name)
	if err != nil {
		return err
	}

	defer file.Close()

	object_upload_input := &s3.PutObjectInput{
		Bucket: &t.s_bucket__name,
		Key:    &t.s_object__name,
		Body:   file,
	}

	_, err = putFile(context.TODO(), client, object_upload_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 5.5 Object 삭제
//---------------------------------------------------------------------------------//
func (t *C_S3) Object__remove(_s__bucket__name, _s__object__name string) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil
	}

	t.s_bucket__name = _s__bucket__name
	t.s_object__name = _s__object__name

	client := s3.NewFromConfig(cfg)

	object_remove_input := &s3.DeleteObjectInput{
		Bucket: &t.s_bucket__name,
		Key:    &t.s_object__name,
	}

	_, err = deleteItem(context.TODO(), client, object_remove_input)
	if err != nil {
		return err
	}
	return nil
}

//---------------------------------------------------------------------------------//

// 5.6 Object 권한 변경
//---------------------------------------------------------------------------------//
// (sample) _s_ob_acl = "private"
func (t *C_S3) Object__change(_s_bucket__name, _s_object__name, _s_object__acl string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket__name = _s_bucket__name
	t.s_object__name = _s_object__name
	t.s_object__acl = _s_object__acl

	object_change_input := &s3.PutObjectAclInput{
		Bucket: &t.s_bucket__name,
		Key:    &t.s_object__name,
		ACL:    s3types.ObjectCannedACL(t.s_object__acl),
	}
	_, err = changeObjectAcl(context.TODO(), client, object_change_input)
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

// 호스트존 생성 API 인터페이스
type ROUTE53CreateHostedZoneAPI interface {
	CreateHostedZone(ctx context.Context,
		params *route53.CreateHostedZoneInput,
		optFns ...func(*route53.Options)) (*route53.CreateHostedZoneOutput, error)
}

// 호스트존 삭제 API 인터페이스
type ROUTE53DeleteHostedZoneAPI interface {
	DeleteHostedZone(ctx context.Context,
		params *route53.DeleteHostedZoneInput,
		optFns ...func(*route53.Options)) (*route53.DeleteHostedZoneOutput, error)
}

// EC2 Isntacne 생성 API 인터페이스
type EC2CreateInstanceAPI interface {
	RunInstances(ctx context.Context,
		params *ec2.RunInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.RunInstancesOutput, error)

	CreateTags(ctx context.Context,
		params *ec2.CreateTagsInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateTagsOutput, error)
}

// EC2 Instance 삭제 API 인터페이스
type EC2TerminateInstancesAPI interface {
	TerminateInstances(ctx context.Context,
		params *ec2.TerminateInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.TerminateInstancesOutput, error)
}

//EC2 Instance 시작 API 인터페이스
type EC2StartInstancesAPI interface {
	StartInstances(ctx context.Context,
		params *ec2.StartInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

// EC2 Instance 중지 API 인터페이스
type EC2StopInstancesAPI interface {
	StopInstances(ctx context.Context,
		params *ec2.StopInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StopInstancesOutput, error)
}

// IP 할당 및 부여 API 인터페이스
type EC2AllocateAddressAPI interface {
	AllocateAddress(ctx context.Context,
		params *ec2.AllocateAddressInput,
		optFns ...func(*ec2.Options)) (*ec2.AllocateAddressOutput, error)

	AssociateAddress(ctx context.Context,
		params *ec2.AssociateAddressInput,
		optFns ...func(*ec2.Options)) (*ec2.AssociateAddressOutput, error)
}

// EC2 Key-pair 생성 API 인터페이스
type EC2CreateKeyPairAPI interface {
	CreateKeyPair(ctx context.Context,
		params *ec2.CreateKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateKeyPairOutput, error)
}

// EC2 Key-pair 삭제 API 인터페이스
type EC2DeleteKeyPairAPI interface {
	DeleteKeyPair(ctx context.Context,
		params *ec2.DeleteKeyPairInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteKeyPairOutput, error)
}

// EC2 Security Group API 인터페이스
type EC2CreateSecurityGroupAPI interface {
	CreateSecurityGroup(ctx context.Context,
		params *ec2.CreateSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.CreateSecurityGroupOutput, error)
}

// EC2 Security Group API 인터페이스
type EC2DeleteSecurityGroupAPI interface {
	DeleteSecurityGroup(ctx context.Context,
		params *ec2.DeleteSecurityGroupInput,
		optFns ...func(*ec2.Options)) (*ec2.DeleteSecurityGroupOutput, error)
}

// EC2 Security Group Inboubd 규칙 추가 API 인터페이스
type EC2AuthorizeSecurityGroupIngressAPI interface {
	AuthorizeSecurityGroupIngress(ctx context.Context,
		params *ec2.AuthorizeSecurityGroupIngressInput,
		optFns ...func(*ec2.Options)) (*ec2.AuthorizeSecurityGroupIngressOutput, error)
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
func changeBucketAcl(c context.Context, api S3PutBucketAclAPI, input *s3.PutBucketAclInput) (*s3.PutBucketAclOutput, error) {
	return api.PutBucketAcl(c, input)
}

// API Object 업로드 함수
func putFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

// API Object 삭제 함수
func deleteItem(c context.Context, api S3DeleteObjectAPI, input *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {
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

// 호스트존 생성 함수
func makeHostzone(c context.Context, api ROUTE53CreateHostedZoneAPI, input *route53.CreateHostedZoneInput) (*route53.CreateHostedZoneOutput, error) {
	return api.CreateHostedZone(c, input)
}

// 호스트존 삭제 함수
func removeHostzone(c context.Context, api ROUTE53DeleteHostedZoneAPI, input *route53.DeleteHostedZoneInput) (*route53.DeleteHostedZoneOutput, error) {
	return api.DeleteHostedZone(c, input)
}

// EC2 Instance 생성 함수
func makeIns(c context.Context, api EC2CreateInstanceAPI, input *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {
	return api.RunInstances(c, input)
}

// EC2 Instance tag 생성 함수
func makeTags(c context.Context, api EC2CreateInstanceAPI, input *ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return api.CreateTags(c, input)
}

// EC2 Instance 삭제 함수
func removeIns(c context.Context, api EC2TerminateInstancesAPI, input *ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return api.TerminateInstances(c, input)
}

// EC2 Instance 시작 함수
func startIns(c context.Context, api EC2StartInstancesAPI, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	return api.StartInstances(c, input)
}

// EC2 Instance 중지 함수
func stopIns(c context.Context, api EC2StopInstancesAPI, input *ec2.StopInstancesInput) (*ec2.StopInstancesOutput, error) {
	return api.StopInstances(c, input)
}

// EC2 Elastic IP 할당
func allocateIP(c context.Context, api EC2AllocateAddressAPI, input *ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
	return api.AllocateAddress(c, input)
}

// EC2 Elastic IP 부여
func associateIP(c context.Context, api EC2AllocateAddressAPI, ipinput *ec2.AssociateAddressInput) (*ec2.AssociateAddressOutput, error) {
	return api.AssociateAddress(c, ipinput)
}

// EC2 Key-pair 생성
func makeKeypair(c context.Context, api EC2CreateKeyPairAPI, input *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {
	return api.CreateKeyPair(c, input)
}

// EC2 Key-pair 삭제
func removeKeypair(c context.Context, api EC2DeleteKeyPairAPI, input *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {
	return api.DeleteKeyPair(c, input)
}

// EC2 Security Group 생성
func createsg(c context.Context, api EC2CreateSecurityGroupAPI, input *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {
	return api.CreateSecurityGroup(c, input)
}

// EC2 Seucurity Group 삭제
func removesg(c context.Context, api EC2DeleteSecurityGroupAPI, input *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {
	return api.DeleteSecurityGroup(c, input)
}

// EC2 Security Group Inbound 추가
func change_sg__insbound(c context.Context, api EC2AuthorizeSecurityGroupIngressAPI, input *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {
	return api.AuthorizeSecurityGroupIngress(c, input)
}

// // EC2 Seucurity Group Inbound 규칙 추가
// func change_sg__insbound(c context.Context, api EC2UpdateSecurityGroupRuleDescriptionsIngressAPI, input *ec2.UpdateSecurityGroupRuleDescriptionsIngressInput) (*ec2.UpdateSecurityGroupRuleDescriptionsIngressOutput, error) {
// 	return api.UpdateSecurityGroupRuleDescriptionsIngress(c, input)
// }

// <------------------------------------------------------------------------->

func Test_tmp2(_t *testing.T) {
	define := C_Network{}
	// define.Hostzone__Create("devtoolstest3.com", "2021-01-25-15:39")
	// define.Hostzone__remove("Z01323092XMDLOD2S2FPX")
	// define.Instance__create("ami-014009fa4a1467d53", "t2.micro", "devtoolstest-group", "testkey", "Name", "test2server", 1)
	// define.Instance__stop("i-00a50ebab0cb0da66")
	// define.Instance__run("i-00a50ebab0cb0da66")
	// define.Instance__remove("i-00a50ebab0cb0da66")
	// define.Instance__associate_ip("i-00a50ebab0cb0da66")
	// define.Keypair__create("devtoolstestkey1211")
	// define.Keypair__remove("devtoolstestkey1211")
	// define.Securitygroup__create("devtoolstestgroup1111", "sdk go test group")
	// define.Securitygroup__remove("devtoolstestgroup1111")
	define.Securitygroup_Inbound("sg-0f71d13740864b02b", "0.0.0.0/32", "tcp", 100, 200)
}
