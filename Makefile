APP = tmdbbot
VERSION = $(if $(TAG),$(TAG),$(if $(BRANCH_NAME),$(BRANCH_NAME),$(shell git symbolic-ref -q --short HEAD || git describe --tags --exact-match)))
YOUR_BOT_TOKEN = 123
YOUR_TMDB_API_KEY = 123

run: dep
	@echo "Run bot..."
	@go run .

dep:
	@echo "Resolve dependencies..."
	@go mod tidy

docker-run: docker-build
	@echo "Running docker container..."
	@docker run --env TMDB_BOT_BOT_TOKEN=${YOUR_BOT_TOKEN} --env TMDB_BOT_TMDB_API_KEY=${YOUR_TMDB_API_KEY} ${APP}:${VERSION}

docker-build:
	@echo "Building docker image..."
	@docker build -t ${APP}:${VERSION} --rm --progress=plain .