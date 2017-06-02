package main

import (
	"encoding/json"
	"fmt"
	"log"

	pb "github.com/tiwariHD/goDrmCdi/grpcImplemetation/commandProto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address    = "localhost:50051"
	defaultCmd = "VERSION"
)

//getJSONStruct prints json ouput on stdin
func getJSONStruct(i interface{}) []byte {
	out, _ := json.MarshalIndent(i, "", "    ")
	//fmt.Printf("%s\n", out)
	return out
}

func checkVersion(c pb.CmdProtoClient) {
	fmt.Println("#VERSION: Empty{}")
	r, err := c.Version(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))
}

func checkInfo(c pb.CmdProtoClient) {
	fmt.Println("#INFO: Empty{}")
	r, err := c.Info(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))
}

func checkAdd(c pb.CmdProtoClient) {
	fmt.Println("#ADD: AddRequest{Version: \"0.0.1\", Request: \"gpu:1\", RequestId: \"1234\"}")
	r, err := c.Add(context.Background(), &pb.AddRequest{Version: "0.0.1",
		Request: "gpu:1", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))

	fmt.Println("ADD: AddRequest{Version: \"0.0.1\", Request: \"gpu:2\", RequestId: \"1234\"}")
	r, err = c.Add(context.Background(), &pb.AddRequest{Version: "0.0.1",
		Request: "gpu:2", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))

	fmt.Println("ADD: AddRequest{Version: \"0.0.1\", Request: \"gpu:1, gpu-memory=2048Mi\", RequestId: \"1234\"}")
	r, err = c.Add(context.Background(), &pb.AddRequest{Version: "0.0.1",
		Request: "gpu:1, gpu-memory=2048Mi", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))

	fmt.Println("ADD: AddRequest{Version: \"0.999\", Request: \"gpu:1\", RequestId: \"1234\"}")
	r, err = c.Add(context.Background(), &pb.AddRequest{Version: "0.999",
		Request: "gpu:1", RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))

	fmt.Println("ADD: AddRequest{Version: \"0.0.1\", Request: \"gpu:1\"}")
	r, err = c.Add(context.Background(), &pb.AddRequest{Version: "0.0.1",
		Request: "gpu:1"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))

}

func checkDel(c pb.CmdProtoClient) {
	fmt.Println("DEL: DelRequest{Version: \"0.0.1\", RequestId: \"1234\"}")
	r, err := c.Del(context.Background(), &pb.DelRequest{Version: "0.0.1",
		RequestId: "1234"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))

	fmt.Println("DEL: DelRequest{Version: \"0.0.1\", RequestId: \"3456\"}")
	r, err = c.Del(context.Background(), &pb.DelRequest{Version: "0.0.1",
		RequestId: "3456"})
	if err != nil {
		log.Fatalf("Failure: %v", err)
	}
	fmt.Printf("%s\n", getJSONStruct(r))
}

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewCmdProtoClient(conn)

	// Contact the server and print out its response.
	checkVersion(c)
	checkInfo(c)
	checkAdd(c)
	checkDel(c)
}
