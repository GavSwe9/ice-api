.PHONY: build clean deploy

build:
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-player-shots/bootstrap get-player-shots/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-game-shots/bootstrap get-game-shots/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-line-stats/bootstrap get-line-stats/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-line-inverse-stats/bootstrap get-line-inverse-stats/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-team-stats/bootstrap get-team-stats/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-season-teams/bootstrap get-season-teams/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-season-teams-map/bootstrap get-season-teams-map/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-line-plays-with/bootstrap get-line-plays-with/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-seasons/bootstrap get-seasons/*
	env GOARCH=arm64 GOOS=linux go build -ldflags="-s -w" -o build/lambda/get-season-team-roster/bootstrap get-season-team-roster/*

zip:
	zip -j build/lambda/get-player-shots.zip build/lambda/get-player-shots/bootstrap
	zip -j build/lambda/get-game-shots.zip build/lambda/get-game-shots/bootstrap
	zip -j build/lambda/get-line-stats.zip build/lambda/get-line-stats/bootstrap
	zip -j build/lambda/get-line-inverse-stats.zip build/lambda/get-line-inverse-stats/bootstrap
	zip -j build/lambda/get-team-stats.zip build/lambda/get-team-stats/bootstrap
	zip -j build/lambda/get-season-teams.zip build/lambda/get-season-teams/bootstrap
	zip -j build/lambda/get-season-teams-map.zip build/lambda/get-season-teams-map/bootstrap
	zip -j build/lambda/get-line-plays-with.zip build/lambda/get-line-plays-with/bootstrap
	zip -j build/lambda/get-seasons.zip build/lambda/get-seasons/bootstrap
	zip -j build/lambda/get-season-team-roster.zip build/lambda/get-season-team-roster/bootstrap

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
