# Benchmarking

## Running the Benchmark

1. Open up the command line interface and navigate to the folder of tests you want to run
   ```
   cd idgenerator
   cd propagator
   cd sampledspan
   cd unsampledspan
   ```

2. To run the benchmark test, run the following command
   ```
   go test -bench=.
   ```

## Benchmark Test Results
```
BenchmarkStartAndEndSampledSpan-12              20230353                49.9 ns/op
BenchmarkStartAndEndNestedSampledSpan-12        11911389                99.7 ns/op
BenchmarkGetCurrentSampledSpan-12               19009425                63.4 ns/op
BenchmarkAddAttributesToSampledSpan-12           8242645                145  ns/op

BenchmarkStartAndEndUnSampledSpan-12            23559434                51.0 ns/op
BenchmarkStartAndEndNestedUnSampledSpan-12      11958974                100  ns/op
BenchmarkGetCurrentUnSampledSpan-12             19014970                63.0 ns/op
BenchmarkAddAttributesToUnSampledSpan-12         8241534                145  ns/op

BenchmarkIDGenerator-12                           116726                9202 ns/op
BenchmarkPropagatorExtract-12                   36464926                28.6 ns/op
BenchmarkPropagatorInject-12                     2430130                 488 ns/op
```
