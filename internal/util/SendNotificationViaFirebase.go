package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
// func getAuthToken() (string, error) {
// 	ctx := context.Background()

// 	// Replace with the path to your service account key file
// 	credsFile := "internal/util/mulvansham-186e2-firebase-adminsdk-7rbje-71755a2bcf.json"

// 	// Read the service account key file
// 	credsData, err := ioutil.ReadFile(credsFile)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to read service account key file: %v", err)
// 	}

// 	// Parse the JSON key file into the appropriate credentials object
// 	creds, err := google.CredentialsFromJSON(ctx, credsData, "https://www.googleapis.com/auth/firebase.messaging")
// 	if err != nil {
// 		return "", fmt.Errorf("failed to create credentials from JSON: %v", err)
// 	}

// 	// Create a token source from the credentials
// 	tokenSource := creds.TokenSource

// 	// Get the token
// 	token, err := tokenSource.Token()
// 	if err != nil {
// 		return "", fmt.Errorf("failed to get token: %v", err)
// 	}

// 	return token.AccessToken, nil
// }

// SendNotificationViaFirebase sends a notification to the specified FCM token
func SendNotificationViaFirebase(token string, title string, body string) bool {
	// authToken, err := getAuthToken()
	// if err != nil {
	// 	log.Printf("Error getting auth token: %v", err)
	// 	return false
	// }

	authToken := "ya29.c.c0ASRK0GYHpXzq8CYcOdk6RzZnZRHkGCf4d-12xP-0hk1uwpF_zyGpv-Sox1EbH0tTTx_HYLejg-zl0xHajEdoJiiVje_DKwe9fK4L0T8TXQFC-lh_9FWKOaRhtAPeXeaxeQjKjniDuq3n6JBgOtA7jFuLhEsjeBPb7037dyHfZ2ySHQEuIfheauN0DZg4hDGZ9euewvpbLwuCG77v0FJ4EPPYWYABKqb8uHTnhQnnuBWCkNVKzWw5a0gz3o3X8UUcCxTtbXvrzB7YloGew86-CQsVDaI41zLlECDHUvpBWiEHQCVvQn_xiMbKKoeaYfiF8GH9R2anVtfayvoOHgoWL5lFHaI73-KY_16NvUIW5W2gF99dH5rca8A_T385CYmZudgBnaJn0SFYqy019-2ta4zkbdUrvMQbhOwfJYejrzyFcR_cpXsr7oqtzbcr9b-ehb08WnSi5tUU86V_Q4do-ugi6MZiOlSzJeo-kds3irWbeOzwMlb8YYIuezqk2rg-ihk9v_jk-h-UwbOvXi_26Fi9riFRuXYnY-axze_jF5mj3boIxqRbudidl6Q83J6qdYJ2u6k1v1pQ5ilOzh9yg2OXfY4OgnzJQRc1w_Yvo46mnq7Ig_r7QnJ__2nq2-kdaIIeZgxbBF_1cwt1djWJabSJ-tcRlzWsZr__rq2oaz667VrnrvZohRy6znfIm6-fJSraJhwyvxWquqc97V9-1Z2lg7Qz0Jb2lbg4u5VFOjzWSgs8eyVp0VtRhjuqw40OlRY-2Xw6OV6g9o00yX7ZpcnX5Si_JZwf8Rlm1wvlVhcImwZkdfZ8U83_-WgZ62_BZ9MazSFI2wkFpVwFB44iWJxB33_Bydem9pJ_qYb8nyo990UiqgVn8i49J46UIW5w1WgS6r3f1xXUqtixbVwgecqtOI_3o2WS9veQi2a2U4ryoWJwF0IsyUvsh8rF4xZIg1pb2gJkRza32jxgks8Q9d350kFnYB_rbuFjBI9j94__zOmUV8vt99F"
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
