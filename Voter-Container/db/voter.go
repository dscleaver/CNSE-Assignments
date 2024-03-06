package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	RedisNilError        = "redis: nil"
	RedisDefaultLocation = "0.0.0.0:6379"
	RedisKeyPrefix       = "voter:"
)

type cache struct {
	client  *redis.Client
	context context.Context
}

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
	cache
}

// constructor for VoterList struct
func NewVoterList() (*VoterList, error) {
	redisUrl := os.Getenv("REDIS_URL")
	//This handles the default condition
	if redisUrl == "" {
		redisUrl = RedisDefaultLocation
	}
	return NewWithCacheInstance(redisUrl)
}

func NewWithCacheInstance(location string) (*VoterList, error) {
	client := redis.NewClient(&redis.Options{
		Addr: location,
	})

	ctx := context.TODO()

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Println("Error connecting to redis" + err.Error())
		return nil, err
	}

	return &VoterList{
		cache: cache{
			client:  client,
			context: ctx,
		},
	}, nil
}

// func isRedisNilError(err error) bool {
// 	return errors.Is(err, redis.Nil) || err.Error() == RedisNilError
// }

func redisKeyFromId(id uint) string {
	return fmt.Sprintf("%s%d", RedisKeyPrefix, id)
}

func (v *VoterList) getAllKeys() ([]string, error) {
	key := fmt.Sprintf("%s*", RedisKeyPrefix)
	return v.client.Keys(v.context, key).Result()
}

func fromJsonString(s string, item *Voter) error {
	err := json.Unmarshal([]byte(s), &item)
	if err != nil {
		return err
	}
	return nil
}

func (v *VoterList) upsertVoter(item *Voter) error {
	log.Println("Adding new Id:", redisKeyFromId(item.VoterId))
	return v.client.JSONSet(v.context, redisKeyFromId(item.VoterId), ".", item).Err()
}

// Helper to return a ToDoItem from redis provided a key
func (v *VoterList) getItemFromRedis(key string, item *Voter) error {
	itemJson, err := v.client.JSONGet(v.context, key, ".").Result()
	if err != nil {
		return err
	}

	return fromJsonString(itemJson, item)
}

func (t *VoterList) doesKeyExist(id uint) bool {
	kc, _ := t.client.Exists(t.context, redisKeyFromId(id)).Result()
	return kc > 0
}

func (v *VoterList) GetAllVoters() ([]Voter, error) {
	keyList, err := v.getAllKeys()
	if err != nil {
		return nil, err
	}
	voters := make([]Voter, len(keyList))

	for idx, key := range keyList {
		err := v.getItemFromRedis(key, &voters[idx])
		if err != nil {
			return nil, err
		}
	}

	return voters, nil
}

func (v *VoterList) AddVoter(voter Voter) error {
	if v.doesKeyExist(voter.VoterId) {
		return fmt.Errorf("Voter with id %d already exists", voter.VoterId)
	}
	return v.upsertVoter(&voter)
}

func (v *VoterList) GetVoter(id uint) (Voter, error) {
	var newVoter Voter
	err := v.getItemFromRedis(redisKeyFromId(id), &newVoter)
	if err != nil {
		return Voter{}, err
	}
	return newVoter, nil
}

func (v *VoterList) DeleteVoter(id uint) error {

	if !v.doesKeyExist(id) {
		return errors.New("no voter for id")
	}

	return v.client.Del(v.context, redisKeyFromId(id)).Err()
}

func (v *VoterList) DeleteAll() error {
	keyList, err := v.getAllKeys()
	if err != nil {
		return err
	}
	if len(keyList) == 0 {
		return nil
	}
	return v.client.Del(v.context, keyList...).Err()
}

func (v *VoterList) UpdateVoter(voter Voter) error {
	if !v.doesKeyExist(voter.VoterId) {
		return errors.New("Voter does not exist")
	}

	return v.upsertVoter(&voter)
}

func (v *Voter) GetVote(pollId uint) (VoterHistory, error) {
	for _, vote := range v.VoteHistory {
		if vote.PollId == pollId {
			return vote, nil
		}
	}
	return VoterHistory{}, errors.New("vote does not exist")
}

func (v Voter) AddVote(newVote VoterHistory) (Voter, error) {
	for _, vote := range v.VoteHistory {
		if vote.PollId == newVote.PollId {
			return Voter{}, errors.New("vote already exists")
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
		return Voter{}, errors.New("vote does not exist")
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
		return Voter{}, errors.New("vote does not exist")
	}

	v.VoteHistory = append(v.VoteHistory[:voteIndex], v.VoteHistory[voteIndex+1:]...)
	return v, nil
}

func (v *Voter) ValidateVotes() error {
	var seen = make(map[uint]struct{})
	for _, vote := range v.VoteHistory {
		_, ok := seen[vote.PollId]
		if ok {
			return errors.New("duplicate vote found")
		}
		seen[vote.PollId] = struct{}{}
		if vote.VoterId != v.VoterId {
			return errors.New("vote does not match voter")
		}
	}

	return nil
}
