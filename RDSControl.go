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

//------------------------------------------------------------------------------//
// 외부 노출 버튼(생성, 삭제)

func (t *C_RDS) Button__DBcreate() error {
	return t.control(TD_Control__Create)
}

func (t *C_RDS) Button__DBdelete() error {
	return t.control(TD_Control__Create)
}

//------------------------------------------------------------------------------//

type TD_Control string

const (
	TD_Control__Create TD_Control = "Create"
	TD_Control__Delete TD_Control = "Delete"
)

// 에러 정의
const (
	def__s_dbname_null string = "DB 이름 공백"
	def__s_dbid_null   string = "DB ID 공백"
	def__s_dbpwd_null  string = "DB PASSWD 공백"
)

// DB 속성 구조체
type C_RDS struct {
	TD_DB__Name    string
	TD_DB__UserID  string
	TD_DB__UserPWD string
}

// DB INPUT(string) 데이터 "" 입력 시 에러
func (t *C_RDS) control(_td_ctl TD_Control) error {
	switch _td_ctl {
	case TD_Control__Create:
		if t.TD_DB__Name == "" {
			return errors.New(def__s_dbname_null)
		} else if t.TD_DB__UserID == "" {
			return errors.New(def__s_dbid_null)
		} else if t.TD_DB__UserPWD == "" {
			return errors.New(def__s_dbid_null)
		}
		fmt.Println("Create DB")
	case TD_Control__Delete:
		if t.TD_DB__Name == "" {
			return errors.New(def__s_dbname_null)
		}
		fmt.Println("Delete DB")
	default:
		return errors.New("Error")
	}
	return nil
}

// AWS CreateDB API Define
type RDSCreateDBClusterAPI interface {
	CreateDBCluster(ctx context.Context,
		params *rds.CreateDBClusterInput,
		optFns ...func(*rds.Options)) (*rds.CreateDBClusterOutput, error)
}

// DB 생성 함수
func makeDBCluster(c context.Context, api RDSCreateDBClusterAPI, makeinput *rds.CreateDBClusterInput) (*rds.CreateDBClusterOutput, error) {
	return api.CreateDBCluster(c, makeinput)
}

// AWS DeleteDB API Define
type RDSDeleteDBClusterAPI interface {
	DeleteDBCluster(ctx context.Context,
		params *rds.DeleteDBClusterInput,
		optFns ...func(*rds.Options)) (*rds.DeleteDBClusterOutput, error)
}

// DB 삭제 함수
func removeDBCluster(c context.Context, api RDSDeleteDBClusterAPI, deleteinput *rds.DeleteDBClusterInput) (*rds.DeleteDBClusterOutput, error) {
	return api.DeleteDBCluster(c, deleteinput)
}

// 내부 함수
func Test__rdstest(_t *testing.T) {

	// AWS Config 파일 로드 및 접속 세션 구성
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-2"))
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := rds.NewFromConfig(cfg)

	// // 구조체 호출
	Property := &C_RDS{}
	Property.TD_DB__Name = "test111"
	Property.TD_DB__UserID = "master"
	Property.TD_DB__UserPWD = "devttools1234"

	// DB 생성 INPUT 값 설정
	makeinput := &rds.CreateDBClusterInput{
		DBClusterIdentifier: &Property.TD_DB__Name,
		Engine:              aws.String("aurora"),
		EngineMode:          aws.String("serverless"),
		MasterUsername:      &Property.TD_DB__UserID,
		MasterUserPassword:  &Property.TD_DB__UserPWD,
		ScalingConfiguration: &types.ScalingConfiguration{
			AutoPause:             aws.Bool(true),
			MinCapacity:           aws.Int32(1),
			MaxCapacity:           aws.Int32(32),
			SecondsBeforeTimeout:  aws.Int32(300),
			SecondsUntilAutoPause: aws.Int32(300),
			TimeoutAction:         aws.String("ForceApplyCapacityChange"),
		},
	}

	// DB 삭제 INPUT 값 설정
	deleteinput := &rds.DeleteDBClusterInput{
		DBClusterIdentifier: aws.String("devtoolstest113"),
		SkipFinalSnapshot:   *aws.Bool(true),
	}

	//DB 생성 (makeDBCluster 함수 실행)
	_, err = makeDBCluster(context.TODO(), client, makeinput)
	if err != nil {
		fmt.Println("Create DB Cluster Error")
		fmt.Println(err)
		return
	}

	// DB 삭제 (removeDBCluster 함수 실행)
	_, err = removeDBCluster(context.TODO(), client, deleteinput)
	if err != nil {
		fmt.Println("Create DB Cluster Error")
		fmt.Println(err)
		return
	}
}
