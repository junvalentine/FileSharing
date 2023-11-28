package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

/*
	Struct type for the Peers.
*/
type Peer struct {
	PeerID    int
	files     []string
	fileloc   []string
	peers     []int
	numFiles  int
	numPeers  int
	directory string
	Port      string
	mu        sync.Mutex
}

/*
	A lightweight data type for the SwarmMaster and Peers to
	hold relevant information about the Peers connected
	to it, including their port and the files they posses.
*/
type PeerInfo struct {
	PeerID      int
	Port        string
	Files       [100]string
	// Fileloc		[100]string
	numFiles    int
	isConnected bool
}


func (m *Peer) Welcome() {
	fmt.Printf("Welcome to the File-Sharing Application\n")
}
/*
	Method used to make Remote Procedure Calls (RPCs)
	Adopted from provided lab code.
*/
func call(rpcname string, args interface{}, reply interface{}, port string) bool {
	c, err := rpc.DialHTTP("tcp", port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

/*
	Method for the Peers to make RPC calls to the SwarmMaster.
*/
func serverCall(rpcname string, args interface{}, reply interface{}) bool {
	c, err := rpc.DialHTTP("tcp", "192.168.32.101:1337")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()

	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}
	fmt.Println(err)
	return false
}

/*
	Creates a server for the Peer so that other Peers can connect.
*/
func (p *Peer) peerServer(port string) {
	rpc.Register(p)
	serv := rpc.NewServer()
	serv.Register(p)
	oldMux := http.DefaultServeMux
	mux := http.NewServeMux()
	http.DefaultServeMux = mux
	serv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	http.DefaultServeMux = oldMux
	l, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	go http.Serve(l, mux)
}

/*
	Method called to create a new Peer.
*/
func MakePeer(directory string, port string) *Peer {
	p := Peer{}

	// p.PeerID = id
	p.directory = directory 
	p.files = make([]string, 100)
	p.fileloc = make([]string, 100)
	p.Port = port
	p.numFiles = 0
	p.peers = make([]int, 100)
	p.numPeers = 0

	p.peerServer(port)
	return &p
}

func (p *Peer) ConnectServer() {
	request := ConnectRequest{}
	reply := ConnectReply{}
	// request.PeerID = p.PeerID
	request.Port = p.Port
	serverCall("Server.ConnectPeer", &request, &reply)
	p.PeerID = reply.PeerID
	if reply.Accepted == true {
		fmt.Printf("Connected to server, PeerID: %v\n", p.PeerID)
	}
}

/*
	Handles incoming connection RPCs (ConnectRequest{}) from other Peers.
*/
// func (p *Peer) AcceptConnect(request *ConnectRequest, reply *ConnectReply) error {
// 	fmt.Printf("Received ConnectRequest from Peer %v\n", request.PeerID)

// 	p.mu.Lock()
// 	defer p.mu.Unlock()

// 	reply.Accepted = true
// 	reply.PeerID = request.PeerID

// 	p.peers[p.numPeers] = request.PeerID
// 	p.numPeers = p.numPeers + 1

// 	fmt.Printf("Connected to Peer: %v\n", request.PeerID)
// 	return nil
// }
func (p *Peer) AcceptConnect(request *ConnectRequest, reply *ConnectReply) error {
    done := make(chan bool) // create a new channel

    go func() { // start a new goroutine
        defer func() { done <- true }() // ensure done is always sent on exit
        fmt.Printf("Received ConnectRequest from Peer %v\n", request.PeerID)
        p.mu.Lock() // acquire the lock before updating shared data
        defer p.mu.Unlock()
        p.peers[p.numPeers] = request.PeerID
        reply.PeerID = request.PeerID
        p.numPeers = p.numPeers + 1
        reply.Accepted = true
        fmt.Printf("Accepted connection from Peer %v\n", request.PeerID)
    }()

    <-done // wait for the goroutine to send a value
    return nil
}
/*
	Connects the Peer to the provided Peer.
*/
func (p *Peer) ConnectPeer(port string, id int) {
	request := ConnectRequest{}
	reply := ConnectReply{}
	request.PeerID = p.PeerID
	request.Port = p.Port
	call("Peer.AcceptConnect", &request, &reply, port)
	if reply.Accepted == false {
		fmt.Printf("Connection refused from Peer %v\n", id)
		return
	}
	p.peers[p.numPeers] = id
	p.numPeers = p.numPeers + 1
	fmt.Printf("Connected to Peer %v\n", id)
}