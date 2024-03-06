package tests

import (
	"log"
	"os"
	"testing"
	"time"

	"drexel.edu/voter/db"
	fake "github.com/brianvoe/gofakeit/v6" //aliasing package name
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

var (
	BASE_API = "http://localhost:1080"

	cli = resty.New()
)

func TestMain(m *testing.M) {

	//SETUP GOES FIRST
	rsp, err := cli.R().Delete(BASE_API + "/voters")

	if rsp.StatusCode() != 200 {
		log.Printf("error clearing database, %v", err)
		os.Exit(1)
	}

	code := m.Run()

	//CLEANUP

	//Now Exit
	os.Exit(code)
}

func newRandVoter(id uint) db.Voter {
	return db.Voter{
		VoterId:     id,
		Name:        fake.RandomString([]string{"name1", "name2", "name3"}),
		Email:       fake.RandomString([]string{"name1@place.com", "name2@place.com", "name3@place.com"}),
		VoteHistory: make([]db.VoterHistory, 0),
	}
}

func newRandVote(voter_id uint, poll_id uint) db.VoterHistory {
	return db.VoterHistory{
		VoterId:  voter_id,
		PollId:   poll_id,
		VoteDate: fake.Date(),
	}
}

func Test_LoadDB(t *testing.T) {
	numLoad := 3
	for i := 0; i < numLoad; i++ {
		item := newRandVoter(uint(i))
		rsp, err := cli.R().
			SetBody(item).
			Post(BASE_API + "/voters")

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	}
}

func Test_LoadDuplicateVoter(t *testing.T) {
	item := newRandVoter(1)

	rsp, err := cli.R().SetBody(item).
		Post(BASE_API + "/voters")

	assert.Nil(t, err)
	assert.Equal(t, 500, rsp.StatusCode())
}

func Test_GetAllVoters(t *testing.T) {
	var items []db.Voter

	rsp, err := cli.R().SetResult(&items).Get(BASE_API + "/voters")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

	assert.Equal(t, 3, len(items))
}

func Test_LoadVoterWithWrongVoterInPolls(t *testing.T) {
	item := newRandVoter(4)
	item.VoteHistory = append(item.VoteHistory, newRandVote(1, 0))

	rsp, err := cli.R().SetBody(item).
		Post(BASE_API + "/voters")

	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_LoadVoterWithDuplicatePolls(t *testing.T) {
	item := newRandVoter(4)
	item.VoteHistory = append(item.VoteHistory, newRandVote(4, 0), newRandVote(4, 0))

	rsp, err := cli.R().SetBody(item).
		Post(BASE_API + "/voters")

	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_DeleteVoter(t *testing.T) {
	var item db.Voter

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/2")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #2 expected")

	rsp, err = cli.R().Delete(BASE_API + "/voters/2")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter not deleted expected")

	rsp, err = cli.R().SetResult(item).Get(BASE_API + "/voters/2")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode(), "expected not found error code")
}

func Test_DeleteNonExistentVoter(t *testing.T) {
	rsp, err := cli.R().Delete(BASE_API + "/voters/4")
	assert.Nil(t, err)
	assert.Equal(t, 500, rsp.StatusCode())
}

func Test_UpdateVoter(t *testing.T) {
	var item, changedItem db.Voter

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/1")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 expected.")

	item.Email = "new@email.com"

	rsp, err = cli.R().SetResult(&changedItem).SetBody(item).Put(BASE_API + "/voters/1")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "expected successful response")

	assert.Equal(t, item, changedItem)

	rsp, err = cli.R().SetResult(&changedItem).Get(BASE_API + "/voters/1")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 expected.")

	assert.Equal(t, item, changedItem)
}

