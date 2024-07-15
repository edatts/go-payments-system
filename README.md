## Go Payments System

This repo contains an example of a simple payments application written in Go. The application consists of two containerized services, one service for authenticating the user and one service for receiving and processing their payment requests.

### Current Features

- Microservices run using Docker Compose
- HTTP REST APIs for registration, logging in, deposits, withdrawals, and payments.
- JWT authenticaion.
- PostgreSQL for backing database.
- Dependency Injection for the data stores.
- Goalng-migrate tool for setting up the database.
- Transactional database queries to maintain data integrity.

### In progres/Future work

- Add message queue (RabbitMQ or Kafka) for more robust request processing.
- Integration tests
- Database triggers for timestamping rows on update.
- Row-level locking for concurrent database access.
- TOML file configuration.
- 2FA support.
- Prometheus metrics.