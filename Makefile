BINTARGET=bin/kable
CLI_VERSION?=nightly
CLI_DATE=`date +%FT%T%z`

# Run tests
test: fmt vet
	go test ./... -coverprofile cover.out

# Run go fmt against code
fmt:
	go fmt ./...

install:
	go install -ldflags="-X github.com/redradrat/kable/cmd.CliVersion=$(CLI_VERSION) -X github.com/redradrat/kable/cmd.CliDate=$(CLI_DATE)"

# Run go vet against code
vet:
	go vet ./...

build:
	go build -ldflags="-X github.com/redradrat/kable/cmd.CliVersion=$(CLI_VERSION) -X github.com/redradrat/kable/cmd.CliDate=$(CLI_DATE)" -o ${BINTARGET}
	chmod +x ${BINTARGET}
