/*
	This file contains the Server structs, functions and RPC handlers.
*/

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
	Server data type for the server.
*/
type Server struct {
	peers    []PeerInfo
	numPeers int
	mu       sync.Mutex
}

/*
	RPC handler for when a Peer wishes to connect
	to the Server.
*/
func (m *Server) ConnectPeer(request *ConnectRequest, reply *ConnectReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = true
	reply.PeerID = m.numPeers
	m.peers[m.numPeers].PeerID = m.numPeers
	m.peers[m.numPeers].Port = request.Port
	m.peers[m.numPeers].isConnected = true
	fmt.Printf("Connected to Peer: %v\n", m.numPeers)

	m.numPeers = m.numPeers + 1
	return nil
}

/*
	Simple function to let us know when the
	Server has successfully been built.
*/
func (m *Server) Welcome() {
	fmt.Printf("Welcome to the File-Sharing Application\n")
}

/*
	RPC handler for when a Peer registers a file in
	the FileShare system to be shareable. This function will
	update the Server's peers data to include the new
	file.
*/
func (m *Server) Register(request *PeerSendFile, reply *ServerReceiveFile) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Accepted = false
	reply.FileName = request.FileName
	reply.Received = true
	for i := 0; i < m.numPeers; i++ {
		if m.peers[i].PeerID == request.PeerID {
			m.peers[i].Files[m.peers[i].numFiles] = request.FileName
			m.peers[i].numFiles++
			// m.peers[i].Fileloc[m.peers[i].numFiles] = request.location
			reply.Accepted = true
			fmt.Printf("Registered %v from Peer %v\n", request.FileName, request.PeerID)
			break
		}
	}
	return nil
}

/*
	RPC handler for when a Peer is in search of a file.
	This function will search the registered files in each
	Peer's file list to find which Peer contains the requested
	file. Then a FindPeerReply RPC will be sent to the requesting
	Peer telling it how to contact the Peer with the desired file.
*/
func (m *Server) SearchFile(request *RequestFileArgs, reply *FindPeerReply) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	reply.Found = false
	reply.File = request.File
	fmt.Printf("Peer %v requested a search for file %v\n", request.PeerID, request.File)
	for i := 0; i < m.numPeers; i++ {
		for j := 0; j < m.peers[i].numFiles; j++ {
			if request.File == m.peers[i].Files[j] {
				reply.Found = true
				reply.PeerID = append(reply.PeerID,m.peers[i].PeerID)
				reply.Port = append(reply.Port,m.peers[i].Port)
				fmt.Printf("Found file %v for Peer %v on Peer %v\n", request.File, request.PeerID, m.peers[i].PeerID)
			}
		}
	}

	if reply.Found == false{
		fmt.Printf("Cannot find a Peer containing file %v for Peer %v\n", request.File, request.PeerID)
	}
	return nil
}

/*
	Starts the server.
*/
func (m *Server) server() {
	rpc.Register(m)
	rpc.HandleHTTP()

	l, e := net.Listen("tcp", ":1337")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

/*
	Creates a new Server
*/
func MakeServer() *Server {
	m := Server{}
	// 10 Peers is arbitrary
	m.peers = make([]PeerInfo, 100)
	m.numPeers = 0
	m.server()
	return &m
}

/* 
	List all the peer that has connected to server
*/
func (m *Server) ListPeers() {
	fmt.Printf("Num      PeerID      Address\n")
	for i := 0; i < m.numPeers; i++ {
		fmt.Printf("%v        %v           %v\n", i+1, m.peers[i].PeerID, "0.0.0.0" + m.peers[i].Port)
	}
}

/* 
	Ping a peer
*/
func (m *Server) PingPeer(peerID int) bool{
	fmt.Printf("Pinging Peer %v\n", peerID)
	check := 0
	for i := 0; i < 3; i++ {
		_ , err := net.Dial("tcp", m.peers[peerID].Port)
		if err != nil {
			check += 1
		}
	}
	if check == 3 {
		fmt.Printf("Peer not live!\n")
		return false
	} else{
		fmt.Printf("%v/3 connection success\n", 3 - check)
		fmt.Printf("Peer live!\n")
		return true
	}
}
/* 
	Discover all file in local repo of a peer
*/
func (m *Server) DiscoverFile(peerID int) {
	fmt.Printf("Discovering file in Peer %v\n", peerID)
	check := m.PingPeer(peerID)
	if check == false{
		return
	}

	request := RequestListFile{}
	reply := ListFileReply{}
	reply.Accepted = false

	call("Peer.ListFileReply", &request, &reply, m.peers[peerID].Port)
	if reply.Accepted == true {
		fmt.Printf("Num      Files\n")
		for i := 0; i < reply.NumFiles; i++ {
			fmt.Printf("%v        %v\n", i+1, reply.File[i])
		}
	}
	return 
}

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