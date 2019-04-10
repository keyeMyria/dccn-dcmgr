package main

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/common"

	//"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	//"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
	"log"
	"time"

	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

	//common_proto "github.com/Ankr-network/dccn-common/protos/common"
//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)

//var addr = "localhost:50051"
var addr = "client-stage.dccn.ankr.network:50051"

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
	userClient := usermgr.NewUserMgrClient(conn)



	req := &usermgr.LoginRequest{}
	req.Email = `yousong.zhang@gmail.com`
	req.Password = "zddzys123"
	//if _, err := userClient.Register(context.Background(), user); err != nil {
	//	log.Fatal(err.Error())
	//} else {
	//	log.Println("Register Ok")
	//}

	var token string
	//var userId string
	if rsp, err := userClient.Login(context.TODO(), req); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("login Success: %s\n", rsp.AuthenticationResult.AccessToken)
		token = rsp.AuthenticationResult.AccessToken
	//	userId = rsp.UserId
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	//task := apiCommon.MockTasks()[0]
	//log.Println("Test CreateTask")
	//if rsp, err := taskClient.CreateTask(tokenContext, &taskmgr.CreateTaskRequest{UserId: userId, Task: &task}); err != nil {
	//	log.Fatal(err.Error())
	//} else {
	//	log.Println(*rsp)
	//}

	// var userTasks []*common_proto.Task
	if rsp, err := dcClient.DataCenterList(tokenContext, &common_proto.Empty{}); err != nil {
		log.Fatal(err.Error())
	} else {
		for i := 0; i < len(rsp.DcList); i++ {
			dc := rsp.DcList[i]
			log.Printf("DataCenterList metrics %s name %s  status %s \n", dc.Name, dc.Name, dc.Status)

		}

	}

	log.Println("END")
}
