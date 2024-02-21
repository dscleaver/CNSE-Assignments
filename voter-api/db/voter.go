package db

import (
	"errors"
	"time"
)

type VoterHistory struct {
	PollId   uint      `json:"poll_id"`
	VoterId  uint      `json:"voter_id"`
	VoteDate time.Time `json:"vote_date"`
}

type Voter struct {
	VoterId     uint           `json:"voter_id"`
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	VoteHistory []VoterHistory `json:"voter_history"`
}

type VoterList struct {
	voters map[uint]Voter //A map of VoterIDs as keys and Voter structs as values
}

// constructor for VoterList struct
func NewVoterList() (*VoterList, error) {
	voterList := &VoterList{
		voters: make(map[uint]Voter),
	}
	return voterList, nil
}

func (v *VoterList) GetAllVoters() ([]Voter, error) {
	var voters []Voter

	for _, voter := range v.voters {
		voters = append(voters, voter)
	}

	return voters, nil
}

func (v *VoterList) AddVoter(voter Voter) error {
	_, ok := v.voters[voter.VoterId]
	if ok {
		return errors.New("Voter already exists.")
	}

	v.voters[voter.VoterId] = voter

	return nil
}

func (v *VoterList) GetVoter(id uint) (Voter, error) {
	voter, ok := v.voters[id]
	if !ok {
		return Voter{}, errors.New("No voter for id.")
	}

	return voter, nil
}

func (v *VoterList) DeleteVoter(id uint) error {
	_, ok := v.voters[id]

	if !ok {
		return errors.New("No voter for id.")
	}

	delete(v.voters, id)

	return nil
}

func (v *VoterList) DeleteAll() error {
	v.voters = make(map[uint]Voter)
	return nil
}

func (v *VoterList) UpdateVoter(voter Voter) error {
	_, ok := v.voters[voter.VoterId]
	if !ok {
		return errors.New("Voter does not exist.")
	}

	v.voters[voter.VoterId] = voter

	return nil
}

func (v *Voter) GetVote(pollId uint) (VoterHistory, error) {
	for _, vote := range v.VoteHistory {
		if vote.PollId == pollId {
			return vote, nil
		}
	}
	return VoterHistory{}, errors.New("Vote does not exist.")
}

func (v Voter) AddVote(newVote VoterHistory) (Voter, error) {
	for _, vote := range v.VoteHistory {
		if vote.PollId == newVote.PollId {
			return Voter{}, errors.New("Vote already exists.")
		}
	}

	v.VoteHistory = append(v.VoteHistory, newVote)
	return v, nil
}

func (v Voter) UpdateVote(updatedVote VoterHistory) (Voter, error) {
	var foundVote bool
	var voteIndex int
	for i, vote := range v.VoteHistory {
		if vote.PollId == updatedVote.PollId {
			foundVote = true
			voteIndex = i
			break
		}
	}

	if !foundVote {
		return Voter{}, errors.New("Vote does not exist.")
	}

	v.VoteHistory[voteIndex] = updatedVote

	return v, nil
}

func (v Voter) DeleteVote(pollId uint) (Voter, error) {
	var foundVote bool
	var voteIndex int
	for i, vote := range v.VoteHistory {
		if vote.PollId == pollId {
			foundVote = true
			voteIndex = i
			break
		}
	}

	if !foundVote {
		return Voter{}, errors.New("Vote does not exist.")
	}

	v.VoteHistory = append(v.VoteHistory[:voteIndex], v.VoteHistory[voteIndex+1:]...)
	return v, nil
}

func (v *Voter) ValidateVotes() error {
	var seen = make(map[uint]struct{})
	for _, vote := range v.VoteHistory {
		_, ok := seen[vote.PollId]
		if ok {
			return errors.New("Duplicate vote found")
		}
		seen[vote.PollId] = struct{}{}
		if vote.VoterId != v.VoterId {
			return errors.New("Vote does not match voter")
		}
	}

	return nil
}
