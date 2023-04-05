package middleware

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
	jwt_user "userService/services/user-service/jwt"
)

func AccessLogInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	md, _ := metadata.FromOutgoingContext(ctx)
	start := time.Now()

	private := map[string]bool{
		"/pb.UserService/RandomPrivateMethod": false,
	}

	var traceId string
	for k, v := range private {
		if method == k && v != true {
			if len(md["authorization"]) > 0 {
				tokenString := md["authorization"][0]
				if tokenString != "" {
					err, _ := jwt_user.CheckJWTToken(tokenString)
					if err != nil {
						return errors.New("your token is invalid")
					}
				} else {
					return errors.New("error authorization")
				}
			} else {
				return errors.New("error authorization")
			}
		}
	}
	traceId = fmt.Sprintf("%d", time.Now().UTC().UnixNano())

	callContext := context.Background()
	mdOut := metadata.Pairs(
		"trace-id", traceId,
	)
	callContext = metadata.NewOutgoingContext(callContext, mdOut)

	err := invoker(callContext, method, req, reply, cc, opts...)

	msg := fmt.Sprintf("Call:%v, traceId: %v, time: %v", method, traceId, time.Since(start))
	log.Println(msg)

	return err
}
