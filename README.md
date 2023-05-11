# Let's implement DDIA in golang

> “What I cannot create, I do not understand.” – Richard Feynman

This repository contains code implementations for 'Designing Data-Intensive Applications (DDIA)' by Martin Kleppmann in
Go. The project aims to provide practical examples and hands-on exercises based on the concepts and
principles discussed in the book, showcasing solutions in multiple programming languages.

It serves as a valuable companion for those who want to delve deeper into the world of designing and building scalable
and maintainable data-intensive applications across different language ecosystems.

## project status

### BTree

* [B-Tree in mem](./btree/b-tree-mem)
* [B-Tree on Disk](./btree/b-tree-on-disk)
* B+Tree in progress

## SkipList

* [SkipList](./skiplist)

Other projects are in progress.


| project | status | info |
  |------------------|---------------|----------------------------------------------------------------------|
  | bitcask | rust ✅, go ❌ | currently no multi-thread supported |
  | delay-queue | partial done | rust impl according to beanstalkd, based on tokio and priority queue |
  | lsm | working hard | |
  | gossip | working hard | | |
  | raft | working hard | | |
  | binlog | working hard | | |
  | ntp | working hard | | |
  | column-storage | working hard | | |

If you are interested in this project or would like me to prioritize implementing a certain component, please don't
hesitate to submit an issue. If you encounter any problems with the code, you are also welcome to submit a PR with the
appropriate modifications.
