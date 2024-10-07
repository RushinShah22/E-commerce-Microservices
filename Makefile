all: setup broker-micro products-micro orders-micro users-micro server

products-micro: 
	cd ./services/products && \
	docker compose down && \
	docker compose build && \
	docker compose up -d

orders-micro: 
	cd ./services/orders && \
	docker compose down && \
	docker compose build && \
	docker compose up -d

users-micro: 
	cd ./services/users && \
	docker compose down && \
	docker compose build  && \
	docker compose up -d

broker-micro: 
	cd ./services/broker && \
	docker compose down && \
	docker compose build  && \
	docker compose up -d && \
	go run . localhost:9092 products 2 1 && \
	go run . localhost:9092 orders 1 1 && \
	go run . localhost:9092 users 1 1

stop-services:
	cd ./services/broker && \
	docker compose down
	cd ./services/orders && \
	docker compose down
	cd ./services/products && \
	docker compose down
	cd ./services/users && \
	docker compose down
	cd ./gateway && \
	docker compose down

server:
	cd ./gateway && \
	docker compose down && \
	docker compose build && \
	docker compose up -d

setup:
	bash ./setup.sh

