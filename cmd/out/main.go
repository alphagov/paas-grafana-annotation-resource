package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/out"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func main() {
	inBytes, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	var req types.OutRequest
	err = json.Unmarshal(inBytes, &req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(
			os.Stderr,
			"Usage: out /path/to/working/directory ; %d args given\n",
			len(os.Args)-1,
		)
	}

	path := os.Args[1]
	env := make(map[string]string, 0)

	for _, pair := range os.Environ() {
		split := strings.SplitN(pair, "=", 2)
		if len(split) != 2 {
			env[split[0]] = "true"
		} else {
			env[split[0]] = split[1]
		}
	}

	resp, err := out.Out(req, env, path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	outBytes, err := json.Marshal(resp)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, string(outBytes))
	os.Exit(0)
}
