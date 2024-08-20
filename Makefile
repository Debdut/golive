build:
	go mod tidy && go mod vendor
	go build -o ./bin/golive && chmod +X ./bin/*
	echo "Executable Ready in ./bin/golive"

docker:
	docker build -t antsankov/golive:latest .

format: 
	gofmt -l -s -w .

# Compiled with Arch Go package guidelines and removed whitespace.
cross-compile:
	env GOOS=linux GOARCH=arm go build -o ./release/golive-linux-arm32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=arm64 go build -o ./release/golive-linux-arm64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=darwin GOARCH=amd64 go build -o ./release/golive-mac-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=386 go build -o ./release/golive-linux-x32 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=linux GOARCH=amd64 go build -o ./release/golive-linux-x64 -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=windows GOARCH=386 go build -o ./release/golive-windows-x32.exe -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=windows GOARCH=amd64 go build -o ./release/golive-windows-x64.exe -ldflags "-s -w" -trimpath -mod=readonly
	env GOOS=darwin GOARCH=arm64 go build -o ./release/golive-mac-arm64 -ldflags "-s -w" -trimpath -mod=readonly

m1:
	env GOOS=darwin GOARCH=arm64 go build -o ./release/golive-mac-arm64 -ldflags "-s -w" -trimpath -mod=readonly

test:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out	

clean:
	rm -rf ./bin/*
	rm -rf ./release/*
	rm -rf ./vendor/*

run:
	./bin/golive
