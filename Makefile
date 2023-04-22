BINARY_NAME=goscrapy

build:
	$(info Building the app for local testing)
	go mod tidy
	cd $(PWD)/cmd/api && go build -o ${BINARY_NAME} ./...
	mv ./cmd/api/${BINARY_NAME} .

run:	build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}
