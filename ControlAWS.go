package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// 구조체 정의
// -----------------------------------------------------------------------------//

// RDS DB 구조체
type TD__RDS struct {
	_s_db_name string
	_s_id      string
	_s_pw      string
}

// S3 구조체
type TD__S3 struct {
	_s_bk_name string
	_s_ob_name string
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
	rds := &TD__RDS{}
	s3 := &TD__S3{}

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

	// Create Bucket func
	err = s3.Bucket__craete(Def_bucket_name)
	if err != nil {
		_t.Error(err)
		return
	}

	// Delete Bucket func
	err = s3.Bucket__remove(Def_bucket_name)
	if err != nil {
		_t.Error(err)
		return
	}

}

// ------------------------------------------------------------------------------ //

// 실행 함수 모음
// ------------------------------------------------------------------------------ //

// DB 생성 함수
func (t *TD__RDS) DB__create(_s_db_name, _s_id, _s_pw string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := rds.NewFromConfig(cfg)

	// DB 생성 INPUT 값 설정
	db_make_input := &rds.CreateDBClusterInput{
		DBClusterIdentifier: aws.String(t._s_db_name),
		Engine:              aws.String("aurora"),
		EngineMode:          aws.String("serverless"),
		MasterUsername:      aws.String(t._s_id),
		MasterUserPassword:  aws.String(t._s_pw),
		ScalingConfiguration: &types.ScalingConfiguration{
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
func (t *TD__RDS) DB__remove(_s_db_name string) error {
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := rds.NewFromConfig(cfg)

	// DB 삭제 INPUT 값 설정
	db_remove_input := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String(t._s_db_name),
		SkipFinalSnapshot:   *aws.Bool(true),
	}

	// DB 삭제 (removeDBCluster 함수 실행)
	_, err = removeDBCluster(context.TODO(), client, db_remove_input)
	if err != nil {
		return fmt.Errorf("Create DB Cluster Error - %v", err)
	}
	return nil
}

// DB 조회 함수
func (t *TD__RDS) DB__getinfo(_s_db_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := rds.NewFromConfig(cfg)

	// DB 조회 INPUT 값 설정
	db__find_input := &rds.DescribeDBClustersInput{
		DBClusterIdentifier: aws.String(t._s_db_name),
	}

	// DB 조회
	result, err := FindDBCluster(context.TODO(), client, db__find_input)
	if err != nil {
		return fmt.Errorf("Find DB Cluster Error - %v", err)
	}
	fmt.Println("all done")
	output := result.DBClusters

	fmt.Printf("%v", output)
	return nil
}

// Bucket 생성
func (t *TD__S3) Bucket__craete(_s_bk_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	bucket__make_input := &s3.CreateBucketInput{
		Bucket: aws.String(_s_bk_name),
	}

	_, err = MakeBucket(context.TODO(), client, bucket__make_input)

	if err != nil {
		return err
	}

	return nil
}

// Bucket 삭제
func (t *TD__S3) Bucket__remove(_s_bk_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)

	bucket__remove_input := &s3.DeleteBucketInput{
		Bucket: aws.String(t._s_bk_name),
	}

	_, err = RemoveBucket(context.TODO(), client, bucket__remove_input)
	if err != nil {
		return err
	}
	return nil
}

// Object 업로드
func (t *TD__S3) Object__uplaod(_s_ob_name string) error {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		return fmt.Errorf("configuration error - %v", err)
	}

	client := s3.NewFromConfig(cfg)
	file, err := os.Open(t._s_ob_name)

	if err != nil {
		fmt.Println("Unable to open file " + t._s_ob_name)
	}

	defer file.Close()

	object__upload_input := &s3.PutObjectInput{
		Bucket: aws.String("devtoolstest113"),
		Key:    &t._s_ob_name,
		Body:   file,
	}

	_, err = PutFile(context.TODO(), client, object__upload_input)
	if err != nil {
		return err
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

// Object 업로드 API 정의
type S3PutObjectAPI interface {
	PutObject(ctx context.Context,
		params *s3.PutObjectInput,
		optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
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
func FindDBCluster(c context.Context, api RDSDescribeDBClustersAPI, input *rds.DescribeDBClustersInput) (*rds.DescribeDBClustersOutput, error) {
	return api.DescribeDBClusters(c, input)
}

// API Bucket 생성 함수
func MakeBucket(c context.Context, api S3CreateBucketAPI, input *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {
	return api.CreateBucket(c, input)
}

// API Bucket 삭제 함수
func RemoveBucket(c context.Context, api S3DeleteBucketAPI, input *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {
	return api.DeleteBucket(c, input)
}

// API Object 업로드 함수
func PutFile(c context.Context, api S3PutObjectAPI, input *s3.PutObjectInput) (*s3.PutObjectOutput, error) {
	return api.PutObject(c, input)
}

// <------------------------------------------------------------------------->

func Test_tmp2(_t *testing.T) {

	define := TD__S3{}
	define.Bucket__craete("devtoolstest12111")
}
