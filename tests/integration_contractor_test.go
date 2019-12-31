package tests

import (
	"context"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func (info *testHarness) TestCreateContractor(t *testing.T) {
	t.Run("CreateContractor", func(t *testing.T) {
		expectedResponse := api.Contractor{
			Company: "Schimmel Group",
			Email:   "dstehr@swift.com",
			Phone:   "555-555-5555",
			Contact: "David",
			Jobs: []*api.Job{
				&api.Job{
					Name: "Presidential Paintdown",
					Address: &api.Address{
						Street:  "1600 Pennsylvania Ave NW",
						Street2: "ATTN: Secret Service",
						City:    "Washington",
						State:   "District of Columbia",
						Zipcode: "20500",
					},
				},
			},
		}

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Schimmel Group",
			Email:   "dstehr@swift.com",
			Phone:   "555-555-5555",
			Contact: "David",
			Jobs: []*api.Job{
				&api.Job{
					Name: "Presidential Paintdown",
					Address: &api.Address{
						Street:  "1600 Pennsylvania Ave NW",
						Street2: "ATTN: Secret Service",
						City:    "Washington",
						State:   "District of Columbia",
						Zipcode: "20500",
					},
				},
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)

		require.NotNil(t, createResponse)
		require.NoError(t, err)
		require.NotEmpty(t, createResponse)
		require.NotNil(t, createResponse.Contractor.Id)
		require.Equal(t, expectedResponse.Company, createResponse.Contractor.Company)
		require.NotNil(t, createResponse.Contractor.Jobs)
		require.Equal(t, expectedResponse.Jobs[0].Address, createResponse.Contractor.Jobs[0].Address)
	})
}

func (info *testHarness) TestGetContractor(t *testing.T) {
	t.Run("GetContractor", func(t *testing.T) {

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Beer, Bergnaum and Simonis",
			Email:   "fortiz@morissette.com",
			Phone:   "555-555-5555",
			Contact: "David",
			Jobs: []*api.Job{
				&api.Job{
					Name: "Presidential Paintdown",
					Address: &api.Address{
						Street:  "1600 Pennsylvania Ave NW",
						Street2: "ATTN: Secret Service",
						City:    "Washington",
						State:   "District of Columbia",
						Zipcode: "20500",
					},
				},
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		getContractorRequest := &api.GetContractorRequest{
			Id: createResponse.Contractor.Id,
		}

		getResponse, err := info.client.GetContractor(ctx, getContractorRequest)

		require.NotNil(t, getResponse)
		require.NoError(t, err)
		require.NotEmpty(t, getResponse)
		require.NotNil(t, getResponse.Contractor.Id)
		require.Equal(t, getResponse.Contractor.Company, createContractorRequest.Company)
		require.NotNil(t, getResponse.Contractor.Jobs)
		require.Equal(t, getResponse.Contractor.Jobs[0].Address, createResponse.Contractor.Jobs[0].Address)
	})
}

func (info *testHarness) TestListContractors(t *testing.T) {
	t.Run("ListContractors", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		listResponse, err := info.client.ListContractors(ctx, &api.ListContractorsRequest{})

		require.NotNil(t, listResponse)
		require.NoError(t, err)
		require.NotEmpty(t, listResponse.Contractors)
	})
}

func (info *testHarness) TestUpdateContractor(t *testing.T) {
	t.Run("UpdateContractor", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Original Company",
			Jobs: []*api.Job{
				&api.Job{
					Name: "Presidential Palace",
					Address: &api.Address{
						Zipcode: "20501",
					},
				},
			},
		}

		createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		updateContractorRequest := &api.UpdateContractorRequest{
			Id:      createResponse.Contractor.Id,
			Company: "Changed Company",
			Jobs: []*api.Job{
				&api.Job{
					Name: "Zombieland",
					Address: &api.Address{
						Zipcode: "20501",
					},
				},
			},
		}

		updateResponse, err := info.client.UpdateContractor(ctx, updateContractorRequest)

		require.NotNil(t, updateResponse)
		require.NoError(t, err)
		require.NotEmpty(t, updateResponse)
		require.Equal(t, updateContractorRequest.Company, updateResponse.Contractor.Company)
		require.Equal(t, updateContractorRequest.Jobs[0].Name, updateResponse.Contractor.Jobs[0].Name)
		require.Equal(t, updateContractorRequest.Jobs[0].Address.Zipcode, updateResponse.Contractor.Jobs[0].Address.Zipcode)
	})
}
