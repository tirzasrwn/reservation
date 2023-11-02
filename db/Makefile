docker_build:
	docker build . -t reservation-postgres
docker_start:
	docker run --name reservation-postgres -p 5432:5432 -v ./data/postgres/:/var/lib/postgresql/data -d reservation-postgres
docker_stop:
	docker stop reservation-postgres
	docker rm reservation-postgres
