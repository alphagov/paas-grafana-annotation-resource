package in

import (
	"io/ioutil"
	"path"

	"github.com/govuk-pay/paas-grafana-annotation-resource/pkg/types"
)

func In(req types.InRequest, dirPath string) (types.InOutResponse, error) {
	err := ioutil.WriteFile(
		path.Join(dirPath, "id"),
		[]byte(req.Version.ID),
		0644,
	)

	if err != nil {
		return types.InOutResponse{}, err
	}

	return types.InOutResponse{Version: req.Version}, nil
}
