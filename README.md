# Pandora Core

**Pandora Core** is the central backend of the Open Source Pandora project, designed for efficient and secure management of API Keys for services. Its primary purpose is to validate, log, and control the consumption of requests made through these API Keys, providing an organized and straightforward management of clients, projects, services, and environments.

**Pandora Core** exposes two main interfaces:

* A RESTful HTTP API for administrative tasks and general configuration.
* An optimized gRPC-based service for rapid API Key validation from services.

It uses PostgreSQL as its persistent storage system, ensuring data consistency and convenient access to logged information.

Thanks to its architecture and clean design, **Pandora Core** offers an efficient and scalable solution, ideal for securely managing controlled access to multiple APIs or services.

👉 See [DEVELOPMENT.md](./DEVELOPMENT.md) for developer setup, debugging and contributing instructions.
