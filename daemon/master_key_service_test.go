package daemon

import (
	"context"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"testing"
)

func TestDaemon_UpdateAllAndListMasterKeys(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		mks := types.NewMasterKeyServiceClient(conn)
		_, err := mks.UpdateAllMasterKeys(context.Background(), &types.UpdateAllMasterKeysRequest{
			MasterKeys: []*types.MasterKey{
				{Fingerprint: "a", PublicKey: "b"},
				{Fingerprint: "c", PublicKey: "d"},
				{Fingerprint: "e", PublicKey: "f"},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		res, err := mks.ListMasterKeys(context.Background(), &types.ListMasterKeysRequest{})
		if err != nil {
			t.Fatal(err)
		}
		for _, k := range res.MasterKeys {
			if k.Fingerprint == "a" {
				if k.PublicKey != "b" {
					t.Fatal("failed 1")
				}
			} else if k.Fingerprint == "c" {
				if k.PublicKey != "d" {
					t.Fatal("failed 2")
				}
			} else if k.Fingerprint == "e" {
				if k.PublicKey != "f" {
					t.Fatal("failed 3")
				}
			} else {
				t.Fatal("what ?")
			}
		}
		_, err = mks.UpdateAllMasterKeys(context.Background(), &types.UpdateAllMasterKeysRequest{
			MasterKeys: []*types.MasterKey{
				{Fingerprint: "b", PublicKey: "b"},
				{Fingerprint: "c", PublicKey: "c"},
				{Fingerprint: "f", PublicKey: "f"},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		res, err = mks.ListMasterKeys(context.Background(), &types.ListMasterKeysRequest{})
		if err != nil {
			t.Fatal(err)
		}
		for _, k := range res.MasterKeys {
			if k.Fingerprint == "b" {
				if k.PublicKey != "b" {
					t.Fatal("failed 1")
				}
			} else if k.Fingerprint == "c" {
				if k.PublicKey != "c" {
					t.Fatal("failed 2")
				}
			} else if k.Fingerprint == "f" {
				if k.PublicKey != "f" {
					t.Fatal("failed 3")
				}
			} else {
				t.Fatal("what ?")
			}
		}
	})
}
