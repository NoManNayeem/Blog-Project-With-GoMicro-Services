# Blogging Website Backend with Go (Microservices)

## Project Objectives

### 1. **User Management Microservice**
   - **Authentication (Login/Register)**: 
     - Users can register an account with a unique username and password.
     - Login system to authenticate users using JWT (JSON Web Tokens).
   - **Profile Management**: 
     - Users can create and update their profiles (with fields like Full Name, Bio, etc.).
   - **Role-Based Access Control**:
     - Two roles: `Writer` and `Admin`.
     - **Writer**: Can create, read, update, and delete (CRUD) their own blogs.
     - **Admin**: Can manage all users, their blogs, and assign roles.

### 2. **Blog Microservice**
   - **CRUD Operations on Blogs**:
     - Writers can create, edit, update, and delete their own blogs.
     - Admins have full access to all blogs from any writer.
   - **Blog Management**: 
     - Simple text-based blog creation.
     - Ensure that blogs are only editable by the correct user based on their role.

### 3. **Microservices Architecture**
   - Develop two separate Go services:
     - **User Management Service**: Responsible for authentication, role management, and profile management.
     - **Blog Service**: Handles the CRUD operations for blogs.
   - Use Go standard library without any frameworks.
   - Each service runs independently and communicates via REST API.

### 4. **Security and Session Management**
   - Implement JWT-based authentication for stateless session management.
   - Protect blog operations (CRUD) based on user roles (Writer/Admin).
   - Ensure secure password storage using hashing algorithms (e.g., bcrypt).
