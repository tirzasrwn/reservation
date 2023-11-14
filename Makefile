start:
	go build -o reservation cmd/web/*.go && ./reservation

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

