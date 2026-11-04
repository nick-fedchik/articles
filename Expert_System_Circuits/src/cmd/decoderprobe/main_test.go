package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestProbeUsesExplicitGrammarAndNoPrivateCorpus(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		if r.URL.Path != "/api/generate" || r.Method != http.MethodPost {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		var request struct {
			Format struct {
				Properties map[string]struct {
					Enum []string `json:"enum"`
				} `json:"properties"`
			} `json:"format"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		labels := request.Format.Properties["modality"].Enum
		if len(labels) != 4 && len(labels) != 5 {
			t.Error("missing grammar")
		}
		if err := json.NewEncoder(w).Encode(map[string]any{
			"response": `{"modality":"OBLIGATION"}`, "done": true, "done_reason": "stop",
		}); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	if err := run(context.Background(), server.URL, "synthetic-test-model", io.Discard); err != nil {
		t.Fatal(err)
	}
	if calls != 2*len(fixtures) {
		t.Fatalf("got %d calls", calls)
	}
}

func TestProbeRejectsMalformedResults(t *testing.T) {
	for _, response := range []string{
		`{"response":"{\"modality\":\"UNKNOWN\"}","done":true,"done_reason":"stop"}`,
		`{"response":"{\"modality\":\"OBLIGATION\"} garbage","done":true,"done_reason":"stop"}`,
		`{"response":"{\"modality\":\"OBLIGATION\"}","done":true,"done_reason":"length"}`,
		`{"response":"{\"modality\":\"OBLIGATION\"}","done":true,"done_reason":"stop","total_duration":-1}`,
		`{"response":"{\"modality\":\"OBLIGATION\"}","done":true,"done_reason":"stop"} trailing`,
		`invalid`,
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, err := io.Copy(w, strings.NewReader(response))
			if err != nil {
				t.Error(err)
			}
		}))
		_, err := classify(context.Background(), server.Client(), server.URL, "test", fixtures[0], false)
		server.Close()
		if err == nil {
			t.Fatalf("accepted malformed result: %s", response)
		}
	}
	if err := run(context.Background(), "https://example.com", "test", io.Discard); err == nil {
		t.Fatal("external endpoint accepted")
	}
}
