PREFIX?=	/usr/local
GO_CMD=		${PREFIX}/bin/go
GO_BIN?=	ports-webhook
GO_FLAGS?=	-ldflags="-s -w"

default:
	${GO_CMD} build ${GO_FLAGS} -o ${GO_BIN}

build: default

clean:
	${GO_CMD} clean

mod:
	${GO_CMD} mod tidy -v
	${GO_CMD} mod verify

mod-update:
	${GO_CMD} get -u -v

update: mod-update mod

fmt:
	find . -type f -name "*.go" -exec gofmt -w {} \;

.PHONY: default build clean mod mod-update update fmt
