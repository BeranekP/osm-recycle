BINARY_NAME=recyko2

build:
	env GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}
	env GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}

