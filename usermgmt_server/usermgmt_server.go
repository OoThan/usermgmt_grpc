package main

import (
	"context"
	pb "example.com/go-usermgmt_grpc/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
)

const (
	port = ":50051"
)

func NewUserManagementServer() *UserManagementServer {
	return &UserManagementServer{}
}

type UserManagementServer struct {
	pb.UnimplementedUserManagementServer
}

func (server *UserManagementServer) Run() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserManagementServer(s, server)
	log.Printf("Server listening at %v", lis.Addr())
	return s.Serve(lis)
}

func (server *UserManagementServer) CreateNewUser(ctx context.Context, in *pb.NewUser) (*pb.User, error) {
	log.Printf("Received : %v", in.GetName())
	readBytes, err := ioutil.ReadFile("users.json")
	var userList *pb.UserList = &pb.UserList{}
	var userId int32 = int32(rand.Intn(10000))
	createdUser := &pb.User{
		Name: in.GetName(),
		Age:  in.GetAge(),
		Id:   userId,
	}

	if err != nil {
		if os.IsNotExist(err) {
			log.Print("File not found. Creating a new file")
			userList.Users = append(userList.Users, createdUser)
			jsonBytes, err := protojson.Marshal(userList)
			if err != nil {
				log.Fatalf("JSON Marshaling failed: %v", err)
			}
			if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
				log.Fatalf("Failed write to file: %v", err)
			}
			return createdUser, nil
		} else {
			log.Fatalf("Error Reading file : %v", err)
		}
	}

	if err := protojson.Unmarshal(readBytes, userList); err != nil {
		log.Fatalf("Failed to pass user list: %v", err)
	}
	userList.Users = append(userList.Users, createdUser)
	jsonBytes, err := protojson.Marshal(userList)
	if err != nil {
		log.Fatalf("JSON Marchaling failed: %v", err)
	}
	if err := ioutil.WriteFile("users.json", jsonBytes, 0664); err != nil {
		log.Fatalf("Failed write to file : %v", err)
	}
	return createdUser, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	jsonBytes, err := ioutil.ReadFile("users.json")
	if err != nil {
		log.Fatalf("Failed read from file: %v", err)
	}
	var userList *pb.UserList = &pb.UserList{}
	if err := protojson.Unmarshal(jsonBytes, userList); err != nil {
		log.Fatalf("Unmarshaling failed : %v", err)
	}
	return userList, nil
}

func main() {
	var userManagementServer *UserManagementServer = NewUserManagementServer()
	if err := userManagementServer.Run(); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}
}
