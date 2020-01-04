GO_SRC_DIRS := $(shell \
	find . -name "*.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)
GO_TEST_DIRS := $(shell \
	find . -name "*_test.go" -not -path "./vendor/*" | \
	xargs -I {} dirname {}  | \
	uniq)

cats-v1_BUILD_DATE_TIME=$(shell date -u "+%Y.%m.%d %H:%M:%S %Z")
cats-v1_VERSION ?= UNSET
cats-v1_BRANCH ?= UNSET
cats-v1_COMMIT ?= UNSET

format: check-gofmt test

build: go-build

go-build:
	@echo "Building for native..."
	@CGO_ENABLED=0 go build -i -ldflags='-X "github.com/waikco/cats-v1/api.version=$(CATS-V1_VERSION)" -X "github.com/waikco/cats-v1/api.buildDateTime=$(CATS-V1_BUILD_DATE_TIME)" -X "github.com/waikco/cats-v1/api.branch=$(CATS-V1_BRANCH)" -X "github.com/waikco/cats-v1/api.revision=$(CATS-V1_COMMIT)"' -o ./builds/cats-v1 .

go-build-mac:
	@echo "Building for mac"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -i -ldflags='-X "github.com/waikco/cats-v1/api.version=$(CATS-V1_VERSION)" -X "github.com/waikco/cats-v1/api.buildDateTime=$(CATS-V1_BUILD_DATE_TIME)" -X "github.com/waikco/cats-v1/api.branch=$(CATS-V1_BRANCH)" -X "github.com/waikco/cats-v1/api.revision=$(CATS-V1_COMMIT)"' -o ./builds/cats-v1-mac .

check-gofmt: $(GO_SRC_DIRS)
	@echo "Checking formatting..."
	@FMT="0"; \
	for pkg in $(GO_SRC_DIRS); do \
		OUTPUT=`gofmt -l $$pkg/*.go`; \
		if [ -n "$$OUTPUT" ]; then \
			echo "$$OUTPUT"; \
			FMT="1"; \
		fi; \
	done ; \
	if [ "$$FMT" -eq "1" ]; then \
		echo "Problem with formatting in files above."; \
		exit 1; \
	else \
		echo "Success - way to run gofmt!"; \
	fi

test: $(GO_TEST_DIRS)
	@for dir in $^; do \
		pushd ./$$dir > /dev/null ; \
		go test -v ; \
		popd > /dev/null ; \
	done;

