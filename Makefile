build:
	cd /app && go mod vendor && GOFLAGS="-mod=vendor" go build -o app

build-docker:
	docker-compose exec -T golang make build
	docker-compose restart golang

build-local:
	CGO_ENABLED=0 GOOS=linux GOFLAGS="-mod=vendor" /usr/local/Cellar/go@1.14/1.14.12/bin/go build -a -installsuffix cgo -o app
