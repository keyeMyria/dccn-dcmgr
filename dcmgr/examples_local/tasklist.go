package main

//import (
//	"context"
//	//"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
//	"log"
//	"time"
//
//	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/metadata"
//
//	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
//
////	common_proto "github.com/Ankr-network/dccn-common/protos/common"
////	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
//)
//
//var addr = "127.0.0.1:50052"
////var addr = "client-dev.dccn.ankr.network:50051"
//
//func main() {
//
//	log.SetFlags(log.LstdFlags | log.Lshortfile)
//	conn, err := grpc.Dial(addr, grpc.WithInsecure())
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	defer func(conn *grpc.ClientConn) {
//		if err := conn.Close(); err != nil {
//			log.Println(err.Error())
//		}
//	}(conn)
//
//	taskClient := taskmgr.NewTaskMgrClient(conn)
//	userClient := usermgr.NewUserMgrClient(conn)
//
//	req := &usermgr.LoginRequest{}
//	req.Email = "yousong1@gmail.com"
//	req.Password = "111111c"
//
//	//var userId string
//	if rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: req.Email, Password: req.Password}); err != nil {
//		log.Fatal(err.Error())
//	} else {
//		log.Printf("response %+v \n", rsp)
//		//log.Printf("login Success: id : %s name : %s , email %s  refresh_token : %s  access_token %s \n", rsp.User.Id, rsp.User.Attributes.Name, rsp.User.Email ,rsp.AuthenticationResult.RefreshToken, rsp.AuthenticationResult.AccessToken)
//		//token = rsp.Token
//		//userId = rsp.UserId
//		//_ := rsp.AuthenticationResult.RefreshToken
//		access_token := rsp.AuthenticationResult.AccessToken
//
//		//log.Printf("refresh_token  %s  access_token %s", refresh_token, access_token)
//
//		md := metadata.New(map[string]string{
//			"token": access_token,
//		})
//
//		log.Printf("get access_token after login %s \n", access_token)
//		ctx := metadata.NewOutgoingContext(context.Background(), md)
//
//		tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
//		defer cancel()
//
//		filter := taskmgr.TaskFilter{TaskId: ""}
//
//		if rsp, err := taskClient.TaskList(tokenContext, &taskmgr.TaskListRequest{TaskFilter: &filter}); err != nil {
//
//			//log.Println("detail create %+v " + rsp)
//			log.Fatal(err)
//		} else {
//
//			for i := 0; i < len(rsp.Tasks); i++ {
//				log.Printf("task list %+v   %+v \n", rsp.Tasks[i], rsp.Tasks[i].Attributes)
//			}
//		}
//
//	}
//}
