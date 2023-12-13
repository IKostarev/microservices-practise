up:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d postgres users todo
	docker-compose -f docker-compose.yaml up --build migrate_todo migrate_users
	docker-compose -f docker-compose.yaml ps

down:
	docker-compose -f docker-compose.yaml down -v
