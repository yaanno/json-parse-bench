package main

import "testing"

func BenchmarkProcessJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ProcessJSONStream("large-file.json", func(SimplifiedData) error {
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParallelProcessJSON(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := ParallelProcessJSONStream("large-file.json", 4, func(SimplifiedData) error {
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
	}
}
