package main

import (
	"encoding/json"
	"errors"
	fmt "fmt"
	ztchooks "github.com/zerotier/ztchooks"
	zap "go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

// WebHook Secret is provided through the environment via a GCP Secret
var Psk = os.Getenv("ZEROTIER_ONE_WEBHOOK_SECRET")

type SimpleMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var ErrUnhandledHook = errors.New("unhandled hook type")
var ErrUnknownHookType = errors.New("unknown hook type")

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

func EventCatcher(w http.ResponseWriter, req *http.Request) {
	// read post body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}
	Logger.Infof("EventCatcher Body: %s", string(body))

	// get signature header from request.  If signature is empty, signature verification
	// is skipped
	signature := req.Header.Get("X-ZTC-Signature")
	if signature != "" {
		if err := ztchooks.VerifyHookSignature(Psk, signature, body, ztchooks.DefaultTolerance); err != nil {
			Logger.Errorf("error verifying Zero Tier One Hook Signature: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Signature Verification Failed"))
			return
		}
	}

	if err := processPayload(body, Logger); err != nil && err != ErrUnhandledHook {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - payload processing failed"))
		return
	}
}

func processPayload(payload []byte, logger *zap.SugaredLogger) error {
	hType, err := ztchooks.GetHookType(payload)
	logger.Debugf("Got Hook Type: %s", hType)

	if err != nil {
		return err
	}
	//
	switch hType {
	case ztchooks.NETWORK_JOIN:
		println(ztchooks.NETWORK_JOIN)
		var nmj ztchooks.NewMemberJoined
		if err := json.Unmarshal(payload, &nmj); err != nil {
			return err
		}

		// ... do something with NewMemberJoined data
	case ztchooks.NETWORK_AUTH:
		println(ztchooks.NETWORK_AUTH)
		var na ztchooks.NetworkMemberAuth
		if err := json.Unmarshal(payload, &na); err != nil {
			return err
		}

		// ... do something with NetworkMemberAuth data
	case ztchooks.NETWORK_DEAUTH:
		println(ztchooks.NETWORK_DEAUTH)
		var nd ztchooks.NetworkMemberDeauth
		if err := json.Unmarshal(payload, &nd); err != nil {
			return err
		}

		// ... do something with NetworkMemberDeauth data
	case ztchooks.NETWORK_CREATED:
		println(ztchooks.NETWORK_CREATED)
		var nc ztchooks.NetworkCreated
		if err := json.Unmarshal(payload, &nc); err != nil {
			return err
		}
		if err := PersistNetworkCreatedEvent(&nc); err != nil {
			return err
		}
	//
	// Continue with cases you wish to handle as needed
	//
	case ztchooks.HOOK_TYPE_UNKNOWN:
		return ErrUnknownHookType
	default:
		return ErrUnhandledHook
	}

	return nil
}
