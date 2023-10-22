up:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d postgres users
	docker-compose -f docker-compose.yaml up --build migrate
	docker-compose -f docker-compose.yaml ps

down:
	docker-compose -f docker-compose.yaml down -v
