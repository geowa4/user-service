SERVICE_NAME=user-service
SRC_FILES=./...
DIST_DIR=dist
DIST_BIN=$(DIST_DIR)/$(SERVICE_NAME)
DB_SETUP_FILE=$(DIST_DIR)/dbsetup
GOBUILD_OPTS=CGO_ENABLED=0 GOOS=linux GOARCH=amd64

all: docker

install:
	go get -u golang.org/x/tools/cmd/vet github.com/golang/lint/golint golang.org/x/tools/cmd/oracle golang.org/x/tools/cmd/goimports golang.org/x/tools/cmd/cover github.com/nsf/gocode github.com/cespare/reflex
	go get -u github.com/gorilla/mux gopkg.in/jackc/pgx.v2 gopkg.in/inconshreveable/log15.v2 github.com/stretchr/testify github.com/DATA-DOG/go-sqlmock github.com/kisielk/sqlstruct

clean:
	-rm -f $(DIST_BIN)
	-docker kill $(SERVICE_NAME)-web
	-docker rm $(SERVICE_NAME)-web

realclean: clean
	-rm -f $(DB_SETUP_FILE)
	-docker kill $(SERVICE_NAME)-db
	-docker rm $(SERVICE_NAME)-db

fmt:
	go fmt $(SRC_FILES)

lint:
	golint $(SRC_FILES)

vet:
	go vet $(SRC_FILES)

test: vet
	go test -cover $(SRC_FILES)

$(DIST_BIN): *.go
	$(GOBUILD_OPTS) go build -o $@ .

dist: lint test $(DIST_BIN)

docker: dist
	docker build --no-cache --rm -t $(SERVICE_NAME) $(DIST_DIR)

$(DB_SETUP_FILE):
	docker run --name $(SERVICE_NAME)-db -v $(CURDIR)/db:/docker-entrypoint-initdb.d -e POSTGRES_USER=$(SERVICE_NAME) -e POSTGRES_PASSWORD="" -d postgres
	touch $(DB_SETUP_FILE)

dev: $(DB_SETUP_FILE) docker
	docker run --name $(SERVICE_NAME)-web -p 8080:8080 --link $(SERVICE_NAME)-db:db -d $(SERVICE_NAME)

watch: clean dev
	reflex -r '\.go$$' make clean dev

.PHONY: install clean realclean fmt lint vet test dist docker dev watch
