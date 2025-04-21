# 📚 Go Backend Book Manager API

A clean and modular RESTful API built using **Go (Golang)** for managing books and users. This backend project demonstrates core backend development concepts including CRUD operations, clean architecture, token-based authentication, and database interaction using PostgreSQL and GORM.

---

## Features

- Create, read, update, and delete (CRUD) operations for books
- User registration and login with JWT authentication
- Secure route protection for authorized access
- Clean architecture and layered folder structure
- PostgreSQL database with GORM ORM
- RESTful API endpoints

---

## Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin Gonic (HTTP router)
- **ORM:** GORM
- **Authentication:** JWT (JSON Web Token)
- **Database:** PostgreSQL
- **Architecture:** Clean Architecture (Handlers, Services, Repositories, Models)

---

## Project Structure
Go_Backend_Book_Manager_Project/ 
├── cmd/ # Main application entry 
├── configs/ # Configuration files 
├── controllers/ # HTTP handlers (controllers) 
├── models/ # Data models 
├── repositories/ # Database operations 
├── routes/ # Route definitions 
├── services/ # Business logic 
├── utils/ # Utility functions (e.g., JWT handling) 
├── go.mod / go.sum # Go module definitions 
└── main.go # Main server file

## API Endpoints
**Auth**:
    POST /api/register – Register a new user
    POST /api/login – Login and receive JWT token
**Books**
    GET /api/books – Get all books (requires token)
    GET /api/books/:id – Get book by ID
    POST /api/books – Create new book (requires token)
    PUT /api/books/:id – Update book (requires token)
    DELETE /api/books/:id – Delete book (requires token)
