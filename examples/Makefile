export PORT ?= 5555 

serve: bin/serve
	PORT=$(PORT) bin/serve

bin/serve: serve.go
	go build -o $@ $^
