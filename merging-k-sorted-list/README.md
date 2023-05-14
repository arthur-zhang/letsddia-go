# Merging K sorted list

This project shows how to merge K sorted files using min-heap

It's useful in many scenario especially if the size of sorted file is huge or the number of files is huge.

We can use small memory to do this very well.

## demo explain

we create 10 files in /tmp/k_way_merging

```
$ ls -l /tmp/k_way_merging/
total 80
-rw-r--r--  1 arthur  wheel  28 May 12 11:18 0000.dat
-rw-r--r--  1 arthur  wheel  29 May 12 11:18 0001.dat
-rw-r--r--  1 arthur  wheel  27 May 12 11:18 0002.dat
-rw-r--r--  1 arthur  wheel  29 May 12 11:18 0003.dat
-rw-r--r--  1 arthur  wheel  28 May 12 11:18 0004.dat
-rw-r--r--  1 arthur  wheel  28 May 12 11:18 0005.dat
-rw-r--r--  1 arthur  wheel  28 May 12 11:18 0006.dat
-rw-r--r--  1 arthur  wheel  29 May 12 11:18 0007.dat
-rw-r--r--  1 arthur  wheel  29 May 12 11:18 0008.dat
-rw-r--r--  1 arthur  wheel  28 May 12 11:18 0009.dat
```

each file contains sorted integers like this

```shell
$ cat 0002.dat
5
5
15
26
47
77
77
78
78
84
```