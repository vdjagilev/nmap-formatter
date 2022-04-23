# NMAP-Formatter

![build status](https://github.com/vdjagilev/nmap-formatter/actions/workflows/go.yml/badge.svg)
[![codecov](https://codecov.io/gh/vdjagilev/nmap-formatter/branch/main/graph/badge.svg?token=8WSYXRKMFA)](https://codecov.io/gh/vdjagilev/nmap-formatter)
[![Maintainability](https://api.codeclimate.com/v1/badges/7836d3a52439fb1affa0/maintainability)](https://codeclimate.com/github/vdjagilev/nmap-formatter/maintainability)

---

A tool that allows you to convert NMAP XML output to html/csv/json/markdown.

## Table of Contents

- [NMAP-Formatter](#nmap-formatter)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
    - [Flags](#flags)
      - [Output Related Flags](#output-related-flags)
	- [Installation](#installation)
		- [Using Go](#using-go)
    - [Docker](#docker)
    - [Download Binary](#download-binary)
    - [Compile](#compile)
  - [Example](#example)
  - [Use as a library](#use-as-a-library)

## Usage

```
nmap-formatter [path-to-nmap.xml] [html|csv|md|json] [flags]
```

Convert XML output to nicer HTML

```
nmap-formatter [path-to-nmap.xml] html > some-file.html
```

or Markdown

```
nmap-formatter [path-to-nmap.xml] md > some-markdown.md
```

or JSON

```
nmap-formatter [path-to-nmap.xml] json
```

it can be also combined with a `jq` tool, for example, list all the found ports and count them:

```
nmap-formatter [nmap.xml] json | jq -r '.Host[]?.Ports?.Port[]?.PortID' | sort | uniq -c
```

```
    1 "22"
    2 "80"
    1 "8080"
```

another example where only those hosts are selected, which have port where some http service is running:

```
nmap-formatter [nmap.xml] json | jq '.Host[]? | . as $host | .Ports?.Port[]? | select(.Service.Name== "http") | $host.HostAddress.Address' | uniq -c
```

```
    1 "192.168.1.1"
    1 "192.168.1.2"
    2 "192.168.1.3"
```

In this case `192.168.1.3` has 2 http services running (for example on ports 80 and 8080)`.

Another example where it is needed to display only filtered ports:

```
nmap-formatter [nmap.xml] json | jq '.Host[]?.Ports?.Port[]? | select(.State.State == "filtered") | .PortID'
```

Display host IP addresses that have filtered ports:

```
nmap-formatter [nmap.xml] json | jq '.Host[]? | . as $host | .Ports?.Port[]? | select(.State.State == "filtered") | .PortID | $host.HostAddress.Address'
```

### Flags

* `-f, --file [filename]` outputs result to the file (by default output goes to STDOUT)
* `--help` display help message
* `--version` display version (also can be used: `./nmap-formatter version`)

#### Output Related Flags

* `--skip-down-hosts` skips hosts that are down
  * Applicable in: `html`, `md`, `csv`
  * Default: `true`
* `--skip-summary` skips summary table
  * Applicable in: `html`, `md`
  * Default: `false`
* `--skip-traceroute` skips traceroute information
  * Applicable in: `html`
  * Default: `false`
* `--skip-metrics` skips metrics information
  * Applicable in: `html`
  * Default: `false`
* `--skip-port-scripts` skips port scripts information in ports table
  * Applicable in: `html`, `md`
  * Default: `false`
* `--json-pretty` pretty-prints JSON
  * Applicable in: `json`
  * Default: `true`

## Installation

### Using Go

```
go install github.com/vdjagilev/nmap-formatter@latest
```

### Docker

No installation needed, just run `docker run`:

```
docker run -v /path/to/xml/file.xml:/opt/file.xml ghcr.io/vdjagilev/nmap-formatter:latest /opt/file.xml json
```

### Download Binary

Choose version from Release page and download it:

```
curl https://github.com/vdjagilev/nmap-formatter/releases/download/v0.3.2/nmap-formatter-linux-amd64.tar.gz --output nmap-formatter.tar.gz -L
tar -xzvf nmap-formatter.tar.gz
./nmap-formatter --help
```

### Compile

```
git clone git@github.com:vdjagilev/nmap-formatter.git
cd nmap-formatter
go mod tidy
go build
# or 
go run . path/to/nmap.xml html
```

## Example

Example of HTML generated output from (https://nmap.org/book/output-formats-xml-output.html)

```
nmap-formatter basic-example.xml html
```

![Basic HTML Example](docs/images/basic-example-html.png)

## Use as a library

How to parse nmap results using golang

```go
package main

import (
	"encoding/xml"
	"os"

	"github.com/vdjagilev/nmap-formatter/formatter"
)

func main() {
	var nmap formatter.NMAPRun
	var config formatter.Config = formatter.Config{}

	// Read XML file that was produced by nmap (with -oX option)
	content, err := os.ReadFile("example.xml")
	if err != nil {
		panic(err)
	}
	// Unmarshal XML and map structure(s) fields accordingly
	if err = xml.Unmarshal(content, &nmap); err != nil {
		panic(err)
	}

	// Output data to console stdout
	// You can use any other io.Writer implementation
	// for example: os.OpenFile("file.json", os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.ModePerm)
	config.Writer = os.Stdout
	// Formatting data to JSON, you can use:
	// CSVOutput, MarkdownOutput, HTMLOutput as well
	config.OutputFormat = formatter.JSONOutput

	// Setting formatter data/options
	templateData := formatter.TemplateData{
		NMAPRun: nmap, // NMAP output data itself
		OutputOptions: formatter.OutputOptions{
			JSONPrettyPrint: true, // Additional option to prettify JSON
		},
	}

	// New formatter instance
	formatter := formatter.New(&config)
	if formatter == nil {
		// Not json/markdown/html/csv
		panic("wrong formatter provided")
	}

	// Attempt to format the data
	if err = formatter.Format(&templateData); err != nil {
		// html template could not be parsed or some other issue occured
		panic(err)
	}
}

```