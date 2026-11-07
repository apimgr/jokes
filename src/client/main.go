package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const defaultServer = "https://jokes.apimgr.us"

var (
	Version      string
	CommitID     string
	BuildDate    string
	BuildEpoch   string
	OfficialSite string
)

type apiEnvelope struct {
	Type    string          `json:"type"`
	Value   json.RawMessage `json:"value"`
	OK      bool            `json:"ok"`
	Error   string          `json:"error"`
	Message string          `json:"message"`
}

type joke struct {
	ID         int      `json:"id"`
	Type       string   `json:"type"`
	Joke       string   `json:"joke"`
	Setup      string   `json:"setup"`
	Delivery   string   `json:"delivery"`
	Categories []string `json:"categories"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("jokes-cli", flag.ContinueOnError)
	flags.SetOutput(stderr)
	var showHelp, showVersion bool
	server := flags.String("server", resolveServer(), "Jokes API base URL")
	category := flags.String("category", "", "Filter jokes by category")
	search := flags.String("search", "", "Search joke text and categories")
	id := flags.Int("id", 0, "Get a joke by ID")
	count := flags.Int("count", 1, "Number of random jokes to fetch (1-100)")
	format := flags.String("format", "text", "Output format: text or json")
	flags.BoolVar(&showHelp, "help", false, "Show help and exit")
	flags.BoolVar(&showHelp, "h", false, "Show help and exit")
	flags.BoolVar(&showVersion, "version", false, "Show version and exit")
	flags.BoolVar(&showVersion, "v", false, "Show version and exit")
	debug := flags.Bool("debug", false, "Log request details to stderr")
	color := flags.String("color", "auto", "Color output: auto, yes, or no")
	flags.Usage = func() {
		name := filepath.Base(os.Args[0])
		fmt.Fprintf(stderr, "%s queries a public Jokes API.\n\nUsage: %s [options]\n\n", name, name)
		flags.PrintDefaults()
	}
	if len(args) > 0 && args[0] == "help" {
		flags.Usage()
		return nil
	}
	if err := flags.Parse(args); err != nil {
		return err
	}
	if showHelp {
		flags.Usage()
		return nil
	}
	if showVersion {
		version := Version
		if version == "" {
			version = "1.0.1"
		}
		fmt.Fprintf(stdout, "%s %s\n", filepath.Base(os.Args[0]), version)
		return nil
	}
	if *color != "auto" && *color != "yes" && *color != "no" {
		return errors.New("--color must be auto, yes, or no")
	}
	if _, set := os.LookupEnv("NO_COLOR"); set {
		*color = "no"
	}
	_ = color
	if *count < 1 || *count > 100 {
		return errors.New("--count must be between 1 and 100")
	}
	if *id < 0 {
		return errors.New("--id must be greater than zero")
	}
	if *format != "text" && *format != "json" {
		return errors.New("--format must be text or json")
	}
	base, err := url.Parse(strings.TrimRight(*server, "/"))
	if err != nil || (base.Scheme != "https" && base.Scheme != "http") || base.Host == "" || base.User != nil {
		return errors.New("--server must be an absolute HTTP or HTTPS URL without credentials")
	}
	if base.Scheme != "https" && base.Hostname() != "localhost" && base.Hostname() != "127.0.0.1" && base.Hostname() != "::1" {
		return errors.New("--server must use HTTPS outside localhost")
	}

	path := "/api/v1/jokes/random"
	query := url.Values{}
	switch {
	case *id > 0:
		path = "/api/v1/jokes/" + strconv.Itoa(*id)
	case strings.TrimSpace(*search) != "":
		path = "/api/v1/jokes/search"
		query.Set("q", strings.TrimSpace(*search))
		query.Set("limit", strconv.Itoa(*count))
	case *count > 1:
		path = "/api/v1/jokes/random/" + strconv.Itoa(*count)
	}
	if *category != "" {
		query.Set("category", *category)
	}
	endpoint := *base
	endpoint.Path = strings.TrimRight(endpoint.Path, "/") + path
	endpoint.RawQuery = query.Encode()
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")
	version := Version
	if version == "" {
		version = "1.0.1"
	}
	request.Header.Set("User-Agent", "jokes-cli/"+version)
	if *debug {
		fmt.Fprintf(stderr, "GET %s\n", endpoint.Redacted())
	}

	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request Jokes API: %w", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read Jokes API response: %w", err)
	}
	var envelope apiEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("decode Jokes API response: %w", err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if envelope.Message != "" {
			return fmt.Errorf("jokes API returned %s: %s", response.Status, envelope.Message)
		}
		return fmt.Errorf("jokes API returned %s: %s", response.Status, strings.TrimSpace(string(body)))
	}
	if *format == "json" {
		var value any
		if err := json.Unmarshal(envelope.Value, &value); err != nil {
			return fmt.Errorf("decode joke data: %w", err)
		}
		encoded, err := json.MarshalIndent(value, "", "  ")
		if err != nil {
			return fmt.Errorf("format joke data: %w", err)
		}
		_, err = fmt.Fprintln(stdout, string(encoded))
		return err
	}
	var jokes []joke
	if len(envelope.Value) > 0 && envelope.Value[0] == '[' {
		if err := json.Unmarshal(envelope.Value, &jokes); err != nil {
			return fmt.Errorf("decode jokes: %w", err)
		}
	} else if strings.Contains(string(envelope.Value), `"jokes"`) {
		var collection struct {
			Jokes []joke `json:"jokes"`
		}
		if err := json.Unmarshal(envelope.Value, &collection); err != nil {
			return fmt.Errorf("decode jokes: %w", err)
		}
		jokes = collection.Jokes
	} else {
		var single joke
		if err := json.Unmarshal(envelope.Value, &single); err != nil {
			return fmt.Errorf("decode joke: %w", err)
		}
		jokes = []joke{single}
	}
	for index, item := range jokes {
		if item.Type == "twopart" {
			fmt.Fprintln(stdout, item.Setup)
			fmt.Fprintln(stdout, item.Delivery)
		} else {
			fmt.Fprintln(stdout, item.Joke)
		}
		if index < len(jokes)-1 {
			fmt.Fprintln(stdout)
		}
	}
	if len(jokes) == 0 {
		return errors.New("no jokes matched the request")
	}
	return nil
}

func resolveServer() string {
	if server := strings.TrimSpace(os.Getenv("JOKES_SERVER_PRIMARY")); server != "" {
		return server
	}
	if OfficialSite != "" {
		return OfficialSite
	}
	return defaultServer
}
