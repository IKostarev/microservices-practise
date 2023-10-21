up:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d users
	docker-compose -f docker-compose.yaml ps

down:
	docker-compose -f docker-compose.yaml down -v
