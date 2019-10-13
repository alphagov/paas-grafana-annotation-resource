.PHONY: test
test:
	ginkgo -r pkg

.PHONY: integration-test
integration-test:
	ginkgo -v integration

.PHONY: local
local:
	docker-compose -f docker-compose-with-concourse.yml up --detach

.PHONY: examples
examples: local
	fly -t local login -k -c http://localhost:8080 -u admin -p admin
	fly -t local sync
	fly -t local set-pipeline -p simple-pipeline -c examples/simple-pipeline.yml

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
