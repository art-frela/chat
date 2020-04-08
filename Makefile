APPNAME = "chat"
DOCKERSPACE = "artfrela"
# h - help
h help:
	@echo "h help 	- this help"
	@echo "build 	- build and the app"
	@echo "run 	- run the app"
	@echo "clean 	- clean app trash"
	@echo "test 	- run all tests"
	@echo "docker 	- build docker image and pull to dockerhub"
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

# docker build
docker:
	docker build . -t $(DOCKERSPACE)/$(APPNAME):latest
	docker push $(DOCKERSPACE)/$(APPNAME):latest
.PHONY: docker