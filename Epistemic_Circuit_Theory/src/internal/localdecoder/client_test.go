package localdecoder

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFullVocabularyProtocol(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var request map[string]any
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Error(err)
		}
		var result any
		switch r.URL.Path {
		case "/tokenize":
			result = map[string]any{"tokens": []int{1, 2}}
		case "/completion":
			if request["post_sampling_probs"] != false || request["n_probs"] != float64(3) {
				t.Error("did not request the full pre-sampling distribution")
			}
			result = map[string]any{
				"generation_settings": map[string]any{"n_probs": 3, "post_sampling_probs": false},
				"completion_probabilities": []any{map[string]any{
					"id": 1, "top_logprobs": []any{
						map[string]any{"id": 0, "bytes": []int{79}, "logprob": -1.0},
						map[string]any{"id": 1, "bytes": []int{78}, "logprob": -1.0},
						map[string]any{"id": 2, "bytes": []int{120}, "logprob": -1.5},
					},
				}},
			}
		default:
			t.Error("unexpected native endpoint")
		}
		if err := json.NewEncoder(w).Encode(result); err != nil {
			t.Error(err)
		}
	}))
	defer server.Close()
	client, err := New(server.URL)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Tokenize(context.Background(), "public", false); err != nil {
		t.Fatal(err)
	}
	got, err := client.Next(context.Background(), "public", Grammar([]string{"YES", "NO"}), 3)
	if err != nil || len(got.Positions[0].Top) != 3 {
		t.Fatal(got, err)
	}
}

func TestRejectsExternalAddressesAndIncompleteResponses(t *testing.T) {
	for _, endpoint := range []string{"https://example.com", "http://127.0.0.1/path", "http://user:pass@localhost"} {
		if _, err := New(endpoint); err == nil {
			t.Fatal("invalid address accepted")
		}
	}
	for _, body := range []string{
		`not json`,
		`{} trailing`,
		`{"completion_probabilities":[],"generation_settings":{"n_probs":3}}`,
		`{"truncated":true,"completion_probabilities":[{}],"generation_settings":{"n_probs":3}}`,
	} {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte(body)); err != nil {
				t.Error(err)
			}
		}))
		client, err := New(server.URL)
		if err != nil {
			t.Fatal(err)
		}
		_, err = client.Next(context.Background(), "public", "", 3)
		server.Close()
		if err == nil {
			t.Fatal("incomplete distribution accepted")
		}
	}
	if !strings.Contains(Grammar([]string{`a"b`, "NO"}), `\"`) {
		t.Fatal("grammar string was not escaped")
	}
}
