package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/basecoat/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var JobID string
var FormulaID string
var APIKey string

var opts []grpc.DialOption

func init() {

	creds, err := credentials.NewClientTLSFromFile("../localhost.crt", "")
	if err != nil {
		utils.StructuredLog(utils.LogLevelFatal, "failed to get certificates", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))

	hash, err := utils.HashPassword([]byte("test"))
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	storage, err := storage.InitStorage()
	if err != nil {
		log.Fatalf("could not connect to storage: %v", err)
	}

	err = storage.CreateUser("test", &api.User{
		Name: "test",
		Hash: string(hash),
	})
	if err != nil {
		if err == utils.ErrEntityExists {
			log.Printf("could not create user: %v\n", err)
			return
		}
		log.Fatalf("could not create user: %v", err)
	}
}

func TestCreateAPIToken(t *testing.T) {

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	createAPITokenRequest := &api.CreateAPITokenRequest{
		User:     "test",
		Password: "test",
		Duration: 300,
	}

	createResponse, err := basecoatClient.CreateAPIToken(context.Background(), createAPITokenRequest)
	if err != nil {
		t.Fatalf("could not create Token: %v", err)
	}

	APIKey = createResponse.Key
}

func TestCreateJob(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	createJobRequest := &api.CreateJobRequest{
		Name:    "The White House",
		Street:  "1600 Pennsylvania Ave NW",
		Street2: "ATTN: Secret Service",
		City:    "Washington",
		State:   "District of Columbia",
		Zipcode: "20500",
		Notes:   "Some sample notes",
		Contact: &api.Contact{
			Name: "President Barack Obama",
			Info: "Barryo@gmail.com",
		},
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	createResponse, err := basecoatClient.CreateJob(ctx, createJobRequest)
	if err != nil {
		t.Fatalf("could not create Job: %v", err)
	}

	getResponse, err := basecoatClient.GetJob(ctx, &api.GetJobRequest{
		Id: createResponse.Id,
	})

	newJob := getResponse.Job
	if newJob.Id == "" {
		t.Fatalf("failed to get ID for new job: %v", getResponse)
	}
	if newJob.Name != createJobRequest.Name {
		t.Fatalf("failed to get correct name for new job; expected %s; got %s", createJobRequest.Name, newJob.Name)
	}
	if newJob.Contact.Name != createJobRequest.Contact.Name {
		t.Fatalf("failed to get correct contact for new job; expected %s; got %s", createJobRequest.Contact, newJob.Contact)
	}

	JobID = newJob.Id
}

func TestUpdateJob(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	updateJobRequest := &api.UpdateJobRequest{
		Id:      JobID,
		Street2: "ATTN: Russia",
		Contact: &api.Contact{
			Name: "President Donald J. Trump",
			Info: "djt@gmail.com",
		},
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.UpdateJob(ctx, updateJobRequest)
	if err != nil {
		t.Fatalf("could not update Job: %v", err)
	}

	getResponse, err := basecoatClient.GetJob(ctx, &api.GetJobRequest{
		Id: JobID,
	})

	updatedJob := getResponse.Job
	if updatedJob.Street2 != updateJobRequest.Street2 {
		t.Fatalf("failed to get correct street2 for updated job; expected %s; got %s", updateJobRequest.Street2, updatedJob.Street2)
	}
	if updatedJob.Contact.Name != updateJobRequest.Contact.Name {
		t.Fatalf("failed to get correct contact for updated job; expected %s; got %s", updateJobRequest.Contact, updatedJob.Contact)
	}
}

func TestCreateFormula(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	createFormulaRequest := &api.CreateFormulaRequest{
		Name:   "antique copper",
		Number: "1001",
		Notes:  "Some writing here that we don't care about for testing purposes",
		Jobs:   []string{JobID},
		Bases: []*api.Base{
			&api.Base{
				Name: "Color 0",
			},
		},
		Colorants: []*api.Colorant{
			&api.Colorant{
				Name: "Color 1",
			},
			&api.Colorant{
				Name: "Color 2",
			},
		},
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	createResponse, err := basecoatClient.CreateFormula(ctx, createFormulaRequest)
	if err != nil {
		t.Fatalf("could not create formula: %v", err)
	}

	getResponse, err := basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
		Id: createResponse.Id,
	})

	newFormula := getResponse.Formula
	if newFormula.Id == "" {
		t.Fatalf("failed to get ID for new Formula: %v", getResponse)
	}
	if newFormula.Name != createFormulaRequest.Name {
		t.Fatalf("failed to get correct name for new formula; expected %s; got %s", createFormulaRequest.Name, newFormula.Name)
	}
	if newFormula.Bases[0].Name != createFormulaRequest.Bases[0].Name {
		t.Fatalf("failed to get correct base name for new formula; expected %s; got %s", createFormulaRequest.Bases[0].Name, newFormula.Bases[0].Name)
	}

	FormulaID = newFormula.Id
}

func TestUpdateFormula(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	updateFormulaRequest := &api.UpdateFormulaRequest{
		Id:     FormulaID,
		Name:   "antique copper",
		Number: "1001",
		Notes:  "Some writing here that we don't care about for testing purposes",
		Jobs:   []string{JobID},
		Colorants: []*api.Colorant{
			&api.Colorant{
				Name: "Color 5",
			},
		},
		Bases: []*api.Base{
			&api.Base{
				Name: "Color 0",
			},
		},
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.UpdateFormula(ctx, updateFormulaRequest)
	if err != nil {
		t.Fatalf("could not update formula: %v", err)
	}

	getResponse, err := basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
		Id: FormulaID,
	})

	updatedFormula := getResponse.Formula
	if updatedFormula.Id == "" {
		t.Fatalf("failed to get ID for new Formula: %v", getResponse)
	}
	if updatedFormula.Notes != updateFormulaRequest.Notes {
		t.Fatalf("failed to get correct notes for updated formula; expected %s; got %s", updateFormulaRequest.Name, updatedFormula.Name)
	}
	if updatedFormula.Bases[0].Name != "Color 0" {
		t.Fatalf("failed to get correct base name for new formula; expected %s; got %s", updateFormulaRequest.Bases[0].Name, updatedFormula.Bases[0].Name)
	}
}

