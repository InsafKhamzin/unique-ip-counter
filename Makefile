run:
	go build .
	./lightspeed $(path)

test:
	go test -v .