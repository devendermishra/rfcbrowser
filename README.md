# An RFC browser
This project provides the following:
1. UI to view RFC (list RFCs and can be searched)
2. Backend service to handle with RFC website

Both are hosted at localhost. UI is at localhost:3000 and service is at localhost:8080.

## Why this project?
This is to improve search and experience in reading RFC. By using an UI, customization can be added (like search). By using a backend service, different operations can be performed (like storing, indexing).

# Different components

## UI
UI is written in react (using typescript).

## Backend service
It is written in golang using gin and sqlite3 for storage.

## How to run?
### Dependencies required
node: 16.20+
npm: 9.8+
yarn: 1.22+
go: 1.21+ (ensure `GOPATH` are setup)

### Run backend service

    > go build
    > ./rfc-service

Service will start `localhost:8080`

### Run UI

> yarn start

UI will run at `localhost:3000`.
