# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gbchain android ios gbchain-cross swarm evm all test clean
.PHONY: gbchain-linux gbchain-linux-386 gbchain-linux-amd64 gbchain-linux-mips64 gbchain-linux-mips64le
.PHONY: gbchain-linux-arm gbchain-linux-arm-5 gbchain-linux-arm-6 gbchain-linux-arm-7 gbchain-linux-arm64
.PHONY: gbchain-darwin gbchain-darwin-386 gbchain-darwin-amd64
.PHONY: gbchain-windows gbchain-windows-386 gbchain-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest

gbchain:
	build/env.sh go run build/ci.go install ./cmd/gbchain
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gbchain\" to launch gbchain."

swarm:
	build/env.sh go run build/ci.go install ./cmd/swarm
	@echo "Done building."
	@echo "Run \"$(GOBIN)/swarm\" to launch swarm."

all:
	build/env.sh go run build/ci.go install

android:
	build/env.sh go run build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gbchain.aar\" to use the library."

ios:
	build/env.sh go run build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/gbchain.framework\" to use the library."

test: all
	build/env.sh go run build/ci.go test

lint: ## Run linters.
	build/env.sh go run build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gbchain-cross: gbchain-linux gbchain-darwin gbchain-windows gbchain-android gbchain-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/sipe-*

#CGO_ENABLED=0 GOOS=linux GOARCH=amd64
gbchain-linux: gbchain-linux-386 gbchain-linux-amd64 gbchain-linux-arm gbchain-linux-mips64 gbchain-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-*

gbchain-linux-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gbchain
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep 386

gbchain-linux-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gbchain
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep amd64

gbchain-linux-arm: gbchain-linux-arm-5 gbchain-linux-arm-6 gbchain-linux-arm-7 gbchain-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep arm

gbchain-linux-arm-5:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gbchain
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep arm-5

gbchain-linux-arm-6:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gbchain
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep arm-6

gbchain-linux-arm-7:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gbchain
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep arm-7

gbchain-linux-arm64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gbchain
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep arm64

gbchain-linux-mips:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gbchain
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep mips

gbchain-linux-mipsle:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gbchain
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep mipsle

gbchain-linux-mips64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gbchain
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep mips64

gbchain-linux-mips64le:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gbchain
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-linux-* | grep mips64le

gbchain-darwin: gbchain-darwin-386 gbchain-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-darwin-*

gbchain-darwin-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gbchain
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-darwin-* | grep 386

gbchain-darwin-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gbchain
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-darwin-* | grep amd64

gbchain-windows: gbchain-windows-386 gbchain-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-windows-*

gbchain-windows-386:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gbchain
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-windows-* | grep 386

gbchain-windows-amd64:
	build/env.sh go run build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gbchain
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gbchain-windows-* | grep amd64
