package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	formatExports        = "exports"
	formatDotenv         = "dotenv"
	formatInputParameter = "parameter"
	formatInputJson      = "json"
)

func main() {
	if os.Getenv("AWS_ENV_PATH") == "" {
		log.Println("aws-env running locally, without AWS_ENV_PATH")
		return
	}

	recursivePtr := flag.Bool("recursive", false, "recursively process parameters on path")
	format := flag.String("format", formatExports, "output format")
	formatInput := flag.String("formatInput", formatInputParameter, "input format")
	flag.Parse()

	if *format == formatExports || *format == formatDotenv {
	} else {
		log.Fatal("Unsupported format option. Must be 'exports' or 'dotenv'")
	}

	sess := CreateSession()
	client := CreateClient(sess)

	ExportVariables(client, os.Getenv("AWS_ENV_PATH"), *recursivePtr, *format, *formatInput, "")
}

func ExportVariables(client *ssm.SSM, path string, recursive bool, format string, inputFormat string, nextToken string) {
	input := &ssm.GetParametersByPathInput{
		Path:           &path,
		WithDecryption: aws.Bool(true),
		Recursive:      aws.Bool(recursive),
	}

	if nextToken != "" {
		input.SetNextToken(nextToken)
	}

	output, err := client.GetParametersByPath(input)

	if err != nil {
		log.Panic(err)
	}

	for _, element := range output.Parameters {
		name := *element.Name
		value := *element.Value
		switch inputFormat {
		case formatInputJson:
			OutputParameterByJsonInputType(format, name, value)
		case formatInputParameter:
			OutputParameter(path, name, value, format)
		}
	}

	if output.NextToken != nil {
		ExportVariables(client, path, recursive, format, inputFormat, *output.NextToken)
	}
}

func OutputParameter(path string, name string, value string, format string) {
	env := strings.Replace(strings.Trim(name[len(path):], "/"), "/", "_", -1)
	value = strings.Replace(value, "\n", "\\n", -1)
	FormatOutput(format, env, value)
}

func OutputParameterByJsonInputType(outputFormat string, path string, value string) {
	var jsonMap map[string]interface{}
	err := json.Unmarshal([]byte(value), &jsonMap)
	if err != nil {
		return
	}

	for key, value := range jsonMap {
		FormatOutput(outputFormat, key, value)
	}
}

func FormatOutput(format string, env string, value interface{}) {
	switch format {
	case formatExports:
		FormatExport(env, value)
	case formatDotenv:
		FormatEnv(env, value)
	}
}

func FormatEnv(env string, value interface{}) (int, error) {
	return fmt.Printf("%s=\"%s\"\n", env, value)
}

func FormatExport(env string, value interface{}) (int, error) {
	return fmt.Printf("export %s='%s'\n", env, value)
}

func CreateSession() *session.Session {
	return session.Must(session.NewSession())
}

func CreateClient(sess *session.Session) *ssm.SSM {
	return ssm.New(sess)
}
