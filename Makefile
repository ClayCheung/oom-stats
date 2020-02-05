build:
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o oom-stats -a ./cmd/entry.go

image: build
	docker build -t clayz95/oom-stats .

push: image
	docker push clayz95/oom-stats:latest

image-tar: image
	docker save -o oom-stats.tar clayz95/oom-stats