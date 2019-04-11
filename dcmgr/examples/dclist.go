package main

import (
	"context"
	"fmt"
	"github.com/Ankr-network/dccn-common/protos/common"
	//"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	"log"
	"time"

	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	//	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)
var addr = "localhost:50056"
//var addr = "client-dev.dccn.ankr.network:50051"

//var addr = "afcac29ea274711e99cb106bbae7419f-1982485008.us-west-1.elb.amazonaws.com:50051"

//func parseError(err string) string{
//
//}

func main() {

	log.SetFlags(log.LstdFlags | log.Llongfile)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)

	dcClient := dcmgr.NewDCAPIClient(conn)
	//userClient := usermgr.NewUserMgrClient(conn)
	//
	//req := &usermgr.LoginRequest{}
	//req.Email = "yousong1@gmail.com"
    //req.Password = "111111c"
	//
	//
	//
	////var userId string
	//if rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: req.Email, Password: req.Password}); err != nil {
	//	if err == ankr_default.ErrPasswordError {
	//		log.Printf("password error  %s", err.Error())
	//	}
	//
	//
	//	a := err.Error()
	//
	//	log.Printf(">>>>>%d  %s <<<<", len(a), a)
	//} else {
	//	log.Printf("response %+v \n", rsp)
	//	//log.Printf("login Success: id : %s name : %s , email %s  refresh_token : %s  access_token %s \n", rsp.User.Id, rsp.User.Attributes.Name, rsp.User.Email ,rsp.AuthenticationResult.RefreshToken, rsp.AuthenticationResult.AccessToken)
	//	//token = rsp.Token
	//	//userId = rsp.UserId
	//	refresh_token := rsp.AuthenticationResult.RefreshToken
	//	access_token := rsp.AuthenticationResult.AccessToken

	//	log.Printf("refresh_token  %s  access_token %s", refresh_token, access_token)

		md := metadata.New(map[string]string{
			"token": "",
		})

		//log.Printf("get access_token after login %s \n", access_token)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// var userTasks []*common_proto.Task
		if rsp, err := dcClient.DataCenterList(tokenContext, &common_proto.Empty{}); err != nil {
			log.Fatal(err.Error())
		} else {

			fmt.Print("this is a new test")

			for i := 0; i < len(rsp.DcList); i++ {
				d := rsp.DcList[i]
				fmt.Printf("task list id %s name %s lat %s lng %s cournty %s \n", d.Id, d.Name, d.GeoLocation.Lng, d.GeoLocation.Lng, d.GeoLocation.Country)
			}

		}
	//}
}
