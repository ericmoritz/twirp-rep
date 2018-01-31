package main

import (
	"github.com/ericmoritz/twirp-rep/internal/repservice"
	pb "github.com/ericmoritz/twirp-rep/rpc/rep"
	usersRPC "github.com/ericmoritz/twirp-users/rpc/users"
	"net/http"
	"os"
	"fmt"
)

func main() {
	var bind = ":8081"
	if port := os.Getenv("PORT"); port != "" {
		bind = ":"+port
	}

	var userServiceURL = "http://localhost:8080"
	if url := os.Getenv("USER_SERVICE_URL"); url != "" {
		userServiceURL = url
	}
	usersClient := usersRPC.NewUsersProtobufClient(
		userServiceURL,
		&http.Client{},
	)

	server, err := repservice.New(
		"./.rep.db",
		usersClient,
	)
	if err != nil {
		panic(err)
	}

	handler := pb.NewRepServer(server, nil)
	fmt.Printf("Listening on %s\n", bind)
	err = http.ListenAndServe(bind, handler)
	if err != nil {
		panic(err)
	}
}
