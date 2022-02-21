package main

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

type TD_Control string

const (
	TD_Control__Create TD_Control = "devtoolstest1234"
	TD_Control__Delete TD_Control = "devtoolstest123111"
)

// 에러 정의
const (
	def__s_already_delete  string = "삭제 대상 DB 없음"
	def__s_already_create  string = "DB 이름 중복"
	def__s_name_type_error string = "DB Type 오류"
)

type C_RDS struct {
	td_rds_control TD_Control
}

func (t *C_RDS) Init() error {
	t.td_rds_control = TD_Control__Create

	return nil
}

func (t *C_RDS) control(_td_control TD_Control) error {
	switch _td_control {
	case TD_Control__Create:
		if TD_Control__Create == "" {
			return errors.New(def__s_name_type_error)
		}
		fmt.Println("DB 생성")
		// Create DB

	case TD_Control__Delete:
		if TD_Control__Delete == "" {
			return errors.New(def__s_name_type_error)
		}
		fmt.Println("DB 삭제")
		// Delete DB
	}

	t.td_rds_control = _td_control
	return nil
}

// AWS CreateDB API Define
type RDSCreateDBClusterAPI interface {
	CreateDBCluster(ctx context.Context,
		params *rds.CreateDBClusterInput,
		optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
}

func makeDBCluster(c context.Context, api RDSCreateDBClusterAPI, makeinput *rds.CreateDBClusterInput) (*rds.CreateDBClusterOutput, error) {
	return api.CreateDBCluster(c, makeinput)
}

// AWS DeleteDB API Define
type RDSDeleteDBClusterAPI interface {
	DeleteDBCluster(ctx context.Context,
		params *rds.DeleteDBClusterInput,
		optFns ...func(*rds.Options)) (*rds.DeleteDBClusterOutput, error)
}

func removeDBCluster(c context.Context, api RDSDeleteDBClusterAPI, deleteinput *rds.DeleteDBClusterInput) (*rds.DeleteDBClusterOutput, error) {
	return api.DeleteDBCluster(c, deleteinput)
}

func Test__rdstest(_t *testing.T) {
	var err error
	aws_control := &C_RDS{}
	err = aws_control.Init()
	if err != nil {
		fmt.Printf("err - %v", err)
		return
	}
	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := rds.NewFromConfig(cfg)

	// DB 생성 INPUT 값 설정
	makeinput := &rds.CreateDBClusterInput{
		DBClusterIdentifier: aws.String("devtestool1231"),
		Engine:              aws.String("aurora"),
		EngineMode:          aws.String("serverless"),
		MasterUsername:      aws.String("master"),
		MasterUserPassword:  aws.String("devtoolstest123"),
		ScalingConfiguration: &types.ScalingConfiguration{
			AutoPause:             aws.Bool(true),
			MinCapacity:           aws.Int32(1),
			MaxCapacity:           aws.Int32(32),
			SecondsBeforeTimeout:  aws.Int32(300),
			SecondsUntilAutoPause: aws.Int32(300),
			TimeoutAction:         aws.String("ForceApplyCapacityChange"),
		},
	}
	//DB 생성 함수
	_, err = makeDBCluster(context.TODO(), client, makeinput)
	if err != nil {
		fmt.Println("Create DB Cluster Error")
		fmt.Println(err)
		return
	}

	// DB 삭제 INPUT 값 설정
	deleteinput := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String("devtoolstest113"),
		SkipFinalSnapshot:   *aws.Bool(true),
	}

	// DB 삭제 함수
	_, err = removeDBCluster(context.TODO(), client, deleteinput)
	if err != nil {
		fmt.Println("Create DB Cluster Error")
		fmt.Println(err)
		return
	}

}

func (t *C_RDS) makeDBCluster() error {
	return t.control(TD_Control__Create)
}

func (t *C_RDS) deleteDBCluster() error {
	return t.control(TD_Control__Delete)
}
