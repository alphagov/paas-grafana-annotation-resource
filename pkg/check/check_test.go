package check_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/check"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Check Suite")
}

var _ = Describe("Check", func() {
	Context("with no versions", func() {
		It("should return no versions", func() {
			req := types.CheckRequest{Version: nil}

			resp, err := check.Check(req)

			Expect(err).NotTo(HaveOccurred())

			Expect(resp).To(HaveLen(0))
		})
	})

	Context("with a version", func() {
		It("should return that version", func() {
			version := types.ResourceVersion{ID: "📊"}

			req := types.CheckRequest{Version: &version}

			resp, err := check.Check(req)

			Expect(err).NotTo(HaveOccurred())

			Expect(resp).To(HaveLen(1))

			Expect(resp[0].ID).To(Equal("📊"))
		})
	})
})
