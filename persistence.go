package main

import (
	firestore "cloud.google.com/go/firestore"
	context "context"
	ztchooks "github.com/zerotier/ztchooks"
	os "os"
)

var ctx = context.Background()

// parameters for Firestore are gotting from the running Container Environment
var gcp_project_name = os.Getenv("GCP_PROJECT")
var database_name = os.Getenv("FIRESTORE_DB_NAME")
var events_collection_name = os.Getenv("FIRESTORE_EVENTS_COLLECTION_NAME")
var test_document_path = os.Getenv("TEST_DOCUMENT_PATH")

func EventStoreClient() (*firestore.Client, error) {
	client, err1 := firestore.NewClientWithDatabase(ctx, gcp_project_name, database_name)
	if err1 != nil {
		Logger.Errorf("unable to connect to Firestore: %v", err1)
	}
	return client, err1
}

func GetDocument(client *firestore.Client, doc_path string) (*firestore.DocumentSnapshot, error) {
	ref := client.Doc(doc_path)
	doc, err := ref.Get(ctx)
	if err != nil {
		Logger.Errorf("Error fetching document from Firestore: %s", err)
	}
	return doc, err
}

func FetchTestDocument(client *firestore.Client) (*firestore.DocumentSnapshot, error) {
	doc, err := GetDocument(client, test_document_path)
	if err != nil {
		Logger.Errorf("Error fetching Aventine test document from Firestore: %s", err)
	}
	return doc, err
}

func PersistNetworkCreatedEvent(event *ztchooks.NetworkCreated) error {
	client, _ := EventStoreClient()
	_, _, err := client.Collection(events_collection_name).Add(ctx, event)
	if err != nil {
		if err != nil {
			Logger.Errorf("Unable to persist Network Created event to Firestore: %s", err)
		}
	}
	return err
}
