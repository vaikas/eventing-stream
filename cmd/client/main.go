package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/vaikas/eventing-stream/proto"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:10000"
	defaultName = "world"
)

func main() {
	events := []pb.CloudEvent{
		{Id: "first"},
		{Id: "second"},
		{Id: "third"},
		{Id: "fourth"},
		{Id: "fifth"},
		{Id: "sixth"},
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewStreamClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	stream, err := c.Stream(ctx)
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a event : %v", err)
			}
			log.Printf("Got Cloud Event %+vs ID: %s", in, in.Id)
		}
	}()
	for _, event := range events {
		if err := stream.Send(&event); err != nil {
			log.Fatalf("Failed to send a event: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}
