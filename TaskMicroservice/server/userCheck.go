package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"taskService/models"
)

// CheckUserExists calls the UserMicroservice to check if the user with the given email exists
// and returns the user details if found, a boolean indicating existence, and an error if any.
func CheckUserExists(email string) (models.User, bool, error) {
	var userResponse models.UserResponce

	// Perform HTTP GET request to UserMicroservice
	resp, err := http.Get("http://localhost:9091/api/v1/users/" + email)
	if err != nil {
		return userResponse.Message, false, fmt.Errorf("error calling UserMicroservice: %w", err)
	}
	defer resp.Body.Close()

	// Check if status code is 200 OK
	if resp.StatusCode != http.StatusOK {
		return userResponse.Message, false, fmt.Errorf("UserMicroservice returned status: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return userResponse.Message, false, fmt.Errorf("error reading UserMicroservice response: %w", err)
	}

	// Unmarshal the response into a UserResponse model
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return userResponse.Message, false, fmt.Errorf("error unmarshalling response: %w", err)
	}

	fmt.Println("body ", string(body))

	// Check if the email from the response matches the requested email
	if userResponse.Message.EmailId == email {
		return userResponse.Message, true, nil
	}

	return userResponse.Message, false, nil
}
