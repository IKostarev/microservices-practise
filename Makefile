up:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d postgres users
	docker-compose -f docker-compose.yaml up --build migrate
	docker-compose -f docker-compose.yaml ps

down:
	docker-compose -f docker-compose.yaml down -v

generate-users:
	protoc -I protos/users \
		--go_out=protos/users \
		--go_opt=paths=source_relative \
		--go-grpc_out=protos/users \
		--go-grpc_opt=paths=source_relative \
		protos/users/users.proto
