package out

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/tags"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

const defaultTemplate = "${BUILD_ID} ${ATC_EXTERNAL_URL}/teams/${BUILD_TEAM_NAME}/pipelines/${BUILD_PIPELINE_NAME}/jobs/${BUILD_JOB_NAME}/builds/${BUILD_NAME}"

func pathToIDFile(dirPath string) string {
	return path.Join(dirPath, "id")
}

func Out(
	req types.OutRequest,
	env map[string]string,
	dirPath string) (types.InOutResponse, error) {

	err := req.Source.Validate()

	if err != nil {
		return types.InOutResponse{}, err
	}

	if req.Params.Path == nil {
		// when we are not given a path, we are creating a new annotation
		return outCreate(req, env, dirPath)
	}

	return outUpdate(req, env, dirPath)
}

func addGrafanaAPIHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
}

func addGrafanaAuth(req *http.Request, username string, password string) {
	req.SetBasicAuth(username, password)
}

func outCreate(
	req types.OutRequest,
	env map[string]string,
	dirPath string,
) (types.InOutResponse, error) {

	combinedTags := tags.CombineTags(req.Source.Tags, req.Params.Tags)

	actualTemplate := defaultTemplate
	if req.Source.Template != nil {
		actualTemplate = *req.Source.Template
	}
	if req.Params.Template != nil {
		actualTemplate = *req.Params.Template
	}

	text := os.Expand(actualTemplate, func(varName string) string {
		if val, ok := req.Params.Env[varName]; ok {
			return val
		}

		if val, ok := req.Source.Env[varName]; ok {
			return val
		}

		if val, ok := env[varName]; ok {
			return val
		}

		return "nil"
	})

	currentTime := time.Now().Unix() * int64(1000)

	requestBody := types.GrafanaCreateAnnotationRequest{
		Time:    currentTime,
		TimeEnd: currentTime + 1,
		Tags:    combinedTags,
		Text:    text,
	}

	requestBodyBytes, err := json.Marshal(requestBody)

	if err != nil {
		return types.InOutResponse{}, err
	}

	url := fmt.Sprintf("%s/api/annotations", req.Source.URL)

	httpReq, err := http.NewRequest(
		"POST", url, bytes.NewReader(requestBodyBytes),
	)

	if err != nil {
		return types.InOutResponse{}, err
	}

	addGrafanaAPIHeaders(httpReq)
	addGrafanaAuth(httpReq, req.Source.Username, req.Source.Password)

	httpResp, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return types.InOutResponse{}, err
	}

	respBodyBytes, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return types.InOutResponse{}, err
	}

	if httpResp.StatusCode != 200 {
		return types.InOutResponse{}, fmt.Errorf(
			fmt.Sprintf(
				"Expected 200 from Grafana API, %d received: %s",
				httpResp.StatusCode, string(respBodyBytes),
			),
		)
	}

	var parsedResponse types.GrafanaCreateAnnotationResponse
	err = json.Unmarshal(respBodyBytes, &parsedResponse)

	if err != nil {
		return types.InOutResponse{}, err
	}

	annotationID := fmt.Sprintf("%d", parsedResponse.ID)

	err = ioutil.WriteFile(pathToIDFile(dirPath), []byte(annotationID), 0644)

	if err != nil {
		return types.InOutResponse{}, err
	}

	return types.InOutResponse{
		Version: types.ResourceVersion{
			ID: annotationID,
		},
		Metadata: []types.ResourceMetadataPair{
			{Name: "id", Value: annotationID},
			{Name: "tags", Value: tags.FormatTags(combinedTags)},
			{Name: "text", Value: text},
		},
	}, nil
}

func outUpdate(
	req types.OutRequest,
	env map[string]string,
	dirPath string,
) (types.InOutResponse, error) {

	idBytes, err := ioutil.ReadFile(pathToIDFile(dirPath))

	if err != nil {
		return types.InOutResponse{}, err
	}

	requestBody := types.GrafanaCreateAnnotationRequest{
		TimeEnd: time.Now().Unix() * int64(1000),
	}

	requestBodyBytes, err := json.Marshal(requestBody)

	if err != nil {
		return types.InOutResponse{}, err
	}

	annotationID := string(idBytes)
	url := fmt.Sprintf("%s/api/annotations/%s", req.Source.URL, annotationID)

	httpReq, err := http.NewRequest(
		"PATCH", url, bytes.NewReader(requestBodyBytes),
	)

	if err != nil {
		return types.InOutResponse{}, err
	}

	addGrafanaAPIHeaders(httpReq)
	addGrafanaAuth(httpReq, req.Source.Username, req.Source.Password)

	httpResp, err := http.DefaultClient.Do(httpReq)

	if err != nil {
		return types.InOutResponse{}, err
	}

	respBodyBytes, err := ioutil.ReadAll(httpResp.Body)

	if err != nil {
		return types.InOutResponse{}, err
	}

	if httpResp.StatusCode != 200 {
		return types.InOutResponse{}, fmt.Errorf(
			fmt.Sprintf(
				"Expected 200 from Grafana API, %d received: %s",
				httpResp.StatusCode, string(respBodyBytes),
			),
		)
	}

	return types.InOutResponse{
		Version: types.ResourceVersion{
			ID: annotationID,
		},
		Metadata: []types.ResourceMetadataPair{
			{Name: "id", Value: annotationID},
		},
	}, nil
}
