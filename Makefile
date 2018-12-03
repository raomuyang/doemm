.PHONY: vendor
vendor:
	@echo "Prepare environment via govendor"
	govendor init
	govendor add +e


.PHONY: prepare
prepare:
	@echo "Prepare environment via go get"
	go get -t ./...

test: prepare
	@echo "Test project"
	go test -race ./...
	go vet ./...

.PHONY: compile
compile: prepare
	@echo "Compile all"
	mkdir -p target
	GOOS=windows GOARCH=amd64 go build  -o target/doemm-win.exe
	GOOS=linux GOARCH=amd64 go build  -o target/doemm-linux
	GOOS=darwin GOARCH=amd64 go build  -o target/doemm-mac
	chmod a+x target/*

.PHONY: clean
clean:
	rm -r target

.PHONY: install
install: prepare
	go install github.com/raomuyang/doemm
