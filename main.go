package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"gopkg.in/yaml.v3"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
}

func run() error {
	var (
		dataFile   string
		dataFormat string
	)

	flag.StringVar(&dataFile, "d", dataFile, "Path to data file which will be passed into the template. Data format is determined from the file extension (.json, .yaml) or can be specified using the -f flag.")
	flag.StringVar(&dataFormat, "f", dataFormat, "Specify data format (json, yaml).")

	flag.Parse()

	if flag.NArg() != 1 {
		return fmt.Errorf("usage: %s [-d data.json] [-f (json|yaml)] TEMPLATE-FILE", os.Args[0])
	}

	templateFile := flag.Arg(0)

	var (
		data map[string]any
		err  error
	)
	if dataFile != "" {
		data, err = loadData(dataFile, dataFormat)
		if err != nil {
			return fmt.Errorf("error loading data: %w", err)
		}
	}

	rawTemplate, err := os.ReadFile(templateFile)
	if err != nil {
		return fmt.Errorf("error reading template file: %w", err)
	}

	tmpl, err := template.New("").Funcs(sprig.FuncMap()).Parse(string(rawTemplate))
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}
	return nil
}

func loadData(filename string, dataFormat string) (map[string]any, error) {
	if dataFormat == "" {
		dataFormat = filepath.Ext(filename)
	}

	var unmarshalFn func([]byte, any) error
	switch dataFormat {
	case "json":
		unmarshalFn = json.Unmarshal
	case "yaml":
		unmarshalFn = yaml.Unmarshal
	case "":
		return nil, fmt.Errorf("could not determine data formta form file extension. use -f to specify it")
	default:
		return nil, fmt.Errorf("unknown data format '%s'. allowed formats are json and yaml", dataFormat)
	}

	rawData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	data := map[string]any{}
	err = unmarshalFn(rawData, &data)
	if err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}
	return data, nil
}
