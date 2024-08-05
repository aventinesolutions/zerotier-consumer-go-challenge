package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// WebHook Secret is provided through the environment via a GCP Secret
var Psk = os.Getenv("ZEROTIER_ONE_WEBHOOK_SECRET")

type SimpleMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

/* just say "hello!" with a JSON response */
func HelloWorld(w http.ResponseWriter, req *http.Request) {
	Logger.Info("someone wants to say hello")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SimpleMessage{
		Type:    "hello",
		Message: "Hello, Aventine Solutions!",
	})
}

/* "liveness" for orchestration with a JSON response */
func Liveness(w http.ResponseWriter, req *http.Request) {
	Logger.Debug("Liveness Check")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SimpleMessage{
		Type:    "livez",
		Message: "true",
	})
}

/* "readiness" for orchestration with a JSON response */
func Readiness(w http.ResponseWriter, req *http.Request) {
	Logger.Debug("Readiness Check")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SimpleMessage{
		Type:    "readyz",
		Message: "true",
	})
}

/* check that the ZeroTier One Webhook Token is set correctly */
func CheckToken(w http.ResponseWriter, req *http.Request) {
	Logger.Debug("Check Token")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SimpleMessage{
		Type:    "token_set",
		Message: fmt.Sprintf("%t", len(Psk) == 64),
	})
}

func CheckFirestore(w http.ResponseWriter, req *http.Request) {
	Logger.Debug("Check Firestore Events database")
	client, _ := EventStoreClient()
	defer client.Close()
	doc, _ := FetchTestDocument(client)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(SimpleMessage{
		Type:    "test_firestore_document",
		Message: fmt.Sprintf("%+v", doc.Data()),
	})
}
