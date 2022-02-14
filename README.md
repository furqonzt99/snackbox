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

Each endpoint manipulates or displays information related to the partner whose Token is provided with the request:

- User apply as Partner : `POST /partners/submission`
- Upload legal document as Partner : `POST "/partners/submission/upload"`
- Get Partner data & their Product : `GET /partners/:id`
- Get report PDF as Partner : `GET /partners/report`
- Get all Partner information by Admin : `GET /partners/submission`
- Accept request User as Partner by Admin : `PUT /partners /submission/:id/accept`
- Reject request User as Partner by Admin : `PUT /partners/submission/:id/reject`

### Product Related

Each endpoint manipulates or displays information related to the Product whose Token is provided with the request:

- Add a Product by Partner : `POST /products`
- Update data Product by Partner : `PUT /products/:id`
- Delete Product by Partner : `DELETE /products/:id`
- Get All Product : `GET /products`
- Upload image Product : PUT /products/:id/image`
