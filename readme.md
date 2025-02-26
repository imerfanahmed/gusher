# Gusher - A Go-Based Real-Time WebSocket Server (Pusher Alternative)

Welcome to Gusher, a lightweight, Go-based real-time WebSocket server designed as an alternative to Pusher Channels. Gusher enables developers to build scalable, real-time applications with features like channel subscriptions, event broadcasting, and app management, mimicking Pusher’s functionality while offering flexibility and control through Go.

## Table of Contents
- [Gusher - A Go-Based Real-Time WebSocket Server (Pusher Alternative)](#gusher---a-go-based-real-time-websocket-server-pusher-alternative)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [Installation](#installation)
    - [Prerequisites](#prerequisites)
    - [Steps](#steps)
  - [Usage](#usage)
    - [Running the Server](#running-the-server)
    - [WebUI Dashboard](#webui-dashboard)
    - [API Endpoints](#api-endpoints)
  - [Configuration](#configuration)
  - [Contributing](#contributing)
    - [Reporting Issues](#reporting-issues)
  - [License](#license)
  - [Support](#support)

## Overview
Gusher is a Go-based server that provides real-time communication through WebSockets, supporting Pusher-like protocols for subscribing to channels, broadcasting events, and managing applications. It includes a WebUI for app management, channel debugging, and event logging, making it ideal for developers building chat applications, live updates, or collaborative tools.

The project leverages Go’s concurrency model, Redis for caching, and MySQL for persistence, ensuring high performance and scalability.

## Features
- **Real-Time WebSockets**: Supports Pusher-compatible events like `pusher:subscribe`, `client_event`, `channel_occupied`, `channel_vacated`, `member_added`, and `member_removed`.
- **App Management**: Create, view, and manage applications with unique keys, secrets, and IDs.
- **Channel Debugging**: Monitor active channels, occupancy, and subscriptions in real-time.
- **Event Logging**: Comprehensive logging of all events with timestamps, filterable by channel or event type.
- **WebUI Dashboard**: Interactive web interface for app management, channel debugging, and event monitoring.
- **Scalability**: Integrates with Redis for caching and MySQL for persistent storage.

## Installation

### Prerequisites
- **Go**: Version 1.18 or higher (`go version`).
- **MySQL**: Version 8.0 or higher for database storage.
- **Redis**: Version 6.0 or higher for caching.

### Steps
1. **Clone the Repository**:
   ```bash
   git clone https://github.com/imerfanahmed/gusher.git
   cd gusher
   ```

2. **Install Dependencies**:
   Run the following to fetch Go dependencies:
   ```bash
   go mod tidy
   ```

3. **Set Up Database**:
   - Install MySQL and create a database:
     ```sql
     CREATE DATABASE gusher;
     ```
   - Apply migrations (run from project root):
     ```bash
     go run cmd/main.go --migrate
     ```
   - Ensure the `apps` and `webhooks` tables are created (see `internal/database/migrations`).

4. **Configure Redis**:
   - Install Redis and start the server:
     ```bash
     redis-server
     ```
   - No configuration changes are typically needed, but ensure it’s accessible on `localhost:6379`.

5. **Set Environment Variables**:
   Create a `.env` file in the project root:
   ```env
   DB_DSN=user:password@tcp(localhost:3306)/gusher
   REDIS_HOST=localhost
   REDIS_PORT=6379
   HOST=localhost
   PORT=8080
   ```

## Usage

### Running the Server
Start the Gusher WebSocket server:
```bash
go run cmd/main.go
```
- The server will listen on `localhost:8080` by default.
- Access the WebUI at `http://localhost:8080` (once implemented in your Go server or hosted separately).

### WebUI Dashboard
Gusher includes a WebUI for managing apps, debugging channels, and logging events:
- **Open in Browser**: Navigate to `http://localhost:8080` (or the hosted WebUI URL).
- **Features**:
  - **Apps Tab**: Create and list apps with keys, secrets, and IDs.
  - **Channels Tab**: View active channels, filter by name, and see occupancy.
  - **Events Tab**: Real-time event logger with filtering by channel or event type.
  - **Debug Console**: Trigger custom events (e.g., `client_event`) on channels for testing.

### API Endpoints
Gusher supports REST and WebSocket APIs. Example endpoints (to be implemented in your Go server):
- **GET `/apps`**: List all applications.
- **POST `/apps`**: Create a new app with `{key, secret, id}`.
- **GET `/channels`**: List active channels and their occupancy.
- **WS `/app/{key}`**: WebSocket endpoint for real-time connections (e.g., `ws://localhost:8080/app/app_key`).

## Configuration
Customize Gusher via environment variables in `.env`:
- `DB_DSN`: MySQL connection string (e.g., `user:password@tcp(localhost:3306)/gusher`).
- `REDIS_HOST/REDIS_PORT`: Redis server details (default: `localhost:6379`).
- `HOST/PORT`: Server host and port (default: `localhost:8080`).

## Contributing
We welcome contributions to Gusher! To contribute:
1. Fork the repository:
   ```bash
   git clone https://github.com/imerfanahmed/gusher.git
   cd gusher
   git checkout -b feature/your-feature
   ```
2. Make changes and commit:
   ```bash
   git commit -m "Add feature: Your description"
   ```
3. Push to your fork and submit a pull request:
   ```bash
   git push origin feature/your-feature
   ```

### Reporting Issues
- File issues on GitHub:
  - Go to [https://github.com/imerfanahmed/gusher/issues](https://github.com/imerfanahmed/gusher/issues).
  - Include:
    - Description of the issue.
    - Steps to reproduce.
    - Expected vs. actual behavior.

## License
Gusher is licensed under the MIT License. See the `LICENSE` file for details.

## Support
- **Issues**: Use [GitHub Issues](https://github.com/imerfanahmed/gusher/issues).
- **Community**: Join discussions on [GitHub Discussions](https://github.com/imerfanahmed/gusher/discussions) (to be set up).
- **Documentation**: Refer to this README and in-code comments for usage.