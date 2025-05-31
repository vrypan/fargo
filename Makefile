SNAPCHAIN_VER := $(shell cat SNAPCHAIN_VERSION 2>/dev/null || echo "unset")
FARGO_VERSION := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

BINS = fargo
PROTO_FILES := $(wildcard proto/*.proto)
FARGO_SOURCES := $(wildcard */*.go)

# Colors for output
GREEN = \033[0;32m
NC = \033[0m

all: fargo

# Compile .proto files, touch stamp file
.proto-bindings: $(PROTO_FILES) .proto
	@echo -e "$(GREEN)Compiling .proto files...$(NC)"
	protoc --proto_path=proto --go_out=. \
	    $(shell cd proto; ls | xargs -I \{\} echo -n '--go_opt=M'{}=farcaster/" " '--go-grpc_opt=M'{}=farcaster/" " ) \
		--go-grpc_out=. \
		proto/*.proto
	@touch $@

.proto: SNAPCHAIN_VERSION
	@echo -e "$(GREEN)Downloading proto files (Hubble v$(SNAPCHAIN_VER))...$(NC)"
	curl -s -L "https://codeload.github.com/farcasterxyz/snapchain/tar.gz/refs/tags/v$(SNAPCHAIN_VER)" \
	    | tar -zxvf - -C . --strip-components 2 "snapchain-$(SNAPCHAIN_VER)/src/proto/"
	@touch $@

proto-clean:
	@echo -e "$(GREEN)Deleting protobuf definitions...$(NC)"
	rm -f proto/*.proto .proto
	@echo -e "$(GREEN)Deleting protobuf bindings...$(NC)"
	rm -f farcaster/*.pb.go farcaster/*.pb.gw.go .proto-bindings

clean:
	@echo -e "$(GREEN)Deleting fargo binary...$(NC)"
	rm -f $(BINS)

.PHONY: all proto-clean clean local release-notes tag tag-minor tag-major releases

fargo: .proto-bindings $(FARGO_SOURCES)
	@echo -e "$(GREEN)Building fargo ${FARGO_VER} $(NC)"
	go build \
	-ldflags "-w -s" \
	-ldflags "-X github.com/vrypan/fargo/config.FARGO_VERSION=${FARGO_VERSION}" \
	-o fargo

release-notes:
	# Automatically generate release_notes.md
	./bin/generate_release_notes.sh

tag:
	./bin/auto_increment_tag.sh patch

tag-minor:
	./bin/auto_increment_tag.sh minor

tag-major:
	./bin/auto_increment_tag.sh major

releases:
	goreleaser release --clean
