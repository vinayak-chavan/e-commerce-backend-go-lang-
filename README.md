# My E-commerce Project

## Description

This project is a basic e-commerce application built with Go, Gin, and GORM. It supports user registration, login, cart management, and order processing.

## Prerequisites

- Go
- Git
- PostgreSQL

## Setup Instructions

1. **Clone the Repository**

   ```sh
   git clone https://github.com/vinayak-chavan/e-commerce-backend-go-lang-
   cd e-commerce
   ```
2. **Install Dependencies**

   ```sh
   go mod tidy
   ```

3. **Set Up Environment Variables**

Change .env file in the root directory:

DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name
DB_PORT=5432
DB_SSLMODE=disable
DB_TIMEZONE=Asia/Kolkata
PORT=8000
SECRET=ThisIsSecretKey

4. **Run the Project**

   ```sh
   go run main.go
   ```
