package standalone_storage

import (
	"path/filepath"

	"github.com/Connor1996/badger"
	"github.com/pingcap-incubator/tinykv/kv/config"
	"github.com/pingcap-incubator/tinykv/kv/storage"
	"github.com/pingcap-incubator/tinykv/kv/util/engine_util"
	"github.com/pingcap-incubator/tinykv/proto/pkg/kvrpcpb"
)

// StandAloneStorage is an implementation of `Storage` for a single-node TinyKV instance. It does not
// communicate with other nodes and all data is stored locally.
type StandAloneStorage struct {
	// Your Data Here (1).
	engine *engine_util.Engines
	config *config.Config
}

func NewStandAloneStorage(conf *config.Config) *StandAloneStorage {
	// Your Code Here (1).
	dbPath := conf.DBPath
	kvPath := filepath.Join(dbPath, "kv")
	raftPath := filepath.Join(dbPath, "raft")

	kvDB := engine_util.CreateDB("kv", false)
	raftDB := engine_util.CreateDB("raft", false)
	return &StandAloneStorage{
		engine: engine_util.NewEngines(kvDB, raftDB, kvPath, raftPath),
		config: conf,
	}
}

func (s *StandAloneStorage) Start() error {
	// Your Code Here (1).
	return nil
}

func (s *StandAloneStorage) Stop() error {
	// Your Code Here (1).
	return s.engine.Close()
}

func (s *StandAloneStorage) Reader(ctx *kvrpcpb.Context) (storage.StorageReader, error) {
	// Your Code Here (1).
	txn := s.engine.Kv.NewTransaction(false)
	return &StandAloneReader{
		txn: txn,
	}, nil
}

func (s *StandAloneStorage) Write(ctx *kvrpcpb.Context, batch []storage.Modify) error {
	// Your Code Here (1).
	writeBatch := &engine_util.WriteBatch{}
	for _, m := range batch {
		switch m.Data.(type) {
		case storage.Put:
			writeBatch.SetCF(m.Cf(), m.Key(), m.Value())
		case storage.Delete:
			writeBatch.DeleteCF(m.Cf(), m.Key())
		}
	}
	return s.engine.WriteKV(writeBatch)
}

type StandAloneReader struct {
	txn *badger.Txn
}

func (reader *StandAloneReader) GetCF(cf string, key []byte) ([]byte, error) {
	txn := reader.txn
	val, err := engine_util.GetCFFromTxn(txn, cf, key)
	if err == badger.ErrKeyNotFound {
		return nil, nil
	}
	return val, err
}

func (reader *StandAloneReader) IterCF(cf string) engine_util.DBIterator {
	txn := reader.txn
	return engine_util.NewCFIterator(cf, txn)
}

func (reader *StandAloneReader) Close() {
	reader.txn.Discard()
}
