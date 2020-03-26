APPNAME = "chat"
# h - help
h help:
	@echo "h help 	- this help"
	@echo "build 	- build and the app"
	@echo "run 	- run the app"
	@echo "clean 	- clean app trash"
	@echo "test 	- run all tests"
.PHONY: h

# build - build the app
build:
	go build -o $(APPNAME)
.PHONY: build

# run - build and run the app
run: build
	./$(APPNAME) -p=8080
.PHONY: run

clean:
	rm ./$(APPNAME)
.PHONY: clean

# test - run all tests
test:
	go test -cover ./...
.PHONY: test
