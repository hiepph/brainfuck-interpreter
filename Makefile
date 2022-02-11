all:
	mkdir -p bin/
	go build -o bin/bf

clean:
	rm -rf bin/*
