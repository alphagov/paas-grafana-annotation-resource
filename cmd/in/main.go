package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/in"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func main() {
	inBytes, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	var req types.InRequest
	err = json.Unmarshal(inBytes, &req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	if len(os.Args) != 2 {
		fmt.Fprintf(
			os.Stderr,
			"Usage: in /path/to/working/directory ; %d args given\n",
			len(os.Args)-1,
		)
	}

	path := os.Args[1]

	resp, err := in.In(req, path)

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
