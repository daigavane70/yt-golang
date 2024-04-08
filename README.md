# FamPay Backend Assignment - YouTube API

## Project Overview

This project is part of the FamPay Backend Assignment. The goal is to create an API that fetches the latest videos from YouTube based on a predefined search query, stores the video data in a database, and provides endpoints for retrieving and searching videos.

## Instructions
1. Clone the repository: 
```
git clone https://github.com/daigavane70/yt-golang.git
```

2. Navigate to the project
```
cd go-sprint
```

3. Build and run docker
```
docker build -t go-sprint . 
docker run --env-file .env -p 8080:8080 go-sprint   
```

## Development Environment

- Language: Go (Golang)
- Database: MySql

## Project Structure

```plaintext
go-sprint/
├── cmd/
│   └── main/
│       └── main.go
├── pkg/
│   ├── config/
│   │   └── config.go
│   ├── controllers/
│   │   ├── test-controller.go
│   │   └── video-controller.go
│   ├── entities/
│   │   ├── Config.go
│   │   └── Videos.go
│   ├── models/
│   │   ├── CommonResponse.go
│   │   └── Youtube.go
│   ├── routes/
│   │   └── routes.go
│   ├── services/
│   │   └── youtube-services.go
│   └── common/
│       ├── constants/
│       │   └── constants.go
│       ├── logger/
│       │   └── logger.go
│       └── utils/
│           └── utils.go
├── .env
├── go.mod
├── go.sum
├── README.md
└── Dockerfile
```

## Api Endpoints

#### GET /video?searchKeyword=keyword
Example - http://localhost:8080/video?search=iplRcb&page=58&pageSize=1