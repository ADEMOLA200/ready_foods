### Danas Food App

This is the README file for the Danas Food App. Below you will find an overview of the application structure, functionalities, and how to run it.

---

### Overview

The Danas Food App is a web application built to facilitate food ordering. It provides endpoints for creating, retrieving, updating, and deleting food orders. The application is developed using Go programming language and utilizes Fiber framework for building web APIs.

### Application Structure

The application consists of the following main components:

- **`controller` Package**: Contains handlers for different API endpoints related to food orders.
  
- **`main` Package**: Contains the entry point of the application and setup for middleware and routes.

- **`database` Package**: Handles database connection and initialization.

- **`middleware` Package**: Provides middleware functions for logging, error handling, and CORS.

- **`routes` Package**: Defines API routes and registers them with the Fiber app.

### Features

- **Order Management**: Users can create, retrieve, update, and delete food orders.
  
- **Authentication Middleware**: Middleware for user authentication (currently commented out) can be easily integrated.

- **Logging Middleware**: Middleware to log HTTP requests and responses for debugging and monitoring purposes.

- **Error Handling Middleware**: Middleware to handle errors and return appropriate HTTP responses.

- **Cross-Origin Resource Sharing (CORS) Middleware**: Middleware to enable Cross-Origin Resource Sharing for allowing requests from other origins.

### Running the Application

To run the Danas Food App, follow these steps:

1. Ensure you have Go installed on your system.

2. Clone the repository:

   ```
   git clone https://github.com/ADEMOLA200/danas-food.git
   ```

3. Navigate to the project directory:

   ```
   cd danas-food
   ```

4. Install dependencies:

   ```
   go mod tidy
   ```

5. Configure your database connection details in the `database.ConnectDB()` function in `main.go`.

6. Run the application:

   ```
   go run main.go
   ```

The application will start and listen for requests on port 3000 by default.

### API Documentation

The API endpoints provided by the Danas Food App can be found in the `controller` package. Each endpoint is associated with a specific handler function responsible for processing incoming requests.

---

This README provides an overview of the Danas Food App, its structure, features, and instructions for running the application. For detailed documentation on specific API endpoints and functionalities, please refer to the comments and code implementation in the respective packages and files.