package raft

type RequestVoteArgs struct {
	// Candidate's Term
	Term uint64

	// Candidate's ID
	CandidateID string

	// Index of candidate's last log entry
	LastLogIndex uint64

	// Term of candidate's last log entry
	LastLogTerm uint64
}

type RequestVoteResp struct {
	// For candidate to update himself
	Term uint64

	// Is candidate received vote
	VoteGranted bool
}

type AppendEntriesArgs struct {
	// TODO
}

type AppendEntriesResp struct {
	// TODO
}
