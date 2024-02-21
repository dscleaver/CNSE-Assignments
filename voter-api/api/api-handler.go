package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"drexel.edu/voter/db"
	"github.com/gofiber/fiber/v2"
)

// The api package creates and maintains a reference to the data handler
// this is a good design practice
type VoterAPI struct {
	db         *db.VoterList
	bootTime   time.Time
	totalCalls uint64
	errors     map[int]uint64
}

func New() (*VoterAPI, error) {
	dbHandler, err := db.NewVoterList()
	if err != nil {
		return nil, err
	}

	return &VoterAPI{db: dbHandler, bootTime: time.Now(), errors: make(map[int]uint64)}, nil
}

func (v *VoterAPI) HandleStats(c *fiber.Ctx) error {
	v.totalCalls++
	err := c.Next()

	var e *fiber.Error
	if errors.As(err, &e) {
		v.errors[e.Code]++
	}
	return err
}

func (v *VoterAPI) ListAllVoters(c *fiber.Ctx) error {
	voterList, err := v.db.GetAllVoters()
	if err != nil {
		log.Println("Error getting all voters: ", err)
		return fiber.NewError(http.StatusNotFound,
			"Error Getting All Voters")
	}
	if voterList == nil {
		voterList = make([]db.Voter, 0)
	}
	return c.JSON(voterList)
}

func (v *VoterAPI) DeleteAllVoters(c *fiber.Ctx) error {
	if err := v.db.DeleteAll(); err != nil {
		log.Println("Error deleting all voters: ", err)
		return fiber.NewError(http.StatusNotFound,
			"Error Deleting All Voters")
	}

	return c.Status(http.StatusOK).SendString("Delete OK")
}

func (v *VoterAPI) withVoter(c *fiber.Ctx, run func(voter db.Voter) error) error {
	param := struct {
		ID uint `params:"id"`
	}{}

	err := c.ParamsParser(&param) // "{"id": 111}"
	if err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	voter, err := v.db.GetVoter(param.ID)
	if err != nil {
		log.Println("Voter not found: ", err)
		return fiber.NewError(http.StatusNotFound)
	}

	return run(voter)
}

func (v *VoterAPI) GetVoter(c *fiber.Ctx) error {
	return v.withVoter(c, func(voter db.Voter) error {
		return c.JSON(voter)
	})
}

func (v *VoterAPI) AddVoter(c *fiber.Ctx) error {
	var voter db.Voter

	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := voter.ValidateVotes(); err != nil {
		log.Println("Voter is not valid.")
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.AddVoter(voter); err != nil {
		log.Println("Error adding voter: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

func (v *VoterAPI) UpdateVoter(c *fiber.Ctx) error {
	param := struct {
		ID uint `params:"id"`
	}{}

	if err := c.ParamsParser(&param); err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	var voter db.Voter

	if err := c.BodyParser(&voter); err != nil {
		log.Println("Error binding JSON: ", err)
		return fiber.NewError(http.StatusBadRequest)
	}

	if voter.VoterId != param.ID {
		log.Println("Voter does not match id parameter.")
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := voter.ValidateVotes(); err != nil {
		log.Println("Voter is not valid.")
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.UpdateVoter(voter); err != nil {
		log.Println("Error updating voter: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.JSON(voter)
}

func (v *VoterAPI) DeleteVoter(c *fiber.Ctx) error {
	param := struct {
		ID uint `params:"id"`
	}{}

	if err := c.ParamsParser(&param); err != nil {
		return fiber.NewError(http.StatusBadRequest)
	}

	if err := v.db.DeleteVoter(param.ID); err != nil {
		log.Println("Error deleting voter: ", err)
		return fiber.NewError(http.StatusInternalServerError)
	}

	return c.Status(http.StatusOK).SendString("Delete OK")
}

func (v *VoterAPI) GetAllVotes(c *fiber.Ctx) error {
	return v.withVoter(c, func(voter db.Voter) error {
		return c.JSON(voter.VoteHistory)
	})
}

func (v *VoterAPI) AddVote(c *fiber.Ctx) error {
	return v.withVoter(c, func(voter db.Voter) error {
		var vote db.VoterHistory

		if err := c.BodyParser(&vote); err != nil {
			log.Println("Error binding JSON: ", err)
			return fiber.NewError(http.StatusBadRequest)
		}

		if voter.VoterId != vote.VoterId {
			log.Println("Voter id does not match voter_id.")
			return fiber.NewError(http.StatusBadRequest)
		}

		voter, err := voter.AddVote(vote)
		if err != nil {
			log.Println("Error adding vote to voter: ", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		if err := v.db.UpdateVoter(voter); err != nil {
			log.Println("Error updating voter with vote: ", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		return c.JSON(vote)
	})
}

func (v *VoterAPI) GetVote(c *fiber.Ctx) error {
	return v.withVoter(c, func(voter db.Voter) error {
		param := struct {
			ID uint `params:"pollid"`
		}{}

		if err := c.ParamsParser(&param); err != nil {
			return fiber.NewError(http.StatusBadRequest)
		}

		vote, err := voter.GetVote(param.ID)
		if err != nil {
			log.Println("Error getting vote: ", err)
			return fiber.NewError(http.StatusNotFound)
		}

		return c.JSON(vote)
	})
}

func (v *VoterAPI) UpdateVote(c *fiber.Ctx) error {
	return v.withVoter(c, func(voter db.Voter) error {
		param := struct {
			ID uint `params:"pollid"`
		}{}

		if err := c.ParamsParser(&param); err != nil {
			return fiber.NewError(http.StatusBadRequest)
		}

		var vote db.VoterHistory

		if err := c.BodyParser(&vote); err != nil {
			log.Println("Error binding JSON: ", err)
			return fiber.NewError(http.StatusBadRequest)
		}

		if vote.PollId != param.ID {
			log.Println("Poll Id does not match pollid parameter")
			return fiber.NewError(http.StatusBadRequest)
		}

		if voter.VoterId != vote.VoterId {
			log.Println("Voter does not match voter_id.")
			return fiber.NewError(http.StatusBadRequest)
		}

		voter, err := voter.UpdateVote(vote)
		if err != nil {
			log.Println("Error updating vote: ", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		if err := v.db.UpdateVoter(voter); err != nil {
			log.Println("Error updating voter: ", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		return c.JSON(vote)
	})
}

func (v *VoterAPI) DeleteVote(c *fiber.Ctx) error {
	return v.withVoter(c, func(voter db.Voter) error {
		param := struct {
			ID uint `params:"pollid"`
		}{}

		if err := c.ParamsParser(&param); err != nil {
			return fiber.NewError(http.StatusBadRequest)
		}

		voter, err := voter.DeleteVote(param.ID)
		if err != nil {
			log.Println("Error deleting vote: ", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		if err := v.db.UpdateVoter(voter); err != nil {
			log.Println("Error updating voter: ", err)
			return fiber.NewError(http.StatusInternalServerError)
		}

		return c.Status(http.StatusOK).SendString("Delete OK")

	})
}

// implementation of GET /health. It is a good practice to build in a
// health check for your API.  Below the results are just hard coded
// but in a real API you can provide detailed information about the
// health of your API with a Health Check
func (v *VoterAPI) HealthCheck(c *fiber.Ctx) error {
	return c.Status(http.StatusOK).
		JSON(fiber.Map{
			"status":             "ok",
			"version":            "1.0.0",
			"uptime":             time.Now().Sub(v.bootTime),
			"total_calls":        v.totalCalls,
			"errors_encountered": v.errors,
		})
}
