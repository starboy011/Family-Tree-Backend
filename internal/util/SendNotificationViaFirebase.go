package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/oauth2/google"
)

const (
	ProjectID   = "mulvansham-186e2"
	FCMEndpoint = "https://fcm.googleapis.com/v1/projects/%s/messages:send"
)

// Message represents the structure of the notification payload for FCM
type Message struct {
	Message struct {
		Token        string `json:"token"`
		Notification struct {
			Title string `json:"title"`
			Body  string `json:"body"`
		} `json:"notification"`
	} `json:"message"`
}

// getAuthToken fetches an OAuth 2.0 token for authenticating requests
func getAuthToken() (string, error) {
	ctx := context.Background()

	// Replace with the path to your service account key file
	credsFile := "internal/util/mulvansham-186e2-firebase-adminsdk-7rbje-9fa98f8714.json"

	// Read the service account key file
	credsData, err := ioutil.ReadFile(credsFile)
	if err != nil {
		return "", fmt.Errorf("failed to read service account key file: %v", err)
	}

	// Parse the JSON key file into the appropriate credentials object
	creds, err := google.CredentialsFromJSON(ctx, credsData, "https://www.googleapis.com/auth/firebase.messaging")
	if err != nil {
		return "", fmt.Errorf("failed to create credentials from JSON: %v", err)
	}

	// Create a token source from the credentials
	tokenSource := creds.TokenSource

	// Get the token
	token, err := tokenSource.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %v", err)
	}

	return token.AccessToken, nil
}

// SendNotificationViaFirebase sends a notification to the specified FCM token
func SendNotificationViaFirebase(token string, title string, body string) bool {
	authToken, err := getAuthToken()
	if err != nil {
		log.Printf("Error getting auth token: %v", err)
		return false
	}

	message := Message{
		Message: struct {
			Token        string `json:"token"`
			Notification struct {
				Title string `json:"title"`
				Body  string `json:"body"`
			} `json:"notification"`
		}{
			Token: token,
			Notification: struct {
				Title string `json:"title"`
				Body  string `json:"body"`
			}{
				Title: title,
				Body:  body,
			},
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshalling message: %v", err)
		return false
	}

	url := fmt.Sprintf(FCMEndpoint, ProjectID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonMessage))
	if err != nil {
		log.Printf("Error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	// Create an HTTP client to perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Received non-OK response code: %d", resp.StatusCode)
	}
	return true
}
