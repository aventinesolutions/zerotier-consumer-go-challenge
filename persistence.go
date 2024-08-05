package main

import (
	firestore "cloud.google.com/go/firestore"
	context "context"
)

const (
	GCP_PROJECT   string = "aventine-k8s"
	DATABASE_NAME string = "zerotier-events"
	TEST_DOCUMENT string = "aventine/OLU9N1NkR0EYwYpGeBXi"
)

var ctx = context.Background()

func EventStoreClient() (*firestore.Client, error) {
	client, err1 := firestore.NewClientWithDatabase(ctx, GCP_PROJECT, DATABASE_NAME)
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
	doc, err := GetDocument(client, TEST_DOCUMENT)
	if err != nil {
		Logger.Errorf("Error fetching Aventine test document from Firestore: %s", err)
	}
	return doc, err
}
