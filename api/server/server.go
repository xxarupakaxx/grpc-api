package main

import (
	"fmt"
	"github.com/xxarupakaxx/grpc-api/api/gen/api"
	"google.golang.org/grpc"
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

	server := grpc.NewServer()
  	api.RegisterPancakeBakerServiceServer(server,handler.NewBakerHandler())
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
