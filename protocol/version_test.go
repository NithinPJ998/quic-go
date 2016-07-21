package protocol

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Version", func() {
	It("converts tags to numbers", func() {
		Expect(VersionTagToNumber('Q' + '1'<<8 + '2'<<16 + '3'<<24)).To(Equal(VersionNumber(123)))
		Expect(VersionTagToNumber('Q' + '0'<<8 + '3'<<16 + '4'<<24)).To(Equal(Version34))
	})

	It("converts number to tag", func() {
		Expect(VersionNumberToTag(VersionNumber(123))).To(Equal(uint32('Q' + '1'<<8 + '2'<<16 + '3'<<24)))
		Expect(VersionNumberToTag(Version34)).To(Equal(uint32('Q' + '0'<<8 + '3'<<16 + '4'<<24)))
	})

	It("has proper tag list", func() {
		Expect(SupportedVersionsAsTags).To(Equal([]byte("Q031Q032Q033Q034")))
	})

	It("has proper version list", func() {
		Expect(SupportedVersionsAsString).To(Equal("34,33,32,31"))
	})

	It("recognizes supported versions", func() {
		Expect(IsSupportedVersion(0)).To(BeFalse())
		Expect(IsSupportedVersion(SupportedVersions[0])).To(BeTrue())
	})
})
