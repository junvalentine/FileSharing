/*
	This file contains the functions and RPC handlers that Peers
	use to handle files.
	RequestFile():
		- Requests a file from a Peer using a RequestFileArgs RPC.
	ServeFile():
		- Handles RequestFileArgs RPCs from Peers and returns the requested
		  file (if possible) to the requesting Peer using a RequestFileReply RPC.
	RegisterFile():
		- Peers use this function to register a file in the system. This means
		  to make the file publicly shareable with other peers.
	saveFile():
		- Private function that Peers use to save a file to 'disk' once obtained from
		  another Peer.
*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
	Requests a given file from a given Peer.
*/
func (p *Peer) RequestFile(port string, id int, file string) bool {
	requestFileArgs := RequestFileArgs{}
	requestFileReply := RequestFileReply{}

	requestFileArgs.PeerID = p.PeerID
	requestFileArgs.File = file
	call("Peer.ServeFile", &requestFileArgs, &requestFileReply, port)

	if requestFileReply.FileExists == false {
		fmt.Printf("Did not receive %v from Peer %v, the file does not exist\n", file, id)
		return false
	}

	fmt.Printf("Received %v from Peer %v\n", requestFileReply.File, id)
	save := saveFile(requestFileReply.File, requestFileReply.FileContents, p.PeerID, p.directory)
	return save
}

/*
	Handles file request RPCs (RequestFileArgs{}) from other Peers.
*/
// func (p *Peer) ServeFile(request *RequestFileArgs, reply *RequestFileReply) error {
// 	for i := 0; i < p.numFiles; i++ {
// 		if p.files[i] == request.File {
// 			reply.FileExists = true
// 			reply.File = p.files[i]
// 			location := p.fileloc[i]
// 			f, err := ioutil.ReadFile(location + request.File)
// 			if err != nil {
// 				fmt.Printf("Error reading file: %v\n", err)
// 			}
// 			data := string(f)
// 			reply.FileContents = data
// 			reply.PeerID = request.PeerID
// 			fmt.Printf("Served file %v to Peer %v\n", request.File, request.PeerID)
// 			return nil
// 		}
// 	}
// 	reply.FileExists = false
// 	reply.ErrorMessage = "File not found on the Server\n"
// 	reply.File = request.File
// 	reply.PeerID = request.PeerID
// 	fmt.Printf("Peer %v requested %v, but the file does not exist\n", request.PeerID, request.File)
// 	return nil
// }
func (p *Peer) ServeFile(request *RequestFileArgs, reply *RequestFileReply) error {
    done := make(chan bool) // create a new channel

    go func() { // start a new goroutine
        defer func() { done <- true }() // ensure done is always sent on exit
        p.mu.Lock() // acquire the lock before accessing shared data
        defer p.mu.Unlock() // release the lock after accessing shared data
        for i := 0; i < p.numFiles; i++ {
            if p.files[i] == request.File {
                reply.FileExists = true
                reply.File = p.files[i]
                location := p.fileloc[i]
                f, err := ioutil.ReadFile(location + request.File)
                if err != nil {
                    fmt.Printf("Error reading file: %v\n", err)
                }
                // data := string(f)
                reply.FileContents = f
                reply.PeerID = request.PeerID
                fmt.Printf("Served file %v to Peer %v\n", request.File, request.PeerID)
                return
            }
        }
        reply.FileExists = false
        reply.ErrorMessage = "File not found on the Server\n"
        reply.File = request.File
        reply.PeerID = request.PeerID
        fmt.Printf("Peer %v requested %v, but the file does not exist\n", request.PeerID, request.File)
    }()

    <-done // wait for the goroutine to send a value
    return nil
}

/*
	Registers a file that a Peer has on disk into the FileShare system.
*/
func (p *Peer) RegisterFile(fileName string, location string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.files[p.numFiles] = fileName
	p.fileloc[p.numFiles] = location
	p.numFiles = p.numFiles + 1

	request := PeerSendFile{}
	reply := ServerReceiveFile{}
	request.FileName = fileName
	request.PeerID = p.PeerID
	// request.location = location

	serverCall("Server.Register", &request, &reply)
	fmt.Printf("Registered file %v\n", fileName)
	return nil
}

/*
	Asks the Server for a particular file.
	The Server will search the network of Peers
	and find the Peer with the requested file, and
	then send the connection details back to the
	requesting Peer.
*/
func (p *Peer) SearchForFile(fileName string) error {
	p.mu.Lock()

	request := RequestFileArgs{}
	reply := FindPeerReply{}
	request.File = fileName
	request.PeerID = p.PeerID
	serverCall("Server.SearchFile", &request, &reply)

	if reply.Found {
		fmt.Printf("Num      PeerID\n")
		for i := 0; i < len(reply.PeerID); i++ {
			fmt.Printf("%v        %v\n", i+1, reply.PeerID[i])
		}

		fmt.Printf("Please choose a PeerID to connect to: ")
		var id int
		fmt.Scanf("%d", &id)

		p.ConnectPeer(reply.Port[id], reply.PeerID[id])
		save := p.RequestFile(reply.Port[id], reply.PeerID[id], reply.File)
		p.mu.Unlock()
		if save == true{
			p.RegisterFile(reply.File, p.directory)
		}
	} else{
		fmt.Printf("File %v not found\n", fileName)
		p.mu.Unlock()
	}
	return nil
}

/*
	Saves a newly received file to the Peer's repository.
*/

func saveFile(fileName string, fileContents []byte, id int, directory string) bool{
	filePath, _ := filepath.Abs(directory + fileName)
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating the file: %v\n", err)
		return false
	}

	l, err := f.Write(fileContents)
	if err != nil {
		fmt.Printf("Error writing the file: %v %v\n", err, l)
		return false
	}
	fmt.Printf("Saved file successfully %v\n", fileName)
	return true
}

func (p* Peer) ListFileReply(request *RequestListFile, reply *ListFileReply) error{
	reply.File = p.files
	reply.PeerID = p.PeerID
	reply.NumFiles = p.numFiles
	reply.Accepted = true

	return nil
}