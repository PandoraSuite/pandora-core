# Pandora Core

**Pandora Core** is the central backend of the Open Source Pandora project, designed for efficient and secure management of API Keys for services. Its primary purpose is to validate, log, and control the consumption of requests made through these API Keys, providing an organized and straightforward management of clients, projects, services, and environments.

**Pandora Core** exposes two main interfaces:

* A RESTful HTTP API for administrative tasks and general configuration.
* An optimized gRPC-based service for rapid API Key validation from services.

It uses PostgreSQL as its persistent storage system, ensuring data consistency and convenient access to logged information.

Thanks to its architecture and clean design, **Pandora Core** offers an efficient and scalable solution, ideal for securely managing controlled access to multiple APIs or services.

## :zap: Features
* :lock: API Key Validation & Quota Control: Secure your APIs and manage usage effectively.
* :bar_chart: Request Logging and Usage Metrics: Gain insights into APIs consumption.
* :puzzle_piece: Multi-Project, Multi-Environment Architecture: Designed for versatile deployment and management across various setups.
* :globe_with_meridians: RESTful API for Admin Tasks: Easy administration and configuration via a standard HTTP interface.
* :high_voltage: High-Performance gRPC Service for Inline Validations: Fast and efficient validation.
* :card_file_box: PostgreSQL as Persistent: Reliable and scalable data storage.

## :rocket: Developer Setup

Ready to dive in? For a full guide on setting up your development environment, running the project, and debugging:

:point_right: See our comprehensive DEVELOPMENT.md guide.

## :compass: Project Status

Pandora Core is under active development. We're continuously working to enhance its capabilities.

## :handshake: Contributing

We welcome community contributions! Your ideas and efforts are highly valued.
Please read CONTRIBUTING.md to learn how to open issues, submit pull requests, and participate in the development process.

## :shield: Security

Found a security vulnerability? We take security seriously.

Please report it privately by following our disclosure guidelines in SECURITY.md. Do not open a public GitHub issue or pull request for security concerns.

## License

This project is licensed under the terms of the MIT license.
