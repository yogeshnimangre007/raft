package raft

// Log entries are replicated to all members of the Raft cluster
// and form the heart of the replicated state machine.
type Log struct {
	Index uint64
	Term uint64
	Data []byte
}

// LogStore is used to provide an interface for storing
// and retrieving logs in a durable fashion
type LogStore interface {
	// Returns the last index written. 0 for no entries.
	LastIndex() uint64

	// Returns the last term written. 0 for no entries.
	LastTerm() uint64

	// Gets a log entry at a given index
	GetLog(index uint64, log *Log) error

	// Stores a log entry
	StoreLog(log *Log) error

	// Deletes a range of log entries. The range is inclusive.
	DeleteRange(min, max uint64) error
}

