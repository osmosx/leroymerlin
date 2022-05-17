run:
	go run main.go

build:
	 go build -o bin/util main.go

run-bin:
	@make build
	./bin/util
