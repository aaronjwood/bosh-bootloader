package compute_test

import (
	"errors"

	"github.com/genevieve/leftovers/gcp/compute"
	"github.com/genevieve/leftovers/gcp/compute/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Image", func() {
	var (
		client *fakes.ImagesClient
		name   string

		image compute.Image
	)

	BeforeEach(func() {
		client = &fakes.ImagesClient{}
		name = "banana"

		image = compute.NewImage(client, name)
	})

	Describe("Delete", func() {
		It("deletes the image", func() {
			err := image.Delete()
			Expect(err).NotTo(HaveOccurred())

			Expect(client.DeleteImageCall.CallCount).To(Equal(1))
			Expect(client.DeleteImageCall.Receives.Image).To(Equal(name))
		})

		Context("when the client fails to delete", func() {
			BeforeEach(func() {
				client.DeleteImageCall.Returns.Error = errors.New("the-error")
			})

			It("returns the error", func() {
				err := image.Delete()
				Expect(err).To(MatchError("ERROR deleting image banana: the-error"))
			})
		})
	})

	Describe("Name", func() {
		It("returns the name", func() {
			Expect(image.Name()).To(Equal(name))
		})
	})
})
