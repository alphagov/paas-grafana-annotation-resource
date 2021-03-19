package check

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
	"github.com/grafana-tools/sdk"
)

func Check(req types.CheckRequest) (types.CheckResponse, error) {
	err := req.Source.Validate()

	if err != nil {
		return types.CheckResponse{}, err
	}

	params := []sdk.GetAnnotationsParams{
		sdk.WithAnnotationType(),
	}
	for _, tag := range req.Source.Tags {
		params = append(params, sdk.WithTag(tag))
	}

	var client *sdk.Client
	if req.Source.APIToken != "" {
		client = sdk.NewClient(req.Source.URL, req.Source.APIToken, http.DefaultClient)
	} else {
		client = sdk.NewClient(req.Source.URL, fmt.Sprintf("%s:%s", req.Source.Username, req.Source.Password), http.DefaultClient)
	}

	annotations, err := client.GetAnnotations(context.Background(), params...)
	if err != nil || len(annotations) == 0 {
		return types.CheckResponse{}, err
	}

	sort.SliceStable(annotations, func(i, j int) bool {
		return annotations[i].Time < annotations[j].Time
	})

	if req.Version == nil {
		return []types.ResourceVersion{
			{ID: fmt.Sprint(annotations[len(annotations)-1].ID)},
		}, nil
	}

	var existing *sdk.AnnotationResponse
	for _, annotation := range annotations {
		if fmt.Sprint(annotation.ID) == req.Version.ID {
			existing = &annotation
			break
		}
	}

	var versions []types.ResourceVersion
	for _, annotation := range annotations {
		if annotation.Time >= existing.Time {
			versions = append(versions, types.ResourceVersion{ID: fmt.Sprint(annotation.ID)})
		}
	}

	return versions, nil
}
