package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"drexel.edu/voter/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Global variables to hold the command line flags to drive the todo CLI
// application
var (
	hostFlag string
	portFlag uint
)

func processCmdLineFlags() {

	flag.StringVar(&hostFlag, "h", "0.0.0.0", "Listen on all interfaces")
	flag.UintVar(&portFlag, "p", 1080, "Default Port")

	flag.Parse()
}

// main is the entry point for our todo API application.  It processes
// the command line flags and then uses the db package to perform the
// requested operation
func main() {
	processCmdLineFlags()

	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())

	apiHandler, err := api.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	//HTTP Standards for "REST" APIS
	//GET - Read/Query
	//POST - Create
	//PUT - Update
	//DELETE - Delete

	app.Use("/voters", apiHandler.HandleStats)
	app.Get("/voters", apiHandler.ListAllVoters)
	app.Post("/voters", apiHandler.AddVoter)
	app.Put("/voters/:id", apiHandler.UpdateVoter)
	app.Delete("/voters", apiHandler.DeleteAllVoters)
	app.Delete("/voters/:id", apiHandler.DeleteVoter)
	app.Get("/voters/:id", apiHandler.GetVoter)
	app.Get("/voters/:id/polls", apiHandler.GetAllVotes)
	app.Post("/voters/:id/polls", apiHandler.AddVote)
	app.Get("/voters/:id/polls/:pollid", apiHandler.GetVote)
	app.Put("/voters/:id/polls/:pollid", apiHandler.UpdateVote)
	app.Delete("/voters/:id/polls/:pollid", apiHandler.DeleteVote)

	app.Get("/health", apiHandler.HealthCheck)

	serverPath := fmt.Sprintf("%s:%d", hostFlag, portFlag)
	log.Println("Starting server on ", serverPath)
	app.Listen(serverPath)
}
