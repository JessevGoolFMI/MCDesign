package TestDesign

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// fetchRandomNumber makes a GET request to the random number API and returns the random number and any error encountered.
func fetchRandomNumber() (int, error) {
	// Define the URL of the API endpoint
	url := "https://www.randomnumberapi.com/api/v1.0/random?min=1&max=100"

	// Make the HTTP GET request
	response, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("error making the request: %w", err)
	}
	defer response.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading the response body: %w", err)
	}

	// Unmarshal the JSON response into a RandomNumberResponse struct
	var numbers []uint8
	err = json.Unmarshal(body, &numbers)
	if err != nil {
		return 0, fmt.Errorf("error unmarshalling the JSON: %w", err)
	}

	// Return the random numbers
	return int(numbers[0]), nil
}
