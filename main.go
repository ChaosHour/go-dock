package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"text/template"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

// Define the flags
var (
	readYAMLFile = flag.String("f", "", "Path to the YAML file")
	help         = flag.Bool("h", false, "Show this help message")
)

// Define the colors
var (
	red    = color.New(color.FgRed).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

// Define the YAML structure. This is a list of installed dependencies
type YAML struct {
	InstalledDependencies []string `yaml:"installed_dependencies"`
}

// Define the yaml Data
var yamlData YAML

// Define the Dockerfile template
const dockerfileTemplate = `FROM ubuntu:20.04

LABEL maintainer="Kurt Larsen <kurt_lv at cox dot net>"

ARG DEBIAN_FRONTEND=noninteractive
WORKDIR /opt
ENV LANG en_US.utf8

RUN TZ=America/Los_Angeles && \
    apt update && \
    apt install -y software-properties-common git curl p7zip-full wget whois locales python3 python3-pip upx psmisc iproute2 && \
    add-apt-repository -y ppa:longsleep/golang-backports && \
    apt update && \
    localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8

RUN {{ .Command }} && \
    rm -rf /var/lib/apt/lists/*
`

func init() {
	// Parse the flags
	flag.Parse()

	// Check if the YAML file is specified
	if *readYAMLFile == "" {
		fmt.Println(red("[!]"), "Please specify the path to the YAML file using the -f flag")
		os.Exit(1)
	}

	// Check if the YAML file exists
	if _, err := os.Stat(*readYAMLFile); os.IsNotExist(err) {
		fmt.Println(red("[!]"), "YAML file does not exist:", *readYAMLFile)
		os.Exit(1)
	}

	// Check if the YAML file is readable
	if _, err := os.Open(*readYAMLFile); err != nil {
		fmt.Println(red("[!]"), "Error reading YAML file:", err)
		os.Exit(1)
	}

	// Read the contents of the YAML file
	contents, err := ioutil.ReadFile(*readYAMLFile)
	if err != nil {
		fmt.Println(red("[!]"), "Error reading file:", err)
		os.Exit(1)
	}

	// Parse the YAML contents of the file
	err = yaml.Unmarshal(contents, &yamlData)
	if err != nil {
		fmt.Println(red("[!]"), "Error unmarshaling YAML:", err)
		os.Exit(1)
	}

	// Define Output using text/template - Kurt check this code.
	output := fmt.Sprintf("apt install -y %s ", join(yamlData.InstalledDependencies, " "))

	// Execute the Dockerfile template with the output of the main program
	var buf bytes.Buffer
	tmpl := template.Must(template.New("dockerfile").Parse(dockerfileTemplate))
	err = tmpl.Execute(&buf, struct{ Command string }{Command: output})
	if err != nil {
		fmt.Println("Error executing Dockerfile template:", err)
		os.Exit(1)
	}

	// Write the generated Dockerfile to disk
	err = ioutil.WriteFile("Dockerfile", buf.Bytes(), 0644)
	if err != nil {
		fmt.Println("Error writing Dockerfile:", err)
		os.Exit(1)
	}

	// Print a success message
	fmt.Println(green("[+]"), "Dockerfile generated successfully")
}

func main() {

	// Loop through the installed dependencies and send them to a buffer
	// print the buffer to stdout
	fmt.Println(yellow("[*]"), "Installed Dependencies:")
	for _, dep := range yamlData.InstalledDependencies {
		fmt.Println("  -", dep)
	}
}

// join concatenates a list of strings with a separator
func join(list []string, sep string) string {
	var buf bytes.Buffer
	for i, s := range list {
		if i > 0 {
			buf.WriteString(sep)
		}
		buf.WriteString(s)
	}
	return buf.String()
}
