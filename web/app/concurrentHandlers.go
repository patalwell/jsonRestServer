package main

/*
Pat Alwell
Software Engineer - Cloud

Build a rest API that allows someone to add, modify, view, and delete users.
It must also connect to an external service, such as a stock or weather
feed, and give information about a user’s chosen cities or stocks.

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

//Information about our users
type User struct {
	Id        int         `json:"Id"`
	Email     string      `json:"Email"`
	FirstName string      ` json:"Firstname"`
	LastName  string      `json:"Lastname"`
	Stocks    []string    `json:"Stocks"`
	Portfolio []StockData `json:"Portfolio"`
}

//Slice to store our users
var allUsers []User

type StockData struct {
	Symbol string
	Open   struct {
		Price float64 `json:"price"`
		Time  int64   `json:"time"`
	} `json:"open"`
	Close struct {
		Price float64 `json:"price"`
		Time  int64   `json:"time"`
	} `json:"close"`
	High float64 `json:"high"`
	Low  float64 `json:"low"`
}

type Result struct {
	Error    error
	Response StockData
}

//create a function to manage request for stock info
func GetStockData(done <-chan interface{}, stockTickers ...string) <-chan Result {

	//channel for results
	results := make(chan Result)

	go func() {

		defer close(results)

		for _, stock := range stockTickers {
			var result Result
			//update Results with Symbol prior to query
			result.Response.Symbol = stock

			url := "https://investors-exchange-iex-trading.p.rapidapi.com/stock/" + stock + "/ohlc"

			req, _ := http.NewRequest("GET", url, nil)


			req.Header.Add("x-rapidapi-host", os.Args[1])
			req.Header.Add("x-rapidapi-key", os.Args[2])

			rsp, _ := http.DefaultClient.Do(req)
			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				results <- Result{Error: err, Response: StockData{}}
			}
			defer rsp.Body.Close()
			err = json.Unmarshal(body, &result.Response)
			result = Result{Error: err, Response: result.Response}
			select {
			case <-done:
				return
			case results <- result:
			}
		}
	}()
	return results
}

//Return all users
func ShowAllUsers(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/user/all" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		// If it's not, use the w.WriteHeader() method to send a 405 status
		// code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	if allUsers != nil {
		js, err := json.Marshal(allUsers)
		log.Println("GET /user/all", allUsers)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//return json
		w.Header().Set("Content-Type","application/json")
		w.Write(js)
		return

	}
	http.Error(w, "There aren't any users!", 400)

}

//Show user by id
func ShowUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		// If it's not, use the w.WriteHeader() method to send a 405 status
		// code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	//get id from the path args
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Id must be an integer!", 405)
		return
	}

	for _, user := range allUsers {
		if user.Id == id {
			js, err := json.Marshal(user)
			log.Println("GET /user?id=", allUsers)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//return json
			w.Header().Set("Content-Type","application/json")
			w.Write(js)
			return
		}
	}

	http.Error(w, "User with id:"+r.URL.Query().Get("id")+" Not Found!", 400)

}

//Create a user with JSON
func CreateUser(w http.ResponseWriter, r *http.Request) {

	// Use r.Method to check whether the request is using POST or not. Note that
	// http.MethodPost is a constant equal to the string "POST".
	if r.Method != http.MethodPost {
		// If it's not, use the w.WriteHeader() method to send a 405 status
		// code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	//Check headers
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting Header with Content-Type application/json", http.StatusUnsupportedMediaType)
		return
	}

	//get json from client
	clientRequest := json.NewDecoder(r.Body)

	//Decode the json into our User type
	var newUser User
	err := clientRequest.Decode(&newUser)
	if err != nil {
		http.Error(w, "Unable to Marshall Request!", http.StatusInternalServerError)
	}


	//If the user exists send a bad request
	for _, v := range allUsers {
		if v.Id == newUser.Id {
			http.Error(w, "Users already exists!", http.StatusBadRequest)
			return
		}
	}

	//close our concurrent requests
	done := make(chan interface{})
	log.Println("POST /user/create ", newUser)

	//iterate over our result channel for stock data
	for result := range GetStockData(done, newUser.Stocks...) {
		if result.Error != nil {
			log.Println(result.Error.Error())
			continue
		}

		//update our portfolio
		newUser.Portfolio = append(newUser.Portfolio, result.Response)

	}

	allUsers = append(allUsers, newUser)
	//sort after appending
}

//Alter user by id
func EditUser(w http.ResponseWriter, r *http.Request) {

	// Use r.Method to check whether the request is using POST or not. Note that // http.MethodPost is a constant equal to the string "POST".
	if r.Method != http.MethodPut {
		// If it's not, use the w.WriteHeader() method to send a 405 status // code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.
		w.Header().Set("Allow", http.MethodPut)
		http.Error(w, "Method Not Allowed", 405)
		return
	}

	//Check headers
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Expecting Header with Content-Type application/json", http.StatusUnsupportedMediaType)
		return
	}

	//Get id from query params
	queryParam := r.URL.Query().Get("id")

	//get id from the path args
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, "Id must be an integer!", 405)
		return
	}

	//get Json data from client for new entry
	clientRequest := json.NewDecoder(r.Body)

	//un-marshall the bytes into a user
	var userEntry User
	err = clientRequest.Decode(&userEntry)

	for idx, user := range allUsers {
		if user.Id == id {
			user.FirstName = userEntry.FirstName
			user.LastName = userEntry.LastName
			user.Email = userEntry.Email

			//clear previous portfolio
			//for stockIdx := range user.Stocks {
			//		user.Portfolio = append(user.Portfolio[:stockIdx], user.Portfolio[stockIdx+1:]...)
			//	}
			user.Portfolio = []StockData{}
			//update the new stocks
			user.Stocks = userEntry.Stocks
			log.Println("Here are the new stock entries", user.Stocks)

			//close our concurrent requests
			done := make(chan interface{})
			log.Println("PUT /user/edit ", userEntry)

			//iterate over our result channel for stock data
			for result := range GetStockData(done, user.Stocks...) {
				if result.Error != nil {
					log.Println(result.Error.Error())
					continue
				}
				//update our portfolio
				user.Portfolio = append(user.Portfolio, result.Response)
			}

			//add new entry to userList
			allUsers = append(allUsers[:idx], user)
			js, err := json.Marshal(allUsers)
			if err != nil {
				http.Error(w, "Something went wrong while encoding allUsers", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type","application/json")
			w.Write(js)
		}
	}

	log.Println("Edited user with Id ", queryParam, userEntry)

}

//Remove a user by id
func DeleteUser(w http.ResponseWriter, r *http.Request) {

	// Use r.Method to check whether the request is using POST or not. Note that // http.MethodPost is a constant equal to the string "POST".
	if r.Method != http.MethodDelete {
		// If it's not, use the w.WriteHeader() method to send a 405 status // code and the w.Write() method to write a "Method Not Allowed"
		// response body. We then return from the function so that the
		// subsequent code is not executed.
		w.Header().Set("Allow", http.MethodDelete)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	//get query param
	queryParam := r.URL.Query().Get("id")

	//get id from the path args
	id, err := strconv.Atoi(queryParam)
	if err != nil {
		http.Error(w, "Id must be an integer!", 405)
		return
	}

	for idx, user := range allUsers {
		if user.Id == id {
			//remove the user
			allUsers = append(allUsers[:idx], allUsers[idx+1:]...)

			//marshall user for output
			js, err := json.Marshal(user)
			log.Println("DELETE /user/delete?id=", queryParam, " ", string(js))

			//write our response to the client
			_, err = fmt.Fprintf(w,"Deleted User %d",user.Id )
			if err != nil {
				log.Println(err.Error())
				http.Error(w, "Something went wrong while deleting a users", http.StatusInternalServerError)
				return
			}
			return
		}
	}
	//If the user isn't in our list
	http.Error(w, "User doesn't exist!", http.StatusBadRequest)

}

//Health check
func Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {

	http.HandleFunc("/", Ping)
	http.HandleFunc("/user/all",ShowAllUsers)
	http.HandleFunc("/user",ShowUser)
	http.HandleFunc("/user/create",CreateUser)
	http.HandleFunc("/user/edit",EditUser)
	http.HandleFunc("/user/delete",DeleteUser)

	//log.Println(os.Getenv("RAPID_API_HOST"))
	//log.Println( os.Getenv("RAPID_API_KEY"))

	log.Println("Starting Server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}