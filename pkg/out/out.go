package out

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/tags"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
	"github.com/grafana-tools/sdk"
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

	var client *sdk.Client
	if req.Source.APIToken != "" {
		client = sdk.NewClient(req.Source.URL, req.Source.APIToken, http.DefaultClient)
	} else {
		client = sdk.NewClient(req.Source.URL, fmt.Sprintf("%s:%s", req.Source.Username, req.Source.Password), http.DefaultClient)
	}

	if req.Params.Path == nil {
		// when we are not given a path, we are creating a new annotation
		return outCreate(client, req, env, dirPath)
	}

	return outUpdate(client, req, env, dirPath)
}

func outCreate(
	client *sdk.Client,
	req types.OutRequest,
	env map[string]string,
	dirPath string,
) (types.InOutResponse, error) {
	combinedTags := tags.CombineTags(req.Source.Tags, req.Params.Tags)
	text := getText(req, env)
	currentTime := time.Now().Unix() * int64(1000)
	status, err := client.CreateAnnotation(context.Background(), sdk.CreateAnnotationRequest{
		Time: currentTime,
		Tags: combinedTags,
		Text: text,
	})

	if err != nil {
		return types.InOutResponse{}, err
	}

	annotationID := fmt.Sprintf("%d", *status.ID)
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
	client *sdk.Client,
	req types.OutRequest,
	env map[string]string,
	dirPath string,
) (types.InOutResponse, error) {
	idBytes, err := ioutil.ReadFile(
		pathToIDFile(path.Join(dirPath, *req.Params.Path)),
	)
	if err != nil {
		return types.InOutResponse{}, err
	}

	annotationID, err := strconv.ParseUint(string(idBytes), 10, 64)
	if err != nil {
		return types.InOutResponse{}, err
	}

	combinedTags := tags.CombineTags(req.Source.Tags, req.Params.Tags)
	text := getText(req, env)
	_, err = client.PatchAnnotation(context.Background(), uint(annotationID), sdk.PatchAnnotationRequest{
		TimeEnd: time.Now().Unix() * int64(1000),
		Tags:    tags.CombineTags(req.Source.Tags, req.Params.Tags),
		Text:    getText(req, env),
	})

	if err != nil {
		return types.InOutResponse{}, err
	}

	return types.InOutResponse{
		Version: types.ResourceVersion{
			ID: string(idBytes),
		},
		Metadata: []types.ResourceMetadataPair{
			{Name: "id", Value: string(idBytes)},
			{Name: "tags", Value: tags.FormatTags(combinedTags)},
			{Name: "text", Value: text},
		},
	}, nil
}

func getText(req types.OutRequest, env map[string]string) string {
	actualTemplate := defaultTemplate
	if req.Params.Template != nil {
		actualTemplate = *req.Params.Template
	}
	return os.Expand(actualTemplate, func(varName string) string {
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
}
