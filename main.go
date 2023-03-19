package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// You are an AI programming assistant. User will send all files from a Git repository, separated by =#=#= headers, followed by a change request. Implement the requested change and output the modified files in the same format.

// var (
// 	openAICreds openai.Credentials
// )

func main() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Ltime)

	//env.Var("OPENAI_API_KEY", required, envloader.StringVar(&openAICreds.APIKey), "OpenAI API key")

	var (
		envFile string
		outFile string
	)
	flag.Usage = usage
	flag.StringVar(&envFile, "conf", "", "load environment variables from this file")
	flag.StringVar(&outFile, "o", "", "file name to save results to instead of printing to stdout")
	flag.Parse()

	if envFile != "" && envFile != "none" {
		loadEnv(envFile)
	}

	ign := newIgnorer()

	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}

	var buf strings.Builder
	for _, rootDir := range args {
		loadFiles(&buf, rootDir, ign.ShouldIgnore)
	}
	buf.WriteString("=#=#= END\n\n")

	output := buf.String()
	if outFile != "" {
		ensure(os.WriteFile(outFile, []byte(output), 0644))
	} else {
		fmt.Println(output)
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
