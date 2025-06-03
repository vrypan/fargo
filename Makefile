SNAPCHAIN_VER := $(shell cat SNAPCHAIN_VERSION 2>/dev/null || echo "unset")
FARGO_VERSION := $(shell git describe --tags 2>/dev/null || echo "v0.0.0")

BINS = fargo
PROTO_FILES := $(wildcard proto/*.proto)
FARGO_SOURCES := $(wildcard */*.go)

# Colors for output
GREEN = \033[0;32m
NC = \033[0m

all: fargo

clean:
	@echo -e "$(GREEN)Deleting fargo binary...$(NC)"
	rm -f $(BINS)

.PHONY: all clean local release-notes tag tag-minor tag-major releases

fargo: $(FARGO_SOURCES)
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
