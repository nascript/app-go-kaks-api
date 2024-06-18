# WEB API APPS

## Project Description: Product Management CRUD Web API

### BASE URL 
```http://13.236.121.17:4040/api/v1/swagger/index.html```

### Overview

This project is a web API designed to perform CRUD (Create, Read, Update, Delete) operations on product data. The API is built using Go (Golang) and follows a clean architecture pattern to ensure scalability,, maintainability, and testability. It leverages various AWS services for deployment and cloud management, providing a robust solution for product data management.

### Key Features

1. **Create Product**: Allows users to add new product entries to the database. Each product entry includes essential details such as product name, stock quantity, and timestamps for creation and updates.

2. **Read Products**: Provides endpoints to retrieve product data. Users can fetch all products or a specific product by its ID.

3. **Update Product**: Enables users to update existing product information, including product name and stock levels.

4. **Delete Product**: Allows users to remove product entries from the database.

### Technical Stack

- **Programming Language**: Go (Golang)
- **Web Framework**: Fiber
- **Database**: MongoDB
- **Cloud Services**: AWS EC2
- **Version Control**: GitHub for source code management and continuous integration
- **Deployment**: GitHub Action to Ubuntu AWS

## Features

- Create Data
- Update Data
- Read All Data
- Read One Data
- Delete Data

## VIEW (UI)

![Alt text](./img/Screenshot%202024-06-15%20at%2011.02.19 AM.png)

## DOCUMENTATION

### Create Data

![Alt text](./img/Screenshot%202024-06-15%20at%2010.58.15 AM.png)

### Read All Data

![Alt text](./img/Screenshot%202024-06-15%20at%2010.57.42 AM.png)

### Read One Data

![Alt text](./img/Screenshot%202024-06-15%20at%2010.58.23 AM.png)

### Edit Data

![Alt text](./img/Screenshot%202024-06-15%20at%2010.58.35 AM.png)

### Delete One Data

![Alt text](./img/Screenshot%202024-06-15%20at%2010.58.42 AM.png)

## Cloud Database Using MongoDB

![Alt text](./img/Screenshot%202024-06-15%20at%2012.19.46 PM.png)

## Folder Structure

```
└── 📁kaks-apps
    └── .env
    └── 📁.idea
        └── .gitignore
        └── CRUD_Hexagonal.iml
        └── modules.xml
        └── vcs.xml
    └── README.md
    └── 📁api
        └── 📁product
            └── adapter.go
            └── handler.go
    └── 📁cmd
        └── main.go
    └── 📁docs
        └── docs.go
        └── swagger.json
        └── swagger.yaml
    └── 📁domain
        └── 📁product
            └── product.go
            └── service.go
    └── go.mod
    └── go.sum
    └── 📁infrastructure
        └── mongo.db.go
        └── opentelemetry.go
    └── 📁repository
        └── 📁product
            └── repository.go
    └── 📁service
        └── 📁product
            └── adapter.go
    └── 📁utils
        └── http.go
        └── utils.go
        └── validator.go
```
