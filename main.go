package main

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/goccy/go-yaml"
)

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("not enough parameters\nusage: %s [pattern] [file]", path.Base(os.Args[0]))
	}

	raw_path := os.Args[1]
	path, err := yaml.PathString("$." + raw_path)
	if err != nil {
		return err
	}

	var data []byte
	var input string

	if len(os.Args) > 2 {
		input = os.Args[2]
		data, err = os.ReadFile(os.Args[2])
		if err != nil {
			return err
		}
	} else {
		data, err = io.ReadAll(os.Stdin)
		input = "-"
		if err != nil {
			return err
		}
	}

	node, err := path.ReadNode(strings.NewReader(string(data)))
	if err != nil {
		return err
	}
	token := node.GetToken()

	var str []byte
	if node.Type().String() == "Sequence" || node.Type().String() == "Mapping" {
		str, err = yaml.YAMLToJSON([]byte(node.String()))
		if err != nil {
			return err
		}
	} else {
		str = []byte(node.String())
	}

	fmt.Printf("%s:%d:%d:%s: %s\n", input, token.Position.Line, token.Position.Column, raw_path, str)

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
