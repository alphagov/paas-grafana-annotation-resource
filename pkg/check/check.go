package check

import (
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func Check(req types.CheckRequest) (types.CheckResponse, error) {
	if req.Version == nil {
		return types.CheckResponse{}, nil
	}

	return types.CheckResponse{*req.Version}, nil
}
