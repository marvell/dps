BUILDDIR := ./build

# target: all - Run tests and generate binary
all: test build

# target: help - Display targets
help:
	@egrep "^# target:" [Mm]akefile | sort - |sed 's/# target://'

# target: clean - Cleans build artifacts
clean:
	@echo " --> Cleaning build artifacts..."
	go clean
	rm -rf ${BUILDDIR}
	@echo

# target: test - Runs CLI tests
test:
	@echo " --> Testing packages..."
	go test .
	@echo

binaries:
	@echo " --> Making binaries..."
	GOARCH=amd64 GOOS=darwin go build -o ${BUILDDIR}/dps_darwin_amd64
	GOARCH=amd64 GOOS=linux go build -o ${BUILDDIR}/dps_linux_amd64

# target: build - Build CLI binary
build: clean binaries
	@echo

.PHONY: all help clean build
