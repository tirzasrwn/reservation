# reservation

## Setup
```sh
go mod tidy
```

## Running

### Web Application

```sh
go run $(ls cmd/web/*.go | grep -v _test.go)
# Or just run run.sh script.
chmod +x run.sh
./run.sh
```
### Test

```sh
# Go to inside the directory part you want to test
go test -v
go test -cover # To look pecentage of covarage test.
go test -coverprofile=coverage.out && go tool cover -html=coverage.out # Same as cover but with more detail and html format.

```
Example
```sh
# From root level project, go to internal/handlers/
go test -v
```

### Setup Postgressql LInux Debian

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

### Install soda
https://gobuffalo.io/documentation/database/soda/  
```sh
go install github.com/gobuffalo/pop/v6/soda@latest
soda generate fizz CreateUserTable
soda migrate # Up migration.
soda migrate down # Down migration.
soda generate fizz CreateReservationTable
```

