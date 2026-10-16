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

**Note**: The UI is still under development. Some features, such as adding task dependencies, checking off tasks, or adding other users to domains, are not yet accessible through the UI.

## Features

- User registration and login

- Creation and management of domains and tasks

- Adding users to domains

- Task assignment to domain members

- Task dependencies, with circular dependencies prevented

- Identification of actionable tasks based on completed dependencies


## Technical Implementation

- RESTful API built with Go

- PostgreSQL database with relational constraints

- JWT-based authentication and protected endpoints

- Database migrations managed with Goose

- Dependency cycle detection using depth-first search

- Comprehensive error handling and logging

- Front-end in HTML/CSS/Javascript


