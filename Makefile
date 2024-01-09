up:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d postgres rabbitmq jaeger
	docker-compose -f docker-compose.yaml up --build migrate-users
	docker-compose -f docker-compose.yaml up -d --build --force-recreate users-service
	docker-compose -f docker-compose.yaml up -d --build --force-recreate gateway-service
	docker-compose -f docker-compose.yaml up -d --build --force-recreate notifications-service
	docker-compose -f docker-compose.yaml ps

down:
	docker-compose -f docker-compose.yaml down -v

generate-users:
	protoc -I api/protos/users \
		--go_out=api/protos/users \
		--go_opt=paths=source_relative \
		--go-grpc_out=api/protos/users \
		--go-grpc_opt=paths=source_relative \
		api/protos/users/users.proto

	mkdir -p users/pkg/grpc_stubs/users
	cp -r api/protos/users/* users/pkg/grpc_stubs/users

	mkdir -p gateway/pkg/grpc_stubs/users
	cp -r api/protos/users/* gateway/pkg/grpc_stubs/users
