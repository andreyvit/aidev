package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andreyvit/openai"
)

var (
	openAICreds openai.Credentials
)

//go:embed prompt.txt
var systemPrompt string

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	var (
		envFile   string
		rootDirs  []string
		include   []string
		exclude   []string
		unexclude []string
		codeFile  string
		replay    bool
		prompt    string
		model     string = openai.ModelChatGPT4
		slices    []string
	)
	flag.Usage = usage
	flag.StringVar(&envFile, "conf", "", "load environment variables from this file")
	flag.StringVar(&codeFile, "C", "", "file name to save combined code to (- for stdout, copy for clipboard)")
	flag.StringVar(&prompt, "p", "", "prompt to execute")
	flag.BoolVar(&replay, "replay", false, "replay response from response file (if any) instead of obtaining new one")
	flag.Var((*stringList)(&slices), "s", "specify a slice to use (can specify multiple times)")
	flag.Var((*stringList)(&rootDirs), "d", "add code directory (defaults to ., can specify multiple times)")
	flag.Var((*stringList)(&include), "i", "include only this glob pattern (can specify multiple times)")
	flag.Var((*stringList)(&exclude), "x", "exclude this glob pattern (can specify multiple times, in case of conflict with -i longest pattern wins)")
	flag.Var((*stringList)(&unexclude), "u", "un-exclude this glob pattern (can specify multiple times, in case of conflict always wins over -x/ignore)")
	flag.Var(&choiceFlag[string]{&model, openai.ModelChatGPT4}, "gpt4", "use GPT-4")
	flag.Var(&choiceFlag[string]{&model, openai.ModelChatGPT4With32k}, "gpt4-32k", "use GPT-4 32k")
	flag.Var(&choiceFlag[string]{&model, openai.ModelChatGPT35Turbo}, "gpt35", "use GPT 3.5")
	flag.Parse()

	if envFile != "" && envFile != "none" {
		loadEnv(envFile)
	}

	if codeFile == "" {
		codeFile = os.Getenv("AIDEV_SAVE_CODE")
	}
	var (
		respFile   string = os.Getenv("AIDEV_SAVE_RESP")
		promptFile string = os.Getenv("AIDEV_SAVE_PROMPT")
	)

	openAICreds = openai.Credentials{
		APIKey:         needEnv("OPENAI_API_KEY"),
		OrganizationID: os.Getenv("OPENAI_ORG"),
	}

	if len(rootDirs) == 0 {
		rootDirs = []string{"."}
	}

	ign := newIgnorer(&TreeConfig{
		Includes:   include,
		Excludes:   exclude,
		Unexcludes: unexclude,
	}, slices)

	items, ignored := loadFiles(rootDirs, ign.ShouldIgnore)
	if len(ignored) > 0 {
		log.Printf("%d ignored:", len(ignored))
		for _, path := range ignored {
			log.Printf("\t%s", path)
		}
	}
	log.Printf("%d files matched:", len(items))
	for _, item := range items {
		log.Printf("\t%s", item.relPath)
	}

	code := formatItems(items)
	if codeFile != "" {
		ensure(saveText(codeFile, code))
	}

	opt := openai.DefaultChatOptions()
	opt.Model = model
	opt.MaxTokens = 2048
	opt.Temperature = 0.7

	limit := openai.MaxTokens(opt.Model)
	log.Printf("Code tokens: %d, max for %s: %d.", openai.TokenCount(code, opt.Model), opt.Model, limit)

	if prompt == "" {
		fmt.Fprintf(os.Stderr, "Prompt? (end with EOF)\n")
		prompt = strings.TrimSpace(string(must(io.ReadAll(os.Stdin))))
		if prompt == "" {
			log.Printf("Empty prompt, nothing to do.")
			os.Exit(0)
		}
	}

	chat := []openai.Msg{
		openai.SystemMsg(systemPrompt),
		openai.UserMsg(fmt.Sprintf("%s\n\n%s", strings.TrimSpace(code), strings.TrimSpace(prompt))),
	}

	tokens := openai.ChatTokenCount(chat, opt.Model)
	log.Printf("Prompt tokens: %d, with completions: %d, max for %s: %d.", tokens, tokens+opt.MaxTokens, opt.Model, limit)
	if tokens+opt.MaxTokens > limit {
		log.Printf("WARNING: prompt exceeds %s capacity.", opt.Model)
		// os.Exit(1)
	}

	httpClient := &http.Client{
		Timeout: 5 * time.Minute,
	}

	var response string

	if replay && respFile != "" && respFile != "-" {
		response = string(mustSkippingOSNotExists(os.ReadFile(respFile)))
	}

	if response == "" {
		for {
			log.Printf("Talking to %s...", opt.Model)
			if promptFile != "" {
				ensure(saveText(promptFile, chat[0].Content+"\n\n=====\n\n"+chat[1].Content))
			}
			start := time.Now()
			msg, usage, err := openai.Chat(context.Background(), chat, openai.Options{
				Model:            model,
				MaxTokens:        2048,
				Temperature:      0.7,
				TopP:             0,
				N:                0,
				BestOf:           0,
				Stop:             []string{},
				PresencePenalty:  0,
				FrequencyPenalty: 0,
			}, httpClient, openAICreds)
			if err != nil {
				log.Printf("** ERROR: %v", err)
				var e *openai.Error
				if errors.As(err, &e) && retriable(e) {
					log.Println("Will retry in 5 seconds...")
					time.Sleep(5 * time.Second)
					continue
				}
				os.Exit(1)
			}
			elapsed := time.Since(start)

			response = msg[0].Content
			if respFile != "" {
				saveText(respFile, response)
			}

			cost := openai.Cost(usage.PromptTokens, usage.CompletionTokens, opt.Model)
			log.Printf("OpenAI %s time: %.0f sec, cost: %v (prompt = %d [vs estimated %d], completion = %d)", opt.Model, elapsed.Seconds(), cost, usage.PromptTokens, tokens, usage.CompletionTokens)
			break
		}
	}

	log.Printf("len(response) = %d", len(response))

	respItems, unfinished := parseItems(response)
	if unfinished {
		log.Printf("WARNING: output is not finished")
	}

	log.Printf("%d files updated:", len(respItems))
	for _, item := range respItems {
		log.Printf("\t%s", item.relPath)
	}

	for _, item := range respItems {
		fn := filepath.Join(rootDirs[0], item.relPath)
		// ext := filepath.Ext(fn)
		// fn = fn[:len(fn)-len(ext)] + ".draft" + ext
		fn = fn + ".draft"

		ensure(os.MkdirAll(filepath.Dir(fn), 0755))
		ensure(os.WriteFile(fn, item.content, 0644))
	}

	// dir := director.New()
	// defer dir.Wait()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()
	// gracefulshutdown.InterceptShutdownSignals(cancel)

	// app := setupApp(dataDir, appOpt)
	// defer app.Close()

}

func usage() {
	base := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s [options]\n\n", base)

	fmt.Printf("Options:\n")
	flag.PrintDefaults()

	fmt.Printf("\nMost options are set via environment variables. Run %s -print-env for a list.\n", base)
}

func retriable(err *openai.Error) bool {
	return err.IsNetwork || err.StatusCode == http.StatusTooManyRequests || (err.StatusCode >= 500 && err.StatusCode <= 599)
}
