package types

type ResourceVersion struct {
	ID string `json:"id"`
}

type ResourceSource struct {
	URL string `json:"url"`

	Username string `json:"username"`
	Password string `json:"password"`

	APIToken string `json:"api_token"`

	Tags []string          `json:"tags"`
	Env  map[string]string `json:"env"`
}

type ResourceParams struct {
	Path *string `json:"path"`

	Template *string `json:"template"`

	Tags []string          `json:"tags"`
	Env  map[string]string `json:"env"`
}

type ResourceMetadataPair struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CheckRequest struct {
	Source  ResourceSource   `json:"source"`
	Version *ResourceVersion `json:"version"`
}

type CheckResponse = []ResourceVersion

type InRequest struct {
	Source  ResourceSource  `json:"source"`
	Version ResourceVersion `json:"version"`
}

type OutRequest struct {
	Source ResourceSource `json:"source"`
	Params ResourceParams `json:"params"`
}

type InOutResponse struct {
	Version  ResourceVersion        `json:"version"`
	Metadata []ResourceMetadataPair `json:"metadata"`
}

type GrafanaCreateAnnotationRequest struct {
	Time int64    `json:"time"`
	Tags []string `json:"tags"`
	Text string   `json:"text"`
}

type GrafanaUpdateAnnotationRequest struct {
	TimeEnd int64    `json:"timeEnd"`
	Tags    []string `json:"tags"`
	Text    string   `json:"text"`
}

type GrafanaCreateAnnotationResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
}
