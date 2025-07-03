# Inventory Management API

This is a simple inventory management API built with Go and Gin. It provides endpoints for managing users, products, categories, and stocks.

## Features

* User authentication (registration, login, refresh token)
* Role-based authorization (admin, user)
* CRUD operations for users, products, and categories
* Stock management

## Getting Started

### Prerequisites

* Go (version 1.24.4 or later)
* PostgreSQL

### Installation

1. Clone the repository:
   ```sh
   git clone https://github.com/esaprakoso/post-api.git
   ```
2. Install dependencies:
   ```sh
   go mod tidy
   ```
3. Create a `.env` file in the root directory and add the following environment variables:
   ```
   DB_HOST=localhost
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=inventory
   DB_PORT=5432
   JWT_SECRET=your_jwt_secret
   ```
4. Run the application:
   ```sh
   go run main.go
   ```

The application will be running at `http://localhost:3000`.

## API Endpoints

All endpoints are prefixed with `/api`.

### Auth

* `POST /auth/register`: Register a new user.
* `POST /auth/login`: Login a user.
* `POST /auth/refresh`: Refresh a user's token.

### Users

* `GET /users`: Get all users (admin only).
* `GET /users/:id`: Get a user by ID (admin only).
* `PATCH /users/:id`: Update a user by ID (admin only).

### Profile

* `GET /profile`: Get the current user's profile.
* `PATCH /profile`: Update the current user's profile.
* `PATCH /profile/password`: Update the current user's password.



### Products

* `GET /products`: Get all products.
* `GET /products/:id`: Get a product by ID.
* `POST /products`: Create a new product (admin only).
* `PUT /products/:id`: Update a product by ID (admin only).
* `DELETE /products/:id`: Delete a product by ID (admin only).

### Stocks

* `GET /stocks`: Get all stocks.
* `POST /stocks`: Create or update stock (admin only).

### Categories

* `GET /categories`: Get all categories.
* `GET /categories/:id`: Get a category by ID.
* `POST /categories`: Create a new category (admin only).
* `PUT /categories/:id`: Update a category by ID (admin only).
* `DELETE /categories/:id`: Delete a category by ID (admin only).

## Database Schema

The database schema consists of the following tables:

* `users`

* `products`
* `categories`
* `stocks`

For more details, see the `models` directory.
