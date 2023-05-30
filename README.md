# Let's implement DDIA in golang

> “What I cannot create, I do not understand.” – Richard Feynman

This repository contains code implementations for 'Designing Data-Intensive Applications (DDIA)' by Martin Kleppmann in
Go. The project aims to provide practical examples and hands-on exercises based on the concepts and
principles discussed in the book, showcasing solutions in multiple programming languages.

It serves as a valuable companion for those who want to delve deeper into the world of designing and building scalable
and maintainable data-intensive applications across different language ecosystems.

## Project status

### Algorithms

* [K Way Merging File](./merging-k-sorted-list)
* [Rolling Hash (Rabin-Karp Algorithm)](./algorithms/rabin-karp)
* [Cuckoo Hashing](./algorithms/cuckoo-hashing)
* [Consistent Hash(Hash Ring)](./hash-ring)
* [Bloom Filter](./bloom-filter)
* [Boyer–Moore majority vote algorithm](./algorithms/boyer-moore-majority)
* [Count-Min Sketch](./algorithms/count-min-sketch)
* [HyperLogLog](./algorithms/hyperloglog)
* [TDigest](./t-digest)
* [Run Length Encoding](./algorithms/run-length-encoding)
* [SimHash](./algorithms/simhash)
* [Roaring Bitmap](./algorithms/roaring-bitmap)

### Storage

* [Bitcask](./bitcask)
* [LSM](./lsm)
* [Inverted Index](./inverted-index)

### BTree

* [B-Tree in mem](./btree/b-tree-mem)
* [B-Tree on Disk](./btree/b-tree-on-disk)
* B+Tree in progress

### SkipList

* [SkipList](./skiplist)

### Commons

* [byteorder](./commons-io/byteorder)  provide rust like byteorder crate for golang
* [file_utils](./commons-io/file_utils) provide Java commons-io like file_utils for golang

### Upcoming

* Gossip
* Raft
* Wal
* Delay queue(Rust done, Go in progress)
* Ntp
* Sql parser
* Column-storage

If you are interested in this project or would like me to prioritize implementing a certain component, please don't
hesitate to submit an issue. If you encounter any problems with the code, you are also welcome to submit a PR with the
appropriate modifications.
