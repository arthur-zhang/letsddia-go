package lsm

import (
	"fmt"
	"skiplist"
	"testing"
)

func TestMultiIterar(t *testing.T) {
	sl1 := skiplist.NewSkipList[Key, Value]()
	sl2 := skiplist.NewSkipList[Key, Value]()
	sl3 := skiplist.NewSkipList[Key, Value]()
	for i := 0; i < 10; i++ {
		str := Key(fmt.Sprintf("%2d", i))
		kv := KVPair{
			SeqId: 0,
			Op:    0,
			Key:   str,
			Value: Value{byte(i)},
		}
		kvByte := kv.encode()
		if i%3 == 0 {
			sl1.Insert(&str, (*Value)(&kvByte))
		} else if i%3 == 1 {
			sl2.Insert(&str, (*Value)(&kvByte))
		} else {
			sl3.Insert(&str, (*Value)(&kvByte))
		}
	}
	sl1.Display()

	iters := []Iterator[Key, Value]{sl1.Iterator(), sl2.Iterator()}
	multi := NewMultiIterator(iters)
	multi.SeekToFirst()
	for multi.Valid() {
		key := *multi.Key()
		t.Logf("%s\n", key)
		multi.Next()
	}
}
