Rendezvous
==========

`rendezvous` is a Go implementation of [rendezvous
hashing][wikipedia] (also known as highest random weight hashing).

This implementation is not currently Go-routine safe at all.

[wikipedia]: http://en.wikipedia.org/wiki/Rendezvous_hashing

Benchmarks
----------

(All benchmarks run on a MacBook Pro (Retina), 2.3HGz Intel Core i7.)

1 thread:

    BenchmarkHashGet_5nodes        449 ns/op
    BenchmarkHashGet_10nodes       784 ns/op
    BenchmarkHashGetN3_5_nodes    1079 ns/op
    BenchmarkHashGetN5_5_nodes    1220 ns/op
    BenchmarkHashGetN3_10_nodes   1785 ns/op
    BenchmarkHashGetN5_10_nodes   1924 ns/op

2 threads:

    TODO

