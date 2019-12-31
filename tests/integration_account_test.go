package tests

import (
	"context"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func (info *testHarness) TestCreateAccount(t *testing.T) {
	t.Run("CreateAccount", func(t *testing.T) {

		createAccountRequest := &api.CreateAccountRequest{
			Id:       "test",
			Password: "test",
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		createResponse, err := info.client.CreateAccount(ctx, createAccountRequest)

		require.NotNil(t, createResponse)
		require.NoError(t, err)
	})
}

func (info *testHarness) TestGetAccount(t *testing.T) {
	t.Run("GetAccount", func(t *testing.T) {

		expectedResponse := &api.GetAccountResponse{
			Account: &api.Account{
				Id: "test",
			},
		}

		getAccountRequest := &api.GetAccountRequest{
			Id: "test",
		}

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		getResponse, err := info.client.GetAccount(ctx, getAccountRequest)

		require.NotNil(t, getResponse)
		require.NoError(t, err)
		require.Equal(t, getResponse.Account.Id, expectedResponse.Account.Id)
	})
}

func (info *testHarness) TestListAccounts(t *testing.T) {
	t.Run("ListAccounts", func(t *testing.T) {

		md := metadata.Pairs("Authorization", "Bearer "+info.adminkey)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		listResponse, err := info.client.ListAccounts(ctx, &api.ListAccountsRequest{})

		require.NotNil(t, listResponse)
		require.NoError(t, err)
		require.NotEmpty(t, listResponse.Accounts)
		require.Equal(t, listResponse.Accounts["test"].Id, "test")
	})
}

// func (info *testHarness) TestCreateContractor(t *testing.T) {

// 	expectedResponse := api.Contractor{
// 		Company: "Schimmel Group",
// 		Email:   "dstehr@swift.com",
// 		Phone:   "555-555-5555",
// 		Contact: "David",
// 		Jobs: []*api.Job{
// 			&api.Job{
// 				Name: "Presidential Paintdown",
// 				Address: &api.Address{
// 					Street:  "1600 Pennsylvania Ave NW",
// 					Street2: "ATTN: Secret Service",
// 					City:    "Washington",
// 					State:   "District of Columbia",
// 					Zipcode: "20500",
// 				},
// 			},
// 		},
// 	}

// 	createContractorRequest := &api.CreateContractorRequest{
// 		Company: "Schimmel Group",
// 		Email:   "dstehr@swift.com",
// 		Phone:   "555-555-5555",
// 		Contact: "David",
// 		Jobs: []*api.AddJob{
// 			&api.AddJob{
// 				Name: "Presidential Paintdown",
// 				Address: &api.Address{
// 					Street:  "1600 Pennsylvania Ave NW",
// 					Street2: "ATTN: Secret Service",
// 					City:    "Washington",
// 					State:   "District of Columbia",
// 					Zipcode: "20500",
// 				},
// 			},
// 		},
// 	}

// 	md := metadata.Pairs("Authorization", "Bearer "+info.key)
// 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// 	createResponse, err := info.client.CreateContractor(ctx, createContractorRequest)
// 	if err != nil {
// 		t.Fatalf("could not create Contractor: %v", err)
// 	}

// 	getResponse, err := info.client.GetContractor(ctx, &api.GetContractorRequest{
// 		Id: createResponse.Id,
// 	})
// 	if err != nil {
// 		t.Fatalf("could not get Contractor: %v", err)
// 	}

// 	assert.NotNil(t, getResponse)
// 	assert.NoError(t, err)
// 	assert.NotEmpty(t, getResponse)

// 	getResponse.Contractor.Id = ""
// 	for _, job := range getResponse.Contractor.Jobs {
// 		job.Id = ""
// 	}

// 	expectedJSON, err := protobufToJSON(&expectedResponse)
// 	assert.NoError(t, err)
// 	gotJSON, err := protobufToJSON(getResponse.GetContractor())
// 	assert.NoError(t, err)

// 	assert.Equal(t, expectedJSON, gotJSON)
// }

// // func TestUpdateJob(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	updateJobRequest := &api.UpdateJobRequest{
// // 		Id:      JobID,
// // 		Street2: "ATTN: Russia",
// // 		Contact: &api.Contact{
// // 			Name: "President Donald J. Trump",
// // 			Info: "djt@gmail.com",
// // 		},
// // 	}

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	_, err = basecoatClient.UpdateJob(ctx, updateJobRequest)
// // 	if err != nil {
// // 		t.Fatalf("could not update Job: %v", err)
// // 	}

// // 	getResponse, err := basecoatClient.GetJob(ctx, &api.GetJobRequest{
// // 		Id: JobID,
// // 	})
// // 	if err != nil {
// // 		t.Fatalf("could not get Job: %v", err)
// // 	}

// // 	updatedJob := getResponse.Job
// // 	if updatedJob.Street2 != updateJobRequest.Street2 {
// // 		t.Fatalf("failed to get correct street2 for updated job; expected %s; got %s", updateJobRequest.Street2, updatedJob.Street2)
// // 	}
// // 	if updatedJob.Contact.Name != updateJobRequest.Contact.Name {
// // 		t.Fatalf("failed to get correct contact for updated job; expected %s; got %s", updateJobRequest.Contact, updatedJob.Contact)
// // 	}
// // }

// // func TestCreateFormula(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	createFormulaRequest := &api.CreateFormulaRequest{
// // 		Name:   "antique copper",
// // 		Number: "1001",
// // 		Notes:  "Some writing here that we don't care about for testing purposes",
// // 		Jobs:   []string{JobID},
// // 		Bases: []*api.Base{
// // 			&api.Base{
// // 				Name: "Color 0",
// // 			},
// // 		},
// // 		Colorants: []*api.Colorant{
// // 			&api.Colorant{
// // 				Name: "Color 1",
// // 			},
// // 			&api.Colorant{
// // 				Name: "Color 2",
// // 			},
// // 		},
// // 	}

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	createResponse, err := basecoatClient.CreateFormula(ctx, createFormulaRequest)
// // 	if err != nil {
// // 		t.Fatalf("could not create formula: %v", err)
// // 	}

// // 	getResponse, err := basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
// // 		Id: createResponse.Id,
// // 	})
// // 	if err != nil {
// // 		t.Fatalf("could not get Formula: %v", err)
// // 	}

// // 	newFormula := getResponse.Formula
// // 	if newFormula.Id == "" {
// // 		t.Fatalf("failed to get ID for new Formula: %v", getResponse)
// // 	}
// // 	if newFormula.Name != createFormulaRequest.Name {
// // 		t.Fatalf("failed to get correct name for new formula; expected %s; got %s", createFormulaRequest.Name, newFormula.Name)
// // 	}
// // 	if newFormula.Bases[0].Name != createFormulaRequest.Bases[0].Name {
// // 		t.Fatalf("failed to get correct base name for new formula; expected %s; got %s", createFormulaRequest.Bases[0].Name, newFormula.Bases[0].Name)
// // 	}

// // 	FormulaID = newFormula.Id
// // }

// // func TestUpdateFormula(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	updateFormulaRequest := &api.UpdateFormulaRequest{
// // 		Id:     FormulaID,
// // 		Name:   "antique copper",
// // 		Number: "1001",
// // 		Notes:  "Some writing here that we don't care about for testing purposes",
// // 		Jobs:   []string{JobID},
// // 		Colorants: []*api.Colorant{
// // 			&api.Colorant{
// // 				Name: "Color 5",
// // 			},
// // 		},
// // 		Bases: []*api.Base{
// // 			&api.Base{
// // 				Name: "Color 0",
// // 			},
// // 		},
// // 	}

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	_, err = basecoatClient.UpdateFormula(ctx, updateFormulaRequest)
// // 	if err != nil {
// // 		t.Fatalf("could not update formula: %v", err)
// // 	}

// // 	getResponse, err := basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
// // 		Id: FormulaID,
// // 	})
// // 	if err != nil {
// // 		t.Fatalf("could not get formula: %v", err)
// // 	}

// // 	updatedFormula := getResponse.Formula
// // 	if updatedFormula.Id == "" {
// // 		t.Fatalf("failed to get ID for new Formula: %v", getResponse)
// // 	}
// // 	if updatedFormula.Notes != updateFormulaRequest.Notes {
// // 		t.Fatalf("failed to get correct notes for updated formula; expected %s; got %s", updateFormulaRequest.Name, updatedFormula.Name)
// // 	}
// // 	if updatedFormula.Bases[0].Name != "Color 0" {
// // 		t.Fatalf("failed to get correct base name for new formula; expected %s; got %s", updateFormulaRequest.Bases[0].Name, updatedFormula.Bases[0].Name)
// // 	}
// // }

// // func TestJobState(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	getResponse, err := basecoatClient.GetJob(ctx, &api.GetJobRequest{
// // 		Id: JobID,
// // 	})
// // 	if err != nil {
// // 		t.Fatal("failed to get job")
// // 	}

// // 	job := getResponse.Job

// // 	if len(job.Formulas) == 0 {
// // 		t.Fatal("Formula list for job is empty")
// // 	} else if job.Formulas[0] != FormulaID {
// // 		t.Fatalf("Incorrect Formula ID; expected %s; got %s", FormulaID, job.Formulas[0])
// // 	}
// // }

// // func TestDeleteJob(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	deleteJobRequest := &api.DeleteJobRequest{
// // 		Id: JobID,
// // 	}

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	_, err = basecoatClient.DeleteJob(ctx, deleteJobRequest)
// // 	if err != nil {
// // 		t.Fatalf("could not delete job: %v", err)
// // 	}

// // 	_, err = basecoatClient.GetJob(ctx, &api.GetJobRequest{
// // 		Id: JobID,
// // 	})

// // 	if status.Code(err) != codes.NotFound {
// // 		t.Fatalf("failed to remove job: %s", JobID)
// // 	}
// // }

// // func TestFormulaState(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	getResponse, err := basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
// // 		Id: FormulaID,
// // 	})
// // 	if err != nil {
// // 		t.Fatal("failed to get formula")
// // 	}

// // 	formula := getResponse.Formula
// // 	if len(formula.Jobs) != 0 {
// // 		t.Fatalf("Removed job was not correctly removed from job list: %v", formula.Jobs)
// // 	}
// // }

// // func TestDeleteFormula(t *testing.T) {
// // 	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8080"), opts...)
// // 	if err != nil {
// // 		log.Fatalf("could not connect to basecoat: %v", err)
// // 	}
// // 	defer conn.Close()

// // 	basecoatClient := api.NewBasecoatClient(conn)

// // 	deleteFormulaRequest := &api.DeleteFormulaRequest{
// // 		Id: FormulaID,
// // 	}

// // 	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
// // 	ctx := metadata.NewOutgoingContext(context.Background(), md)

// // 	_, err = basecoatClient.DeleteFormula(ctx, deleteFormulaRequest)
// // 	if err != nil {
// // 		t.Fatalf("could not delete formula: %v", err)
// // 	}

// // 	_, err = basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
// // 		Id: FormulaID,
// // 	})

// // 	if status.Code(err) != codes.NotFound {
// // 		t.Fatalf("failed to remove formula: %s", FormulaID)
// // 	}
// // }
