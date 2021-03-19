package check_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/check"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Out Suite")
}

var _ = Describe("Check", func() {
	const (
		username = "admin"
		password = "password"
	)

	Context("when validating the source", func() {
		It("returns an error when there is no 'url' defined in the source", func() {
			_, err := check.Check(
				types.CheckRequest{Source: types.ResourceSource{}},
			)

			Expect(err).To(MatchError(ContainSubstring(
				"'url' is required in the source definition",
			)))
		})

		It("returns an error when there is no 'username' or 'api_token' defined in the source", func() {
			_, err := check.Check(
				types.CheckRequest{Source: types.ResourceSource{
					URL: "http://grafana",
				}},
			)

			Expect(err).To(MatchError(ContainSubstring(
				"'username' or 'api_token' are required in the source definition",
			)))
		})

		It("returns an error when there is a 'username' but no 'password'", func() {
			_, err := check.Check(
				types.CheckRequest{Source: types.ResourceSource{
					URL:      "http://grafana",
					Username: username,
				}},
			)

			Expect(err).To(MatchError(ContainSubstring(
				"'password' is required in the source definition when 'username' is present",
			)))
		})
	})

	BeforeSuite(func() {
		httpmock.Activate()
	})

	AfterSuite(func() {
		httpmock.DeactivateAndReset()
	})

	BeforeEach(func() {
		httpmock.Reset()

		httpmock.RegisterResponder(
			"GET", fmt.Sprintf("http://%s:%s@grafana/api/annotations", username, password),
			func(req *http.Request) (*http.Response, error) {
				Expect(req.URL.Query()["tags"]).To(
					ConsistOf("s1", "s2"),
					"Tags should be present and correct",
				)

				responder := httpmock.NewStringResponder(
					200,
					`[
						{"id": 123, "time": 1000},
						{"id": 456, "time": 2000},
						{"id": 789, "time": 3000}
					]`,
				)

				return responder(req)
			},
		)
	})

	Context("when no version is provided", func() {
		var (
			req types.CheckRequest
		)

		BeforeEach(func() {
			req = types.CheckRequest{
				Source: types.ResourceSource{
					URL: "http://grafana",

					Username: username,
					Password: password,

					Tags: []string{"s2", "s1"},
					Env: map[string]string{
						"ENV_VAR_SOURCE": "env-var-source",
					},
				},
				Version: nil,
			}
		})

		AfterEach(func() {
			callCount := httpmock.GetTotalCallCount()
			Expect(callCount).To(Equal(1), "Out should make an API call to Grafana")
		})

		It("should return the most recent annotation id", func() {
			resp, err := check.Check(req)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp).To(ConsistOf(
				types.ResourceVersion{ID: "789"},
			))
		})
	})

	Context("when a version is provided", func() {
		var (
			req types.CheckRequest
		)

		BeforeEach(func() {
			req = types.CheckRequest{
				Source: types.ResourceSource{
					URL: "http://grafana",

					Username: username,
					Password: password,

					Tags: []string{"s2", "s1"},
					Env: map[string]string{
						"ENV_VAR_SOURCE": "env-var-source",
					},
				},
				Version: &types.ResourceVersion{
					ID: "456",
				},
			}
		})

		AfterEach(func() {
			callCount := httpmock.GetTotalCallCount()
			Expect(callCount).To(Equal(1), "Out should make an API call to Grafana")
		})

		It("should return all annotations after and including the version provided", func() {
			resp, err := check.Check(req)
			Expect(err).NotTo(HaveOccurred())

			Expect(resp).To(Equal([]types.ResourceVersion{
				{ID: "456"},
				{ID: "789"},
			}))
		})
	})
})
