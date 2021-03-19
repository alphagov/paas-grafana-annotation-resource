package out_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"testing"
	"time"

	"github.com/grafana-tools/sdk"
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/out"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Out Suite")
}

func stringAddress(s string) *string {
	return &s
}

var _ = Describe("Out", func() {
	const (
		username = "admin"
		password = "password"
	)

	var (
		workingDirectory *string

		env map[string]string = map[string]string{
			"BUILD_ID":            "build-id",
			"BUILD_NAME":          "build-name",
			"BUILD_JOB_NAME":      "build-job-name",
			"BUILD_PIPELINE_NAME": "build-pipeline-name",
			"BUILD_TEAM_NAME":     "build-team-name",
			"ATC_EXTERNAL_URL":    "http://concourse-url",
		}
	)

	BeforeSuite(func() {
		httpmock.Activate()
	})

	AfterSuite(func() {
		httpmock.DeactivateAndReset()
	})

	BeforeEach(func() {
		httpmock.Reset()

		dir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		workingDirectory = &dir
	})

	AfterEach(func() {
		if workingDirectory != nil {
			os.RemoveAll(*workingDirectory)
		}

		workingDirectory = nil
	})

	Context("when validating the source", func() {
		It("returns an error when there is no 'url' defined in the source", func() {
			_, err := out.Out(
				types.OutRequest{Source: types.ResourceSource{}},
				map[string]string{},
				"",
			)

			Expect(err).To(MatchError(ContainSubstring(
				"'url' is required in the source definition",
			)))
		})

		It("returns an error when there is no 'username' or 'api_token' defined in the source", func() {
			_, err := out.Out(
				types.OutRequest{Source: types.ResourceSource{
					URL: "http://grafana",
				}},
				map[string]string{},
				"",
			)

			Expect(err).To(MatchError(ContainSubstring(
				"'username' or 'api_token' are required in the source definition",
			)))
		})

		It("returns an error when there is a 'username' but no 'password'", func() {
			_, err := out.Out(
				types.OutRequest{Source: types.ResourceSource{
					URL:      "http://grafana",
					Username: username,
				}},
				map[string]string{},
				"",
			)

			Expect(err).To(MatchError(ContainSubstring(
				"'password' is required in the source definition when 'username' is present",
			)))
		})
	})

	Context("when no id file exists", func() {
		var (
			req types.OutRequest
		)

		BeforeEach(func() {
			httpmock.RegisterResponder(
				"POST", fmt.Sprintf("http://%s:%s@grafana/api/annotations", username, password),
				func(req *http.Request) (*http.Response, error) {
					bodyBytes, err := ioutil.ReadAll(req.Body)
					Expect(err).NotTo(HaveOccurred())

					var requestBody sdk.CreateAnnotationRequest
					err = json.Unmarshal(bodyBytes, &requestBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(requestBody.Text).To(
						Equal("build-id env-var-source env-var-param"),
						"Text interpolation should work",
					)

					Expect(requestBody.Tags).To(
						ConsistOf("p1", "p2", "s1", "s2"),
						"Tags should be present and correct",
					)

					Expect(requestBody.Time).To(BeNumerically(
						"==", time.Now().Unix()*int64(1000), 1500,
					), "Time should approximately be now")

					responder := httpmock.NewStringResponder(
						200,
						`{ "message":"Annotation added", "id": 12345 }`,
					)

					return responder(req)
				},
			)

			req = types.OutRequest{
				Source: types.ResourceSource{
					URL: "http://grafana",

					Username: username,
					Password: password,

					Tags: []string{"s2", "s1"},
					Env: map[string]string{
						"ENV_VAR_SOURCE": "env-var-source",
					},
				},

				Params: types.ResourceParams{
					Template: stringAddress(
						"${BUILD_ID} ${ENV_VAR_SOURCE} ${ENV_VAR_PARAM}",
					),

					Tags: []string{"p2", "p1"},
					Env: map[string]string{
						"ENV_VAR_PARAM": "env-var-param",
					},
				},
			}
		})

		AfterEach(func() {
			callCount := httpmock.GetTotalCallCount()
			Expect(callCount).To(Equal(1), "Out should make an API call to Grafana")
		})

		It("should return the created id within the version", func() {
			resp, err := out.Out(req, env, *workingDirectory)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Version.ID).To(Equal("12345"))
		})

		It("should return metadata", func() {
			resp, err := out.Out(req, env, *workingDirectory)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Metadata).To(HaveLen(3))

			Expect(resp.Metadata).To(ContainElement(
				types.ResourceMetadataPair{Name: "id", Value: "12345"},
			))
			Expect(resp.Metadata).To(ContainElement(
				types.ResourceMetadataPair{Name: "tags", Value: "p1, p2, s1, s2"},
			))
			Expect(resp.Metadata).To(ContainElement(
				types.ResourceMetadataPair{
					Name:  "text",
					Value: "build-id env-var-source env-var-param",
				},
			))
		})

		It("should create an 'id' file containing the version", func() {
			_, err := out.Out(req, env, *workingDirectory)
			Expect(err).NotTo(HaveOccurred())

			idPath := path.Join(*workingDirectory, "id")

			_, err = os.Stat(idPath)
			Expect(err).NotTo(HaveOccurred(), "Out should create a file 'id'")

			idBytes, err := ioutil.ReadFile(idPath)
			Expect(err).NotTo(HaveOccurred(), "Out should create readable file 'id'")

			Expect(string(idBytes)).To(
				Equal("12345"), "Out should write the id to 'id'",
			)
		})
	})

	Context("when an id file exists", func() {
		const (
			annotationID = "12345"
		)

		var (
			req types.OutRequest
		)

		BeforeEach(func() {
			Expect(
				os.Mkdir(path.Join(*workingDirectory, "resource-name"), 0755),
			).NotTo(HaveOccurred(), "Could not 'create resource-name' directory")

			Expect(ioutil.WriteFile(
				path.Join(*workingDirectory, "resource-name", "id"),
				[]byte(annotationID),
				0644,
			)).NotTo(HaveOccurred(), "Could not write id file needed for test")

			httpmock.RegisterResponder(
				"PATCH", fmt.Sprintf("http://%s:%s@grafana/api/annotations/%s", username, password, annotationID),
				func(req *http.Request) (*http.Response, error) {
					bodyBytes, err := ioutil.ReadAll(req.Body)
					Expect(err).NotTo(HaveOccurred())

					var requestBody sdk.PatchAnnotationRequest
					err = json.Unmarshal(bodyBytes, &requestBody)
					Expect(err).NotTo(HaveOccurred())

					Expect(requestBody.Text).To(
						Equal("build-id http://concourse-url/teams/build-team-name/pipelines/build-pipeline-name/jobs/build-job-name/builds/build-name"),
						"Text interpolation should work",
					)

					Expect(requestBody.Tags).To(
						ConsistOf("p1", "p2"),
						"Tags should be present and correct",
					)

					Expect(requestBody.TimeEnd).To(BeNumerically(
						"==", time.Now().Unix()*int64(1000), 1500,
					), "TimeEnd should approximately be now")

					responder := httpmock.NewStringResponder(
						200,
						`{ "message":"Annotation patched"}`,
					)

					return responder(req)
				},
			)

			req = types.OutRequest{
				Source: types.ResourceSource{
					URL: "http://grafana",

					Username: username,
					Password: password,
				},

				Params: types.ResourceParams{
					Path: stringAddress("resource-name"),

					Tags: []string{"p2", "p1"},
					Env: map[string]string{
						"ENV_VAR_PARAM": "env-var-param",
					},
				},
			}
		})

		AfterEach(func() {
			callCount := httpmock.GetTotalCallCount()
			Expect(callCount).To(Equal(1), "Out should make an API call to Grafana")
		})

		It("should return the created id within the version", func() {
			resp, err := out.Out(req, env, *workingDirectory)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Version.ID).To(Equal("12345"))
		})

		It("should return metadata", func() {
			resp, err := out.Out(req, env, *workingDirectory)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp.Metadata).To(HaveLen(3))

			Expect(resp.Metadata).To(ContainElement(
				types.ResourceMetadataPair{Name: "id", Value: "12345"},
			))
			Expect(resp.Metadata).To(ContainElement(
				types.ResourceMetadataPair{Name: "tags", Value: "p1, p2"},
			))
			Expect(resp.Metadata).To(ContainElement(
				types.ResourceMetadataPair{
					Name:  "text",
					Value: "build-id http://concourse-url/teams/build-team-name/pipelines/build-pipeline-name/jobs/build-job-name/builds/build-name",
				},
			))
		})
	})
})
