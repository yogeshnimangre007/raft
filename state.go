package raft

import "sync/atomic"

type Mode uint32

const NULL = -1

const (
	Follower Mode = iota
	Candidate
	Leader
)

type raftState struct {
	// Current mode of the node
	mode Mode

	// Current term of the node
	currentTerm uint64

	// If voted for someone in this term, else -1
	votedFor int64

	// index of latest committed log entry
	commitIndex uint64

	// index of latest applied log entry
	lastAppliedIndex uint64
}

func (rs *raftState) setMode(mode Mode) {
	atomic.StoreUint32((*uint32)(&rs.mode), uint32(mode))
}

func (rs *raftState) getMode() Mode {
	return Mode(atomic.LoadUint32((*uint32)(&rs.mode)))
}
