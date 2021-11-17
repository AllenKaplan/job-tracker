go-build:
	go build -o app .

go-run: go-build
	./app

build:
	DOCKER_BUILDKIT=1 docker build -t job-tracker .

docker-run:
	docker run  --rm -d -p 8080:8080 -e PORT='8080' \
		--name job-tracker job-tracker


up: build
	docker-compose up -d #> /dev/null