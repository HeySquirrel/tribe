default: test

test:
	go test -v ./...

prtest:
	go test -tags authenticated -v ./...
