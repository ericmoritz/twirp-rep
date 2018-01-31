package usersservice_test

import (
	"context"
	"testing"
	. "github.com/franela/goblin"
	pb "github.com/ericmoritz/twirp-rep/rpc/rep"
	"github.com/ericmoritz/twirp-rep/internal/repservice"
	usersPb "github.com/ericmoritz/twirp-users/rpc/users"
	"os"
	"net/http"
)

// Test tests the server
func Test(t *testing.T) {
	g := Goblin(t)

	g.Describe("Users API", func() {
		var service pb.Rep
		testDbPath := "/tmp/rep.db"
		usersClient := usersPb.NewUsersProtobufClient(
			"http://localhost:8080",
			&http.Client{},
		)

		g.Before(func() {
			// Delete the db if it exists
			if _, err := os.Stat(testDbPath); err == nil {
				if err := os.RemoveAll(testDbPath); err != nil {
					panic(err)
				}
			}

			if s, err := repservice.New(testDbPath, usersClient); err == nil {
				service = s
			} else {
				panic(err)
			}
		})

		g.It("Happy Case", func() {
			// Register an eric user
			_, err := usersClient.Register(
				context.Background(),
				&usersPb.RegisterReq{Username: "eric", Password: "Shhh"},
			)
			g.Assert(err).Equal(nil)

			// Register a jesiah user
			_, err = usersClient.Register(
				context.Background(),
				&usersPb.RegisterReq{Username: "jesiah", Password: "for the horde"},
			)
			g.Assert(err).Equal(nil)


			// Login with the eric user
			loginResp, err := usersClient.Login(
				context.Background(),
				&usersPb.LoginReq{Username: "eric", Password: "Shhh"},
			)
			g.Assert(err).Equal(nil)

			// Store the session
			session := loginResp.Session

			// Get the current votes for jesiah
			startVotesResp, err := service.Votes(
				context.Background(),
				&pb.VotesReq{
					Username: "jesiah",
				},
			)
			g.Assert(err).Equal(nil)

			// Vote for a Jesiah
			_, err = service.Vote(
				context.Background(),
				&pb.VoteReq{
					Session: session,
					Target: "jesiah",
				},
			)
			g.Assert(err).Equal(nil)

			// Get the more current votes for jesiah
			endVotesResp, err := service.Votes(
				context.Background(),
				&pb.VotesReq{
					Username: "jesiah",
				},
			)
			g.Assert(err).Equal(nil)
			g.Assert(startVotesResp.Count+1).Equal(endVotesResp.Count)
		})



	})
}
