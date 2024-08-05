package main

import (
	"encoding/json"
	"errors"
	ztchooks "github.com/zerotier/ztchooks"
	zap "go.uber.org/zap"
	"io"
	"net/http"
	"os"
)

const (
	ServerAddr = ":4444"
)

// WebHook Secret is provided through the environment via a GCP Secret
var psk = os.Getenv("ZEROTIER_ONE_WEBHOOK_SECRET")

var Logger = zap.NewExample().Sugar()

var ErrUnhandledHook = errors.New("unhandled hook type")
var ErrUnknownHookType = errors.New("unknown hook type")

func eventCatcher(w http.ResponseWriter, req *http.Request) {
	// read post body
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err.Error())
	}

	// get signature header from request.  If signature is empty, signature verification
	// is skipped
	signature := req.Header.Get("X-ZTC-Signature")
	if signature != "" {
		if err := ztchooks.VerifyHookSignature(psk, signature, body, ztchooks.DefaultTolerance); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Signature Verification Failed"))
			return
		}
	}

	if err := processPayload(body); err != nil && err != ErrUnhandledHook {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - payload processing failed"))
		return
	}
}

func processPayload(payload []byte) error {
	hType, err := ztchooks.GetHookType(payload)
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

		// ... do something with NetworkCreated data

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

func main() {
	defer func(logger *zap.SugaredLogger) {
		err := logger.Sync()
		if err != nil {
			panic("unable to defer Zap logging, exiting!")
		}
	}(Logger)
	Logger.Info("starting ZeroTier Consumer Coding Challenge")

	http.HandleFunc("/hello", HelloWorld)
	http.HandleFunc("/livez", Liveness)
	http.HandleFunc("/readyz", Readiness)
	http.HandleFunc("/check_token", CheckToken)
	http.HandleFunc("/check_firestore", CheckFirestore)
	http.HandleFunc("/event", eventCatcher)
	err2 := http.ListenAndServe(ServerAddr, nil)
	if err2 != nil {
		Logger.Errorf("error starting Web Service: %s", err2)
		panic("unable to start Web Service, exiting!")
	}
}
