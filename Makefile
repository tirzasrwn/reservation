start:
	go build -o reservation cmd/web/*.go && ./reservation

docker_build:
	docker build . -t reservation
docker_start:
	docker run --name reservation -p 4545:4545 -d reservation
docker_stop:
	docker stop reservation
	docker rm reservation
