# Go makefile

#export env
#basic information
ProjectAdmin := "cronyadmin"
ProjectNode := "cronynode"

PROJECTBASE 	:= $(shell pwd)
PROJECTBIN 	:= $(PROJECTBASE)/bin
AdminConf := "$(PROJECTBIN)/admin"
NodeConf := "$(PROJECTBIN)/node"
TIMESTAMP   := $(shell /bin/date "+%F %T")

#change to deploy environment
AdminFile := "./admin/cmd/main.go"
NodeFile := "./node/cmd/main.go"

#compile ldflags
LDFLAGS		:= -s -w \
			   -X 'main.BuildGitBranch=$(shell git describe --all)' \
			   -X 'main.BuildGitRev=$(shell git rev-list --count HEAD)' \
			   -X 'main.BuildGitCommit=$(shell git rev-parse HEAD)' \
			   -X 'main.BuildDate=$(shell /bin/date "+%F %T")'


#linux amd64开发调试版本
linux-dev:
	@echo "install linux amd64 dev version"
	@if [ ! -d $(AdminConf)/logs ]; then \
        mkdir -p $(AdminConf)/logs; \
    fi
	@if [ ! -d $(NodeConf)/logs ]; then \
        mkdir -p $(NodeConf)/logs; \
    fi
	@cp -r ./admin/conf $(AdminConf)
	@cp -r ./node/conf $(NodeConf)
	@echo "building project cronyadmin..."
	@CGO_ENABLED=0 GOOS=linux  GOARCH=amd64 go build -v -o $(PROJECTBIN)/$(ProjectAdmin) $(AdminFile)
	@echo "building project cronynode..."
	@CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -v -o $(PROJECTBIN)/$(ProjectNode) $(NodeFile)
	@chmod +x $(PROJECTBIN)/$(ProjectAdmin)
	@chmod +x $(PROJECTBIN)/$(ProjectNode)
	@echo "build success."


#本地开发调试版本，取决于系统
local-dev: clean
	@echo "install local dev version"
	@if [ ! -d $(NodeConf)/logs ]; then \
        mkdir -p $(NodeConf)/logs; \
    fi
	@if [ ! -d $(AdminConf)/logs ]; then \
            mkdir -p $(AdminConf)/logs; \
    fi
	@cp -r ./admin/conf $(AdminConf)
	@cp -r ./node/conf $(NodeConf)
	@echo "building project cronyadmin..."
	@CGO_ENABLED=0 go build -v  -o $(PROJECTBIN)/$(ProjectAdmin) $(AdminFile)
	@echo "building project cronynode..."
	@CGO_ENABLED=0 go build -v  -o $(PROJECTBIN)/$(ProjectNode) $(NodeFile)
	@chmod +x $(PROJECTBIN)/$(PROJECTNAME)
	@chmod +x $(PROJECTBIN)/$(ProjectNode)
	@echo "build success."

gitpush: clean fmt
	git add .
	git commit -m "$(m) changed at $(TIMESTAMP)"
	git push
fmt:
	@go fmt $(PROJECTBASE)/...
	@echo "hello"
	@go mod tidy

clean:
	@#echo $(PROJECTBIN)
	@rm -rf $(PROJECTBIN)/* &>/dev/null
depend:
	go mod download
gitpull: fmt
	git add .
	git commit -m "$(m) changed at $(TIMESTAMP)"
	git pull
.PHONY: fmt clean git