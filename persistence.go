package main

import (
	firestore "cloud.google.com/go/firestore"
	context "context"
	zap "go.uber.org/zap"
)

const (
	GCP_PROJECT   string = "aventine-k8s"
	DATABASE_NAME string = "zerotier-events"
	TEST_DOCUMENT string = "aventine/OLU9N1NkR0EYwYpGeBXi"
)

var ctx = context.Background()

func EventStoreClient(logger *zap.SugaredLogger) (*firestore.Client, error) {
	client, err1 := firestore.NewClientWithDatabase(ctx, GCP_PROJECT, DATABASE_NAME)
	if err1 != nil {
		logger.Errorf("unable to connect to Firestore: %v", err1)
	}
	return client, err1
}

func GetDocument(client *firestore.Client, logger *zap.SugaredLogger, doc_path string) (*firestore.DocumentSnapshot, error) {
	ref := client.Doc(doc_path)
	doc, err := ref.Get(ctx)
	if err != nil {
		logger.Errorf("Error fetching document from Firestore: %s", err)
	}
	return doc, err
}

func FetchTestDocument(client *firestore.Client, logger *zap.SugaredLogger) (*firestore.DocumentSnapshot, error) {
	doc, err := GetDocument(client, logger, TEST_DOCUMENT)
	if err != nil {
		logger.Errorf("Error fetching Aventine test document from Firestore: %s", err)
	}
	return doc, err
}
