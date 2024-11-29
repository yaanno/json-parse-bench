package main

import "testing"

func BenchmarkProcessJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := processJSON("large-file.json")
		if err != nil {
			b.Fatal(err)
		}
	}
}
