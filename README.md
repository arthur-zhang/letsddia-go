# Let's implement DDIA in golang (with some Java and Cpp).

This repository contains code implementations for 'Designing Data-Intensive Applications (DDIA)' by Martin Kleppmann in
Rust, Go, and Java. The project aims to provide practical examples and hands-on exercises based on the concepts and
principles discussed in the book, showcasing solutions in multiple programming languages.

It serves as a valuable companion for those who want to delve deeper into the world of designing and building scalable
and maintainable data-intensive applications across different language ecosystems.

## project status

| project          | status        | info                                                                 |
|------------------|---------------|----------------------------------------------------------------------|
| bitcask          | rust ✅, go ❌  | currently no multi-thread supported                                  |
| in memory b-tree | partial  done | demo implementation of Book [Introduction to Algorithms]             |
| in memory b+tree | partial  done |                                                                      |
| on disk b+tree   | partial  done |                                                                      |
| skiplist         | ✅             |                                                                      |
| delay-queue      | partial  done | rust impl according to beanstalkd, based on tokio and priority queue |
| lsm              | working hard  |                                                                      |
| gossip           | working  hard |                                                                      |                                                                      |
| raft             | working  hard |                                                                      |                                                                      |
| binlog           | working  hard |                                                                      |                                                                      |
| ntp              | working  hard |                                                                      |                                                                      |
| column-storage   | working  hard |                                                                      |                                                                      |
| coming soon      | ...           | ..                                                                   |

If you are interested in this project or would like me to prioritize implementing a certain component, please don't
hesitate to submit an issue. If you encounter any problems with the code, you are also welcome to submit a PR with the
appropriate modifications.
