package server

import (
	"context"

	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// The functions below are Server's Raw API. (implements TinyKvServer).
// Some helper methods can be found in sever.go in the current directory

// RawGet return the corresponding Get response based on RawGetRequest's CF and Key fields
func (server *Server) RawGet(_ context.Context, req *kvrpcpb.RawGetRequest) (*kvrpcpb.RawGetResponse, error) {
	// Your Code Here (1).
	rsp := new(kvrpcpb.RawGetResponse)
	reader, _ := server.storage.Reader(req.Context)
	defer reader.Close()
	val, err := reader.GetCF(req.Cf, req.Key)
	if err != nil {
		rsp.Error = err.Error()
		rsp.NotFound = true
		return rsp, err
	}
	if len(val) == 0 || val == nil {
		rsp.NotFound = true
	}
	rsp.Value = val
	return rsp, nil
}

// RawPut puts the target data into storage and returns the corresponding response
func (server *Server) RawPut(_ context.Context, req *kvrpcpb.RawPutRequest) (*kvrpcpb.RawPutResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be modified
	rsp := new(kvrpcpb.RawPutResponse)
	modify := []storage.Modify{
		{
			Data: storage.Put{
				Cf:    req.Cf,
				Key:   req.Key,
				Value: req.Value,
			},
		},
	}
	err := server.storage.Write(req.Context, modify)
	if err != nil {
		rsp.Error = err.Error()
		return rsp, err
	}
	return rsp, nil
}

// RawDelete delete the target data from storage and returns the corresponding response
func (server *Server) RawDelete(_ context.Context, req *kvrpcpb.RawDeleteRequest) (*kvrpcpb.RawDeleteResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using Storage.Modify to store data to be deleted
	rsp := new(kvrpcpb.RawDeleteResponse)
	modify := []storage.Modify{
		{
			Data: storage.Delete{
				Cf:  req.Cf,
				Key: req.Key,
			},
		},
	}
	err := server.storage.Write(req.Context, modify)
	if err != nil {
		rsp.Error = err.Error()
		return rsp, err
	}
	return rsp, nil
}

// RawScan scan the data starting from the start key up to limit. and return the corresponding result
func (server *Server) RawScan(_ context.Context, req *kvrpcpb.RawScanRequest) (*kvrpcpb.RawScanResponse, error) {
	// Your Code Here (1).
	// Hint: Consider using reader.IterCF
	rsp := new(kvrpcpb.RawScanResponse)
	reader, _ := server.storage.Reader(req.Context)
	defer reader.Close()
	iter := reader.IterCF(req.Cf)
	defer iter.Close()
	iter.Seek(req.StartKey)
	for i := uint32(0); iter.Valid() && i < req.Limit; i += 1 {
		item := iter.Item()
		key := item.Key()
		val, _ := item.Value()
		rsp.Kvs = append(rsp.Kvs,
			&kvrpcpb.KvPair{
				Key:   key,
				Value: val,
			},
		)
		iter.Next()
	}
	return rsp, nil
}
