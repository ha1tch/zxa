package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ha1tch/zxa/internal/zxa_assembler"
)

var version = "v0.1.0-dev"

type Config struct {
	inputFile    string
	outputFile   string
	includePaths []string
	hexOutput    bool
	jsonOutput   bool
	verbose      bool
	z80next      bool
	quiet        bool
}

func printUsage() {
	fmt.Fprintf(os.Stderr, `ZXA - Z80 Cross Assembler %s

Usage: zxa [options] <input.asm>

Options:
`, version)
	flag.PrintDefaults()
}

func parseFlags() (*Config, error) {
	cfg := &Config{}

	// Define flags
	outFile := flag.String("o", "", "output file name (default: input base name)")
	includePath := flag.String("I", "", "include search path (can be specified multiple times)")
	flag.BoolVar(&cfg.hexOutput, "hex", false, "generate hex dump output")
	flag.BoolVar(&cfg.jsonOutput, "json", false, "generate JSON assembly report")
	flag.BoolVar(&cfg.verbose, "v", false, "enable verbose output")
	flag.BoolVar(&cfg.z80next, "next", false, "enable Z80N (ZX Spectrum Next) instructions")
	flag.BoolVar(&cfg.quiet, "q", false, "quiet mode (suppress non-error output)")
	showVersion := flag.Bool("version", false, "show version information")

	// Custom usage message
	flag.Usage = printUsage

	flag.Parse()

	// Show version if requested
	if *showVersion {
		fmt.Printf("ZXA version %s\n", version)
		os.Exit(0)
	}

	// Get input file
	if flag.NArg() != 1 {
		return nil, fmt.Errorf("exactly one input file must be specified")
	}
	cfg.inputFile = flag.Arg(0)

	// Set output file
	if *outFile != "" {
		cfg.outputFile = *outFile
	} else {
		base := strings.TrimSuffix(cfg.inputFile, filepath.Ext(cfg.inputFile))
		cfg.outputFile = base
	}

	// Process include paths
	if *includePath != "" {
		cfg.includePaths = strings.Split(*includePath, string(os.PathListSeparator))
	}

	// Add current directory to include paths if not present
	hasCurrentDir := false
	for _, path := range cfg.includePaths {
		if path == "." {
			hasCurrentDir = true
			break
		}
	}
	if !hasCurrentDir {
		cfg.includePaths = append([]string{"."}, cfg.includePaths...)
	}

	// Verbose and quiet are mutually exclusive
	if cfg.verbose && cfg.quiet {
		return nil, fmt.Errorf("cannot specify both verbose (-v) and quiet (-q) modes")
	}

	return cfg, nil
}

func main() {
	startTime := time.Now()

	// Parse command line flags
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(1)
	}

	// Create assembler options
	opts := zxa_assembler.AssemblerOptions{
		Variant: zxa_assembler.Z80Standard,
	}
	if cfg.z80next {
		opts.Variant = zxa_assembler.Z80Next
	}

	// Create assembler instance
	asm := zxa_assembler.NewAssembler(opts)

	// Configure assembler
	asm.SetHexOutput(cfg.hexOutput)
	asm.SetJSONOutput(cfg.jsonOutput)
	for _, path := range cfg.includePaths {
		asm.AddIncludePath(path)
	}

	// Print configuration if verbose
	if cfg.verbose {
		fmt.Printf("ZXA version %s\n", version)
		fmt.Printf("Input file: %s\n", cfg.inputFile)
		fmt.Printf("Output base: %s\n", cfg.outputFile)
		if len(cfg.includePaths) > 0 {
			fmt.Printf("Include paths:\n")
			for _, path := range cfg.includePaths {
				fmt.Printf("  %s\n", path)
			}
		}
		fmt.Printf("Output formats: binary")
		if cfg.hexOutput {
			fmt.Printf(", hex")
		}
		if cfg.jsonOutput {
			fmt.Printf(", json")
		}
		fmt.Printf("\n")
		if cfg.z80next {
			fmt.Printf("Z80N instructions enabled\n")
		}
		fmt.Printf("\n")
	}

	// Perform assembly
	result, err := asm.Assemble(cfg.inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Assembly failed: %v\n", err)
		os.Exit(1)
	}

	// Write outputs
	if err := result.WriteFiles(cfg.outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
		os.Exit(1)
	}

	// Print statistics unless quiet mode
	if !cfg.quiet {
		if cfg.verbose {
			fmt.Printf("\nAssembly statistics:\n")
			fmt.Printf("  Bytes generated: %d\n", result.Statistics.BytesGenerated)
			fmt.Printf("  Lines processed: %d\n", result.Statistics.LinesProcessed)
			fmt.Printf("  Symbols defined: %d\n", result.Statistics.SymbolsDefined)
			fmt.Printf("  Time taken: %v\n", time.Since(startTime))
		} else {
			fmt.Printf("Assembled %s: %d bytes\n", 
				filepath.Base(cfg.inputFile), 
				result.Statistics.BytesGenerated)
		}
	}

	// Exit with error if assembly had errors
	if !result.Success {
		os.Exit(1)
	}
}