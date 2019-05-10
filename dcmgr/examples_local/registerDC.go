package main

import (
	"context"
	"fmt"
	"log"
	"time"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var addr = "localhost:50051"

//var addr = "client-dev.dccn.ankr.com:50051"

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

	md := metadata.New(map[string]string{
		"token": "",
	})

	//log.Printf("get access_token after login %s \n", access_token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// var userTasks []*common_proto.Task
	if rsp, err := dcClient.RegisterDataCenter(tokenContext, &dcmgr.RegisterDataCenterRequest{UserId:"xxxx123", ClusterName:"demo-cluster"}); err != nil {
		log.Fatal(err.Error())
	} else {

		fmt.Printf("RegisterDataCenter response %+v", rsp)

	}
	//}
}
