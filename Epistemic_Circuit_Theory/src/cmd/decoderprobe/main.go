package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type fixture struct {
	ID, Text, Expected string
}

var fixtures = []fixture{
	{"obligation-en", "The controller shall record each alarm.", "OBLIGATION"},
	{"prohibition-en", "The client shall not send an empty password.", "PROHIBITION"},
	{"description-en", "The controller recorded an alarm yesterday.", "NO_ATOM"},
	{"description-uk", "Учора сервер надіслав повідомлення клієнтові.", "NO_ATOM"},
	{"recommendation-en", "The operator should review the diagnostic log.", "RECOMMENDATION"},
	{"permission-en", "The user may export the report.", "PERMISSION"},
}

type observation struct {
	Fixture        string `json:"fixture"`
	NullAllowed    bool   `json:"null_allowed"`
	Expected       string `json:"expected"`
	Observed       string `json:"observed"`
	Correct        bool   `json:"correct"`
	ServerDuration int64  `json:"server_duration_ns"`
	Tokens         int    `json:"output_tokens"`
}

func classify(ctx context.Context, client *http.Client, endpoint, model string, item fixture, nullAllowed bool) (observation, error) {
	labels := []string{"OBLIGATION", "PROHIBITION", "RECOMMENDATION", "PERMISSION"}
	if nullAllowed {
		labels = append(labels, "NO_ATOM")
	}
	request := map[string]any{
		"model": model, "stream": false, "keep_alive": 0,
		"prompt": "Classify the following sentence by its explicit normative modality. " +
			"SHALL is OBLIGATION, SHALL NOT is PROHIBITION, SHOULD is RECOMMENDATION, MAY is PERMISSION. " +
			"A description of an event is NO_ATOM. Return only a JSON object with the field modality.\nSentence: " + item.Text,
		"format": map[string]any{
			"type":       "object",
			"properties": map[string]any{"modality": map[string]any{"type": "string", "enum": labels}},
			"required":   []string{"modality"}, "additionalProperties": false,
		},
		"options": map[string]any{"temperature": 0, "seed": 42, "num_predict": 64, "num_ctx": 1024, "num_thread": 2, "num_gpu": 0},
	}
	body, err := json.Marshal(request)
	if err != nil {
		return observation{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint+"/api/generate", bytes.NewReader(body))
	if err != nil {
		return observation{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)
	if err != nil {
		return observation{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return observation{}, fmt.Errorf("decoder HTTP status %d", response.StatusCode)
	}
	var envelope struct {
		Response      string `json:"response"`
		Done          bool   `json:"done"`
		DoneReason    string `json:"done_reason"`
		TotalDuration int64  `json:"total_duration"`
		EvalCount     int    `json:"eval_count"`
	}
	envelopeDecoder := json.NewDecoder(io.LimitReader(response.Body, 1<<20))
	if err := envelopeDecoder.Decode(&envelope); err != nil {
		return observation{}, fmt.Errorf("decoder envelope: %w", err)
	}
	if err := envelopeDecoder.Decode(new(any)); err != io.EOF {
		return observation{}, errors.New("decoder envelope has trailing content")
	}
	if !envelope.Done || envelope.DoneReason != "stop" {
		return observation{}, errors.New("decoder did not finish normally")
	}
	if envelope.TotalDuration < 0 || envelope.EvalCount < 0 {
		return observation{}, errors.New("decoder returned negative measurement counters")
	}
	var decoded struct {
		Modality string `json:"modality"`
	}
	decoder := json.NewDecoder(bytes.NewBufferString(envelope.Response))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&decoded); err != nil {
		return observation{}, fmt.Errorf("decoder structured result: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return observation{}, errors.New("decoder structured result has trailing content")
	}
	valid := false
	for _, label := range labels {
		valid = valid || decoded.Modality == label
	}
	if !valid {
		return observation{}, errors.New("decoder returned an out-of-grammar label")
	}
	return observation{item.ID, nullAllowed, item.Expected, decoded.Modality, decoded.Modality == item.Expected,
		envelope.TotalDuration, envelope.EvalCount}, nil
}

func run(ctx context.Context, endpoint, model string, writer io.Writer) error {
	address, err := url.Parse(endpoint)
	if err != nil || address.Scheme != "http" || address.User != nil ||
		(address.Hostname() != "127.0.0.1" && address.Hostname() != "localhost" && address.Hostname() != "::1") ||
		address.Path != "" || address.RawQuery != "" || address.Fragment != "" {
		return errors.New("probe requires a local HTTP endpoint without credentials or URL path")
	}
	if model == "" {
		return errors.New("an explicitly selected installed model is required; the probe never downloads models")
	}
	client := &http.Client{Timeout: 45 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return errors.New("decoder redirects are not allowed")
	}}
	rows := []observation{}
	for _, nullAllowed := range []bool{false, true} {
		for _, item := range fixtures {
			got, err := classify(ctx, client, endpoint, model, item, nullAllowed)
			if err != nil {
				return fmt.Errorf("%s null=%t: %w", item.ID, nullAllowed, err)
			}
			rows = append(rows, got)
		}
	}
	result := struct {
		Model         string        `json:"model"`
		Kind          string        `json:"kind"`
		PreMaskLogits string        `json:"pre_mask_logits"`
		Observations  []observation `json:"observations"`
	}{model, "synthetic local decoder pilot, not an independent corpus evaluation",
		"not provided by this endpoint; R_t and prefix-conditioning hypotheses remain untested", rows}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

func main() {
	endpoint := flag.String("endpoint", "http://127.0.0.1:11434", "local Ollama endpoint")
	model := flag.String("model", "", "explicit installed model name; no automatic download")
	flag.Parse()
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	if err := run(ctx, *endpoint, *model, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
