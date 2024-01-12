up:
	docker-compose -f docker-compose.yaml down -v
	docker-compose -f docker-compose.yaml up -d postgres rabbitmq jaeger
	docker-compose -f docker-compose.yaml up --build migrate-users
	docker-compose -f docker-compose.yaml up --build migrate-todo
	docker-compose -f docker-compose.yaml up -d --build --force-recreate users-service
	docker-compose -f docker-compose.yaml up -d --build --force-recreate todo-service
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

generate-todo:
	protoc -I api/protos/todo \
		--go_out=api/protos/todo \
		--go_opt=paths=source_relative \
		--go-grpc_out=api/protos/todo \
		--go-grpc_opt=paths=source_relative \
		api/protos/todo/todo.proto

	mkdir -p todo/pkg/grpc_stubs/todo
	cp -r api/protos/todo/* todo/pkg/grpc_stubs/todo

	mkdir -p gateway/pkg/grpc_stubs/todo
	cp -r api/protos/todo/* gateway/pkg/grpc_stubs/todo
