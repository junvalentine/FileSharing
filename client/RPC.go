package main

import (
	// "os"
	// "strconv"
)

/*
	Request RPC for Peer's to connect.
*/
type ConnectRequest struct {
	PeerID int
	Port   string
}

/*
	Reply RPC for Peer's to connect.
*/
type ConnectReply struct {
	PeerID   int
	Accepted bool
}

/*
	RPC for a Peer to send a file to the server.
*/
type PeerSendFile struct {
	// location string
	PeerID   int
	FileName string
}

/*
	RPC for the server to confirm it received the file.
*/
type ServerReceiveFile struct {
	FileName string
	Received bool
	Accepted bool
}

/*
	Sent by the Peer to the Server when searching for
	a file in the network using Peer.SearchForFile().
*/
type RequestFileArgs struct {
	PeerID int
	File   string
}

/*
	Used by a peer to send another Peer a file in Peer.RequestFile()
	and Peer.ServeFile().
*/
type RequestFileReply struct {
	PeerID       int
	FileExists   bool
	ErrorMessage string
	File         string
	FileContents []byte
}

/*
	Sent by the Server to a Peer indicating the details
	regarding a Peer that possesses a particular file. Used
	in Peer.SearchForFile() and Server.SearchFile().
*/
type FindPeerReply struct {
	PeerID []int
	Port   []string
	File   string
	Found  bool
}

/*
	RPC for discover 
*/
type RequestListFile struct {
	PeerID 	int
}

type ListFileReply struct {
	File 	[]string
	PeerID 	int
	NumFiles int
	Accepted bool
}