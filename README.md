# Go-API Boilerplate: A Golang-based API Boilerplate

## Welcome to Go-API Boilerplate.

Golang is a highly flexible and straightforward language, providing an exceptional environment for implementing customized APIs. We offer a basic boilerplate with a pre-implemented structure to assist and speed up future developments. Additionally, we already ship it within Docker containers, as well as an authentication system based on JWT, Redis, and MongoDB. Please, follow along the next sections in order to have a better glimpse of this exciting project.

Go-API Boilerplate is compatible with: **Go 1.14+**.

---

## Package guidelines

1. The very first information you need is in the very **next** section.
2. **Installing** is also easy if you wish to read the code and bump yourself into, follow along.
3. Note that there might be some **additional** steps in order to use our solutions.
4. If there is a problem, please do not **hesitate**. Call us.

---

## Getting started: 60 seconds with Go-API Boilerplate

First of all. The code is all commented. Yes, it is commented. Just browse to any file, chose your subpackage, and follow it. We have high-level code for most tasks we could think.

Alternatively, if you wish to learn even more, please take a minute:

Go-API Boilerplate is based on the following structure, and you should pay attention to its tree:

```
- go_api_boilerplate
    - docker-compose-prod.yml
    - docker-compose.yml
    - Dockerfile
    - Dockerfile.prod
    - .env.example
    - requests.json
    - src
        - controllers
        - database
        - middleware
        - models
        - server
        - utils
        - api.go
        - go.mod
        - go.sum
        - reflex.conf
```

### Docker-based Files

Docker-based files, such as `Dockerfile` and `docker-compose.yml`, provides a straightforward way to start a container and use the application just out-of-the-box.

### Environment File

In order to provide a more customizable environment, we ship a `.env.example` file, which should be modified and copied to a `.env` file in order to boot the system.

### Requests

We provide a collection of Postman's requests, which contains all the possible `requests` that the API is capable of invoking.

### Source

Known as ```src```, this is where all the magic happens. Follow along the next items to understand what is happening here.

#### controllers

Every common-knowledge route should be handled by a `controller`. This module offers an easy way to implement customized routes and operations.

#### database

Data should be stored somewhere, right? This module offers methods that allow one to access and use our provided `database` applications, such as MongoDB and Redis.

#### middleware

Before resolving any request, one might need to perform some pre-request operations, correct? The `middleware` provides methods that are invocable before the request processing, such as authorizing a known user.

#### models

When dealing with data, one might need to define how it is structured and which type of information it should encode. The `models` package is the perfect place to define how your system's information should look.

#### server

Starting from zero might burden the development, don't you think? We are based-off the Gin framework, and in an attempt to ease the developer's life, we opted to provide the `server` module, which holds all server-related functions, such as initializing the server itself and the router.

#### utils

Common-based implementations should be available throughout the application and never re-implemented every time they need to be used, correct? We provide a `utils` package, where conventional scope methods are available to the whole application.

#### api.go

This is the `application` entry point component. Only modify it if you know what you are doing.

#### go.mod

Golang now provides a better alternative to download and install `required packages`, such as NPM's package.json.

#### go.sum

Golang also provides a `hash` containing all the pre-installed modules needed by the application.

#### reflex.conf

One new development tool concerns the hot-reloading, where developers do not need to restart the application to update their changes. Thus, we opted to use the `reflex` package in order to check which files have been modified and re-compile the application.

---

## Installation

We believe that everything has to be easy. Not tricky or daunting, Go-API Boilerplate will be the one-to-go package that you will need, from the very first installation to the daily-tasks implementing needs.

### Development

To ease one needs in a development environment, we ship this package in a Docker container. Make sure that ```docker``` and ```docker-compose``` are installed and accessible from the command line.

Make sure that you have adjusted your environment variable needs in the `.env.example` file and have copied it into a `.env` file before attempting to build and launch the container.

Finally, you can build the container by using:

```
docker-compose build
```

After the build process is finished, you can run the container in detached mode:

```
docker-compose up -d
```

If you ever need to perform maintenance or update the repository, please put the container down (ensure to use -v; otherwise it will not replace the build):

```
docker-compose down
```

### Production

To ease your needs in a production environment, we also ship this package in a Docker container. Make sure that ```docker``` and ```docker-compose``` are installed and accessible from the command line, and that your `Dockerfile.prod` and `dockerfile-prod.yml` are the Docker's entry points.

Additionally, make sure that you have adjusted your environment variable needs in the `.env.example` file and have copied it into a `.env` file before attempting to build and launch the container.

Finally, you can build the container by using:

```
docker-compose build
```

After the build process is finished, you can run the container in detached mode:

```
docker-compose up -d
```

If you ever need to perform maintenance or update the repository, please put the container down (ensure to use -v; otherwise it will not replace the build):

```
docker-compose down
```

---

## Environment configuration

Note that sometimes, there is a need for additional implementation. If needed, from here, you will be the one to know all of its details.

### Ubuntu

No specific additional commands needed.

### Windows

No specific additional commands needed.

### MacOS

No specific additional commands needed.

---

## Support

We know that we do our best, but it is inevitable to acknowledge that we make mistakes. If you ever need to report a bug, report a problem, talk to us, please do so! We will be available at our bests at this repository.

---
