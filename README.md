# staff-scheduling

## Description
This app was build to schedule shifts for staff. 

## Tech stack
The following tech stack and designs were used:
- Go
- Postgres DB
- Fiber (API framework)
- JWT token (for auth)
- Swagger/swag (endpoints documentation)
- SQLC (to automate PSQL queries generation)
- Docker
- Domain-Driven Design

### DB Schema:
<img width="1106" alt="image" src="https://github.com/Bruary/staff-scheduling/assets/38393880/d0d60be6-35ef-4ca1-b914-4a0f896f78aa">



## Getting Started

### Prerequisits
1) **Go**: [guide](https://www.geeksforgeeks.org/how-to-install-golang-on-macos/)
2) **Docker Desktop**: [download](https://www.docker.com/products/docker-desktop/)
3) **Postman**: [download](https://www.postman.com/downloads/)
4) **IDE (VSCode)**: [download](https://code.visualstudio.com/download)


Finally, **clone the repo**. Run the below command on the terminal _(on mac press cmd + space, select terminal, then paste the below)_:
```
git clone git@github.com:Bruary/staff-scheduling.git
```

Now, open the repo using VSCode and on the root of the repo, run the below command to install all dependencies:
```
go mod vendor
```

**Note: Make sure Docker Desktop app is running in the background before continuing**


After this, spin up the containers using:
```
docker-compose up
```

**And you should be all set (assuming everything went fine)**

## Testing the app

### Access documentation
On your local browser, go to: http://localhost:3000/docs/ and you will find all endpoints with expected request/response bodies.

### Use postman collection
Postman collection can be shared via email.
