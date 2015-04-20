# This Makefile is mainly useful for cross-compilation
# and upload to the GitHub release sections
# It is based on https://github.com/aktau/github-release/Makefile
# In order to be able to cross-compile you need to have
# built GO for all the compilation targets.
# Under Ubuntu Linux those are already available in the repo:
#  $ sudo apt-get install golang-$GOOS-$GOARCH
# For upload you need to have github-release installed
#  $ go get github.com/aktau/github-release

.PHONY: clean builds release dep install

LAST_TAG := $(shell git describe --abbrev=0 --tags)

USER := xogeny
EXECUTABLE := impact
GITHUB_RELEASE = github-release
BUILD_CMD =  go build

UNIX_EXECUTABLES := \
	darwin/386/$(EXECUTABLE) \
	darwin/amd64/$(EXECUTABLE) \
	linux/386/$(EXECUTABLE) \
	linux/amd64/$(EXECUTABLE) \
#	linux/arm/$(EXECUTABLE) \
	freebsd/386/$(EXECUTABLE) \
	freebsd/amd64/$(EXECUTABLE) \
	freebsd/arm/$(EXECUTABLE) \
	netbsd/386/$(EXECUTABLE) \
	netbsd/amd64/$(EXECUTABLE)

WIN_EXECUTABLES := \
	windows/386/$(EXECUTABLE).exe \
	windows/amd64/$(EXECUTABLE).exe

COMPRESSED_EXECUTABLES=$(UNIX_EXECUTABLES:%=%.tar.xz) $(WIN_EXECUTABLES:%.exe=%.zip)
COMPRESSED_EXECUTABLE_TARGETS=$(COMPRESSED_EXECUTABLES:%=build/%)

UPLOAD_CMD = $(GITHUB_RELEASE) upload -u $(USER) -r $(EXECUTABLE) -t $(LAST_TAG) -n $(subst /,-,$(FILE)) -f build/$(FILE)

all: $(EXECUTABLE)

builds: $(COMPRESSED_EXECUTABLE_TARGETS)

# 386
build/darwin/386/$(EXECUTABLE):
	GOARCH=386 GOOS=darwin $(BUILD_CMD) -o "$@"
build/linux/386/$(EXECUTABLE):
	GOARCH=386 GOOS=linux $(BUILD_CMD) -o "$@"
build/windows/386/$(EXECUTABLE).exe:
	GOARCH=386 GOOS=windows $(BUILD_CMD) -o "$@"
# We don't build the more uncomon platforms for now
#build/freebsd/386/$(EXECUTABLE):
#	GOARCH=386 GOOS=freebsd $(BUILD_CMD) -o "$@"
#build/netbsd/386/$(EXECUTABLE):
#	GOARCH=386 GOOS=netbsd $(BUILD_CMD) -o "$@"


# amd64
build/darwin/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin $(BUILD_CMD) -o "$@"
build/linux/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux $(BUILD_CMD) -o "$@"
build/windows/amd64/$(EXECUTABLE).exe:
	GOARCH=amd64 GOOS=windows $(BUILD_CMD) -o "$@"
# We don't build the more uncomon platforms for now
#build/freebsd/amd64/$(EXECUTABLE):
#	GOARCH=amd64 GOOS=freebsd $(BUILD_CMD) -o "$@"
#build/netbsd/amd64/$(EXECUTABLE):
#	GOARCH=amd64 GOOS=netbsd $(BUILD_CMD) -o "$@"
# arm
#build/freebsd/arm/$(EXECUTABLE):
#	GOARCH=arm GOOS=freebsd $(BUILD_CMD) -o "$@"
#build/linux/arm/$(EXECUTABLE):
#	GOARCH=arm GOOS=linux $(BUILD_CMD) -o "$@"


# compressed artifacts, makes a huge difference (Go executable is ~9MB,
# after compressing ~2MB)
%.tar.xz: %
	tar -Jcvf "$<.tar.xz" -C `dirname "$<"` `basename "$<"`
%.zip: %.exe
	zip -j "$@" "$<"

# git tag -a v$(RELEASE) -m 'release $(RELEASE)'
release: $(COMPRESSED_EXECUTABLE_TARGETS)
	git push && git push --tags
	$(GITHUB_RELEASE) release -u $(USER) -r $(EXECUTABLE) \
		-t $(LAST_TAG) -n $(LAST_TAG) --pre-release || true
	$(foreach FILE,$(COMPRESSED_EXECUTABLES),$(UPLOAD_CMD);)

dep:
	go get -t ../...

$(EXECUTABLE): dep
	$(BUILD_CMD) -o "$@"

clean:
	rm $(EXECUTABLE) || true
	rm -rf build/

# Some standard commands though the original go commands are shorter :-)
test:
	go test ../...

install:
	go install
