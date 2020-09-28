BINTARGET=bin/kable
CLI_VERSION?=nightly
CLI_DATE=`date +%FT%T%z`

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
	go build -ldflags="-X github.com/redradrat/kable/pkg/cmd.cliVersion=$(CLI_VERSION) -X github.com/redradrat/kable/pkg/cmd.cliDate=$(CLI_DATE)" -o ${BINTARGET}
	chmod +x ${BINTARGET}

license:
	licensed cache -c ../.licensed
	licensed status -c ../.licensed