package raft

import (
	"net"
	"sync"
)

type Node struct {
	raftState

	// Configuration for node initialization
	conf *Config

	// Node Id to distinguish it from others
	nodeID string

	// LogStore provides durable log storage
	logs LogStore

	// Transport layer , we are suing
	trans Transport

	// Address of other known nodes
	peers []net.Addr

	// Channel to liste for RPC requests
	rpcCh <-chan RPC

	// Channel to signal the shutdown
	shutdownCh <-chan struct{}

	// routines bucket
	routines sync.WaitGroup
}

func NewNode(logs LogStore, peers []net.Addr, trans Transport,
	shutdownCh <-chan struct{}) (*Node, error) {

	conf := DefaultConfig()
	nodeID := UUID()

	rs := raftState{
		currentTerm:      0,
		votedFor:         NULL,
		commitIndex:      0,
		lastAppliedIndex: 0,
	}
	node := &Node{
		raftState:  rs,
		conf:       conf,
		nodeID:     nodeID,
		logs:       logs,
		trans:      trans,
		peers:      peers,
		rpcCh:      trans.Listen(),
		shutdownCh: shutdownCh,
		routines:   sync.WaitGroup{},
	}

	// bootup the node
	node.goFunc(node.boot)

	return node, nil
}

func (n *Node) goFunc(fn func()) {
	n.routines.Add(1)
	go func() {
		defer n.routines.Done()
		fn()
	}()
}

// boot is a long goroutine to boot up the node
func (n *Node) boot() {

	// Set the node's mode to Follower
	n.setMode(Follower)

	for {
		select {
		case _ = <-n.rpcCh:
			// Got an rpc request
			// TODO
		case <-randomTimeout(n.conf.ElectionTimeout):
			// not heard from leader, become Candidate and start the election
			// TODO
		}
	}
}

// runCandidate is a long goroutine which governs the operation of node
// while in Candidate mode
func (n *Node) runCandidate() {
	// update term
	n.setCurrentTerm(n.getCurrentTerm() + 1)

	// voteCh to collect votes
	voteCh := make(chan RequestVoteResp, len(n.peers))

	// vote for yourself
	voteCh <- RequestVoteResp{
		Term:        n.currentTerm,
		VoteGranted: true,
	}

	// Number of votes needed to win
	voteNeedToWin := (len(n.peers) / 2) + 1

	// Vote collected
	voteCount := 0

	//Ask vote from peers
	askVoteFromPeer := func(peer net.Addr) {
		args := &RequestVoteArgs{
			Term:         n.currentTerm,
			CandidateID:  n.nodeID,
			LastLogIndex: n.logs.LastIndex(),
			LastLogTerm:  n.logs.LastTerm(),
		}

		resp := new(RequestVoteResp)

		n.goFunc(func() {
			err := n.trans.RequestVote(peer, args, resp)
			if err != nil {
				// Didn't get the vote
				// TODO Log
			}

			// send the vote
			voteCh <- *(resp)
		})
	}

	for _, peer := range n.peers {
		askVoteFromPeer(peer)
	}

	// Determine if we won the election or not
	for i := 0; i < len(n.peers); i++ {
		select {
		case rpcResp := <-voteCh:
			if !rpcResp.VoteGranted {
				// Vote request denied
				// Got greater term, revert back to follower
				// TODO revet state
			} else {
				// Vote request granted
				voteCount += 1

				// Check is we won the election
				if voteCount >= voteNeedToWin {
					// We won the election
					// TODO become leader node
				}
			}
		}
	}
}
