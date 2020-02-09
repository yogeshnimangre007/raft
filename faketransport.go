package raft

import (
	"fmt"
	"net"
	"sync"
)

// Implements the Transport interface to allow Raft to be tested
// in-memory without going over a network

type FakeTransAddr struct {
	ID string
}

func (fta *FakeTransAddr) Network() string {
	return "faketransport"
}

func (fta *FakeTransAddr) String() string {
	return fta.ID
}

func NewFakeTransAddr() *FakeTransAddr {
	return &FakeTransAddr{UUID()}
}

type FakeTransport struct {
	sync.RWMutex
	listenerCh chan RPC
	localAddr  FakeTransAddr
	peers      map[string]*FakeTransport
}

func (ft *FakeTransport) Listen() <-chan RPC {
	return ft.listenerCh
}

func (ft *FakeTransport) LocalAddr() FakeTransAddr {
	return ft.localAddr
}

func (ft *FakeTransport) RequestVote(target net.Addr, args *RequestVoteArgs, resp *RequestVoteResp) error {
	ft.Lock()
	toPeer, ok := ft.peers[target.String()]
	ft.Unlock()

	if !ok {
		// Log something here
		fmt.Printf("Failed to connect to peer: %v\n", target.String())
	}

	respCh := make(chan RPCResponse)
	toPeer.listenerCh <- RPC{
		Args:     args,
		RespChan: respCh,
	}

	rpcResp := <-respCh
	if rpcResp.Error != nil {
		return rpcResp.Error
	}

	out := rpcResp.Response.(*RequestVoteResp)
	*resp = *out
	return nil
}

func (ft *FakeTransport) AppendEntries(target net.Addr, args *AppendEntriesArgs, resp *AppendEntriesResp) error {
	// TODO
	return nil
}

// Connect is used to connect this transport to another transport for
// a given peer name. This allows for local routing.
func (ft *FakeTransport) Connect(peer net.Addr, trans *FakeTransport) {
	ft.Lock()
	defer ft.Unlock()
	ft.peers[peer.String()] = trans
}

// Disconnect is used to remove the ability to route to a given peer
func (ft *FakeTransport) Disconnect(peer net.Addr) {
	ft.Lock()
	defer ft.Unlock()
	delete(ft.peers, peer.String())
}
