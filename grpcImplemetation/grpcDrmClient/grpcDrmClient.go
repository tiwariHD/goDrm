package main

import (
	"fmt"
	"log"
	"os"

	pb "github.com/tiwariHD/commandProto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address    = "localhost:50051"
	defaultCmd = "VERSION"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCmdProtoClient(conn)

	// Contact the server and print out its response.
	command := defaultCmd
	if len(os.Args) > 1 {
		command = os.Args[1]
	}
	fmt.Println("#env CDI_COMMAND=VERSION ./goDrmCdi < drm.conf")
	r, err := c.GetReply(context.Background(), &pb.CmdRequest{Command: command})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_COMMAND=INFO ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Command: "INFO"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "ADD", Request: "gpu:1", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:2 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "ADD", Request: "gpu:2", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "DEL", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1,gpu-memory=2048Mi CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "ADD", Request: "gpu:1, gpu-memory=2048Mi", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=DEL CDI_REQUEST_ID=3456 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "DEL", RequestId: "3456"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.999 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 CDI_REQUEST_ID=1234 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.999",
		Command: "ADD", Request: "gpu:1", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=ADD CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "ADD", Request: "gpu:1"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_COMMAND=MYCMD ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Command: "MYCMD"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)

	fmt.Println("#env CDI_VERSION=0.0.1 CDI_REQUEST=gpu:1 ./goDrmCdi < drm.conf")
	r, err = c.GetReply(context.Background(), &pb.CmdRequest{Version: "0.0.1",
		Request: "gpu:1"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", r.Message)
}
