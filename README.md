
# CrawlxyAI

A simple Go solution for web scraping, data processing, and AI-driven interactions. This project brings together powerful tools to make it easy to scrape web pages, turn the scraped data into meaningful insights, and interact with the data through a REST API built with the Gin framework. Keep in mind that there are quite some moving parts

## Table of Contents

- [CrawlxyAI](#crawlxyai)
  - [Table of Contents](#table-of-contents)
  - [Tools and Technologies Used](#tools-and-technologies-used)
  - [Installation and Setup](#installation-and-setup)
    - [Prerequisites](#prerequisites)
    - [Step-by-Step Guide](#step-by-step-guide)
  - [Running the Application](#running-the-application)
    - [Database Management](#database-management)
  - [Project Components](#project-components)
    - [service](#service)
    - [database](#database)
    - [envs](#envs)
    - [llms](#llms)
  - [Web Scraping and Embeddings](#web-scraping-and-embeddings)
    - [Scrapper](#scrapper)
  - [REST API](#rest-api)
    - [API Endpoint](#api-endpoint)


## Tools and Technologies Used

Here's a list of the main tools and technologies we're using in this project:

- **Go**: Our primary programming language.
- **Colly**: A powerful web scraping framework in Go.
- **Goquery**: Makes it easy to parse and manipulate HTML.
- **LangChainGo**: For creating embeddings from HTML documents.
- **Gin**: A fast, minimalistic web framework for Go.
- **Docker**: For containerizing our application and database.
- **PostgreSQL**: Our database of choice. Used the pgstore vector extension.
- **SQLC**: Generates type-safe Go code from SQL queries.
- **Makefile**: Automates our build and deployment processes.
- **Docker Compose**: Manages our multi-container Docker applications.

## Installation and Setup

### Prerequisites

Make sure you have these installed before you begin:

- Docker
- Docker Compose
- Go (latest version)

### Step-by-Step Guide

1. **Clone the repository:**
    ```
    git clone https://github.com/wicfasho/crawl-xy-ai
    cd crawl-xy-ai
    ```

2. **Set up environment variables:**
    Create a `.env` file in the root directory and add the necessary configurations based on the templates in `envs/global.env`, `envs/pgadmin.env`, and `envs/postgres.env`.

3. **Build and run the application using Docker Compose:**
    ```
    docker-compose up --build
    ```

## Running the Application

Once you've got everything set up, the application should be up and running. You can interact with it through various routes and endpoints.

### Database Management

- **Migrations**: Manage your database migrations with the following commands:
    ```
    make migrate-up
    make migrate-down
    ```

- **PgAdmin**: Use PgAdmin to manage your PostgreSQL database. Make sure your `pgadmin.env` file is configured properly.

## Project Components

### service

This directory contains the main backend service code:

- `main.go`: The entry point for the application. It initializes the database, starts the scrapper service, and sets up the API server.
- `database/database.go`: Handles database connections and operations.
- `scrapper/scrapper.go`: Handles web scraping using Colly and Goquery, and generates embeddings using LangChainGo.
- `server/server.go`: Sets up the HTTP server using Gin.
- `api/v1/api.go`: Defines API routes and handlers.
- `sqlc/`: Contains SQLC-generated code for type-safe database interactions.
- `routes/routes.go`: Defines the application routes.
- `handlers/ask/ask.go`: Handles the /ask endpoint.

### database

This directory contains database-related files:

- `migrations/`: SQL files for database schema migrations.
- `postgres/`: Dockerfile and SQL scripts for setting up PostgreSQL.

### envs

This directory contains environment configuration files:

- `pgadmin.env`: Environment variables for PgAdmin.
- `postgres.env`: Environment variables for PostgreSQL.
- `global.env`: Global environment variables for the application.

### llms

This directory contains language models and related files:

- `ollama/models/`: Directory for storing various language model files. A `Modelfile` template for ollama can be used to spin up custom models. Ollama should be running on your host machine.

## Web Scraping and Embeddings

The application uses Colly and Goquery to scrape web pages. It processes HTML documents to extract meaningful content, which is then converted into embeddings using the LangChainGo library. These embeddings are stored in PostgreSQL for further analysis and querying.

### Scrapper

The `scrapper/scrapper.go` file contains the logic for scraping web pages and processing the HTML content. Colly is used for efficient web scraping, and Goquery is used to navigate and manipulate the HTML documents.

## REST API

The application provides a REST API built with the Gin framework. The `localhost:8080/api/v1/ask` endpoint allows users to interact with the stored embeddings and models.

### API Endpoint

- **/ask**: This endpoint handles user queries and returns relevant responses based on the embeddings stored in the database. The handler for this endpoint is defined in `handlers/ask/ask.go`.
