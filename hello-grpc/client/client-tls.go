package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"

	pb "github.com/Washington-Ksea/gRPC-example/hello-grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func main() {
	addr := "localhost:50051"
	cred, err := credentials.NewClientTLSFromFile("server.crt", "")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(cred), grpc.WithUnaryInterceptor(unaryInterceptor))
	if err != nil {
		log.Fatalf("did not conect %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := os.Args[1]

	ctx := context.Background()
	//ctx, cancel := context.WithCancel(ctx)
	md := metadata.Pairs("timestamp", time.Now().Format(time.Stamp))
	ctx = metadata.NewOutgoingContext(ctx, md)
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name}, grpc.Trailer(&md))
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			log.Printf("gRPC Error (message %s)", s.Message())
			for _, d := range s.Details() {
				switch info := d.(type) {
				case *errdetails.RetryInfo:
					log.Printf(" RetryInfo: %v", info)
				}
			}
			os.Exit(1)
		} else {
			log.Fatalf("could not greet %v", err)
		}

	}
	log.Printf("Greeting: %s", r.Message)

}

func unaryInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Printf("before call: %s, request: %+v", method, req)
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("after call %s, response: %+v", method, reply)
	return err
}
