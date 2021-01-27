build:
	go build -o foreach

install: build
	cp ./foreach /usr/local/bin/foreach
