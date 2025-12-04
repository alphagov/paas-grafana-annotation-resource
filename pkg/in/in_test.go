package in_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/in"
	"github.com/alphagov/paas-grafana-annotation-resource/pkg/types"
)

func TestCheck(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "In Suite")
}

var _ = Describe("In", func() {
	var (
		err               error
		resourceDirectory *string
	)

	BeforeEach(func() {
		dir, err := ioutil.TempDir("", "")
		Expect(err).NotTo(HaveOccurred())

		resourceDirectory = &dir
	})

	AfterEach(func() {
		if resourceDirectory != nil {
			os.RemoveAll(*resourceDirectory)
		}

		resourceDirectory = nil
	})

	It("should return the given version", func() {
		version := types.ResourceVersion{ID: "📊"}
		req := types.InRequest{Version: version}

		resp, err := in.In(req, *resourceDirectory)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.Version.ID).To(Equal("📊"))
	})

	It("should not return any metadata", func() {
		version := types.ResourceVersion{ID: "📊"}
		req := types.InRequest{Version: version}

		resp, err := in.In(req, *resourceDirectory)
		Expect(err).NotTo(HaveOccurred())

		Expect(resp.Metadata).To(HaveLen(0))
	})

	It("should create an 'id' file containing the version", func() {
		version := types.ResourceVersion{ID: "📊"}
		req := types.InRequest{Version: version}

		_, err = in.In(req, *resourceDirectory)
		Expect(err).NotTo(HaveOccurred())

		idPath := path.Join(*resourceDirectory, "id")

		_, err = os.Stat(idPath)
		Expect(err).NotTo(HaveOccurred(), "In should create a file 'id'")

		idBytes, err := ioutil.ReadFile(idPath)
		Expect(err).NotTo(HaveOccurred(), "In should create readable file 'id'")

		Expect(string(idBytes)).To(Equal("📊"), "In should write the id to 'id'")
	})
})
