package localdecoder

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	base string
	http *http.Client
}

func New(endpoint string) (*Client, error) {
	u, err := url.Parse(endpoint)
	if err != nil || u.Scheme != "http" || u.User != nil ||
		(u.Hostname() != "127.0.0.1" && u.Hostname() != "::1" && u.Hostname() != "localhost") ||
		u.Path != "" || u.RawQuery != "" || u.Fragment != "" {
		return nil, errors.New("decoder: loopback HTTP address without path or credentials required")
	}
	return &Client{endpoint, &http.Client{Timeout: 45 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error {
		return errors.New("decoder: redirects forbidden")
	}}}, nil
}

func (c *Client) Post(ctx context.Context, path string, body, result any) error {
	if path != "/completion" && path != "/tokenize" {
		return errors.New("decoder: unsupported native endpoint")
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	response, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("decoder: HTTP status %d", response.StatusCode)
	}
	limited := &io.LimitedReader{R: response.Body, N: (16 << 20) + 1}
	decoder := json.NewDecoder(limited)
	if err := decoder.Decode(result); err != nil {
		return fmt.Errorf("decoder JSON: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF || limited.N <= 0 {
		return errors.New("decoder: trailing or oversized JSON response")
	}
	return nil
}

type TokenProbability struct {
	ID      int      `json:"id"`
	Bytes   []byte   `json:"bytes"`
	LogProb *float64 `json:"logprob"`
}

type Completion struct {
	Truncated bool `json:"truncated"`
	Positions []struct {
		ID  int                `json:"id"`
		Top []TokenProbability `json:"top_logprobs"`
	} `json:"completion_probabilities"`
	Settings struct {
		NProbs            int  `json:"n_probs"`
		PostSamplingProbs bool `json:"post_sampling_probs"`
	} `json:"generation_settings"`
}

func (c *Client) Next(ctx context.Context, prompt any, grammar string, vocabulary int) (Completion, error) {
	if vocabulary < 1 || vocabulary > 200000 {
		return Completion{}, errors.New("decoder: vocabulary outside [1,200000]")
	}
	req := map[string]any{
		"prompt": prompt, "n_predict": 1, "n_probs": vocabulary,
		"post_sampling_probs": false, "temperature": 1, "top_k": 0, "top_p": 1, "min_p": 0,
		"repeat_penalty": 1, "frequency_penalty": 0, "presence_penalty": 0,
		"samplers": []string{"temperature"}, "seed": 42, "cache_prompt": true, "return_tokens": true,
	}
	if grammar != "" {
		req["grammar"] = grammar
	}
	var response Completion
	if err := c.Post(ctx, "/completion", req, &response); err != nil {
		return Completion{}, err
	}
	if response.Truncated || len(response.Positions) != 1 || len(response.Positions[0].Top) != vocabulary ||
		response.Settings.NProbs != vocabulary || response.Settings.PostSamplingProbs {
		return Completion{}, errors.New("decoder: full pre-sampling vocabulary not returned")
	}
	return response, nil
}

func (c *Client) Tokenize(ctx context.Context, text string, addSpecial bool) ([]int, error) {
	var response struct {
		Tokens []int `json:"tokens"`
	}
	err := c.Post(ctx, "/tokenize", map[string]any{"content": text, "add_special": addSpecial, "parse_special": true}, &response)
	if err != nil {
		return nil, err
	}
	if len(response.Tokens) == 0 || len(response.Tokens) > 1024 {
		return nil, errors.New("decoder: empty or oversized tokenization")
	}
	return response.Tokens, nil
}

func Grammar(labels []string) string {
	quoted := make([]string, len(labels))
	for i, label := range labels {
		quoted[i] = strconv.Quote(label)
	}
	return "root ::= " + strings.Join(quoted, " | ")
}
