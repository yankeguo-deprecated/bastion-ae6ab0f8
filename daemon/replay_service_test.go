package daemon

import (
	"context"
	"github.com/yankeguo/bastion/types"
	"google.golang.org/grpc"
	"testing"
)

func TestFilenameForSessionID(t *testing.T) {
	if "/var/lib/bastion/replays/1234/4321/1234/1234432112344321" != FilenameForSessionID(0x1234432112344321, "/var/lib/bastion/replays") {
		t.Fatal("failed 1")
	}
	if "/var/lib/bastion/replays/0000/4321/1234/0000432112344321" != FilenameForSessionID(0x0000432112344321, "/var/lib/bastion/replays") {
		t.Fatal("failed 2")
	}
}

func TestDaemon_WriteReadReply(t *testing.T) {
	withDaemon(t, func(t *testing.T, daemon *Daemon, conn *grpc.ClientConn) {
		rs := types.NewReplayServiceClient(conn)

		s, err := rs.WriteReplay(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if err = s.Send(&types.ReplayFrame{
			SessionId: 0x123400004321,
			Type:      1,
			Payload:   []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		}); err != nil {
			t.Fatal(err)
		}
		if err = s.Send(&types.ReplayFrame{
			SessionId: 0x123400004321,
			Type:      2,
			Payload:   []byte{0x05, 0x04, 0x03, 0x02, 0x01},
		}); err != nil {
			t.Fatal(err)
		}
		if err = s.Send(&types.ReplayFrame{
			SessionId: 0x123400004321,
			Type:      3,
			Payload:   []byte{0x04, 0x03, 0x02, 0x01, 0x01, 0x02, 0x03, 0x04},
		}); err != nil {
			t.Fatal(err)
		}
		_, err = s.CloseAndRecv()
		if err != nil {
			t.Fatal(err)
		}

		s2, err := rs.ReadReplay(context.Background(), &types.ReadReplayRequest{SessionId: 0x123400004321})
		if err != nil {
			t.Fatal(err)
		}
		f, err := s2.Recv()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(f)
		if f.Type != 1 || f.SessionId != 0x123400004321 || f.Payload[4] != 0x05 {
			t.Fatal("failed 1")
		}
		f, err = s2.Recv()
		if err != nil {
			t.Fatal(err)
		}
		t.Log(f)
		if f.Type != 2 || f.SessionId != 0x123400004321 || f.Payload[4] != 0x01 {
			t.Fatal("failed 2")
		}
		f, err = s2.Recv()
		t.Log(f, err)
		if err != nil {
			t.Fatal(err)
		}
		if f.Type != 3 || f.SessionId != 0x123400004321 || f.Payload[3] != 0x01 {
			t.Fatal("failed 3")
		}
		_, err = s2.Recv()
		if err == nil {
			t.Fatal("failed 4")
		}
	})
}
