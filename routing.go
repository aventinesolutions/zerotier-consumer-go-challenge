package main

import http "net/http"

func SetRoutes() {
	http.HandleFunc("/hello", HelloWorld)
	http.HandleFunc("/livez", Liveness)
	http.HandleFunc("/readyz", Readiness)
	http.HandleFunc("/check_token", CheckToken)
	http.HandleFunc("/check_firestore", CheckFirestore)
	http.HandleFunc("/event", EventCatcher)
}
