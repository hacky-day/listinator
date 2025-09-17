package typifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
)

type input struct {
	Product string `json:"product"`
}

type output struct {
	Product string `json:"product"`
	Type    string `json:"type"`
	UUID    string `json:"uuid"`
}

type Typifier struct {
	http.Client
	url string
}

func Init() (*Typifier, error) {
	url := os.Getenv("LISTINATOR_TYPIFIER_URL")
	if url == "" {
		return nil, nil
	}
	return &Typifier{
		Client: http.Client{},
		url:    url,
	}, nil
}

func (t *Typifier) GetUUID(entry string) (uuid.UUID, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(input{
		Product: entry,
	}); err != nil {
		return uuid.Nil, fmt.Errorf("unable to jsonify input, %w", err)
	}
	resp, err := t.Post(t.url, "application/json", &buf)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to request, %w", err)
	}
	defer resp.Body.Close()

	var o output
	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		return uuid.Nil, fmt.Errorf("unable to parse output, %w", err)
	}

	u, err := uuid.Parse(o.UUID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("unable to parse uuid in response, %w", err)
	}

	return u, nil
}
