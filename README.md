# [Let's Go Further book by Alex Edwards](https://lets-go-further.alexedwards.net/)

Welcome to the README file for the Greenlight application developed as part of the book "Let's Go Further" by Alex Edwards. This document provides an overview of the application, highlights the incredible value of the book, and showcases the additional features and improvements implemented beyond the step-by-step tutorial.

## About the Book

"Let's Go Further" is an exceptional resource for software engineers who want to delve into advanced patterns for building APIs and web applications in Go. Authored by [Alex Edwards](https://www.alexedwards.net/), this book goes beyond the basics and equips developers with the knowledge to build robust and scalable applications.

See more here: https://lets-go-further.alexedwards.net/

## Additional Features and Improvements


1. **Docker Compose**: Included a Docker Compose file to simplify the setup process and allow running the database locally. 
2. **API Documentation**: [TODO]


## Getting Started

1. **Prerequisites**: Make sure you have Go (version 1.21.0 or higher) installed on your machine.
2. **Clone the Repository**: Run the following command to clone the repository:

    ```shell
    git clone git@github.com:jessicatarra/greenlight.git
    ```

3. **Install Dependencies**: Change to the project directory and use the following command to install the dependencies:

    ```shell
    go mod tidy
    ```

4. **Environment Variables**: Create a `.envrc` file in the project directory and add the following environment variables with their corresponding values:

```shell
export GREENLIGHT_DB_DSN=

export MIGRATION_URL=

export SMTP_HOST=

export SMTP_PASSWORD=

export SMTP_PORT=

export SMTP_SENDER=

export SMTP_USERNAME=

export CORS_TRUSTED_ORIGINS=

```

Make sure to provide the necessary details for each environment variable. Here's a brief explanation of each variable:

GREENLIGHT_DB_DSN: The database connection string for PostgreSQL. Modify it according to your PostgreSQL database setup.

MIGRATION_URL: The URL or path to the migrations directory. Modify it to match the location of your migrations.

SMTP_HOST, SMTP_PASSWORD, SMTP_PORT, SMTP_SENDER, SMTP_USERNAME: SMTP server configuration for sending emails. Update these values with your SMTP server details. I use [Mailtrap](https://mailtrap.io/), very easy to setup.

CORS_TRUSTED_ORIGINS: A space-separated list of trusted origins for Cross-Origin Resource Sharing (CORS). Modify it with the origins that should be allowed to access the API.

5. **Build and Start Containers**: Run the following command to build and start the containers using Docker Compose::

    ```shell
    docker-compose up --build
    ```
   
6. **Run the Application**:  Once the containers are running, open a new terminal window and navigate to the project directory. Run the following command to start the application:

    ```shell
    make run/api
    ```

7. **[TODO] API Documentation**: The API documentation can be accessed at [http://localhost:4000/docs](http://localhost:8080/docs) once the application is running.

8. **(Optional) View the Help Section**: If you want to see the help section of the app, you can run the following command:

    ```shell
   make run/api/help
    ```

Feel free to explore the code and experiment with the additional features and improvements implemented in this version of the Greenlight application.