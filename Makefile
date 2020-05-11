NAME = bookshelf
VERSION = v1.0.0

export GO111MODULE = on

.PHONY: test-db
test-db:
	go test -v db_test.go db_mysql.go bookshelf.go db_memory.go

.PHONY: test-main
test-main:
	go test -v main_test.go bookshelf.go db_memory.go main.go template.go db_mysql.go

.PHONY: test
test:
	make test-db && \
	make test-main

.PHONY: migrate-up
migrate-up:
	docker run -v $(PWD)/migrations:/migrations --network host migrate/migrate \
	-path=/migrations/ -database 'mysql://user:password@tcp(localhost:3306)/default' up $(N) || \
	docker run -v $(PWD)/migrations:/migrations --network host migrate/migrate \
	-path=/migrations/ -database 'mysql://user:password@tcp(localhost:3306)/default' force $(N); \
	docker run -v $(PWD)/migrations:/migrations --network host migrate/migrate \
        -path=/migrations/ -database 'mysql://user:password@tcp(localhost:3306)/default' version || true

.PHONY: migrate-down
migrate-down:
	docker run -v $(PWD)/migrations:/migrations --network host migrate/migrate \
	-path=/migrations/ -database 'mysql://user:password@tcp(localhost:3306)/default' down $(N) || \
	docker run -v $(PWD)/migrations:/migrations --network host migrate/migrate \
        -path=/migrations/ -database 'mysql://user:password@tcp(localhost:3306)/default' force $(N); \
        docker run -v $(PWD)/migrations:/migrations --network host migrate/migrate \
        -path=/migrations/ -database 'mysql://user:password@tcp(localhost:3306)/default' version || true
