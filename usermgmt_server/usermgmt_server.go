package main

import (
	"context"
	pb "example.com/go-usermgmt_grpc/usermgmt"
	"fmt"
	"github.com/jackc/pgx/v5"
	"google.golang.org/grpc"
	"log"
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
	conn *pgx.Conn
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
	createSql := `
	create table if not exists users (
	    id SERIAL PRIMARY KEY,
		name text,
		age int
	);
	`
	_, err := server.conn.Exec(context.Background(), createSql)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Table creation failed  : %v\n", err)
		os.Exit(1)
	}

	createdUser := &pb.User{Name: in.GetName(), Age: in.GetAge()}
	tx, err := server.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("conn.Begin Failed: %v", err)
	}
	_, err = tx.Exec(context.Background(), "insert into users(name, age) values ($1, $2)", createdUser.GetName(), createdUser.GetAge())
	if err != nil {
		log.Fatalf("tx.Exec failed : %v", err)
	}
	tx.Commit(context.Background())
	return createdUser, nil
}

func (server *UserManagementServer) GetUsers(ctx context.Context, in *pb.GetUsersParams) (*pb.UserList, error) {
	var userList *pb.UserList = &pb.UserList{}
	rows, err := server.conn.Query(context.Background(), "select * from users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		user := pb.User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		userList.Users = append(userList.Users, &user)
	}

	return userList, nil
}

func main() {
	databaseURL := "postgres://Luke:UpTech20@@08@localhost:5432/postgres"
	conn, err := pgx.Connect(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Unable to establish connection: %v", err)
	}
	defer conn.Close(context.Background())

	var userManagementServer *UserManagementServer = NewUserManagementServer()
	userManagementServer.conn = conn
	if err := userManagementServer.Run(); err != nil {
		log.Fatalf("Failed to serve : %v", err)
	}
}
