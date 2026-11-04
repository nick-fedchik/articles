package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nick-fedchik/articles/Expert_System_Circuits/src/annotation"
)

func readJSON(path string, value any) error {
	if path == "" {
		return errors.New("an explicit local input file is required")
	}
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := io.ReadAll(io.LimitReader(file, (64<<20)+1))
	if err != nil {
		return err
	}
	if len(data) > 64<<20 {
		return errors.New("annotation input exceeds 64 MiB")
	}
	if err := checkJSONKeys(data); err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return fmt.Errorf("strict annotation JSON: %w", err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return errors.New("annotation input has trailing content or exceeds 64 MiB")
	}
	return nil
}

func checkJSONKeys(data []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	var value func(int) error
	value = func(depth int) error {
		if depth > 64 {
			return errors.New("annotation JSON nesting exceeds 64")
		}
		token, err := decoder.Token()
		if err != nil {
			return err
		}
		delim, composite := token.(json.Delim)
		if !composite {
			return nil
		}
		seen := map[string]bool{}
		for decoder.More() {
			if delim == '{' {
				key, err := decoder.Token()
				if err != nil {
					return err
				}
				name, ok := key.(string)
				if !ok {
					return errors.New("annotation JSON object key is not a string")
				}
				name = strings.ToLower(name)
				if seen[name] {
					return errors.New("annotation JSON contains duplicate case-insensitive keys")
				}
				seen[name] = true
			}
			if err := value(depth + 1); err != nil {
				return err
			}
		}
		_, err = decoder.Token()
		return err
	}
	if err := value(0); err != nil {
		return err
	}
	if _, err := decoder.Token(); err != io.EOF {
		return errors.New("annotation JSON must contain exactly one value")
	}
	return nil
}

func writeNew(path string, value any) error {
	if path == "" {
		return errors.New("an explicit new output path is required")
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	if _, err := file.Write(append(data, '\n')); err != nil {
		file.Close()
		return err
	}
	return file.Close()
}

func output(writer io.Writer, value any) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func run(args []string, writer io.Writer) error {
	if len(args) == 0 {
		return errors.New("commands: prepare, validate, compare, adjudicate, evaluate, demo")
	}
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	corpusPath := flags.String("corpus", "", "explicit local corpus JSON")
	aPath := flags.String("a", "", "first submitted annotation")
	bPath := flags.String("b", "", "second submitted annotation")
	annotationPath := flags.String("annotation", "", "submitted annotation to validate")
	arbiterPath := flags.String("arbiter", "", "explicit human adjudication")
	goldPath := flags.String("gold", "", "adjudicated labels")
	predPath := flags.String("predictions", "", "system predictions for the selected split")
	planPath := flags.String("plan", "", "document-family split plan")
	out := flags.String("out", "", "new output file, or new template directory for prepare")
	split := flags.String("split", "test", "train, calibration or test")
	seed := flags.Int64("seed", 42, "deterministic document-family split seed")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional annotation arguments")
	}
	if args[0] == "demo" {
		demo, err := annotation.SyntheticDemo()
		if err != nil {
			return err
		}
		plan, _, _, err := annotation.Prepare(demo.Corpus, *seed)
		if err != nil {
			return err
		}
		agreement, err := annotation.Compare(demo.Corpus, demo.A, demo.B)
		if err != nil {
			return err
		}
		gold, err := annotation.Adjudicate(demo.Corpus, demo.A, demo.B, demo.Adjudication)
		if err != nil {
			return err
		}
		return output(writer, struct {
			Scope     string
			Plan      annotation.SplitPlan
			Agreement annotation.Agreement
			Gold      annotation.Gold
		}{"synthetic workflow test only; no human annotation or independent corpus was produced", plan, agreement, gold})
	}
	var corpus annotation.Corpus
	if err := readJSON(*corpusPath, &corpus); err != nil {
		return err
	}
	if err := corpus.Validate(); err != nil {
		return err
	}
	switch args[0] {
	case "prepare":
		plan, a, b, err := annotation.Prepare(corpus, *seed)
		if err != nil {
			return err
		}
		if *out == "" {
			return errors.New("prepare requires an explicit new output directory")
		}
		files := []struct {
			name  string
			value any
		}{
			{"split.json", plan}, {"annotator-a.template.json", a}, {"annotator-b.template.json", b},
		}
		for _, file := range files {
			if _, err := os.Lstat(filepath.Join(*out, file.name)); err == nil || !os.IsNotExist(err) {
				return errors.New("prepare refuses to overwrite an existing template or split")
			}
		}
		if err := os.MkdirAll(*out, 0750); err != nil {
			return err
		}
		for _, file := range files {
			if err := writeNew(filepath.Join(*out, file.name), file.value); err != nil {
				return err
			}
		}
		return output(writer, map[string]any{"created": 3, "status": "draft templates only; no independent labels created"})
	case "validate":
		var submitted annotation.Annotation
		if err := readJSON(*annotationPath, &submitted); err != nil {
			return err
		}
		if _, err := annotation.ValidateAnnotation(corpus, submitted); err != nil {
			return err
		}
		return output(writer, map[string]any{"valid_records": len(submitted.Records), "independence_verified": false})
	case "compare", "adjudicate":
		var a, b annotation.Annotation
		if err := readJSON(*aPath, &a); err != nil {
			return err
		}
		if err := readJSON(*bPath, &b); err != nil {
			return err
		}
		if args[0] == "compare" {
			agreement, err := annotation.Compare(corpus, a, b)
			if err != nil {
				return err
			}
			return output(writer, agreement)
		}
		var adjudication annotation.Adjudication
		if err := readJSON(*arbiterPath, &adjudication); err != nil {
			return err
		}
		gold, err := annotation.Adjudicate(corpus, a, b, adjudication)
		if err != nil {
			return err
		}
		return writeNew(*out, gold)
	case "evaluate":
		var plan annotation.SplitPlan
		var gold annotation.Gold
		var predictions annotation.Predictions
		if err := readJSON(*planPath, &plan); err != nil {
			return err
		}
		if err := readJSON(*goldPath, &gold); err != nil {
			return err
		}
		if err := readJSON(*predPath, &predictions); err != nil {
			return err
		}
		report, err := annotation.Evaluate(corpus, plan, gold, predictions, *split)
		if err != nil {
			return err
		}
		return output(writer, report)
	default:
		return errors.New("unknown annotation command")
	}
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
