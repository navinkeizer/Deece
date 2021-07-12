package main

import (
	"io/ioutil"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)
/*
func TestHandler(t *testing.T) {
	//Here, we form a new HTTP request. This is the request that's going to be passed to our handler.
	// The first argument is the method, the second argument is the route, and the third is the request body, which we don't have in this case.
	req, err := http.NewRequest("GET", "", nil)

	// In case there is an error in forming the request, we fail and stop the test
	if err != nil {
		t.Fatal(err)
	}

	// We use Go's httptest library to create an http recorder. This recorder will act as the target of our http request
	// (you can think of it as a mini-browser, which will accept the result of the http request that we make)
	recorder := httptest.NewRecorder()

	// Create an HTTP handler from our handler function. "handler" is the handler function defined in our main.go file that we want to test
	hf := http.HandlerFunc(handler)

	// Serve the HTTP request to our recorder. This is the line that actually executes our the handler that we want to test
	hf.ServeHTTP(recorder, req)

	// Check the status code is what we expect.
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `Hello World!`
	actual := recorder.Body.String()
	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func TestRouter(t *testing.T) {
	// Instantiate the router using the constructor function that
	// we defined previously
	r := newRouter()

	// Create a new server using the "httptest" libraries `NewServer` method
	// Documentation : https://golang.org/pkg/net/http/httptest/#NewServer
	mockServer := httptest.NewServer(r)

	// The mock server we created runs a server and exposes its location in the URL attribute
	// We make a GET request to the "hello" route we defined in the router
	resp, err := http.Get(mockServer.URL + "/hello")

	// Handle any unexpected error
	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be ok, got %d", resp.StatusCode)
	}

	// In the next few lines, the response body is read, and converted to a string
	defer resp.Body.Close()
	// read the body into a bunch of bytes (b)
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	// convert the bytes to a string
	respString := string(b)
	expected := "Hello World!"

	// We want our response to match the one defined in our handler.
	// If it does happen to be "Hello world!", then it confirms, that the
	// route is correct
	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}

}

func TestRouterForNonExistentRoute(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)
	// Most of the code is similar. The only difference is that now we make a //request to a route we know we didn't define, like the `PUT /hello` route.
	resp, err := http.Post(mockServer.URL+"/hello", "", nil)

	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 405 (method not allowed)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Status should be 405, got %d", resp.StatusCode)
	}

	// The code to test the body is also mostly the same, except this time, we need an empty body
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	respString := string(b)
	expected := ""

	if respString != expected {
		t.Errorf("Response should be %s, got %s", expected, respString)
	}

}

func TestStaticFileServer(t *testing.T) {
	r := newRouter()
	mockServer := httptest.NewServer(r)

	// We want to hit the `GET /web/` route to get the index.html file response
	resp, err := http.Get(mockServer.URL + "/web/")

	if err != nil {
		t.Fatal(err)
	}

	// We want our status to be 200 (ok)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Status should be 200, got %d", resp.StatusCode)
	}

	// It isn't wise to test the entire content of the HTML file.
	// Instead, we test that the content-type header is "text/html; charset=utf-8"
	// so that we know that an html file has been served
	contentType := resp.Header.Get("Content-Type")
	expectedContentType := "text/html; charset=utf-8"

	if expectedContentType != contentType {
		t.Errorf("Wrong content type, expected %s, got %s", expectedContentType, contentType)
	}

}

func TestGetCrawlTargetsHandler(t *testing.T) {

	crawlTargets = []CrawlTarget{
		{"k2k4r8oxynrlparmnoh62lhk0ozhsdw8lizrcxjxs2w3jlllrqzi2bm8", "A small harmless crawlTarget"},
	}

	req, err := http.NewRequest("GET", "", nil)

	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(getCrawlTargetHandler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := CrawlTarget{"k2k4r8oxynrlparmnoh62lhk0ozhsdw8lizrcxjxs2w3jlllrqzi2bm8", "A small harmless crawlTarget"}
	b := []CrawlTarget{}
	err = json.NewDecoder(recorder.Body).Decode(&b)

	if err != nil {
		t.Fatal(err)
	}

	actual := b[0]

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}
func TestCrawlTargetsHandler(t *testing.T) {

	crawlTargets = []CrawlTarget{
		{"k2k4r8oxynrlparmnoh62lhk0ozhsdw8lizrcxjxs2w3jlllrqzi2bm8", "A small harmless crawlTarget"},
	}

	form := newCreateCrawlTargetForm()
	req, err := http.NewRequest("POST", "", bytes.NewBufferString(form.Encode()))

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()

	hf := http.HandlerFunc(createCrawlTargetHandler)

	hf.ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := CrawlTarget{"QmeLzcTz6KEguARsZNorsJ7RvWMsaGdgYKyX5MQcFMUevA", "Eagle related content"}

	if err != nil {
		t.Fatal(err)
	}

	actual := crawlTargets[1]

	if actual != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", actual, expected)
	}
}

func newCreateCrawlTargetForm() *url.Values {
	form := url.Values{}
	form.Set("CID", "QmeLzcTz6KEguARsZNorsJ7RvWMsaGdgYKyX5MQcFMUevA")
	form.Set("description", "Eagle related content")
	return &form
}

*/
