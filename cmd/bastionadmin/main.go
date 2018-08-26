package main

import (
	"context"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"io/ioutil"
	"log"
	"os"
)

func newConnection(c *cli.Context) (conn *grpc.ClientConn, err error) {
	if conn, err = grpc.Dial(c.GlobalString("endpoint"), grpc.WithInsecure()); err != nil {
		return
	}
	return
}

func main() {
	// setup log
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
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
					Name:  "unblock",
					Usage: "unblock a user",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "account", Usage: "account name of the user"},
					},
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						us := types.NewUserServiceClient(conn)
						res, err := us.UpdateUser(context.Background(), &types.UpdateUserRequest{
							Account:         c.String("account"),
							UpdateIsBlocked: true,
							IsBlocked:       false,
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
				{
					Name:  "list-keys",
					Usage: "list keys of a user",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "account", Usage: "account of the user"},
					},
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						ks := types.NewKeyServiceClient(conn)
						res, err := ks.ListKeys(context.Background(), &types.ListKeysRequest{Account: c.String("account")})
						if err != nil {
							return err
						}
						for _, k := range res.Keys {
							log.Println(k)
						}
						return nil
					},
				},
				{
					Name:  "add-key",
					Usage: "add key to a user",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "account", Usage: "account of the user"},
						cli.StringFlag{Name: "name", Usage: "name of the key"},
						cli.StringFlag{Name: "fingerprint", Usage: "fingerprint of the key"},
						cli.StringFlag{Name: "file", Usage: "the public key file"},
					},
					Action: func(c *cli.Context) error {
						fp := c.String("fingerprint")
						file := c.String("file")
						if len(file) > 0 {
							buf, err := ioutil.ReadFile(file)
							if err != nil {
								return err
							}
							pk, _, _, _, err := ssh.ParseAuthorizedKey(buf)
							if err != nil {
								return err
							}
							fp = ssh.FingerprintSHA256(pk)
						}
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						ks := types.NewKeyServiceClient(conn)
						res, err := ks.CreateKey(context.Background(), &types.CreateKeyRequest{
							Account:     c.String("account"),
							Name:        c.String("name"),
							Fingerprint: fp,
						})
						if err != nil {
							return err
						}
						log.Println(res.Key)
						return nil
					},
				},
				{
					Name:  "delete-key",
					Usage: "delete a key",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "fingerprint", Usage: "fingerprint of the key"},
					},
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						ks := types.NewKeyServiceClient(conn)
						_, err = ks.DeleteKey(context.Background(), &types.DeleteKeyRequest{
							Fingerprint: c.String("fingerprint"),
						})
						if err != nil {
							return err
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
				{
					Name:  "create",
					Usage: "create a node manually",
					Flags: []cli.Flag{
						cli.StringFlag{Name: "hostname", Usage: "hostname of node"},
						cli.StringFlag{Name: "address", Usage: "address of the node, default port is 22"},
						cli.StringFlag{Name: "user", Usage: "ssh user will be used in bastion", Value: types.NodeUserRoot},
					},
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						ns := types.NewNodeServiceClient(conn)
						res, err := ns.PutNode(context.Background(), &types.PutNodeRequest{
							Hostname: c.String("hostname"),
							Address:  c.String("address"),
							User:     c.String("user"),
							Source:   types.NodeSourceManual,
						})
						log.Println(res.Node)
						return nil
					},
				},
			},
		},
		{
			Name:  "sessions",
			Usage: "session related commands",
			Subcommands: []cli.Command{
				{
					Name:  "seed",
					Usage: "create 111 dummy sessions for testing purpose",
					Action: func(c *cli.Context) error {
						conn, err := newConnection(c)
						if err != nil {
							return err
						}
						defer conn.Close()
						us := types.NewUserServiceClient(conn)
						res1, err := us.ListUsers(context.Background(), &types.ListUsersRequest{})
						if err != nil {
							return err
						}
						if len(res1.Users) == 0 {
							return errors.New("no users")
						}
						ss := types.NewSessionServiceClient(conn)
						for i := 0; i < 111; i++ {
							u := res1.Users[i%len(res1.Users)]
							res2, err := ss.CreateSession(context.Background(), &types.CreateSessionRequest{
								Account:    u.Account,
								Command:    "dummy_command",
								IsRecorded: i%2 == 0,
							})
							if err != nil {
								return err
							}
							ss.FinishSession(context.Background(), &types.FinishSessionRequest{
								Id: res2.Session.Id,
							})
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
