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
		expectedResponse := api.CreateContractorResponse{
			Contractor: &api.Contractor{
				Company: "Schimmel Group",
				Contact: &api.Contact{
					Email: "dstehr@swift.com",
					Phone: "555-555-5555",
				},
			},
		}

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Schimmel Group",
			Contact: &api.Contact{
				Email: "dstehr@swift.com",
				Phone: "555-555-5555",
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)
		require.NotNil(t, createResponse)
		require.NotEmpty(t, createResponse)
		require.NotNil(t, createResponse.Contractor.Id)
		require.Equal(t, expectedResponse.Contractor.Company, createResponse.Contractor.Company)
		require.Equal(t, expectedResponse.Contractor.Contact, createResponse.Contractor.Contact)
	})
}

func (info *testHarness) TestGetContractor(t *testing.T) {
	t.Run("GetContractor", func(t *testing.T) {

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Beer, Bergnaum and Simonis",
			Contact: &api.Contact{
				Email: "fortiz@morissette.com",
				Phone: "555-555-5555",
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
		require.NoError(t, err)
		require.NotNil(t, getResponse)
		require.NotEmpty(t, getResponse)
		require.NotNil(t, getResponse.Contractor.Id)
		require.Equal(t, getResponse.Contractor.Company, createContractorRequest.Company)
		require.Equal(t, getResponse.Contractor.String(), createResponse.Contractor.String())
	})
}

func (info *testHarness) TestListContractors(t *testing.T) {
	t.Run("ListContractors", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		listResponse, err := info.client.ListContractors(ctx, &api.ListContractorsRequest{})
		require.NoError(t, err)
		require.NotNil(t, listResponse)
		require.NotEmpty(t, listResponse.Contractors)
		require.Len(t, listResponse.Contractors, 2)
	})
}

func (info *testHarness) TestUpdateContractor(t *testing.T) {
	t.Run("UpdateContractor", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Original Company",
			Contact: &api.Contact{
				Email: "fortiz@morissette.com",
				Phone: "555-555-5555",
			},
		}

		createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		updateContractorRequest := &api.UpdateContractorRequest{
			Id:      createResponse.Contractor.Id,
			Company: "Changed Company",
			Contact: &api.Contact{
				Email: "barack@morissette.com",
				Phone: "555-555-5555",
			},
		}

		updateResponse, err := info.client.UpdateContractor(ctx, updateContractorRequest)
		require.NoError(t, err)
		require.NotNil(t, updateResponse)
		require.NotEmpty(t, updateResponse)
		require.Equal(t, updateContractorRequest.Company, updateResponse.Contractor.Company)
		require.Equal(t, updateContractorRequest.Contact.Email, updateResponse.Contractor.Contact.Email)
	})
}

func (info *testHarness) TestDeleteContractor(t *testing.T) {
	t.Run("DeleteContractor", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Soon to be deleted Company",
			Contact: &api.Contact{
				Email: "fortiz@morissette.com",
				Phone: "555-555-5555",
			},
		}

		createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		deleteContractorRequest := &api.DeleteContractorRequest{
			Id: createResponse.Contractor.Id,
		}

		deleteResponse, err := info.client.DeleteContractor(ctx, deleteContractorRequest)
		require.NoError(t, err)
		require.NotNil(t, deleteResponse)
		require.Empty(t, deleteResponse)

		_, err = info.client.GetContractor(ctx, &api.GetContractorRequest{
			Id: createResponse.Contractor.Id,
		})
		require.Error(t, err)
	})
}

// test if contractor to job linking works
func (info *testHarness) TestCreateContractorAddJob(t *testing.T) {
	t.Run("CreateContractorAddJob", func(t *testing.T) {

		createJobRequest := &api.CreateJobRequest{
			Name: "Hashicorp",
			Address: &api.Address{
				Street: "400 39th street",
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createJobResponse, err := info.client.CreateJob(ctx, createJobRequest)
		require.NoError(t, err)
		require.NotNil(t, createJobResponse)
		require.NotEmpty(t, createJobResponse)
		require.NotNil(t, createJobResponse.Job.Id)

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Armin Industries",
			Contact: &api.Contact{
				Email: "fortiz@morissette.com",
				Phone: "555-555-5555",
			},
			Jobs: []string{createJobResponse.Job.Id},
		}

		createContractorResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		getJobResponse, err := info.client.GetJob(ctx, &api.GetJobRequest{
			Id: createJobResponse.Job.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, getJobResponse)
		require.NotEmpty(t, getJobResponse)
		require.Len(t, createContractorResponse.Contractor.Jobs, 1)
		require.Equal(t, createContractorResponse.Contractor.Id, getJobResponse.Job.ContractorId)
		require.Equal(t, createContractorResponse.Contractor.Jobs[0], getJobResponse.Job.Id)
	})
}
