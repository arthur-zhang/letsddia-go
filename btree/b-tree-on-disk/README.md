#### A toy persistent B-Tree implementation in Go

This is a persistent B-Tree implementation in Go. It is based on the memory B-Tree implementation
in [b-tree-mem](../b-tree-mem).

#### Details

- To simplify implementation, we use a fixed size for the key and value. The key is 10 bytes and value is 10 bytes too.
- Delete is not implemented yet.
- Copy on write version is on the way

disk layout is as follows:

```
|isLeaf 1-byte|
|numKeys 8-bytes|
|Key#0 10-bytes |
|Value#0 10-bytes |
```

#### Usage

```go
btree := NewBtree("/tmp/btree.db", 2)
arr := []byte{'F', 'S', 'Q', 'K', 'C', 'L', 'H', 'T', 'V', 'W', 'M', 'R', 'N', 'P', 'A', 'B', 'X', 'Y', 'D', 'Z', 'E'}
for _, v := range arr {
item := &Item{
key:   [N]byte{v},
value: [N]byte{v, byte('0'), byte('1')},
}
btree.Insert(item)
}
key := Item{key: [N]byte{'A'}}
item := btree.Search(&key)
assert.True(t, item != nil)
assert.Equal(t, "A", decodeString(item.key))
assert.Equal(t, "A01", decodeString(item.value))

key = Item{key: [N]byte{'a'}}
item = btree.Search(&key)
assert.True(t, item == nil)
```




