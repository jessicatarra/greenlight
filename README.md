# [Let's Go Further book by Alex Edwards](https://lets-go-further.alexedwards.net/)

Welcome to the README file for the Greenlight application developed as part of the book "Let's Go Further" by Alex Edwards. This document provides an overview of the application, highlights the incredible value of the book, and showcases the additional features and improvements implemented beyond the step-by-step tutorial.

## About the Book

"Let's Go Further" is an exceptional resource for software engineers who want to delve into advanced patterns for building APIs and web applications in Go. Authored by [Alex Edwards](https://www.alexedwards.net/), this book goes beyond the basics and equips developers with the knowledge to build robust and scalable applications.

See more here: https://lets-go-further.alexedwards.net/

## Additional Features and Improvements


1. **Docker Compose**: Included a Docker Compose file to simplify the setup process and allow running the database locally. 
2. **API Documentation**: Integrated API documentation using Swagger. You can access the API documentation [here](https://greenlight.tarralva.com/swagger/index.html). It provides detailed information about the available endpoints, request/response schemas, and allows testing the API interactively.
3. **Auto Migration**: Implemented an automated migration process. This feature automatically applies any required database schema changes when the application starts, ensuring that the database is always up to date with the latest changes.
4. **Embedded Migration Files**: Added embedded migration files to the application. These migration files contain SQL scripts that define the necessary database schema changes. See more: [GO EMBED FOR MIGRATIONS](https://oscarforner.com/blog/2023-10-10-go-embed-for-migrations/)


## Getting Started

1. **Prerequisites**: Make sure you have Go (version 1.21.0 or higher) installed on your machine.
2. **Clone the Repository**: Run the following command to clone the repository:

    ```shell
    git clone git@github.com:jessicatarra/greenlight.git
    ```

3. **Install Dependencies**: Change to the project directory and use the following command to install the dependencies:

    ```shell
    go mod download
    ```

4. **Environment Variables**: Create a `.envrc` file in the project directory and add the following environment variables with their corresponding values:

```
export DATABASE_URL=

export SMTP_HOST=

export SMTP_PASSWORD=

export SMTP_PORT=

export SMTP_SENDER=

export SMTP_USERNAME=

export CORS_TRUSTED_ORIGINS=

export JWT_SECRET=

```

Make sure to provide the necessary details for each environment variable. Here's a brief explanation of each variable:

DATABASE_URL: The database connection string for PostgreSQL. Modify it according to your PostgreSQL database setup.

SMTP_HOST, SMTP_PASSWORD, SMTP_PORT, SMTP_SENDER, SMTP_USERNAME: SMTP server configuration for sending emails. Update these values with your SMTP server details. I use [Mailtrap](https://mailtrap.io/), very easy to set up.

CORS_TRUSTED_ORIGINS: A space-separated list of trusted origins for Cross-Origin Resource Sharing (CORS). Modify it with the origins that should be allowed to access the API.

JWT_SECRET: The algorithm ( HS256 ) used to sign the JWT means that the secret is a symmetric key that is known by both the sender and the receiver. [See more](https://jwt.io/)

5. **Build and Start Containers**: Run the following command to build and start the containers using Docker Compose::

    ```shell
    docker-compose up --build
    ```
   
6. **Run the Application**:  Once the containers are running, open a new terminal window and navigate to the project directory. Run the following command to start the application:

    ```shell
    make run/api
    ```

7. **API Documentation**: The API documentation can be accessed at [https://localhost:4000/swagger/index.html](http://localhost:8080/docs) once the application is running.

8. **(Optional) View the Help Section**: If you want to see the help section of the app, you can run the following command:

    ```shell
   make run/api/help
    ```

## Migrations

To manage SQL migrations in this project we’re going to use the migrate command-line tool (which itself is written in Go).

Detailed installation instructions for different operating systems can be found here, but on macOS you should be able to install it with the command:

   ```shell
      brew install golang-migrate
   ```
And on Linux and Windows, the easiest method is to download a pre-built binary and move it to a location on your system path. For example, on Linux:

```shell

   curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz &
   mv migrate.linux-amd64 $GOPATH/bin/migrate

```

## Quality Controlling Code

In order to check, test and tidy up our codebase automatically, run the following command:

```shell
   make audit
```

Also, to successfully run the previous command, you would need to download the following package:

- Use the third-party [staticcheck](https://staticcheck.io/) tool to carry out some [additional static analysis checks](https://staticcheck.dev/docs/checks).


Feel free to explore the code and experiment with the additional features and improvements implemented in this version of the Greenlight application.