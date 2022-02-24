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
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// 구조체 정의
// -----------------------------------------------------------------------------//

type AWSCTRL struct {
	TD__S3  C__S3
	TD__RDS C__RDS
}

// RDS DB 구조체
type C__RDS struct {
	s_db_name string
	s_id      string
	s_pw      string
}

// S3 구조체
type C__S3 struct {
	s_bucket_name string
	s_bucket_acl  string
	s_object_name string
	s_object_acl  string
}

// var err_rdb__config__null = errors.New("RDB - Config - 공백")
// var err_s3__config__null = errors.New("S3 - config - 공백")

// DB 생성, 삭제 속성 값
const (
	Def_db__name    string = "devtoolstest12312"
	Def_db__id      string = "master"
	Def_db__pwd     string = "test112233"
	Def_bucket_name string = "devtoolsbucket11"
	Def_object_name string = "C:\\work_space\\devops\\awsctrl\\src\\test1"
)

// ------------------------------------------------------------------------------//
// 조건 통과 시 명령 함수 실행
func Test_all(_t *testing.T) {
	var err error
	rds := &C__RDS{}
	s3 := &C__S3{}

	// DB 생성
	err = rds.DB__create(Def_db__name, Def_db__id, Def_db__pwd)
	if err != nil {
		_t.Error(err)
		return
	}

	// DB 삭제
	err = rds.DB__remove(Def_db__name)
	if err != nil {
		_t.Error(err)
		return
	}

	// DB 조회
	err = rds.DB__getinfo(Def_db__name)
	if err != nil {
		_t.Error(err)
		return
	}

	// Bucket 생성
	err = s3.Bucket__create(Def_bucket_name)
	if err != nil {
		_t.Error(err)
		return
	}

	// Bucket 삭제
	err = s3.Bucket__remove(Def_bucket_name)
	if err != nil {
		_t.Error(err)
		return
	}

	// Object 업로드
	err = s3.Object__uplaod(Def_bucket_name, Def_object_name)
	if err != nil {
		_t.Error(err)
		return
	}

	// Object 삭제
	err = s3.Object__remove(Def_bucket_name, Def_object_name)
	if err != nil {
		_t.Error(err)
		return
	}
}

// ------------------------------------------------------------------------------ //

// 실행 함수 모음
// ------------------------------------------------------------------------------ //

// DB 생성 함수
func (t *C__RDS) DB__create(_s_db_name, _s_id, _s_pw string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	t.s_db_name = _s_db_name
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

// DB 삭제 함수
func (t *C__RDS) DB__remove(_s_db_name string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	t.s_db_name = _s_db_name

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

// DB 조회 함수
func (t *C__RDS) DB__getinfo(_s_db_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := rds.NewFromConfig(cfg)

	t.s_db_name = _s_db_name

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

// Bucket 생성
func (t *C__S3) Bucket__create(_s_bucket_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket_name = _s_bucket_name

	bucket__make_input := &s3.CreateBucketInput{
		Bucket: aws.String(t.s_bucket_name),
	}

	_, err = makeBucket(context.TODO(), client, bucket__make_input)

	if err != nil {
		return err
	}

	return nil
}

// Bucket 삭제
func (t *C__S3) Bucket__remove(_s_bukcket_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket_name = _s_bukcket_name

	bucket__remove_input := &s3.DeleteBucketInput{
		Bucket: aws.String(t.s_bucket_name),
	}

	_, err = removeBucket(context.TODO(), client, bucket__remove_input)
	if err != nil {
		return err
	}
	return nil
}

// Buekct 권한 변경
func (t *C__S3) Bucket__change(_s_bucket_name, _s_bucket_acl string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket_name = _s_bucket_name
	t.s_bucket_acl = _s_bucket_acl

	bucket__change_input := &s3.PutBucketAclInput{
		Bucket: &t.s_bucket_name,
		ACL:    s3types.BucketCannedACL(t.s_bucket_acl),
	}

	_, err = ChangeBucketAcl(context.TODO(), client, bucket__change_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 업로드
// (sample) _s_ob_name = "C:\\Temp\\test1.txt"
func (t *C__S3) Object__uplaod(_s_bucket_name, _s_object_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket_name = _s_bucket_name
	t.s_object_name = _s_object_name

	file, err := os.Open(t.s_object_name)
	if err != nil {
		return err
	}

	defer file.Close()

	object__upload_input := &s3.PutObjectInput{
		Bucket: &t.s_bucket_name,
		Key:    &t.s_object_name,
		Body:   file,
	}

	_, err = putFile(context.TODO(), client, object__upload_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 삭제
func (t *C__S3) Object__remove(_s_bucket_name, _s_object_name string) error {

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return nil
	}

	t.s_bucket_name = _s_bucket_name
	t.s_object_name = _s_object_name

	client := s3.NewFromConfig(cfg)

	object__remove_input := &s3.DeleteObjectInput{
		Bucket: &t.s_bucket_name,
		Key:    &t.s_object_name,
	}

	_, err = DeleteItem(context.TODO(), client, object__remove_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 권한 변경
// (sample) _s_ob_acl = "private"
func (t *C__S3) Object__change(_s_bucket_name, _s_object_name, _s_object_acl string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	t.s_bucket_name = _s_bucket_name
	t.s_object_name = _s_object_name
	t.s_object_acl = _s_object_acl

	object__change_input := &s3.PutObjectAclInput{
		Bucket: &t.s_bucket_name,
		Key:    &t.s_object_name,
		ACL:    s3types.ObjectCannedACL(t.s_object_acl),
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

// <------------------------------------------------------------------------->

func Test_tmp2(_t *testing.T) {
	define := C__S3{}
	define.Bucket__change("testdevtools1231", "public-read-write")
}
