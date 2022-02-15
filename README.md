![Logo](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-logo.png)

[![Go.Dev reference](https://img.shields.io/badge/echo-reference-blue?logo=go&logoColor=blue)](https://github.com/labstack/echo)
[![Go.Dev reference](https://img.shields.io/badge/gorm-reference-blue?logo=go&logoColor=blue)](https://pkg.go.dev/gorm.io/gorm?tab=doc)
[![Go.Dev reference](https://img.shields.io/badge/aws--s3-reference-orange)](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go/example_code/s3)
[![Go.Dev reference](https://img.shields.io/badge/maroto-reference-blue)](https://pkg.go.dev/github.com/johnfercher/maroto?tab=doc)

# Snackbox

Organizations or offices need to set up a consumption section that is in charge of providing food, where they have to find restaurants and contact them all the time to track their orders to ensure the orders are processed properly.

Snackbox can make it easier for customers to order snacks and rice boxes and do tracking, providing security for customers because the connected partners have passed the verification process and this application will provide employment for partners.

## Features

- JWT Authentication
- Search Product By Nearest Partner Location
- Payment Gateway Integration - Invoice & Disbursment (Xendit)
- AWS S3 Integration
- PDF Export (Maroto)

## Installation

- Clone this repo

```bash
git clone https://github.com/furqonzt99/snackbox.git snackbox
```

- Go to repository folder

```bash
cd snackbox
go mod tidy
```

- Create .env file and add the following environment (you can see the variables from .env.example)

- Run this app

```bash
go run .
```

## Use Case Diagram

![Use Case Diagram](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-ucd.png)

## Entity Relationship Diagram

![Entity Relationship Diagram](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-erd.png)

## High Level Architecture

![High Level Architecture](https://github.com/furqonzt99/snackbox/blob/documentation/documentation/snackbox-hla.png)

## API Documentation

[API Documentation](https://app.swaggerhub.com/apis-docs/furqonzt99/snackbox/1)

## Wire Flow

[Wire Flow](https://whimsical.com/snackbox-UcYKhew5MBhFzJWaCXQbAb)

## Code Structure

## Unit Test

## Tech Stack

## Authors

- [@furqonzt99](https://github.com/furqonzt99) - Product Owner & Developer
- [@yogawahyudi7](https://github.com/yogawahyudi7) - Developer
