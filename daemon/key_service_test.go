package daemon

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"github.com/yankeguo/bastion/types"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"testing"
)

func TestDaemon_CreateListDeleteKeys(t *testing.T) {
	v, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	p, err := ssh.NewPublicKey(v.Public())
	if err != nil {
		t.Fatal(err)
	}
	v2, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		t.Fatal(err)
	}
	p2, err := ssh.NewPublicKey(v2.Public())
	if err != nil {
		t.Fatal(err)
	}
	fp := ssh.FingerprintSHA256(p)
	fp2 := ssh.FingerprintSHA256(p2)

	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		s := types.NewKeyServiceClient(conn)

		res, err := s.CreateKey(context.Background(), &types.CreateKeyRequest{
			Fingerprint: fp,
			Account:     "test",
			Name:        "hello",
		})

		if err != nil {
			t.Fatal(err)
		}
		if res.Key.Account != "test" {
			t.Fatal("failed 1")
		}

		t.Log(res)

		res, err = s.CreateKey(context.Background(), &types.CreateKeyRequest{
			Fingerprint: fp,
			Account:     "test",
			Name:        "hello",
		})

		if err == nil {
			t.Fatal("failed duplicated")
		}

		s.CreateKey(context.Background(), &types.CreateKeyRequest{
			Fingerprint: fp2,
			Account:     "test",
			Name:        "hello",
		})

		res1, err := s.ListKeys(context.Background(), &types.ListKeysRequest{Account: "test"})

		if err != nil {
			t.Fatal(err)
		}

		if len(res1.Keys) != 2 {
			t.Fatal("failed 2")
		}

		if res1.Keys[0].Fingerprint != fp && res1.Keys[1].Fingerprint != fp {
			t.Fatal("failed 3")
		}

		t.Log(res1)

		res2, err := s.GetKey(context.Background(), &types.GetKeyRequest{Fingerprint: fp})
		if err != nil {
			t.Fatal(err)
		}
		if res2.Key.Account != "test" {
			t.Fatal("failed")
		}

		t.Log(res2)

		_, err = s.DeleteKey(context.Background(), &types.DeleteKeyRequest{Fingerprint: fp})
		if err != nil {
			t.Fatal(err)
		}

		res1, err = s.ListKeys(context.Background(), &types.ListKeysRequest{Account: "test"})

		if err != nil {
			t.Fatal(err)
		}

		if len(res1.Keys) != 1 {
			t.Fatal("failed 2")
		}

		if res1.Keys[0].Fingerprint != fp2 {
			t.Fatal("failed 3")
		}

	})
}
