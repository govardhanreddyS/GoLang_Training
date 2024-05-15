Running Specific Benchmarks: You can specify which benchmarks to run by providing a regular expression after the -bench flag. For example, if you only want to run benchmarks containing "Factorial" in their names, you can use: go test -benchmem -run=^$ -bench=Factorial .
Parallelism: By default, Go runs benchmarks sequentially. You can increase parallelism by setting the -cpu flag, followed by the number of parallel executions you desire. For example, to run benchmarks using four CPUs in parallel: go test -benchmem -run=^$ -bench . -cpu=4
Benchmark Timeout: You can set a timeout for each benchmark using the -timeout flag, followed by the duration. For example, to set a timeout of 5 seconds for each benchmark: go test -benchmem -run=^$ -bench . -timeout=5s
Verbosity: You can increase the verbosity of the output by setting the -v flag. This will provide more detailed information about each benchmark as it runs: go test -benchmem -run=^$ -bench . -v
MemProfile: If you want to generate memory profiles during benchmark execution, you can use the -memprofile flag, followed by the name of the file to which the memory profile will be written: go test -benchmem -run=^$ -bench . -memprofile=memprofile.out
CPU Profile: Similarly, you can generate CPU profiles during benchmark execution using the -cpuprofile flag: go test -benchmem -run=^$ -bench . -cpuprofile=cpuprofile.out


go tool pprof command. Here's a breakdown of the output:

flat: The total amount of time spent in the function itself.
flat%: The percentage of total time spent in the function itself.
sum%: The cumulative percentage of total time spent in this function and all its callees.
cum: The total amount of time spent in the function and all its callees.
cum%: The percentage of total time spent in the function and all its callees.
From the output, you can see information about two functions:

bench-test.factorialIterative (inline): This function consumed 46.51% of the total CPU time. It's called inline, meaning it's directly invoked within another function.
bench-test.factorialRecursive: This function consumed 38.14% of the total CPU time.
Additionally, you have information about the benchmarks:

bench-test.BenchmarkFactorialIterative: Consumed 12.09% of the total CPU time.
bench-test.BenchmarkFactorialRecursive: Consumed 1.40% of the total CPU time.
Finally, you have some runtime-related entries like runtime.mcall, which consumed a small percentage of CPU time.

This breakdown helps you identify where your program spends the most time, allowing you to optimize performance where necessary. In this case, it appears that the iterative factorial function consumes more CPU time compared to the recursive one during benchmark execution.