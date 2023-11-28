/*
	This file contains the SwarmMaster structs, functions and RPC handlers.
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
	to the SwarmMaster.
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
	SwarmMaster has successfully been built.
*/
func (m *Server) Welcome() {
	fmt.Printf("Welcome to the File-Sharing Application\n")
}

/*
	RPC handler for when a Peer registers a file in
	the FileShare system to be shareable. This function will
	update the SwarmMaster's peers data to include the new
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