func TestJobState(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	getResponse, err := basecoatClient.GetJob(ctx, &api.GetJobRequest{
		Id: JobID,
	})
	if err != nil {
		t.Fatal("failed to get job")
	}

	job := getResponse.Job

	if len(job.Formulas) == 0 {
		t.Fatal("Formula list for job is empty")
	} else if job.Formulas[0] != FormulaID {
		t.Fatalf("Incorrect Formula ID; expected %s; got %s", FormulaID, job.Formulas[0])
	}
}

func TestDeleteJob(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	deleteJobRequest := &api.DeleteJobRequest{
		Id: JobID,
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.DeleteJob(ctx, deleteJobRequest)
	if err != nil {
		t.Fatalf("could not delete job: %v", err)
	}

	_, err = basecoatClient.GetJob(ctx, &api.GetJobRequest{
		Id: JobID,
	})

	if status.Code(err) != codes.NotFound {
		t.Fatalf("failed to remove job: %s", JobID)
	}
}

func TestFormulaState(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	getResponse, err := basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
		Id: FormulaID,
	})
	if err != nil {
		t.Fatal("failed to get formula")
	}

	formula := getResponse.Formula
	if len(formula.Jobs) != 0 {
		t.Fatalf("Removed job was not correctly removed from job list: %v", formula.Jobs)
	}
}

func TestDeleteFormula(t *testing.T) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	deleteFormulaRequest := &api.DeleteFormulaRequest{
		Id: FormulaID,
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.DeleteFormula(ctx, deleteFormulaRequest)
	if err != nil {
		t.Fatalf("could not delete formula: %v", err)
	}

	_, err = basecoatClient.GetFormula(ctx, &api.GetFormulaRequest{
		Id: FormulaID,
	})

	if status.Code(err) != codes.NotFound {
		t.Fatalf("failed to remove formula: %s", FormulaID)
	}
}
