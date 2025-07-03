# POS (Point of Sale) API

This is a comprehensive Point of Sale (POS) API built with Go and Gin. It provides endpoints for managing users, products, categories, promotions, and orders.

## Features

*   User authentication (registration, login, refresh token)

For more detailed functional requirements, please see the [Functional Requirements Document (FRD)](./docs/FRD.md).
*   Role-based authorization (admin, user)
*   CRUD operations for users, products, and categories
*   Stock management
*   Product and cart promotions
*   Order management

## Getting Started

### Prerequisites

*   Go (version 1.24.4 or later)
*   PostgreSQL

### Installation

1.  Clone the repository:
    ```sh
    git clone https://github.com/esaprakoso/pos-api.git
    ```
2.  Install dependencies:
    ```sh
    go mod tidy
    ```
3.  Create a `.env` file in the root directory and add the following environment variables:
    ```
    DB_HOST=localhost
    DB_USER=your_db_user
    DB_PASSWORD=your_db_password
    DB_NAME=inventory
    DB_PORT=5432
    JWT_SECRET=your_jwt_secret
    ```
4.  Run the application:
    ```sh
    go run main.go
    ```

The application will be running at `http://localhost:3000`.

## API Endpoints

All endpoints are prefixed with `/api`.

### Auth

*   `POST /auth/register`: Register a new user.
*   `POST /auth/login`: Login a user.
*   `POST /auth/refresh`: Refresh a user's token.

### Users

*   `GET /users`: Get all users (admin only).
*   `GET /users/:id`: Get a user by ID (admin only).
*   `PATCH /users/:id`: Update a user by ID (admin only).
*   `DELETE /users/:id`: Delete a user by ID (admin only).

### Profile

*   `GET /profile`: Get the current user's profile.
*   `PATCH /profile`: Update the current user's profile.
*   `PATCH /profile/password`: Update the current user's password.

### Products

*   `GET /products`: Get all products.
*   `GET /products/:id`: Get a product by ID.
*   `POST /products`: Create a new product (admin only).
*   `PUT /products/:id`: Update a product by ID (admin only).
*   `DELETE /products/:id`: Delete a product by ID (admin only).
*   `PATCH /products/:id/stock`: Update product stock (admin only).

### Categories

*   `GET /categories`: Get all categories.
*   `GET /categories/:id`: Get a category by ID.
*   `POST /categories`: Create a new category (admin only).
*   `PUT /categories/:id`: Update a category by ID (admin only).
*   `DELETE /categories/:id`: Delete a category by ID (admin only).

### Product Promotions

*   `GET /product-promotions`: Get all product promotions.
*   `GET /product-promotions/:id`: Get a product promotion by ID.
*   `POST /product-promotions`: Create a new product promotion (admin only).
*   `PUT /product-promotions/:id`: Update a product promotion by ID (admin only).
*   `DELETE /product-promotions/:id`: Delete a product promotion by ID (admin only).

### Cart Promotions

*   `GET /cart-promotions`: Get all cart promotions.
*   `GET /cart-promotions/:id`: Get a cart promotion by ID.
*   `POST /cart-promotions`: Create a new cart promotion (admin only).
*   `PUT /cart-promotions/:id`: Update a cart promotion by ID (admin only).
*   `DELETE /cart-promotions/:id`: Delete a cart promotion by ID (admin only).

### Orders

*   `GET /orders`: Get all orders.
*   `GET /orders/:id`: Get an order by ID.
*   `POST /orders`: Create a new order.

## Database Schema

The database schema consists of the following tables:

*   `users`
*   `products`
*   `categories`
*   `stock_transactions`
*   `product_promotions`
*   `cart_promotions`
*   `orders`
*   `order_items`

For more details, see the `models` directory.