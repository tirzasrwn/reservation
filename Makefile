start:
	go build -o reservation cmd/web/*.go 
	./reservation

test:
	-rm coverage.out coverage.html
	-go test -coverprofile=coverage.out ./...
	-go tool cover -html=coverage.out -o ./coverage.html
	-open coverage.html

restart: docker_db_stop docker_mailhog_stop docker_db_start docker_mailhog_start

docker_webapp_build:
	docker build . -t reservation
docker_webapp_start:
	docker run --name reservation -p 4545:4545 -d reservation
docker_webapp_stop:
	docker stop reservation
	docker rm reservation

docker_db_build:
	docker build ./db/ -t reservation-postgres
docker_db_start:
	docker run --name reservation-postgres -p 5432:5432 -v "./../tmp/db/data/postgres/:/var/lib/postgresql/data:Z" -d reservation-postgres
docker_db_stop:
	docker stop reservation-postgres
	docker rm reservation-postgres

docker_mailhog_build:
	docker build ./mailhog/ -t reservation-mailhog
docker_mailhog_start:
	docker run --name reservation-mailhog -p 1025:1025 -p 8025:8025 -d reservation-mailhog
docker_mailhog_stop:
	docker stop reservation-mailhog
	docker rm reservation-mailhog

