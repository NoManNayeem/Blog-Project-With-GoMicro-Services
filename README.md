
# Blog Project with Go Microservices

This project implements a simple blog system using microservices architecture. There are two main services:
1. **User Management Service**: Handles user registration, login, profile management, and role-based access control.
2. **Blog Service**: Allows users to create, update, delete, and view blogs.

## Services
### 1. User Management Service
- **Endpoints**:
  - `/register`: Register a new user.
  - `/login`: Log in a user and generate a JWT token.
  - `/profile`: Create or update a user's profile.
  - `/admin`: Manage users (Admin only access).
  
- **Roles**:
  - `Writer`: Can manage (CRUD) their own blogs.
  - `Admin`: Can manage all users and their blogs.

### 2. Blog Service
- **Endpoints**:
  - `/blogs`: Get all blogs.
  - `/blogs/create`: Create a new blog (Authenticated users).
  - `/blogs/update`: Update a blog post (Authenticated users).
  - `/blogs/delete`: Delete a blog post (Authenticated users).
  
## Project Structure

```
blogging-backend/
    ├── user-management/
    │    ├── main.go
    │    └── user/
    │        ├── handler.go
    │        ├── model.go
    │        └── service.go
    └── blogs/
         ├── main.go
         └── blog/
             ├── handler.go
             ├── model.go
             └── service.go
```

## Swagger API Documentation

Swagger is used to document both services. After running the services, you can access the Swagger UI to view and test the APIs.

- **Swagger UI for User Management**: `http://localhost:8000/swagger/index.html`
- **Swagger UI for Blog Service**: `http://localhost:8001/swagger/index.html`

## How to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/NoManNayeem/Blog-Project-With-GoMicro-Services.git
   ```

2. Navigate to each service directory and install dependencies:
   
   - For User Management:
     ```bash
     cd user-management
     go mod tidy
     ```

   - For Blog Service:
     ```bash
     cd blogs
     go mod tidy
     ```

3. Generate Swagger documentation:
   
   - In both service directories, run:
     ```bash
     swag init
     ```

4. Run each service:

   - **User Management**: 
     ```bash
     go run main.go
     ```
     This will run the user management service on `http://localhost:8000`.

   - **Blog Service**:
     ```bash
     go run main.go
     ```
     This will run the blog service on `http://localhost:8001`.

## Technologies Used
- **Go**: Language for building microservices.
- **MySQL**: Database for user and blog management.
- **Swagger**: API documentation.
- **JWT**: Authentication and role-based access control.

## License
This project is licensed under the MIT License.

