.PHONY: start-db stop-db

start-db:
	docker run -d --name mariadb-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=rootpassword -e MYSQL_DATABASE=mydatabase mariadb:latest

stop-db:
	docker stop mariadb-container
	docker rm mariadb-container