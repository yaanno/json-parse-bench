# Test Json parsing

Default Json parser.

1. open and load large json file
2. process file using entry ids
3. dump json

```
go run main.go
```

benchmark:

```
go test -bench=. -benchmem
```

Compile and run

```
go build main.go
./main
```
