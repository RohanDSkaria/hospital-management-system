# 🏥 Hospital Management System API

A robust, production-ready REST API for hospital management built with Go, featuring role-based access control, JWT authentication, and comprehensive patient management capabilities.

## 🚀 Features

### 🔐 **Authentication & Authorization**
- **JWT-based authentication** with secure token management
- **Role-based access control (RBAC)** with two distinct user roles:
  - **Receptionist**: Full CRUD operations on patient records
  - **Doctor**: Read and update patient information (no deletion rights)
- **Password hashing** for secure credential storage
- **Middleware-based route protection** with automatic token validation

### 👥 **User Management**
- **User registration** with email validation and role assignment
- **Secure login** with JWT token generation
- **Profile management** with protected endpoint access
- **Input validation** with comprehensive error handling

### 🏥 **Patient Management**
- **Complete CRUD operations** for patient records
- **Comprehensive patient data** including:
  - Personal information (name, date of birth, address, contact)
  - Medical history tracking
  - Registration audit trail (who registered the patient)
- **UUID-based identification** for secure record management
- **Data validation** with proper error responses

### 📚 **API Documentation**
- **Interactive Swagger/OpenAPI documentation** with automatic generation
- **Comprehensive endpoint documentation** including:
  - Request/response schemas
  - Authentication requirements
  - Error code documentation
  - Parameter validation rules
- **Real-time API testing** through Swagger UI

## 🛠️ Technology Stack

### **Backend Framework**
- **Go 1.24.4** - High-performance, concurrent programming language
- **Gin Framework** - Fast HTTP web framework with excellent middleware support
- **GORM** - Feature-rich ORM library for Go

### **Database**
- **PostgreSQL** - Robust, ACID-compliant relational database
- **Auto-migration** for seamless schema management
- **Connection pooling** for optimal performance

### **Authentication & Security**
- **JWT (JSON Web Tokens)** for stateless authentication
- **bcrypt** password hashing for secure credential storage
- **Role-based middleware** for granular access control

### **API Documentation**
- **Swaggo/Swagger** for automatic API documentation generation
- **OpenAPI 2.0** specification compliance
- **Interactive testing interface**

### **Development Tools**
- **Environment configuration** with dotenv support
- **Structured logging** for debugging and monitoring
- **UUID generation** for secure identifier management

## 📁 Project Structure

```
hospital-management-system/
├── api/                    # HTTP handlers and middleware
│   ├── auth_handler.go     # Authentication endpoints
│   ├── patient_handler.go  # Patient management endpoints
│   └── middleware.go       # JWT and role-based middleware
├── cmd/
│   └── server/            # Application entry point
├── internal/              # Private application code
│   ├── auth/              # JWT token management
│   ├── database/          # Database connection and configuration
│   ├── model/             # Data models and GORM definitions
│   ├── repository/        # Data access layer
│   └── service/           # Business logic layer
├── pkg/                   # Public packages
│   └── utils/             # Utility functions
├── docs/                  # Auto-generated API documentation
└── go.mod                 # Go module dependencies
```

## 📖 API Documentation

### Available Endpoints

#### 🔐 Authentication
- `POST /api/v1/register` - Register a new user (receptionist/doctor)
- `POST /api/v1/login` - Authenticate and receive JWT token

#### 🏥 Patient Management

**Receptionist Routes** (Full CRUD access):
- `POST /api/v1/receptionist/patients` - Create new patient
- `GET /api/v1/receptionist/patients` - List all patients
- `GET /api/v1/receptionist/patients/{id}` - Get patient by ID
- `PUT /api/v1/receptionist/patients/{id}` - Update patient
- `DELETE /api/v1/receptionist/patients/{id}` - Delete patient

**Doctor Routes** (Read/Update access):
- `GET /api/v1/doctor/patients` - List all patients
- `GET /api/v1/doctor/patients/{id}` - Get patient by ID
- `PUT /api/v1/doctor/patients/{id}` - Update patient

#### 🏥 Health Check
- `GET /ping` - Server health check

## 🏗️ Architecture Patterns

### **Clean Architecture**
- **Separation of concerns** with distinct layers
- **Dependency injection** for testable components
- **Repository pattern** for data access abstraction
- **Service layer** for business logic encapsulation

### **Security Best Practices**
- **Input validation** at multiple layers
- **SQL injection prevention** through GORM
- **JWT token expiration** and validation
- **Role-based access control** with middleware
- **Secure password handling** with bcrypt

### **API Design**
- **RESTful conventions** for resource management
- **Consistent error handling** with proper HTTP status codes
- **Request/response validation** with structured schemas
- **Comprehensive documentation** with OpenAPI specification

## 📊 Database Schema

### Users Table
```sql
- id (UUID, Primary Key)
- full_name (VARCHAR(255), Not Null)
- email (VARCHAR(255), Unique, Not Null)
- password_hash (TEXT, Not Null)
- role (VARCHAR(20), Not Null) -- 'receptionist' or 'doctor'
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

### Patients Table
```sql
- id (UUID, Primary Key)
- full_name (VARCHAR(255), Not Null)
- date_of_birth (DATE)
- address (TEXT)
- contact_number (VARCHAR(20))
- medical_history (TEXT)
- registered_by_id (UUID, Foreign Key to Users)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
```

## 🔒 Security Features

- **JWT Authentication** with configurable expiration
- **Role-based Authorization** with middleware protection
- **Password Hashing** using bcrypt algorithm
- **Input Sanitization** and validation
- **SQL Injection Prevention** through ORM
- **CORS Protection** (configurable)
- **Environment-based Configuration** for sensitive data