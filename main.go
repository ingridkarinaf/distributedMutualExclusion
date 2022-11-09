package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	accessRequest "github.com/ingridkarinaf/distributedMutualExclusion/grpc"
	grpc "google.golang.org/grpc"
)

type peer struct {
	accessRequest.UnimplementedAccessRequestServer
	id           int32
	requestQueue map[int32]accessRequest.AccessRequestClient
	peers      map[int32]accessRequest.AccessRequestClient
	state		 string
	ctx          context.Context
}

func (p *peer) AccessRequest(ctx context.Context, req *accessRequest.Request) (*accessRequest.Reply, error) {
	id := req.Id
	

	rep := &accessRequest.Reply{Id: p.id}
	log.Printf("Sorting access to critical section")
	return rep, nil

}

func main() {
	//creating peer using terminal argument to create port 
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32) //Takes arguments 0, 1 and 2, see comment X
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:            ownPort,
		requestQueue: make(map[int32]accessRequest.AccessRequestClient),
		peers:       make(map[int32]accessRequest.AccessRequestClient),
		ctx:           ctx,
		state 		  string,
	}

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}
	grpcServer := grpc.NewServer()
	accessRequest.RegisterAccessRequestServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}
	}()

	//This is where we dial to all the other peers
	//Comment X: Hardcoded ports 5000, 5001, 5002, which is why these are the only valid ports
	for i := 0; i < 3; i++ {
		port := int32(5000) + int32(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		fmt.Printf("Trying to dial: %v\n", port)
		conn, err := grpc.Dial(fmt.Sprintf(":%v", port), grpc.WithInsecure(), grpc.WithBlock()) //This is going to wait until it receives the connection
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		defer conn.Close()
		c := accessRequest.NewAccessRequestClient(conn)
		p.peers[port] = c
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		p.sendPingToAll()
	}
}

}
