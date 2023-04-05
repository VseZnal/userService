package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"userService/services/api-gateway/config"
	"userService/services/api-gateway/cors"
	"userService/services/api-gateway/middleware"
	user_service "userService/services/api-gateway/proto/user-service"
)

func main() {
	conf := config.PublicConfig()

	userServiceHost := conf.HostUser
	userServicePort := conf.PortUser

	gatewayHost := conf.HostGateway
	gatewayPort := conf.PortGateway

	userServiceAddress := userServiceHost + ":" + userServicePort

	proxyAddr := gatewayHost + ":" + gatewayPort

	GatewayStart(proxyAddr, userServiceAddress)
}

func GatewayStart(proxyAddr, userAddress string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{grpc.WithUnaryInterceptor(middleware.AccessLogInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := user_service.RegisterUserServiceHandlerFromEndpoint(ctx, mux, userAddress, opts)
	if err != nil {
		log.Fatalln("Failed to connect to User service", err)
	}

	gwServer := &http.Server{
		Addr:    proxyAddr,
		Handler: cors.Cors(mux),
	}

	fmt.Println("starting gateway server at " + proxyAddr)
	log.Fatalln(gwServer.ListenAndServe())

}
