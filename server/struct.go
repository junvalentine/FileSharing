package main

import (
	"sync"
)

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

type PeerInfo struct {
	PeerID      int
	Port        string
	Files       [100]string
	// Fileloc		[100]string
	numFiles    int
	isConnected bool
}