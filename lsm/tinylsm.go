package lsm

type Config struct {
	baseDir          string
	maxMemStoreSize  uint64
	blockSizeUpLimit uint64
}

const OpPut = uint8(0)
const OpDelete = uint8(1)

type TinyLsm struct {
	seqId     uint64
	diskStore *DiskStore
	memStore  *MemStore
	config    *Config
}

func New(config Config) TinyLsm {
	diskStore := NewDiskStore(config.baseDir)
	memStore := NewMemStore(&diskStore, &config)
	return TinyLsm{
		seqId:     0,
		diskStore: &diskStore,
		memStore:  &memStore,
		config:    &config,
	}
}
func (db *TinyLsm) append(op uint8, key Key, value Value) {
	db.seqId += 1
	db.memStore.Add(db.seqId, op, key, value)
}

func (db *TinyLsm) Put(key Key, value Value) {
	db.append(OpPut, key, value)
}

func (db *TinyLsm) Delete(key Key) {
	db.append(OpDelete, key, nil)
}
func (db *TinyLsm) Flush() {
	db.memStore.flush()
}

func (db *TinyLsm) Get(key *Key) (Value, bool) {
	iter := db.Scan(key, key)
	if !iter.Valid() {
		return nil, false
	}
	if *iter.Key() == *key {
		return *iter.Value(), true
	}
	return nil, false
}

func (db *TinyLsm) Close() {
	db.Flush()
	db.diskStore.Close()
}

func (db *TinyLsm) Scan(start, end *Key) MultiIterator {
	diskIter := db.diskStore.Iter()
	memIter := db.memStore.Iter()
	iter := NewMultiIterator(
		[]Iterator[Key, Value]{
			diskIter,
			memIter,
		},
	)

	iter.Seek(start)
	return iter
}
