<p align="center">
    <h1 align="center">E-COMMERCE-MICROSERVICES</h1>
</p>

<p align="center">
	<img src="https://img.shields.io/github/license/RushinShah22/E-commerce-Microservices?style=flat&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/RushinShah22/E-commerce-Microservices?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/RushinShah22/E-commerce-Microservices?style=flat&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/RushinShah22/E-commerce-Microservices?style=flat&color=0080ff" alt="repo-language-count">
    
</p>
<p align="center">
		<em>Built with the tools and technologies:</em>
</p>
<p align="center">
	<img src="https://img.shields.io/badge/YAML-CB171E.svg?style=flat&logo=YAML&logoColor=white" alt="YAML">
	<img src="https://img.shields.io/badge/Docker-2496ED.svg?style=flat&logo=Docker&logoColor=white" alt="Docker">
	<img src="https://img.shields.io/badge/apache-kafka.svg?style=flat&logo=Go&logoColor=white" alt="Kafka">
	<img src="https://img.shields.io/badge/GraphQL-E434AA.svg?style=for-the-badge&logo=graphql&logoColor=white" alt="GraphQL">
</p>

<br>

##### 🔗 Table of Contents

- [📍 Overview](#-overview)
- [👾 Features](#-features)
- [📂 Repository Structure](#-repository-structure)
- [🚀 Getting Started](#-getting-started)
  - [🔖 Prerequisites](#-prerequisites)
  - [🤖 Usage](#-usage)
- [📌 Project Roadmap](#-project-roadmap)

---

## 📍 Overview

<code>❯ The E-commerce Microservices project is designed to provide a scalable and modular architecture for building a comprehensive e-commerce platform. Leveraging the principles of microservices, this application breaks down the e-commerce functionalities into discrete services, each responsible for specific operations, such as product management, order processing, user management, payment handling, and shipping logistics.

❯ This architecture promotes flexibility, allowing each service to be developed, deployed, and scaled independently. The system employs an event-driven approach, utilizing message queues for asynchronous communication between services, which enhances performance and reliability.</code>

---

## 👾 Features

<code>❯ Modular Microservices Architecture: Each functionality (e.g., catalog, orders, users) is encapsulated in its own service, enabling independent development, deployment, and scaling.
❯ Scalability: Easily scale individual services based on demand without affecting the entire system, allowing for optimized resource usage.
❯ Event-Driven Communication: Services communicate asynchronously through message queues, enhancing performance and ensuring resilience against service failures.
❯ API Gateway: A unified entry point for all external client requests, simplifying the client-side experience and allowing for centralized request routing and authentication.
❯ Service Discovery: Automatic detection and registration of services, facilitating dynamic communication between microservices without hardcoding service endpoints.
❯ Robust User Management: The User Service provides secure authentication and profile management, allowing users to create accounts, manage their profiles, and log in securely.
❯ Product Management: The Catalog Service allows for comprehensive management of product listings, including adding new products, updating existing ones, and managing inventory levels.
❯ Order Processing: The Order Service enables users to place, modify, and track their orders seamlessly, providing a smooth purchasing experience.</code>

---

## 📂 Repository Structure

```sh
└── E-commerce-Microservices/
    ├── Makefile
    ├── README.md
    ├── docker-compose.yml
    ├── gateway
    │   ├── docker-compose.yml
    │   ├── gateway-dockerfile.dockerfile
    │   ├── go.mod
    │   ├── go.sum
    │   ├── gqlgen.yml
    │   ├── graph
    │   ├── server.go
    │   └── tools.go
    ├── services
    │   ├── broker
    │   ├── orders
    │   ├── products
    │   └── users
    └── setup.sh
```

---

## 🚀 Getting Started

### 🔖 Prerequisites

<p align="center">
	<img src="https://img.shields.io/badge/Docker-2496ED.svg?style=flat&logo=Docker&logoColor=white" alt="Docker">
	<img src="https://img.shields.io/badge/Go-00ADD8.svg?style=flat&logo=Go&logoColor=white" alt="Go">
</p>

### 🤖 Usage

To run the project, execute the following command:

```sh
❯ git clone https://github.com/RushinShah22/E-commerce-Microservices/
```

```sh
❯ cd ./E-commerce-Microservices
```

```sh
❯ Setup JWT_KEY in gateway docker compose file.
```

```sh
❯ make
```

### 🔌 Ports

```sh
❯ Users -> 3000
❯ Products -> 4000
❯ Orders -> 8000
```

---

## 📌 Project Roadmap

- [x] **`Task 1`**: <strike>Create Microservices for users, products, orders.</strike>
- [x] **`Task 2`**: <strike>Implement kafka as broker.</strike>
- [X] **`Task 3`**: <strike>Implement graphql and authentication.</strike>

---
