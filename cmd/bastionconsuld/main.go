package main

import (
	"context"
	"flag"
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"os"
	"time"
)

var (
	dev      bool
	endpoint string

	lastIndex uint64
)

func main() {
	var err error

	// init logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true})

	flag.StringVar(&endpoint, "endpoint", "127.0.0.1:9777", "endpoint address of bunkerd")
	flag.BoolVar(&dev, "dev", false, "enable dev mode")
	flag.Parse()

	// update logger
	if dev {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	for {
		log.Debug().Msg("update started")
		if err = update(); err != nil {
			log.Error().Err(err).Msg("update failed")
			time.Sleep(time.Second * 10)
		} else {
			log.Debug().Msg("update finished")
		}
	}
}

func update() (err error) {
	// create client
	var cc *api.Client
	if cc, err = api.NewClient(api.DefaultConfig()); err != nil {
		return
	}
	// query nodes
	var cns []*api.Node
	var mt *api.QueryMeta
	if cns, mt, err = cc.Catalog().Nodes(&api.QueryOptions{
		WaitIndex: lastIndex,
		WaitTime:  time.Second * 30,
	}); err != nil {
		return
	}
	lastIndex = mt.LastIndex
	// create grpc connection
	var bcn *grpc.ClientConn
	if bcn, err = grpc.Dial(endpoint, grpc.WithInsecure()); err != nil {
		return
	}
	defer bcn.Close()
	// create node service client
	ns := types.NewNodeServiceClient(bcn)
	// query bunkerd nodes
	var lnr *types.ListNodesResponse
	if lnr, err = ns.ListNodes(context.Background(), &types.ListNodesRequest{}); err != nil {
		return
	}
	// remove missing node
remLoop:
	for _, n := range lnr.Nodes {
		// ignore non-consul hosts
		if n.Source != types.NodeSourceConsul {
			continue
		}
		// check if it's existed in consul catalog
		for _, cn := range cns {
			if cn.Node == n.Hostname {
				continue remLoop
			}
		}
		log.Debug().Str("hostname", n.Hostname).Str("address", n.Address).Msg("remove node")
		if _, err = ns.DeleteNode(context.Background(), &types.DeleteNodeRequest{
			Hostname: n.Hostname,
		}); err != nil {
			log.Error().Str("hostname", n.Hostname).Str("address", n.Address).Err(err).Msg("failed to remove node")
			err = nil
			continue
		}
	}
	// add new node
addLoop:
	for _, cn := range cns {
		// check existed and equal
		for _, n := range lnr.Nodes {
			if n.Hostname == cn.Node && n.User == types.NodeUserRoot && n.Address == cn.Address && n.Source == types.NodeSourceConsul {
				log.Debug().Str("hostname", cn.Node).Str("address", cn.Address).Msg("synced node")
				continue addLoop
			}
		}
		log.Debug().Str("hostname", cn.Node).Str("address", cn.Address).Msg("add node")
		if _, err = ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: cn.Node,
			User:     types.NodeUserRoot,
			Address:  cn.Address,
			Source:   types.NodeSourceConsul,
		}); err != nil {
			log.Error().Str("hostname", cn.Node).Str("address", cn.Address).Err(err).Msg("failed to add node")
			err = nil
			return
		}
	}
	return
}
