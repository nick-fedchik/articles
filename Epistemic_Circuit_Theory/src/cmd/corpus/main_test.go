package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadBothLocalSchemasWithoutPublishingText(t *testing.T) {
	for _, body := range []string{
		`{"clauses":[{"source_ref":"RFC-100#sec-1","text":"Public fixture shall be recorded.","arbitration":{"modality":"OBLIGATION"}}]}`,
		`{"items":[{"source":"RFC 100","text":"Public fixture shall be recorded.","category":"OBLIGATION"}]}`,
	} {
		path := filepath.Join(t.TempDir(), "corpus.json")
		if err := os.WriteFile(path, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}

		rows, digest, err := load(path)
		if err != nil || len(rows) != 1 || len(digest) != 64 || rows[0].Label != "OBLIGATION" {
			t.Fatal(rows, digest, err)
		}
	}

}

func FuzzDecodeCorpus(f *testing.F) {
	f.Add([]byte(`{"items":[{"source":"public","text":"fixture","category":"NO_ATOM"}]}`))
	f.Add([]byte(`{"clauses":[]}`))
	f.Fuzz(func(t *testing.T, data []byte) {
		if len(data) > 1<<20 {
			t.Skip("bounded corpus parser fuzz input")
		}
		rows, digest, err := decode(data)
		if err == nil && (len(rows) == 0 || len(digest) != 64) {
			t.Fatal("successful decode without rows or input fingerprint")
		}
	})
}
