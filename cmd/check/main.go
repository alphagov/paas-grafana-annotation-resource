package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/check"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func main() {
	inBytes, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	var req types.CheckRequest
	err = json.Unmarshal(inBytes, &req)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}

	resp, err := check.Check(req)

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
