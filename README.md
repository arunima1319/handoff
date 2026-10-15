# Handoff 

A dependency-focused collaborative task management app. Built using Go for the back-end and PostgreSQL as the DBMS. 

## Motivation 

Collaborating with family and friends on day-to-day tasks was getting difficult - tasks kept getting postponed and ultimately forgotten because something *else* needed to be done first. 

Handoff ensures these dependencies are kept under control, tasks have clear assignees, and that cyclic dependencies are prevented. 

## Quick Start 

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

- SQL queries generated using sqlc

- Database migrations managed with Goose

- Dependency cycle detection using depth-first search

- Comprehensive error handling and structured logging

- Front-end in HTML/CSS/Javascript


