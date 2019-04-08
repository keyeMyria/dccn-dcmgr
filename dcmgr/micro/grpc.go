package micro2

import (
	"log"
	"net"

	//"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)


type GRPCService struct {
	 listener net.Listener
	 s         *grpc.Server
}

func NewGRPCService()  GRPCService{
	s := GRPCService{}
	lis, err := net.Listen("tcp", config.Listen)
	s.listener = lis
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.s = grpc.NewServer()
	return s

}

//*
//pb.RegisterGreeterServer(s, &server{})


func (s *GRPCService)GetServer() *grpc.Server {
	return s.s
}



func (s *GRPCService) Start() {
	// Register reflection service on gRPC server.
 	reflection.Register(s.s)
	WriteLog("GRPCService Start @ " + config.Listen)
	if err := s.s.Serve(s.listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}


