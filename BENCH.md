# Benchmark commands

## Exemple using 'internal/utils' BenchmarkGetMapKey* functions

1. Run benchmarks

This should be used to identify consuming steps in a specific function.
A second function should be created to compare to the baseline function.

```bash
go test \
	-run=^$ \
	-bench "^BenchmarkGetMapKey([2])?$" internal/utils/*.go \
	-v \
	-cpuprofile cpu.prof \
	-benchmem \
	-memprofile mem.prof  \
	-count 4
```

```bash
go tool pprof internal/utils/cpu.prof
    -> top5
    -> top5 -cum
    -> top10
    -> list GetMapKey
```
    
```bash
go tool pprof internal/utils/mem.prof
```


2. Compare the statistical improvements

This should be used to compare the improvements between the first version of the function
(old) and the latest version (new).

```bash
go test \
	-run=^$ \
	-bench ^BenchmarkGetMapKey$ internal/utils/*.go \
	-count 6 > old.txt
```

```bash
go test \
	-run=^$ \
	-bench ^BenchmarkGetMapKey$ internal/utils/*.go \
	-count 6 > new.txt
```

```bash
benchstat old.txt new.txt
```
