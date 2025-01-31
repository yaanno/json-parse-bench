#!/bin/bash
set -e

# Compile Go implementation
echo "=== Compiling Go Implementation ==="
cd /Users/A200246910/workspace/json-parsers/go
go build -o json-parser main.go

# Compile Rust implementation
echo "=== Compiling Rust Implementation ==="
cd /Users/A200246910/workspace/json-parsers/rust
cargo build --release

# Return to project root
cd /Users/A200246910/workspace/json-parsers

echo "=== Starting Benchmarks ==="

# Node.js Benchmark
echo -e "\n=== Node.js Benchmark ==="
time node node/index.js large-file.json

# Python Standard Benchmark
echo -e "\n=== Python Standard Benchmark ==="
time python3 python/main.py large-file.json

# Python Optimized Benchmark
echo -e "\n=== Python Optimized Benchmark ==="
time python3 python/main_optimized.py large-file.json

# Go Standard Benchmark
echo -e "\n=== Go Standard Benchmark ==="
cd go
./json-parser
echo "=== Go Standard Memory Profile ==="
go tool pprof -text mem_profile.prof | awk '
    /^TYPE/ {next}
    /^Dropped/ {next}
    /^#/ {next}
    /^$/ {next}
    {print $1, $2, $3, $4}
' | head -n 10
cd ..

# Go Parallel Benchmark
echo -e "\n=== Go Parallel Benchmark ==="
cd go
./json-parser -parallel
echo "=== Go Parallel Memory Profile ==="
go tool pprof -text mem_profile.prof | awk '
    /^TYPE/ {next}
    /^Dropped/ {next}
    /^#/ {next}
    /^$/ {next}
    {print $1, $2, $3, $4}
' | head -n 10
cd ..

# Rust Standard Benchmark
echo -e "\n=== Rust Standard Benchmark ==="
time rust/target/release/rust

# Rust Parallel Benchmark
echo -e "\n=== Rust Parallel Benchmark ==="
time rust/target/release/rust --parallel

echo -e "\n=== Benchmark Complete ==="
