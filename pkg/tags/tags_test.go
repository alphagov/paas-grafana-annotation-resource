package tags_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/alphagov/paas-grafana-annotation-resource/pkg/tags"
)

func TestTags(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tags Suite")
}

var _ = Describe("CombineTags", func() {
	It("ensures the tags are unique", func() {
		result := tags.CombineTags(
			[]string{"a", "b", "c"},
			[]string{"b", "c", "d"},
		)

		Expect(result).To(HaveLen(4))
		Expect(result).To(ConsistOf("a", "b", "c", "d"))
	})

	It("ensures the tags are sorted alphabetically", func() {
		result := tags.CombineTags(
			[]string{"g", "o", "v", "u", "k"},
			[]string{"p", "a", "a", "s"},
		)

		Expect(result).To(ConsistOf("a", "g", "k", "o", "p", "s", "u", "v"))
	})
})

var _ = Describe("FormatTags", func() {
	It("returns the empty string for zero tags", func() {
		result := tags.FormatTags([]string{})
		Expect(result).To(Equal(""))
	})

	It("returns the sorted tags, joined by commas, with spaces", func() {
		result := tags.FormatTags([]string{"no", "jenkins", "please"})
		Expect(result).To(Equal("jenkins, no, please"))
	})
})
