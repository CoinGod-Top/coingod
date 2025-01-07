ifndef GOOS
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	GOOS := darwin
else ifeq ($(UNAME_S),Linux)
	GOOS := linux
else
$(error "$$GOOS is not defined. If you are using Windows, try to re-make using 'GOOS=windows make ...' ")
endif
endif

PACKAGES    := $(shell go list ./... | grep -v '/vendor/' | grep -v '/crypto/ed25519/chainkd' | grep -v '/mining/tensority')
PACKAGES += 'core/mining/tensority/go_algorithm'

BUILD_FLAGS := -ldflags "-X core/version.GitCommit=`git rev-parse HEAD`"


#COINGODD_BINARY32 := coingodd-$(GOOS)_386
COINGODD_BINARY64 := coingodd-$(GOOS)_amd64

#COINGODCLI_BINARY32 := coingodcli-$(GOOS)_386
COINGODCLI_BINARY64 := coingodcli-$(GOOS)_amd64

VERSION := $(shell awk -F= '/Version =/ {print $$2}' version/version.go | tr -d "\" ")


#COINGODD_RELEASE32 := coingodd-$(VERSION)-$(GOOS)_386
COINGODD_RELEASE64 := coingodd-$(VERSION)-$(GOOS)_amd64

#COINGODCLI_RELEASE32 := coingodcli-$(VERSION)-$(GOOS)_386
COINGODCLI_RELEASE64 := coingodcli-$(VERSION)-$(GOOS)_amd64

#COINGOD_RELEASE32 := coingod-$(VERSION)-$(GOOS)_386
COINGOD_RELEASE64 := coingod-$(VERSION)-$(GOOS)_amd64

all: test target release-all install

coingodd:
	@echo "Building coingodd to cmd/coingodd/coingodd"
	@go build $(BUILD_FLAGS) -o cmd/coingodd/coingodd cmd/coingodd/main.go

coingodd-simd:
	@echo "Building SIMD version coingodd to cmd/coingodd/coingodd"
	@cd mining/tensority/cgo_algorithm/lib/ && make
	@go build -tags="simd" $(BUILD_FLAGS) -o cmd/coingodd/coingodd cmd/coingodd/main.go

coingodcli:
	@echo "Building coingodcli to cmd/coingodcli/coingodcli"
	@go build $(BUILD_FLAGS) -o cmd/coingodcli/coingodcli cmd/coingodcli/main.go

install:
	@echo "Installing coingodd and coingodcli to $(GOPATH)/bin"
	@go install ./cmd/coingodd
	@go install ./cmd/coingodcli

target:
	mkdir -p $@

binary: target/$(COINGODD_BINARY64) target/$(COINGODCLI_BINARY64)

ifeq ($(GOOS),windows)
release: binary
	cd target && cp -f $(COINGODD_BINARY64) $(COINGODD_BINARY64).exe
	cd target && cp -f $(COINGODCLI_BINARY64) $(COINGODCLI_BINARY64).exe
	cd target && md5sum $(COINGODD_BINARY64).exe $(COINGODCLI_BINARY64).exe >$(COINGOD_RELEASE64).md5
	cd target && zip $(COINGOD_RELEASE64).zip $(COINGODD_BINARY64).exe $(COINGODCLI_BINARY64).exe $(COINGOD_RELEASE64).md5
	cd target && rm -f $(COINGODD_BINARY64) $(COINGODCLI_BINARY64) $(COINGODD_BINARY64).exe $(COINGODCLI_BINARY64).exe $(COINGOD_RELEASE64).md5
else
release: binary
	cd target && md5sum $(COINGODD_BINARY64) $(COINGODCLI_BINARY64) >$(COINGOD_RELEASE64).md5
	cd target && tar -czf $(COINGOD_RELEASE64).tgz $(COINGODD_BINARY64) $(COINGODCLI_BINARY64) $(COINGOD_RELEASE64).md5
	cd target && rm -f $(COINGODD_BINARY64) $(COINGODCLI_BINARY64) $(COINGOD_RELEASE64).md5
endif

release-all: clean
#	GOOS=darwin  make release
	GOOS=linux   make release
	GOOS=windows make release

clean:
	@echo "Cleaning binaries built..."
	@rm -rf cmd/coingodd/coingodd
	@rm -rf cmd/coingodcli/coingodcli
	@rm -rf target
	@rm -rf $(GOPATH)/bin/coingodd
	@rm -rf $(GOPATH)/bin/coingodcli
	@echo "Cleaning temp test data..."
	@rm -rf test/pseudo_hsm*
	@rm -rf blockchain/pseudohsm/testdata/pseudo/
	@echo "Cleaning sm2 pem files..."
	@rm -rf crypto/sm2/*.pem
	@echo "Done."

target/$(COINGODD_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/coingodd/main.go

target/$(COINGODCLI_BINARY64):
	CGO_ENABLED=0 GOARCH=amd64 go build $(BUILD_FLAGS) -o $@ cmd/coingodcli/main.go



test:
	@echo "====> Running go test"
	@go test -tags "network" $(PACKAGES)

benchmark:
	@go test -bench $(PACKAGES)

functional-tests:
	@go test -timeout=5m -tags="functional" ./test

ci: test functional-tests

.PHONY: all target release-all clean test benchmark
