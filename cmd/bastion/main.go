package main

import (
	"github.com/urfave/cli"
	"os"
	"log"
	"google.golang.org/grpc"
	"github.com/yankeguo/bastion/types"
	"context"
)

func newConnection(c *cli.Context) (conn *grpc.ClientConn, err error) {
	if conn, err = grpc.Dial(c.GlobalString("endpoint"), grpc.WithInsecure()); err != nil {
		return
	}
	return
}

func main() {
	// clear date time from log
	log.SetFlags(0)
	// build the app
	app := cli.NewApp()
	app.Name = "bastion"
	app.Description = "bastion command-line interface"
	app.Author = "Yanke Guo <guoyk.cn@gmail.com>"
	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "endpoint", Usage: "bastiond rpc address", Value: "127.0.0.1:9777"},
	}
	app.Commands = []cli.Command{
		{
			Name:  "users",
			Usage: "user related commands",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create a user",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "account", Usage: "account name of the user"},
						cli.StringFlag{Name: "password", Usage: "password of the user"},
						cli.StringFlag{Name: "nickname", Usage: "nickname of the user"},
						cli.BoolFlag{Name: "admin", Usage: "if the user is a admin"},
					},
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						us := types.NewUserServiceClient(conn)
						res, err := us.CreateUser(context.Background(), &types.CreateUserRequest{
							Account:  c.String("account"),
							Password: c.String("password"),
							Nickname: c.String("nickname"),
							IsAdmin:  c.Bool("admin"),
						})
						if err != nil {
							return err
						}
						log.Println(res.User)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list all users",
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						us := types.NewUserServiceClient(conn)
						res, err := us.ListUsers(context.Background(), &types.ListUsersRequest{})
						if err != nil {
							return err
						}
						for _, u := range res.Users {
							log.Println(u)
						}
						return nil
					},
				},
			},
		},
		{
			Name:  "nodes",
			Usage: "node related commands",
			Subcommands: []cli.Command{
				{
					Name:  "list",
					Usage: "list all nodes",
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						ns := types.NewNodeServiceClient(conn)
						res, err := ns.ListNodes(context.Background(), &types.ListNodesRequest{})
						if err != nil {
							return err
						}
						for _, n := range res.Nodes {
							log.Println(n)
						}
						return nil
					},
				},
			},
		},
	}
	// run the app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
