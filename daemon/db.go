package daemon

import (
	"github.com/asdine/storm"
	"github.com/asdine/storm/index"
	"github.com/asdine/storm/q"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	errRecordNotFound = status.Error(codes.InvalidArgument, "record not found")
	errInternal       = status.Error(codes.Internal, "internal error")
)

func errDuplicatedField(key string) error {
	return status.Errorf(codes.AlreadyExists, "duplicated field '%s'", key)
}

func transformStormError(err *error) {
	if err == nil {
	} else if *err == nil {
	} else if *err == storm.ErrNotFound {
		*err = errRecordNotFound
	} else {
		*err = status.Error(codes.Internal, (*err).Error())
	}
	return
}

// DB wraps storm.Node, changes some behaviors when error returned
type DB struct {
	db *storm.DB
}

type Node struct {
	node storm.Node
}

type Query struct {
	query storm.Query
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) Tx(writable bool, cb func(*Node) error) (err error) {
	var node storm.Node
	if node, err = d.db.Begin(writable); err != nil {
		transformStormError(&err)
		return
	}
	defer node.Rollback()
	if err = cb(&Node{node: node}); err != nil {
		return
	}
	if writable {
		if err = node.Commit(); err != nil {
			transformStormError(&err)
			return
		}
	}
	return
}

func (d *DB) One(fieldName string, value interface{}, to interface{}) (err error) {
	err = d.db.One(fieldName, value, to)
	transformStormError(&err)
	return
}

func (d *DB) All(to interface{}, options ...func(*index.Options)) (err error) {
	err = d.db.All(to, options...)
	if err == storm.ErrNotFound {
		err = nil
	} else {
		transformStormError(&err)
	}
	return
}

func (d *DB) Find(fieldName string, val interface{}, out interface{}, options ...func(q *index.Options)) (err error) {
	err = d.db.Find(fieldName, val, out, options...)
	if err == storm.ErrNotFound {
		err = nil
	} else {
		transformStormError(&err)
	}
	return
}

func (d *DB) Save(data interface{}) (err error) {
	err = d.db.Save(data)
	transformStormError(&err)
	return
}

func (d *DB) DeleteStruct(data interface{}) (err error) {
	d.db.DeleteStruct(data)
	transformStormError(&err)
	return
}
func (d *DB) CheckDuplicated(bucket string, keyName string, keyVal interface{}) (err error) {
	var exists bool
	if exists, err = d.db.KeyExists(bucket, keyVal); err != nil {
		transformStormError(&err)
		return
	}
	if exists {
		err = errDuplicatedField(keyName)
		return
	}
	return
}

func (d *DB) Count(data interface{}) (c int, err error) {
	c, err = d.db.Count(data)
	transformStormError(&err)
	return
}

func (d *Node) One(fieldName string, value interface{}, to interface{}) (err error) {
	err = d.node.One(fieldName, value, to)
	transformStormError(&err)
	return
}

func (d *Node) All(to interface{}, options ...func(*index.Options)) (err error) {
	err = d.node.All(to, options...)
	if err == storm.ErrNotFound {
		err = nil
	} else {
		transformStormError(&err)
	}
	return
}

func (d *Node) Find(fieldName string, val interface{}, out interface{}, options ...func(q *index.Options)) (err error) {
	err = d.node.Find(fieldName, val, out, options...)
	if err == storm.ErrNotFound {
		err = nil
	} else {
		transformStormError(&err)
	}
	return
}

func (d *Node) Save(data interface{}) (err error) {
	err = d.node.Save(data)
	transformStormError(&err)
	return
}

func (d *Node) DeleteStruct(data interface{}) (err error) {
	d.node.DeleteStruct(data)
	transformStormError(&err)
	return
}

func (d *Node) CheckDuplicated(bucket string, keyName string, keyVal interface{}) (err error) {
	var exists bool
	if exists, err = d.node.KeyExists(bucket, keyVal); err != nil {
		transformStormError(&err)
		return
	}
	if exists {
		err = errDuplicatedField(keyName)
		return
	}
	return
}

func (d *Node) Count(data interface{}) (c int, err error) {
	c, err = d.node.Count(data)
	transformStormError(&err)
	return
}

func (d *Node) Select(matchers ...q.Matcher) *Query {
	return &Query{
		query: d.node.Select(matchers...),
	}
}

func (q *Query) Count(data interface{}) (i int, err error) {
	i, err = q.query.Count(data)
	transformStormError(&err)
	return
}

func (q *Query) Delete(data interface{}) (err error) {
	err = q.query.Delete(data)
	transformStormError(&err)
	return
}
