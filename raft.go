package raft

import (
	"net"
)

type Node struct {
	raftState

	// Configuration for node initialization
	conf *Config

	// Node Id to distinguish it from others
	nodeID     string

	// LogStore provides durable log storage
	logs       LogStore

	// Transport layer , we are suing
	trans      Transport

	// Address of other known nodes
	peers      []net.Addr

	// Channel to liste for RPC requests
	rpcCh      <-chan RPC

	// Channel to signal the shutdown
	shutdownCh <-chan struct{}
}

func NewNode(logs LogStore, peers []net.Addr, trans Transport,
	shutdownCh <-chan struct{}) (*Node, error) {

	conf := DefaultConfig()
	nodeID := UUID()

	rs := raftState{
		currentTerm: 0,
		votedFor: NULL,
		commitIndex: 0,
		lastAppliedIndex: 0,
	}
	node := &Node{
		raftState: rs,
		conf:       conf,
		nodeID:     nodeID,
		logs:       logs,
		trans:      trans,
		peers:      peers,
		rpcCh:      trans.Listen(),
		shutdownCh: shutdownCh,
	}

	return node, nil
}
