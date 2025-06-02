# Pandora Core

**Pandora Core** is the central backend of the Open Source Pandora project, designed for efficient and secure management of API Keys for services. Its primary purpose is to validate, log, and control the consumption of requests made through these API Keys, providing an organized and straightforward management of clients, projects, services, and environments.

It uses PostgreSQL as its persistent storage system, ensuring data consistency and convenient access to logged information.

Thanks to its architecture and clean design, **Pandora Core** offers an efficient and scalable solution, ideal for securely managing controlled access to multiple APIs or services.

## :package: Production Deployment

> :warning: **Beta Release**: This Docker image is currently in beta. Feel free to try it out and share your feedback with the community.

1. **Pull the Docker image**

   ```bash
   docker pull madpaydev/pandora-core:v0.1.0-beta.1
   ```

2. **Run the container**

   ```bash
   docker run -d \
     --name pandora-core \
     -p 80:80 \
     -p 50051:50051 \
     -e PANDORA_DB_PASSWORD="<postgresql_password>" \
     madpydev/pandora-core:v0.1.0-beta.1
   ```

   **Pandora Core** will start and expose ports **80** (HTTP) and **50051** (gRPC).

   > :rocket: **Tip:** To customize Pandora even further, check out the [Pandora Environment Variables](#gear-pandora-environment-variables) section below!

### :gear: Pandora Environment Variables

* **`PANDORA_DB_PASSWORD`** (required) Set the password for the Pandora database. There is no default—this variable **must** be provided.

* **`PANDORA_DB_NAME`** (optional) Rename the database.
  * Default: `pandora`

* **`PANDORA_DB_USER`** (optional) Change the database username.
  * Default: `pandora`

* **`PANDORA_DIR`** (optional) Specify a custom directory for storing Pandora’s configuration and secrets.
  * Default: `/etc/pandora`

* **`PANDORA_JWT_SECRET`** (optional) Provide a fixed secret key for signing JWT authentication tokens. If omitted, Pandora generates a random one at startup (not recommended for consistent development).
  * Default: (randomly generated on each startup)

* **`PANDORA_HTTP_PORT`** (optional) Change the HTTP server’s listening port.
  * Default: `80`

* **`PANDORA_GRPC_PORT`** (optional) Change the gRPC server’s listening port.
  * Default: `50051`

* **`PANDORA_EXPOSE_VERSION`** (optional) Control whether Pandora reveals its version in HTTP responses.
  * Default: `true`

## :thought_balloon: Use Cases

**Pandora Core** provides centralized API Key generation, validation, and quota management for client projects and environments. It addresses the need to secure and monitor access to your services—particularly microservices and AI agents—by enabling fine-grained control over usage and simplifying authentication workflows.

* **Quickly onboard new customers** to your agents and LLM-based microservices by issuing API keys and applying specific usage limits for each of your customers' environments.

* **Centralize authentication** to all your services, with quota enforcement through generated API keys.
 
* **Monitor consumption** of your services in real time to optimize billing and adapt pricing models.

## :computer: Upcoming Features

* **Admin UI (Beta)**
  A web-based dashboard for managing resources (clients, projects, environments, API keys).
  It will be available in the coming weeks in beta. The UI will allow you to:

  * Visualize active API keys.
  * Configure quotas.
  * Manage clients, projects, environments and revoke/regenerate keys without using the REST API directly.

* **Python SDK (Beta)**
  A Python package to simplify integration with **Pandora Core**. Key features include:

  * **FastAPI integration** for seamless authentication and authorization out of the box.
  * **Automatic request** logging and quota enforcement at the service level.
  * **Utility functions** for client, project, service, environment, and API Key management.
  * **Quick setup** install via `pip install pandora-sdk`.

## :rocket: Developer Setup

Ready to dive in? For a full guide on setting up your development environment, running the project, and debugging:

:point_right: See our comprehensive [DEVELOPMENT.md](./DEVELOPMENT.md) guide.

## :compass: Project Status

**Pandora Core** is under active development.

We're continuously working to enhance its capabilities.

## :handshake: Contributing

We welcome community contributions! Your ideas and efforts are highly valued.

Please read [CONTRIBUTING.md](./CONTRIBUTING.md) to learn how to open issues, submit pull requests, and participate in the development process.

## :shield: Security

Found a security vulnerability? We take security seriously.

Please report it privately by following our disclosure guidelines in [SECURITY.md](./SECURITY.md). Do not open a public GitHub issue or pull request for security concerns.

## License

This project is licensed under the terms of the MIT license.
