package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

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
