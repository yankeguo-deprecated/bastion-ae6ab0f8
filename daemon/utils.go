package daemon

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"path/filepath"
	"time"
)

func now() int64 {
	return time.Now().Unix()
}

func newToken() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		panic(err)
	}
	return hex.EncodeToString(buf)
}

func bcryptGenerate(password string) (string, error) {
	buf, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(buf), err
}

func zerologUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	log.Debug().Str("method", info.FullMethod).Interface("request", req).Interface("response", resp).Err(err).Msg("request handled")
	return
}

func zerologStreamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
	start := time.Now()
	err = handler(srv, stream)
	log.Debug().Str("method", info.FullMethod).Dur("duration", time.Now().Sub(start)).Err(err).Msg("stream request handled")
	return
}

func FilenameForSessionID(id int64, dir string) string {
	buf := make([]byte, 8, 8)
	binary.BigEndian.PutUint64(buf, uint64(id))
	name := hex.EncodeToString(buf)
	ret := make([]string, 0, 5)
	ret = append(ret, dir)
	ret = append(ret, name[:4])
	ret = append(ret, name[4:8])
	ret = append(ret, name[8:12])
	ret = append(ret, name)
	return filepath.Join(ret...)
}
