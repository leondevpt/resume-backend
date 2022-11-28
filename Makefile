GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags || git rev-parse --short=8 HEAD)
INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)

BUILDTIME=$(shell TZ=Asia/Shanghai date +%FT%T%z)
GitTag=$(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
GitBranch=$(shell git rev-parse --abbrev-ref HEAD)
GitCommit=$(shell git rev-parse --short=12 HEAD)
GitTreeState=$(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

REPO = github.com/leondevpt/resume-backend


DOCKER_CMD=docker
DOCKER_BUILD=$(DOCKER_CMD) build
DOCKER_PUSH=$(DOCKER_CMD) push

DOCKER_IMAGE_NAME=resume:$(VERSION)

BUILD_FLAGS := -ldflags "-X '${REPO}/version.GitCommit=$(GitCommit)' \
                         -X '${REPO}/version.GitBranch=$(GitBranch)' \
                         -X '${REPO}/version.GitTag=$(GitTag)' \
                         -X '${REPO}/version.GitTreeState=$(GitTreeState)' \
                         -X '${REPO}/version.Version=$(VERSION)' \
                         -X '${REPO}/version.BuildTime=[$(BUILDTIME)]'"

LDFLAGS := ' -w -s'

.PHONY: api
# generate api proto
api:
	echo "开始编译proto文件生成api"
	buf generate api


.PHONY: build
# build
build:
	echo $(BUILD_FLAGS)
	mkdir -p bin/ && go build $(BUILD_FLAGS)$(LDFLAGS) -v -o ./bin/ ./...


.PHONY: docker-build
# docker-build
docker-build:
	echo "打包 Docker Image - $(DOCKER_IMAGE_NAME)"
	$(DOCKER_BUILD) -t $(DOCKER_IMAGE_NAME) .


##项目依赖的相关工具
.PHONY: dep
dep:
	@go install \
	github.com/bufbuild/buf/cmd/buf \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc \
	google.golang.org/protobuf/cmd/protoc-gen-go \
	github.com/google/wire/cmd/wire