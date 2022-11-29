package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	accessRequest "github.com/ingridkarinaf/distributedMutualExclusion/grpc"
	"google.golang.org/grpc"
)

type peer struct {
	accessRequest.UnimplementedAccessRequestServer
	id           int32
	requestQueue map[int32]accessRequest.AccessRequestClient
	peers        map[int32]accessRequest.AccessRequestClient
	state        string
	ctx          context.Context
	lock 		 chan(bool)
	approvals    map[int32]bool
}

func main() {
	//creating peer using terminal argument to create port
	lock := make(chan bool, 1)
	lock <- true 
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32) //Takes arguments 0, 1 and 2, see comment X
	ownPort := int32(arg1) + 5001
	input_state := "" //Takes arguments 0, 1 and 2, see comment X

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:           ownPort,
		requestQueue: make(map[int32]accessRequest.AccessRequestClient),
		peers:        make(map[int32]accessRequest.AccessRequestClient),
		ctx:          ctx,
		state:        input_state, //wanted / not_wanted / holding
		lock: 		  lock,
		approvals: 	  make(map[int32]bool),
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
		port := int32(5001) + int32(i)

		if port == ownPort {
			continue
		}

		var conn *grpc.ClientConn
		log.Printf("Trying to dial: %v\n", port)
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
		log.Printf("%v wants to drive the car", p.id)
		reply, err := p.sendAccessRequestToAll()
		if err != nil {
			log.Printf("Something went wrong with sendAccessRequestToAll")
		}
		log.Printf("%v got a reply: %s", p.id, reply)
	}
}

func (p *peer) AccessRequest(ctx context.Context, req *accessRequest.Request) (*accessRequest.Reply, error) {
	<-p.lock
	id := req.Id

	if p.state == "holding" || (p.state == "wanted" && p.id < id) || (p.state == "wanted" && p.approvals[id]) {
		log.Printf("Either %v is using the car or it wants it. We have put %v in the queue", p.id, id)
		p.requestQueue[id] = p.peers[id]
		for p.state == "holding" {
			//log.Println(p.id, "is driving the car")
		}

	} else {
		log.Printf("Peer with id %v can drive it ", id)

	}
	p.lock<-true
	return &accessRequest.Reply{Id: p.id}, nil

}

func (p *peer) sendAccessRequestToAll() (*accessRequest.Reply, error) {
	request := &accessRequest.Request{Id: p.id}
	p.state = "wanted"

	log.Println("making access requests")

	for id, peer := range p.peers { //All the peers we will send the request to

		reply, err := peer.AccessRequest(p.ctx, request)
		if err != nil {
			log.Printf("Something went wrong when getting reply from %v", id)
		}
		log.Printf("Got reply from peer: %v \n", reply.Id)
		p.approvals[id] = true

	}

	p.state = "holding"
	log.Printf("%v driving the car \n", p.id)
	time.Sleep(5 * time.Second)
	log.Printf("%v released the car", p.id)
	p.state = "not wanted"
	//If the node has anyone in the queue
	var repl *accessRequest.Reply
	for id, peer := range p.requestQueue {
		_ = peer
		p.approvals[id] = false
		rep := &accessRequest.Reply{Id: p.id}
		repl = rep
		log.Printf("%v sends reply to %v ", p.id, id)
		delete(p.requestQueue, id)
	}

	return repl, nil
}
