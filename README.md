# Handoff 

A dependency-focused collaborative task management app. Built using Go for the back-end and PostgreSQL as the DBMS. 

## Motivation 

Collaborating with family and friends on day-to-day tasks was getting difficult - tasks kept getting postponed and ultimately forgotten because something *else* needed to be done first. 

Handoff ensures these dependencies are kept under control, tasks have clear assignees, and that cyclic dependencies are prevented. 

## Quick Start 

### Running Locally with Docker

**Prerequisite:** [Docker Desktop](https://www.docker.com/products/docker-desktop/) (or Docker Engine with Docker Compose).

1. Clone the repository:

   ```bash
   git clone https://github.com/arunima1319/handoff.git
   cd handoff
   ```

2. Create an environment file:

   ```bash
   cp .env.example .env
   ```

   Set `SECRET` to a secure random value in `.env`.

3. Build and start the application:

   ```bash
   docker compose up --build
   ```

4. Open [http://localhost:8080](http://localhost:8080) in your browser.

Docker Compose starts the PostgreSQL database, applies the Goose migrations, and runs the Go application.

To stop the application, press `Ctrl+C`, or run:

```bash
docker compose down
```

Database data is preserved between restarts using a Docker volume.


## Features

- Task dependency management with cycle detection using depth-first search (DFS)
- Identification of actionable tasks based on completed dependencies
- Authentication and authorization using Argon2id password hashing and JWTs
- RESTful API developed in Go for user, domain, task, and dependency management
- Relational data modelling involving five interconnected tables with database-level constraints to enforce data integrity
- Comprehensive error handling and logging, for easy debugging


