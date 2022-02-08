all:
	mkdir -p bin/
	go build -o bin/bfint

clean:
	rm -rf bin/*
