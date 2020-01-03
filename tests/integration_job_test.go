package tests

import (
	"context"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func (info *testHarness) TestCreateJob(t *testing.T) {
	t.Run("CreateJob", func(t *testing.T) {
		expectedResponse := api.CreateJobResponse{
			Job: &api.Job{
				Name: "Hart Productions",
				Address: &api.Address{
					Street: "400 40th street",
				},
			},
		}

		createJobRequest := &api.CreateJobRequest{
			Name: "Hart Productions",
			Address: &api.Address{
				Street: "400 40th street",
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateJob(ctx, createJobRequest)
		require.NoError(t, err)
		require.NotNil(t, createResponse)
		require.NotEmpty(t, createResponse)
		require.NotNil(t, createResponse.Job.Id)
		require.Equal(t, expectedResponse.Job.Name, createResponse.Job.Name)
		require.Equal(t, expectedResponse.Job.Address, createResponse.Job.Address)
	})
}

func (info *testHarness) TestGetJob(t *testing.T) {
	t.Run("GetJob", func(t *testing.T) {

		createJobRequest := &api.CreateJobRequest{
			Name: "Rajj Productions",
			Address: &api.Address{
				Street: "401 40th street",
			},
		}
		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateJob(ctx, createJobRequest)
		require.NoError(t, err)

		getJobRequest := &api.GetJobRequest{
			Id: createResponse.Job.Id,
		}

		getResponse, err := info.client.GetJob(ctx, getJobRequest)
		require.NoError(t, err)
		require.NotNil(t, getResponse)
		require.NotEmpty(t, getResponse)
		require.NotNil(t, getResponse.Job.Id)
		require.Equal(t, getResponse.Job.Name, createJobRequest.Name)
		require.Equal(t, getResponse.Job.Address, createResponse.Job.Address)
	})
}

func (info *testHarness) TestListJobs(t *testing.T) {
	t.Run("ListJobs", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		listResponse, err := info.client.ListJobs(ctx, &api.ListJobsRequest{})
		require.NoError(t, err)
		require.NotNil(t, listResponse)
		require.NotEmpty(t, listResponse.Jobs)
		require.Len(t, listResponse.Jobs, 3)
	})
}

func (info *testHarness) TestUpdateJob(t *testing.T) {
	t.Run("UpdateJob", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createJobRequest := &api.CreateJobRequest{
			Name: "Destiny Productions",
			Address: &api.Address{
				Street: "402 40th street",
			},
		}

		createResponse, err := info.client.CreateJob(ctx, createJobRequest)
		require.NoError(t, err)

		updateJobRequest := &api.UpdateJobRequest{
			Id:   createResponse.Job.Id,
			Name: "Destiny Stardew",
			Address: &api.Address{
				Street: "403 40th street",
			},
		}

		updateResponse, err := info.client.UpdateJob(ctx, updateJobRequest)
		require.NoError(t, err)
		require.NotNil(t, updateResponse)
		require.NotEmpty(t, updateResponse)
		require.Equal(t, updateJobRequest.Name, updateResponse.Job.Name)
		require.Equal(t, updateJobRequest.Address.Street, updateResponse.Job.Address.Street)
	})
}

func (info *testHarness) TestDeleteJob(t *testing.T) {
	t.Run("DeleteJob", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createJobRequest := &api.CreateJobRequest{
			Name: "Destiny Productions",
			Address: &api.Address{
				Street: "402 40th street",
			},
		}

		createResponse, err := info.client.CreateJob(ctx, createJobRequest)
		require.NoError(t, err)

		deleteJobRequest := &api.DeleteJobRequest{
			Id: createResponse.Job.Id,
		}

		deleteResponse, err := info.client.DeleteJob(ctx, deleteJobRequest)
		require.NoError(t, err)
		require.NotNil(t, deleteResponse)
		require.Empty(t, deleteResponse)

		_, err = info.client.GetJob(ctx, &api.GetJobRequest{
			Id: createResponse.Job.Id,
		})
		require.Error(t, err)
	})
}

// test if job to contractor linking works on create
func (info *testHarness) TestCreateJobAddContractor(t *testing.T) {
	t.Run("CreateJobAddContractor", func(t *testing.T) {

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Edwards Industries",
			Contact: &api.Contact{
				Email: "fortiz@morissette.com",
				Phone: "555-555-5555",
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createContractorResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		expectedResponse := api.CreateJobResponse{
			Job: &api.Job{
				Name: "Hart Productions",
				Address: &api.Address{
					Street: "400 40th street",
				},
				ContractorId: createContractorResponse.Contractor.Id,
			},
		}

		createJobRequest := &api.CreateJobRequest{
			Name: "Hart Productions",
			Address: &api.Address{
				Street: "400 40th street",
			},
			ContractorId: createContractorResponse.Contractor.Id,
		}

		createJobResponse, err := info.client.CreateJob(ctx, createJobRequest)
		require.NoError(t, err)
		require.NotNil(t, createJobResponse)
		require.NotEmpty(t, createJobResponse)
		require.NotNil(t, createJobResponse.Job.Id)
		require.Equal(t, expectedResponse.Job.ContractorId, createJobResponse.Job.ContractorId)

		getContractorResponse, err := info.client.GetContractor(ctx, &api.GetContractorRequest{
			Id: createContractorResponse.Contractor.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, getContractorResponse)
		require.NotEmpty(t, getContractorResponse)
		require.Equal(t, createJobResponse.Job.Id, getContractorResponse.Contractor.Jobs[0])
	})
}

// test if job to contractor linking works
func (info *testHarness) TestUpdateJobRemoveContractor(t *testing.T) {
	t.Run("UpdateJobRemoveContractor", func(t *testing.T) {

		createContractorRequest := &api.CreateContractorRequest{
			Company: "Edwards Industries",
			Contact: &api.Contact{
				Email: "fortiz@morissette.com",
				Phone: "555-555-5555",
			},
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.apikey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createContractorResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
		require.NoError(t, err)

		createJobResponse, err := info.client.CreateJob(ctx, &api.CreateJobRequest{
			Name: "Hart Productions",
			Address: &api.Address{
				Street: "400 40th street",
			},
			ContractorId: createContractorResponse.Contractor.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, createJobResponse)

		updateJobRequest := &api.UpdateJobRequest{
			Id:   createJobResponse.Job.Id,
			Name: "Hart Corp",
			Address: &api.Address{
				Street: "400 40th street",
			},
		}

		updateJobResponse, err := info.client.UpdateJob(ctx, updateJobRequest)
		require.NoError(t, err)
		require.NotNil(t, updateJobResponse)
		require.Equal(t, updateJobRequest.Name, updateJobResponse.Job.Name)
		require.Equal(t, updateJobRequest.ContractorId, updateJobResponse.Job.ContractorId)

		getContractorResponse, err := info.client.GetContractor(ctx, &api.GetContractorRequest{
			Id: createContractorResponse.Contractor.Id,
		})
		require.NoError(t, err)
		require.NotNil(t, getContractorResponse)
		require.Len(t, getContractorResponse.Contractor.Jobs, 0)
	})
}
