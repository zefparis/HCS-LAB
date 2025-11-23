package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/corehuman/hcs-lab-api/internal/hcs"
)

const version = "1.0.0-hcs-lab"

func main() {
	// Define command line flags
	var (
		u3Only   = flag.Bool("u3-only", false, "Only compute and output U3 code")
		u4Only   = flag.Bool("u4-only", false, "Only compute and output U4 code")
		pretty   = flag.Bool("pretty", false, "Pretty print JSON output")
		rawJSON  = flag.Bool("raw-json", false, "Print only JSON to stdout (no extra text)")
		showHelp = flag.Bool("help", false, "Show help information")
		showVer  = flag.Bool("version", false, "Show version information")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] input.json\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Generate HCS codes from an input profile\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExample:\n")
		fmt.Fprintf(os.Stderr, "  %s input.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --u3-only profile.json\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --pretty --raw-json input.json\n", os.Args[0])
	}

	flag.Parse()

	// Handle help and version flags
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVer {
		fmt.Printf("hcsgen version %s\n", version)
		os.Exit(0)
	}

	// Check for input file argument
	args := flag.Args()
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "Error: input file required\n")
		flag.Usage()
		os.Exit(1)
	}

	inputFile := args[0]

	// Read input file
	inputData, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(1)
	}

	// Parse input JSON
	var input hcs.InputProfile
	if err := json.Unmarshal(inputData, &input); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing input JSON: %v\n", err)
		os.Exit(1)
	}

	// Create generator
	generator, err := hcs.NewGenerator()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing generator: %v\n", err)
		os.Exit(1)
	}

	// Set generation options
	opts := &hcs.GeneratorOptions{
		U3Only: *u3Only,
		U4Only: *u4Only,
	}

	// Generate HCS codes
	output, err := generator.GenerateWithOptions(&input, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating HCS codes: %v\n", err)
		os.Exit(1)
	}

	// Prepare output file names
	baseName := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	outputJSONFile := baseName + "_output.json"
	outputHCSFile := baseName + "_output.hcs"

	// Write JSON output
	var jsonData []byte
	if *pretty {
		jsonData, err = json.MarshalIndent(output, "", "  ")
	} else {
		jsonData, err = json.Marshal(output)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling output: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(outputJSONFile, jsonData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output.json: %v\n", err)
		os.Exit(1)
	}

	// Write HCS file (codes only)
	var hcsContent []string
	if output.CodeU3 != "" {
		hcsContent = append(hcsContent, output.CodeU3)
	}
	if output.CodeU4 != "" {
		hcsContent = append(hcsContent, output.CodeU4)
	}
	if output.CodeU5 != "" {
		hcsContent = append(hcsContent, output.CodeU5)
	}
	hcsData := []byte(strings.Join(hcsContent, "\n"))
	if err := os.WriteFile(outputHCSFile, hcsData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output.hcs: %v\n", err)
		os.Exit(1)
	}

	// Output to stdout
	if *rawJSON {
		// Print only JSON
		fmt.Print(string(jsonData))
	} else {
		// Print the HCS codes with labels
		if output.CodeU3 != "" {
			fmt.Printf("HCS-U3: %s\n", output.CodeU3)
		}
		if output.CodeU4 != "" {
			fmt.Printf("HCS-U4: %s\n", output.CodeU4)
		}
		if output.CodeU5 != "" {
			fmt.Printf("HCS-U5: %s\n", output.CodeU5)
			if output.ChineseProfile != nil {
				fmt.Printf("\nChinese BaZi Profile detected:\n")
				fmt.Printf("  Four Pillars: %s | %s | %s | %s\n",
					output.ChineseProfile.YearPillar,
					output.ChineseProfile.MonthPillar,
					output.ChineseProfile.DayPillar,
					output.ChineseProfile.HourPillar)
				fmt.Printf("  Day Master: %s (Strength: %.0f%%)\n",
					output.ChineseProfile.DayMaster,
					output.ChineseProfile.DayMasterStrength*100)
				fmt.Printf("  Yin/Yang Balance: %.0f%% Yang\n",
					output.ChineseProfile.YinYangBalance*100)
			}
		}
		fmt.Printf("\nCHIP: %s\n", output.Chip)
		fmt.Printf("\nOutput written to:\n")
		fmt.Printf("  - %s (full JSON)\n", outputJSONFile)
		fmt.Printf("  - %s (codes only)\n", outputHCSFile)
	}
}
