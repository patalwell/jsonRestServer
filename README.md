# Json REST Server

This repository contains a sample concurrent JSON REST server `concurrentHandlers.go`. 
The server has concurrency patterns for retrieving results from an external API service using multiple go routines.

The server also allows clients to create users by exchanging the appropriate headers and JSON payloads. 

Users are composed of an id, an email, a first name, a last name, a list of stocks, and a portfolio of those stocks. 
Users are only persisted for the duration of the application's runtime.

```json
   { 
      "Id":3,
      "Email":"pat.alwell@gmail.com",
      "Firstname":"Pat",
      "Lastname":"Alwell",
      "Stocks":["MSFT"],
      "Portfolio":[ 
         { 
            "Symbol":"MSFT",
            "open":{ 
               "price":145,
               "time":1572964200946
            },
            "close":{ 
               "price":144.46,
               "time":1572987600495
            },
            "high":145.02,
            "low":143.905
         }
      ]
}
 ```

## Requirements:

At least go version `go1.12.6 `

## Usage:

1. Pull the package `go get github.com/patalwell/jsonRestServer`

2. Compile the application:
    1. For Linux Based Systems: `GOOS=linux GOARCH=amd64 go build -o main web/app/concurrentHandlers.go`
    2. For Mac Based Systesms: `GOOS=darwin GOARCH=amd64 go build -o main web/app/concurrentHandlers.go`
3. Run the executable ` ./main {API-Host} {API-Key}` where API Host and API key are your credentials for rapid api: https://rapidapi.com/search/free

The server executes on the host's internal network by default e.g. `http://localhost:8080`. 
Note: You may need to proxy the host if you are running the server on an external host.

## Tests:

There are unit tests for this project located in the `/web/app` directory. These tests should be run in a sequence that makes sense for the operations, 
e.g. the current functions should not be moved for the sake of sequence. You can run a test by navigating to the test directory and
issuing the `go test` command.

`go test web/app/concurrentHandlers_test.go`

## CRUD Methods

#### Create a User

To create a user post JSON to the `/user/create` endpoint in the form of:
```shell

curl -X POST http://localhost:8080/user/create -H "Content-Type: application/json" \
-d '{"Id":3, "Email":"pat.alwell@gmail.com","Firstname":"Pat", "Lastname":"Alwell", "Stocks":["MSFT","HSIC","DGX"]}'

```
Where the JSON payload contains the Id, Email, Firstname, Lastname, and Stocks of a given user.

#### Get all users

To get all the users you've created make a request to the `/user/all` endpoint in the form of:
```shell

curl -X GET http://localhost:8080/user/all

```
You should receive a JSON array similar to the following:

```json
      
[{ 
   "Id":3,
   "Email":"pat.alwell@gmail.com",
   "Firstname":"Pat",
   "Lastname":"Alwell",
   "Stocks":[ 
      "MSFT"
   ],
   "Portfolio":[ 
      { 
         "Symbol":"MSFT",
         "open":{ 
            "price":144.96,
            "time":1572877800906
         },
         "close":{ 
            "price":144.55,
            "time":1572901200577
         },
         "high":145,
         "low":144.16
      }
   ]
},
{ 
   "Id":35,
   "Email":"ben.franklin@gmail.com",
   "Firstname":"Ben",
   "Lastname":"Franklin",
   "Stocks":[ 
      "HSIC"
   ],
   "Portfolio":[ 
       { 
          "Symbol":"HSIC",
          "open":{ 
             "price":63.43,
             "time":1572877800907
          },
          "close":{ 
             "price":65.73,
             "time":1572901200454
          },
          "high":66.02,
          "low":63.43
       }
   ]
}]
``` 
#### Get A User by Id

To get a single user by their id make a request to `/user?id={user_id}` in the form of:

```shell

curl -X GET http://localhost:8080/user?id=3

```

You should receive a single json record similar to the following:

```json
{ 
   "Id":3,
   "Email":"pat.alwell@gmail.com",
   "Firstname":"Pat",
   "Lastname":"Alwell",
   "Stocks":[ 
      "MSFT"
   ],
   "Portfolio":[ 
      { 
         "Symbol":"MSFT",
         "open":{ 
            "price":144.96,
            "time":1572877800906
         },
         "close":{ 
            "price":144.55,
            "time":1572901200577
         },
         "high":145,
         "low":144.16
      }
   ]
}
```

#### Edit an Existing User

To edit an existing user submit a PUT request to the `/user/edit?id={user_id}` endpoint in the form of:

```shell
curl -X PUT http://localhost:8080/user/edit?id=3 -H "Content-Type: application/json" \
-d '{ "Email":"jdoe@gmail.com", "Firstname":"John", "Lastname":"Doe", "Stocks":["DGX"]}'

```
Updates will only be applied to the Email, Firstname, Lastname, and Stocks fields. 
User Id's are immutable for the sake of indexing and portfolio data is updated based on the time of the query.

You should receive a json Array with the edits:

```json
[ 
   { 
      "Id":3,
      "Email":"jdoe@gmail.com",
      "Firstname":"John",
      "Lastname":"Doe",
      "Stocks":[ 
         "DGX"
      ],
      "Portfolio":[ 
         { 
            "Symbol":"DGX",
            "open":{ 
               "price":101.57,
               "time":1572964200414
            },
            "close":{ 
               "price":100.82,
               "time":1572987825665
            },
            "high":102.14,
            "low":100.76
         }
      ]
   }
]

```

#### Deleting an Existing User

To delete and existing user submit a DELETE request to the `/user/delete?id={user_id}` endpoint in the form of:

```shell

curl -X DELETE http://localhost:8080/user/delete?id=3

```
You should receive notification that the user has successfully been deleted.

```json

Deleted User { 
   "Id":3,
   "Email":"jdoe@gmail.com",
   "Firstname":"John",
   "Lastname":"Doe",
   "Stocks":[ 
      "DGX"
   ],
   "Portfolio":[ 
      { 
         "Symbol":"DGX",
         "open":{ 
            "price":101.57,
            "time":1572964200414
         },
         "close":{ 
            "price":100.82,
            "time":1572987825665
         },
         "high":102.14,
         "low":100.76
      }
   ]
}

```