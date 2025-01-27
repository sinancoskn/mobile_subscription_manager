# Subscription Management Platform

This project is a comprehensive **subscription management system** designed to handle receipt validation, subscription processing, and event handling for mobile apps. It supports integration with external receipt validation APIs, such as the Apple App Store or Google Play Store, and includes mock services for local testing and debugging.

## **Features**

- **Purchase API**: Handles subscription-related operations such as receipt validation, subscription status updates, and database interaction.
- **Event Processor**: A background worker system that processes subscription events, such as renewals, cancellations, and refunds, via RabbitMQ.
- **Mock Receipt API**: Simulates third-party receipt validation APIs for testing and development environments.
- **PostgreSQL Integration**: A robust database schema with support for partitioned tables and advanced queries.
- **RabbitMQ Integration**: Message queue support for asynchronous processing of subscription events.
- **Dockerized Services**: Simplified local and production deployment using Docker and Docker Compose.

---

## **Project Structure**

```plaintext
project-root/
├── purchase-api/         # PHP Phalcon API module for subscription management
├── event-processor/      # Golang worker module for subscription event processing
├── mock-receipt-api/     # Mock service for simulating receipt validation APIs
├── init/                 # Initialization scripts for PostgreSQL and RabbitMQ
├── docker-compose.yaml   # Docker Compose configuration for the project
```

### **Services**

1. **Purchase API (`purchase-api/`)**
   - Built with PHP Phalcon.
   - Provides RESTful endpoints for:
     - Validating purchase receipts.
     - Managing subscriptions and their statuses.
     - Interfacing with the PostgreSQL database.
   - Environment variables:
     ```plaintext
     POSTGRESQL_HOST=postgres
     POSTGRESQL_PORT=5432
     POSTGRESQL_DB=default_db
     POSTGRESQL_USER=root
     POSTGRESQL_PASS=root_password
     RABBITMQ_HOST=rabbitmq
     RABBITMQ_PORT=5672
     STORE_API_HOST=http://localhost:1234
     TOKEN_SECRET=your_secret
     ```

2. **Event Processor (`event-processor/`)**
   - Built in Go.
   - **Components**:
     - **Worker Manager**:
       - Listens to HTTP requests for triggers and creates actions and batches for chunked operations.
       - Manages the lifecycle of actions and batches, ensuring workers can process them efficiently.
       - Provides a monitoring interface accessible at `http://127.0.0.1:9090/` for tracking actions and batches.
     - **Worker**:
       - Registers with the Worker Manager using a unique ID.
       - Checks for pending batches, locks a batch, and processes it using the mock API.
       - Handles batch processing and status updates to ensure reliable task execution.
     - **Callback (Optional)**:
       - Listens to RabbitMQ for subscription events.
       - Handles third-party webhook calls to notify external systems about subscription updates or events.
   - Environment variables:
     ```plaintext
     DATABASE_DSN=host=postgres user=root password=root_password dbname=default_db port=5432
     RABBITMQ_URL=amqp://admin:admin@rabbitmq:5672/
     MANAGER_PORT=9090
     MANAGER_HOST=event-processor
     ```

3. **Mock Receipt API (`mock-receipt-api/`)**
   - Simulates third-party receipt validation services (Apple and Google).
   - Useful for local development and testing.

4. **PostgreSQL and RabbitMQ (`init/`)**
   - **PostgreSQL**:
     - Contains `init.sql` for setting up the database schema, including partitioned tables for `devices` and `subscriptions`.
     - Supports seed data via additional SQL files (e.g., `01_apps.sql`, `02_webhooks.sql`).
   - **RabbitMQ**:
     - Includes `rabbitmq-definitions.json` for pre-configuring queues, exchanges, and bindings.
     - Sample queue: `subscription_events`.

---

## **Setup and Deployment**

### **Prerequisites**
- Docker and Docker Compose installed on your system.
- Basic knowledge of Docker and containerized environments.

### **Steps to Run the Project**

1. **Clone the Repository**
   ```bash
   git clone <repository-url>
   cd project-root
   ```

2. **Build and Start Services**
   ```bash
   docker-compose up --build
   ```

3. **Verify Services**
   - Access RabbitMQ Management UI: `http://localhost:15672/` (Username: `admin`, Password: `admin`).
   - Access Purchase API: `http://localhost:8080/`.
   - Access Mock Receipt API: `http://localhost:1234/`.
   - Event Processor Worker Manager Monitoring: `http://127.0.0.1:9090/`.

4. **Testing**
   - Use tools like `Postman` or `curl` to interact with the APIs.
   - Check RabbitMQ queues and PostgreSQL tables for processed events.

---

## **Environment Variables**

### **RabbitMQ**
```plaintext
RABBITMQ_DEFAULT_USER=admin
RABBITMQ_DEFAULT_PASS=admin
RABBITMQ_URL=amqp://admin:admin@rabbitmq:5672/
```

### **PostgreSQL**
```plaintext
POSTGRESQL_HOST=postgres
POSTGRESQL_PORT=5432
POSTGRESQL_DB=default_db
POSTGRESQL_USER=root
POSTGRESQL_PASS=root_password
```

### **Event Processor**
```plaintext
DATABASE_DSN=host=postgres user=root password=root_password dbname=default_db port=5432
MANAGER_HOST=event-processor
MANAGER_PORT=9090
```

---

## **Key Features**

1. **Partitioned Database Tables**
   - Optimized PostgreSQL schema for managing large datasets efficiently.
   - Tables like `devices` and `subscriptions` are partitioned by hash for faster queries.

2. **Scalable Event Processing**
   - RabbitMQ for queue-based asynchronous processing.
   - Event processor supports worker scaling for high throughput.

3. **Modular Architecture**
   - Independent services for API handling (`purchase-api`), event processing (`event-processor`), and testing (`mock-receipt-api`).

---

## **Future Improvements**
- Add monitoring and alerting (e.g., Prometheus, Grafana).
- Integrate CI/CD pipelines for automated testing and deployment.
- Extend mock APIs to simulate edge cases.

---

## **Contributing**
Feel free to submit issues and pull requests. Contributions are welcome!

---

## **License**
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

