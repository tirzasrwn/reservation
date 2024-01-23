# reservation

## About

An implementation of web application in go. Using room reservation for study case.

## Setup

```sh
go mod tidy
```

## Requirement

- make
- go
- docker

## Running

### Web Application

```sh
go run ./cmd/web/
# Or just use the Makefile
make
```

### Test

- Manual command

```sh
# Go to inside the directory part you want to test
go test -v
# To look pecentage of covarage test.
go test -cover
# Same as cover but with more detail and html format.
go test -coverprofile=coverage.out && go tool cover -html=coverage.out

```

- Using make

```sh
make test
```

Example

```sh
# From root level project, go to internal/handlers/
go test -v
```

### Setup Postgresql Linux Debian

```sh
sudo apt update
sudo apt install postgresql
sudo apt install postgresql-client

# Setting up new password
sudo passwd postgres

user:~$ sudo -i -u postgres
postgres@user:~$ psql
postgres=# ALTER USER postgres PASSWORD 'mynewpassword';
```

Make it open to public.

```sh
sudo nano /etc/postgresql/9.3/main/postgresql.conf
# edit: listen_addresses = '*'
# save
sudo nano /etc/postgresql/9.3/main/pg_hba.conf
# add: host    all             all             192.168.7.0/24          md5
# save
sudo systemctl restart postgresql
# Check from another device using nmap
nmap <ip-adress> -p5432
```

### Setup Postgresql using Docker

```sh
make docker_db_build
make docker_db_start
make docker_db_stop
```

### Running database migration

[Read this documentation.](./migrations/readme.md)

### Setup Mailhog using Docker

```sh
make docker_mailhog_build
make docker_mailhog_start
make docker_mailhog_stop
```

## Account

- login page

  ```
  email: admin@admin.com
  password: admin
  ```
