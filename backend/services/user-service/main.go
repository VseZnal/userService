package main

import (
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"userService/services/user-service/config"
	user_service "userService/services/user-service/proto"
	pb "userService/services/user-service/proto/user-service"
)

func main() {
	conf := config.UserConfig()

	userServiceHost := conf.HostUser
	userServicePort := conf.PortUser

	userServiceAddress := fmt.Sprintf("%s:%s", userServiceHost, userServicePort)

	err := user_service.Init()
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	lis, err := net.Listen("tcp", userServiceAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

	pb.RegisterUserServiceServer(
		server,
		&user_service.Server{},
	)

	log.Printf("server listening at %v", lis.Addr())

	err = server.Serve(lis)

	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
