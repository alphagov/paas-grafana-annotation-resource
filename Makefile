.PHONY: compile
compile: compile-check compile-in compile-out

.PHONY: compile-check
compile-check:
	go build -o bin/check cmd/check/main.go

.PHONY: compile-in
compile-in:
	go build -o bin/in cmd/in/main.go

.PHONY: compile-out
compile-out:
	go build -o bin/out cmd/out/main.go
