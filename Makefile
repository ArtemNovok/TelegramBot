build:
	@echo stopping running containers and building new...
	docker-compose down -v
	docker-compose up --build -d
	@echo containers are running! 
stop:
	@echo stopping containers...
	docker-compose down -v
	@echo containers are stopped!

