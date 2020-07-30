# Rate My Bistro: Crawler

## 0 Prerequisites

 * go
 * access to cgm.ag network

## 1 Build
To build the project run in the project root:
```go
go build
```

## 2 Run
To run the project run in the project root:
```go
go run Main.go
```

## 3 Test
To test a feature navigate into the feature directory and run the test:
```go
cd crawler && go test
```

## 4 Contribution
Before you start changing things, read the following infos:

Please document any new code
Express changes in semantic commit messages
Align your changes with the existing coding style
Better ask first and then start changing
Use Templates
You found a bug somewhere in the code?

--> Open an Issue

You fixed a bug somewhere in the code?

--> Open a pull request

You got an awesome idea to improve the project?

You hate your Bistro as much as I do and want to speed up development?

--> The best way to support me in this project starts with a direct contact. Just send me an email and we will figure out a way on how to split up work :)

--> ansgar.sa@gmail.com
--> rouven@himmelstein.info

## TODO
 * Parse a specific date
 * jobs - introduces sequential job queue 
 * REST feature - provides a rest api to control the crawler  (openAPI swagger)