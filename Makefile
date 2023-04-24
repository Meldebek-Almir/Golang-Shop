postgres:
	docker run --name postgres_database \
  -p 5432:5432 \
  -e POSTGRES_PASSWORD=mysecretpassword \
  -v ${PWD}/postgres-data:/var/lib/postgresql/data \
  -d postgres 

stop:
	docker stop postgres_database
	docker rm postgres_database

