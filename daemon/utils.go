package daemon

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"time"

	"github.com/olivere/elastic"
	"github.com/rs/zerolog/log"
	"github.com/yankeguo/bastion/types"
	"github.com/yankeguo/bastion/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
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

type ReplayIndice struct {
	SessionId int64     `json:"session_id"`
	Timestamp uint32    `json:"timestamp"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ReplaySubmitter struct {
	SessionId int64
	CreatedAt time.Time
	Index     string
	EsClient  *elastic.Client
	timestamp uint32
	cache     string
	batch     []ReplayIndice
}

func NewReplaySubmitter(createdAt time.Time, sessionId int64, esClient *elastic.Client) (r *ReplaySubmitter) {
	r = &ReplaySubmitter{
		SessionId: sessionId,
		CreatedAt: createdAt,
		Index:     fmt.Sprintf("%s%04d-%02d-%02d", types.ReplayElasticsearchIndexPrefix, createdAt.Year(), createdAt.Month(), createdAt.Day()),
		EsClient:  esClient,
		batch:     []ReplayIndice{},
	}
	return
}

func (r *ReplaySubmitter) submitBatch() (err error) {
	if len(r.batch) == 0 {
		return
	}
	bulk := r.EsClient.Bulk()
	for _, d := range r.batch {
		bulk = bulk.Add(elastic.NewBulkIndexRequest().Index(r.Index).Type("_doc").Doc(d))
	}
	_, err = bulk.Do(context.Background())
	r.batch = []ReplayIndice{}
	return
}

func (r *ReplaySubmitter) Add(f types.ReplayFrame) (err error) {
	if f.Type != types.ReplayFrameTypeStderr && f.Type != types.ReplayFrameTypeStdout {
		return
	}
	s := utils.ExtractReadableString(f.Payload)
	if len(s) == 0 {
		return
	}
	r.cache = r.cache + s
	if (f.Timestamp - r.timestamp) > 1000 {
		r.batch = append(r.batch, ReplayIndice{
			SessionId: r.SessionId,
			Timestamp: r.timestamp,
			Content:   r.cache,
			CreatedAt: r.CreatedAt,
		})
		r.timestamp = f.Timestamp
		r.cache = ""
		if len(r.batch) > 100 {
			if err = r.submitBatch(); err != nil {
				return
			}
		}
	}
	return
}

func (r *ReplaySubmitter) Close() (err error) {
	if len(r.cache) > 0 {
		r.batch = append(r.batch, ReplayIndice{
			SessionId: r.SessionId,
			Timestamp: r.timestamp,
			Content:   r.cache,
			CreatedAt: r.CreatedAt,
		})
	}
	if err = r.submitBatch(); err != nil {
		return
	}
	return
}
