package compute_test

import (
	"errors"

	"github.com/genevievelesperance/leftovers/gcp/compute"
	"github.com/genevievelesperance/leftovers/gcp/compute/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	gcpcompute "google.golang.org/api/compute/v1"
)

var _ = Describe("Subnetworks", func() {
	var (
		client  *fakes.SubnetworksClient
		logger  *fakes.Logger
		regions map[string]string

		subnetworks compute.Subnetworks
	)

	BeforeEach(func() {
		client = &fakes.SubnetworksClient{}
		logger = &fakes.Logger{}
		regions = map[string]string{"https://region-1": "region-1"}

		subnetworks = compute.NewSubnetworks(client, logger, regions)
	})

	Describe("List", func() {
		var filter string

		BeforeEach(func() {
			logger.PromptCall.Returns.Proceed = true
			client.ListSubnetworksCall.Returns.Output = &gcpcompute.SubnetworkList{
				Items: []*gcpcompute.Subnetwork{{
					Name:   "banana-subnetwork",
					Region: "https://region-1",
				}},
			}
			filter = "banana"
		})

		It("lists, filters, and prompts for subnetworks to delete", func() {
			list, err := subnetworks.List(filter)
			Expect(err).NotTo(HaveOccurred())

			Expect(client.ListSubnetworksCall.CallCount).To(Equal(1))
			Expect(client.ListSubnetworksCall.Receives.Region).To(Equal("region-1"))

			Expect(logger.PromptCall.Receives.Message).To(Equal("Are you sure you want to delete subnetwork banana-subnetwork?"))

			Expect(list).To(HaveLen(1))
			Expect(list).To(HaveKeyWithValue("banana-subnetwork", "region-1"))
		})

		Context("when the client fails to list subnetworks", func() {
			BeforeEach(func() {
				client.ListSubnetworksCall.Returns.Error = errors.New("some error")
			})

			It("returns the error", func() {
				_, err := subnetworks.List(filter)
				Expect(err).To(MatchError("Listing subnetworks for region region-1: some error"))
			})
		})

		Context("when the subnetwork name does not contain the filter", func() {
			It("does not add it to the list", func() {
				list, err := subnetworks.List("grape")
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.CallCount).To(Equal(0))
				Expect(list).To(HaveLen(0))
			})
		})

		Context("when it is the default subnetwork", func() {
			BeforeEach(func() {
				client.ListSubnetworksCall.Returns.Output = &gcpcompute.SubnetworkList{
					Items: []*gcpcompute.Subnetwork{{
						Name:   "default",
						Region: "https://region-1",
					}},
				}
			})

			It("does not add it to the list", func() {
				list, err := subnetworks.List(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(logger.PromptCall.CallCount).To(Equal(0))
				Expect(list).To(HaveLen(0))
			})
		})

		Context("when the user says no to the prompt", func() {
			BeforeEach(func() {
				logger.PromptCall.Returns.Proceed = false
			})

			It("does not add it to the list", func() {
				list, err := subnetworks.List(filter)
				Expect(err).NotTo(HaveOccurred())

				Expect(list).To(HaveLen(0))
			})
		})
	})

	Describe("Delete", func() {
		var list map[string]string

		BeforeEach(func() {
			list = map[string]string{"banana-subnetwork": "region-1"}
		})

		It("deletes subnetworks", func() {
			subnetworks.Delete(list)

			Expect(client.DeleteSubnetworkCall.CallCount).To(Equal(1))
			Expect(client.DeleteSubnetworkCall.Receives.Region).To(Equal("region-1"))
			Expect(client.DeleteSubnetworkCall.Receives.Subnetwork).To(Equal("banana-subnetwork"))

			Expect(logger.PrintfCall.Messages).To(Equal([]string{"SUCCESS deleting subnetwork banana-subnetwork\n"}))
		})

		Context("when the client fails to delete a subnetwork", func() {
			BeforeEach(func() {
				client.DeleteSubnetworkCall.Returns.Error = errors.New("some error")
			})

			It("logs the error", func() {
				subnetworks.Delete(list)

				Expect(logger.PrintfCall.Messages).To(Equal([]string{"ERROR deleting subnetwork banana-subnetwork: some error\n"}))
			})
		})
	})
})
