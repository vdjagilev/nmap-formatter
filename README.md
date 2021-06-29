# NMAP-Formatter

A tool that allows you to convert NMAP XML output to html/csv/json/markdown.

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

### Flags

* `--down-hosts` includes hosts that are down in the template (html, md)
* `-f, --file [filename]` outputs result to the file (by default output goes to STDOUT)
* `--help` display help message

## Installation

### Using Go

```
go install github.com/vjdagilev/nmap-formatter@latest
```

### Download Binary

Choose version from Release page and download it:

```
curl https://github.com/vdjagilev/nmap-formatter/releases/download/v0.0.2/nmap-formatter-linux-amd64.tar.gz --output nmap-formatter.tar.gz -L
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
