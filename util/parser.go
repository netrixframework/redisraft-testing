package util

import (
	"encoding/json"
	"errors"

	"github.com/netrixframework/netrix/types"
)

var AppendEntriesRequestType string = "append_entries_request"
var AppendEntriesResponseType string = "append_entries_response"
var RequestVoteRequestType string = "request_vote_request"
var RequestVoteResponseType string = "request_vote_response"

type RedisRaftMessageParser struct{}

type RedisRaftMessage struct {
	Type                  string `json:"type"`
	AppendEntriesRequest  `json:"-"`
	AppendEntriesResponse `json:"-"`
	RequestVoteRequest    `json:"-"`
	RequestVoteResponse   `json:"-"`
}

func (r RedisRaftMessage) String() string {
	return ""
}

func (r RedisRaftMessage) Clone() types.ParsedMessage {
	return nil
}

func (r RedisRaftMessage) Marshal() ([]byte, error) {
	switch r.Type {
	case AppendEntriesRequestType:
		return json.Marshal(r.AppendEntriesRequest)
	case AppendEntriesResponseType:
		return json.Marshal(r.AppendEntriesResponse)
	case RequestVoteRequestType:
		return json.Marshal(r.RequestVoteRequest)
	case RequestVoteResponseType:
		return json.Marshal(r.RequestVoteResponse)
	}
	return nil, errors.New("invalid message type")
}

var _ types.MessageParser = &RedisRaftMessageParser{}

func (*RedisRaftMessageParser) Parse(data []byte) (types.ParsedMessage, error) {
	var rMsg RedisRaftMessage
	if err := json.Unmarshal(data, &rMsg); err != nil {
		return nil, err
	}
	switch rMsg.Type {
	case AppendEntriesRequestType:
		var ar AppendEntriesRequest
		if err := json.Unmarshal(data, &ar); err != nil {
			return nil, err
		}
		rMsg.AppendEntriesRequest = ar
	case AppendEntriesResponseType:
		var ar AppendEntriesResponse
		if err := json.Unmarshal(data, &ar); err != nil {
			return nil, err
		}
		rMsg.AppendEntriesResponse = ar
	case RequestVoteRequestType:
		var rr RequestVoteRequest
		if err := json.Unmarshal(data, &rr); err != nil {
			return nil, err
		}
		rMsg.RequestVoteRequest = rr
	case RequestVoteResponseType:
		var rr RequestVoteResponse
		if err := json.Unmarshal(data, &rr); err != nil {
			return nil, err
		}
		rMsg.RequestVoteResponse = rr
	default:
		return nil, errors.New("invalid message type")
	}

	return rMsg, nil
}

type Entry struct {
	Term    int64
	ID      int
	Session int64
	Type    int
	Data    string
}

type AppendEntriesRequest struct {
	LeaderID     int `json:"leader_id"`
	Term         int64
	PrevLogIndex int64 `json:"prev_log_idx"`
	PrevLogTerm  int64 `json:"prev_log_term"`
	LeaderCommit int64 `json:"leader_commit"`
	MessageID    int64 `json:"msg_id"`
	Entries      []Entry
}

type AppendEntriesResponse struct {
	Term         int64
	Success      int
	CurrentIndex int64 `json:"current_idx"`
	MessageID    int64 `json:"msg_id"`
}

type RequestVoteRequest struct {
	Prevote      int
	Term         int64
	CandidateID  int64 `json:"candidate_id"`
	LastLogIndex int64 `json:"last_log_idx"`
	LastLogTerm  int64 `json:"last_log_term"`
}

type RequestVoteResponse struct {
	Term        int64
	Prevote     int
	RequestTerm int64
	VoteGranted int
}
