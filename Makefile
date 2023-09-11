.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-player-shots get-player-shots/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-line-stats get-line-stats/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-line-inverse-stats get-line-inverse-stats/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-team-stats get-team-stats/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-season-teams get-season-teams/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-line-plays-with get-line-plays-with/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-seasons get-seasons/*
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/get-season-team-roster get-season-team-roster/*

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
