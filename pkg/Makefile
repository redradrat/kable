BINTARGET=bin/kable

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Run go fmt against code
fmt:
	go fmt ./...

install: build
	cp bin/kable ~/tempapps/kable

# Run go vet against code
vet:
	go vet ./...

build:
	go build -o ${BINTARGET}
	chmod +x ${BINTARGET}