package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/xxarupakaxx/grpc-api/api/gen/api"
	"github.com/xxarupakaxx/grpc-api/api/handler"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

func main() {
	port := 50051
	lis,err := net.Listen("tcp",fmt.Sprintf(":%d",port))
	if err != nil {
		log.Fatalf("failed to listen :%v",err)
	}
	zapLogger ,err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	grpc_zap.ReplaceGrpcLogger(zapLogger)


	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_zap.UnaryServerInterceptor(zapLogger),
				),
			),
		)

  	api.RegisterPancakeBakerServiceServer(server,handler.NewBakerHandler())
	api.RegisterImageUploadServiceServer(server,handler.NewImageUploadHandler())

	reflection.Register(server)

	go func() {
		log.Printf("start gRPC server port :%v",port)
		err = server.Serve(lis)
		if err != nil {
			log.Fatalf("failed to server: %v",err)
			return
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit,os.Interrupt)

	<- quit
	log.Println("stopping gRPC server...")
	server.GracefulStop()
}

func auth(ctx context.Context) (context.Context, error) {
	token ,err := grpc_auth.AuthFromMD(ctx,"bearer")
	if err != nil {
		return nil, err
	}

	if token != "hi/mi/tsu" {
		return nil,grpc.Errorf(codes.Unauthenticated,"invalid bearer token")

	}

	return context.WithValue(ctx,"UserName","God"),nil
}