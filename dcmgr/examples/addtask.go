package main

import (
	"context"
//	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"

	"log"
	"time"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	//usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

//	common_proto "github.com/Ankr-network/dccn-common/protos/common"
//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)
var addr = "appmgr:50051"



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

	taskClient := taskmgr.NewAppMgrClient(conn)

		md := metadata.New(map[string]string{
			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTUwMTUzMTgsImp0aSI6IjQ4NTQ5YjQxLWUzNjYtNGIxMi05NTc3LTU0M2Y5NTE5Y2JlZiIsImlzcyI6ImFua3IubmV0d29yayJ9.A0p3KyxIKZHAZb_buPgadKj3d40Rlw_hSpsFBrNLjuw",
		})

		ctx := metadata.NewOutgoingContext(context.Background(), md)
//
		tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
//
		app := common_proto.App{}
		app.Name = "task"
		//task.Type = common_proto.TaskType_DEPLOYMENT
		app.Attributes = &common_proto.AppAttributes{}
		//task.Attributes.R = 1
		//t := common_proto.APP_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image:"nginx:1.12"}}
		//task.TypeData = &t



		if rsp, err := taskClient.CreateApp(tokenContext, &taskmgr.CreateAppRequest{App: &app}); err != nil {

			//log.Println("detail create %+v " + rsp)
			log.Fatal(err)
		} else {
			log.Println("create task successfully : taskid   " + rsp.AppId)
		}

//	}

}
