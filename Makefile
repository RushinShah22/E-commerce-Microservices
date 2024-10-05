all: broker-micro products-micro orders-micro users-micro

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
	go run . localhost:9092 products 2 1