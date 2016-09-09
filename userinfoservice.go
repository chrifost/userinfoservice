//Simple HTTP server

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
)

// FYI - Remember to make variables start with Uppercase so they can be accessed
// outside the package
type userRecord []struct {
	// These fields use the "json: tag" to specify which field they map to
	UserID   float64 `json:"id"`
	User     string  `json:"name"`
	UserName string  `json:"username"`
	// These fields are mapped direclty by name (note the different case)
	Email string
	Phone string
	// As these fields can be nullable, we use a pointer to a string rather
	// than a string
	Website string
}

var info *log.Logger

// Log the basic HTTP info
func logHTTPInfo(req *http.Request) {

	info.Println("-->HTTP Request Details<--")
	if nil != req {
		info.Println("Method:      ", req.Method)
		info.Println("Protocol:    ", req.Proto)
		info.Println("Scheme:      ", req.URL.Scheme)
		info.Println("Host:        ", req.Host)
		info.Println("URL:         ", req.URL)
		info.Println("Header Info: ")
		// Print the Header
		for key, value := range req.Header {
			info.Println(key, ":", value)
		}

	} else {
		info.Println("HTTP Request Empty")
	}
	info.Println("<--HTTP Request Details-->")
}

func getContent(url string) ([]byte, error) {
	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Send the request via a client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Print the HTTP Header
	//printHeader(*resp)

	// Defer the closing of the body
	defer resp.Body.Close()

	//Read the cotnent into a byte array
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// We are done, return the bytes
	return body, nil
}

func getUserRecord(id string) (*userRecord, error) {
	content, err := getContent(fmt.Sprintf("https://jsonplaceholder.typicode.com/users?id=%s", id))
	if err != nil {
		info.Println("URL Error:", err)
		return nil, err
	}

	// Print the full response
	info.Println("Raw JSON Response:\n", string(content))

	//Fill the record with the data from the JSON
	var record userRecord
	err = json.Unmarshal(content, &record)
	if err != nil {
		info.Println("JSON Error:", err)
		return nil, err
	}

	return &record, nil
}

func (record userRecord) IsEmpty() bool {
	return reflect.DeepEqual(record, userRecord{})
}

func handler(w http.ResponseWriter, req *http.Request) {

	/* Explored an query parse approach found here:
	   http://blog.charmes.net/2015/07/parsing-http-query-string-in-go.html

	   However, for initial work decided to use  simple URL Parse approach
	*/

	var userID string

	logHTTPInfo(req)

	// All we care about are the query params to try and parse them
	m, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		info.Println("Error parsing req.URL.RawQuery")
	}

	// Print the map
	info.Println("Recieved query params:")
	for key, value := range m {
		switch key {
		case "userId":
			userID = value[0]
			if userID == "" {
				info.Println("Error parsing userId: ", err)
			}
		}
	}

	info.Println("userId:", userID)

	//Get user data based on provided Id
	record, err := getUserRecord(userID)
	if err != nil {
		info.Println("Error getUserRecord()")
	}

	// Check to see if we recieved a user record
	x := record
	if x.IsEmpty() {
		info.Println("Empty Respsone")
	}

	info.Printf("User Record return:\n %+v", record)
}

func initLogging(infoHandle io.Writer) {

	info = log.New(infoHandle, "INFO: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}

// Bit bucket for ""/favicon.ico" path
func faviconHandler(w http.ResponseWriter, req *http.Request) {}

func main() {

	// Intialize a basic logger
	initLogging(os.Stdout)
	info.Println("Starting User Info Service...")

	// Setup HTTP handlers.  Add one for "/favicon.ico" just to dump browser request for this
	http.HandleFunc("/favicon.ico", faviconHandler)
	http.HandleFunc("/", handler)

	// Start HTTP server with generic handler
	log.Fatal(http.ListenAndServe(":8080", nil))
	//curl 'http://localhost:8080?limit=42&dryrun=true'
}
