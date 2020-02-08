package raft

import "net"

// RPCResponse captures both a response and a potential error
type RPCResponse struct {
	Response interface{}
	Error error
}

// RPC has a service to invoke, and provides a Reponse mechanism
type RPC struct {
	Service interface{}
	RespChan chan<- RPCResponse
}

// Respond is used to respond with a response, error or both
func (r *RPC) Respond(resp interface{}, err error) {
	r.RespChan <- RPCResponse{resp, err}
}

// Transport provides an interface for network trasnport
// to allow raft node to communicate with others
type Transport interface {
	// Listen returns a channel that can be used to
	// consume and respond to RPC requests.
	Listen() <-chan RPC

	// AppendEntries sends the appropriate RPC to the target node
	AppendEntries(target net.Addr, args *AppendEntriesArgs, resp *AppendEntriesResp) error

	// RequestVote sends the appropriate RPC to the target node
	RequestVote(target net.Addr, args *RequestVoteArgs, resp *RequestVoteResp) error
}
