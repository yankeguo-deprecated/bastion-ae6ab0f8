package main

import (
	"context"
	"flag"
	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"time"
)

var (
	endpoint string
	verbose  bool

	lastIndex uint64
)

func main() {
	flag.StringVar(&endpoint, "endpoint", "127.0.0.1:9777", "endpoint address of bunkerd")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose mode")
	flag.Parse()

	for {
		if err := update(); err != nil {
			time.Sleep(time.Second * 10)
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
	// find hostnames to remove
	removes := make([]string, 0)
	for _, bn := range lnr.Nodes {
		// ignore non-consul hosts
		if bn.Source != types.NodeSourceConsul {
			continue
		}
		// check if it's existed in consul catalog
		var found bool
		for _, cn := range cns {
			if cn.Node == bn.Hostname {
				found = true
				break
			}
		}
		// not found, mark to remove
		if !found {
			removes = append(removes, bn.Hostname)
		}
	}
	// add all hosts from consul catalog
	for _, cn := range cns {
		if verbose {
			log.Print("will add:", cn.Node, cn.Address)
		}
		if _, err = ns.PutNode(context.Background(), &types.PutNodeRequest{
			Hostname: cn.Node,
			User:     types.NodeUserRoot,
			Address:  cn.Address,
			Source:   types.NodeSourceConsul,
		}); err != nil {
			return
		}
	}
	// delete all hosts not existed any more
	for _, n := range removes {
		if verbose {
			log.Print("will remove:", n)
		}
		if _, err = ns.DeleteNode(context.Background(), &types.DeleteNodeRequest{
			Hostname: n,
		}); err != nil {
			return
		}
	}
	return
}
