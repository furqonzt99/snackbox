# SNACKBOX

a rest-ful api project

Run project with:

```
go run main.go
```

## Stack-tech

- [x] RESTful API Using Go, Echo, Gorm, MySQL
- [x] AWS for service api

## Open Endpoints

Open endpoints require no Authentication.

- Register : `POST /register`
- Login : `POST /login `

## Endpoints that require Authentication

Closed endpoints require a valid Token to be included in the header of the request. A Token can be acquired from the Login view above.

User related

Each endpoint manipulates or displays information related to the User whose Token is provided with the request:

- Get user profile data by User ID : `GET /profile`
- Update user data by User ID : `PUT /users`
- Delete user data by User ID : `DELETE /users`

### Partner Related