func Test_UpdateNonExistentVoter(t *testing.T) {
	item := newRandVoter(4)

	rsp, err := cli.R().SetBody(item).Put(BASE_API + "/voters/1")
	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_UpdateVoterWithDuplicatePolls(t *testing.T) {

	var item db.Voter

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/0")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 expected.")

	item.VoteHistory = append(item.VoteHistory, newRandVote(2, 0), newRandVote(2, 0))

	rsp, err = cli.R().SetBody(item).Put(BASE_API + "/voters/0")
	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_UpdateVoterWithWrongVoterInVotes(t *testing.T) {

	var item db.Voter

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/0")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 expected.")

	item.VoteHistory = append(item.VoteHistory, newRandVote(3, 0))

	rsp, err = cli.R().SetBody(item).Put(BASE_API + "/voters/0")
	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_LoadVotes(t *testing.T) {
	numLoad := 3
	for i := 0; i < numLoad; i++ {
		item := newRandVote(1, uint(i))
		rsp, err := cli.R().
			SetBody(item).
			Post(BASE_API + "/voters/1/polls")

		assert.Nil(t, err)
		assert.Equal(t, 200, rsp.StatusCode())
	}
}

func Test_LoadDuplicateVote(t *testing.T) {
	item := newRandVote(1, 0)

	rsp, err := cli.R().SetBody(item).Post(BASE_API + "/voters/1/polls")
	assert.Nil(t, err)
	assert.Equal(t, 500, rsp.StatusCode())
}

func Test_LoadVoteToWrongVoter(t *testing.T) {
	item := newRandVote(1, 0)

	rsp, err := cli.R().SetBody(item).Post(BASE_API + "/voters/0/polls")
	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())

}
func Test_LoadVoteToNonExistentVoter(t *testing.T) {
	item := newRandVote(1, 0)

	rsp, err := cli.R().SetBody(item).Post(BASE_API + "/voters/4/polls")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode())
}

func Test_GetAllVotes(t *testing.T) {
	var items []db.VoterHistory

	rsp, err := cli.R().SetResult(&items).Get(BASE_API + "/voters/1/polls")

	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode())

	assert.Equal(t, 3, len(items))
}

func Test_GetAllVotesForMissingVoter(t *testing.T) {
	rsp, err := cli.R().Get(BASE_API + "/voters/4/polls")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode())
}

func Test_UpdateVote(t *testing.T) {

	var vote, changedVote db.VoterHistory

	url := BASE_API + "/voters/1/polls/1"

	rsp, err := cli.R().SetResult(&vote).Get(url)
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 with poll #1 expected")

	vote.VoteDate = time.Now()

	rsp, err = cli.R().SetBody(vote).SetResult(&changedVote).Put(url)
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "Expect successful change.")

	assert.Equal(t, vote.VoteDate.UnixMicro(), changedVote.VoteDate.UnixMicro())

	rsp, err = cli.R().SetResult(&changedVote).Get(url)
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 expected with poll #1 expected")

	assert.Equal(t, vote.VoteDate.UnixMicro(), changedVote.VoteDate.UnixMicro())

}

func Test_UpdateVoteOnMissingVoter(t *testing.T) {
	vote := newRandVote(4, 0)

	rsp, err := cli.R().SetBody(vote).Put(BASE_API + "/voters/4/polls/0")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode())
}

func Test_UpdateVoteWithWrongVoter(t *testing.T) {
	vote := newRandVote(2, 0)

	rsp, err := cli.R().SetBody(vote).Put(BASE_API + "/voters/1/polls/0")
	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_UpdateVoteWithWrongPollId(t *testing.T) {
	vote := newRandVote(1, 2)

	rsp, err := cli.R().SetBody(vote).Put(BASE_API + "/voters/1/polls/0")
	assert.Nil(t, err)
	assert.Equal(t, 400, rsp.StatusCode())
}

func Test_UpdateVoteWithMissingPollId(t *testing.T) {
	vote := newRandVote(1, 4)

	rsp, err := cli.R().SetBody(vote).Put(BASE_API + "/voters/1/polls/4")
	assert.Nil(t, err)
	assert.Equal(t, 500, rsp.StatusCode())
}

func Test_DeleteVote(t *testing.T) {
	var item db.VoterHistory

	rsp, err := cli.R().SetResult(&item).Get(BASE_API + "/voters/1/polls/2")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "voter #1 with poll #2 expected")

	rsp, err = cli.R().Delete(BASE_API + "/voters/1/polls/2")
	assert.Nil(t, err)
	assert.Equal(t, 200, rsp.StatusCode(), "vote not deleted expected")

	rsp, err = cli.R().SetResult(&item).Get(BASE_API + "/voters/1/polls/2")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode(), "expected not found error code")

	var voter db.Voter
	_, err = cli.R().SetResult(&voter).Get(BASE_API + "/voters/1")
	assert.Nil(t, err)
	assert.Equal(t, 2, len(voter.VoteHistory))
}

func Test_DeleteOnMissingVoter(t *testing.T) {
	rsp, err := cli.R().Delete(BASE_API + "/voters/4/polls/2")
	assert.Nil(t, err)
	assert.Equal(t, 404, rsp.StatusCode())
}

func Test_DeleteOnMissingVote(t *testing.T) {
	rsp, err := cli.R().Delete(BASE_API + "/voters/1/polls/4")
	assert.Nil(t, err)
	assert.Equal(t, 500, rsp.StatusCode())
}
