package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"time"

	pb "example.com/go-usermgmt_grpc/usermgmt"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect : %v", err)
	}
	defer conn.Close()

	c := pb.NewUserManagementClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var newUsers = make(map[string]int32)
	newUsers["Alice"] = 24
	newUsers["Luke"] = 24
	for name, age := range newUsers {
		r, err := c.CreateNewUser(ctx, &pb.NewUser{Name: name, Age: age})
		if err != nil {
			log.Fatalf("Could not create user: %v", err)
		}
		log.Printf(`UserDetails:
			Name: %s
			Age: %d
			Id: %d`, r.GetName(), r.GetAge(), r.GetId())
		rr, err := c.GetUsers(ctx, &pb.GetUsersParams{})
		if err != nil {
			log.Fatalf("Could not retriece users: %v", err)
		}
		log.Print("\nUser List: \n")
		log.Printf("r.GetUsers(): %v\n", rr.GetUsers())
	}
}
