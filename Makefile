# This Makefile is mainly useful for cross-compilation
# In order to be able to cross-compile you need to have
# built GO for all the compilation targets.
# Under Ubuntu Linux those are already available in the repo:
#  sudo apt-get install golang-$GOOS-$GOARCH

.PHONY: cross-compile clean install

cross-compile:
	mkdir -p build
# Linux 32bit:
	GOOS=linux GOARCH=386 go build -o build/gimpact-linux32 client.go
# Linux 64bit:
	GOOS=linux GOARCH=amd64 go build -o build/gimpact-linux64 client.go
# Windows 32bit:
	GOOS=windows GOARCH=386 go build -o build/gimpact-win32.exe client.go
# Windows 64bit:
	GOOS=windows GOARCH=amd64 go build -o build/gimpact-win64.exe client.go
# MacOSX 64bit:
	GOOS=darwin GOARCH=amd64 go build -o build/gimpact-darwin64 client.go

clean:
	rm -rf build/*

# Some standard commands though the original go commands are shorter :-)
test:
	go test

install:
	go install
