package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/basecoat/utils"
	"github.com/icrowley/fake"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var JobIDs []string
var APIKey string
var opts []grpc.DialOption

func init() {
	creds, err := credentials.NewClientTLSFromFile("./localhost.crt", "")
	if err != nil {
		log.Fatalf("failed to get certificates: %v", err)
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

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	createAPITokenRequest := &api.CreateAPITokenRequest{
		User:     "test",
		Password: "test",
		Duration: 14400,
	}

	createResponse, err := basecoatClient.CreateAPIToken(context.Background(), createAPITokenRequest)
	if err != nil {
		log.Fatalf("could not create Token: %v", err)
	}

	APIKey = createResponse.Key
}

func generateColor() string {
	colorNames := []string{
		"Alluring Light",
		"Antigua Blue",
		"Aroma",
		"Baby Frog",
		"Bazooka Pink",
		"Biscay",
		"Blue Beauty",
		"Bluewash",
		"Bright Cerulean",
		"Bungalow Gold",
		"Can Can",
		"Cathay Spice",
		"Cherub",
		"City Loft",
		"Community",
		"Cradle Pink",
		"Daddy-O",
		"Deep Orange-coloured Brown",
		"Dirty White",
		"Dusky Citron",
		"Elusive",
		"Eyre",
		"First Frost",
		"Fozzie Bear",
		"Gainsboro",
		"Gluten",
		"Grape Haze",
		"Greige",
		"Haute Red",
		"Hoeth Blue",
		"Ice Mist",
		"Irresistible Beige",
		"Kaitoke Green",
		"Lahmian Medium",
		"Liberal Lilac",
		"Light Spring Burst",
		"Lola",
		"Maiden's Blush",
		"Mayan Treasure",
		"Midnight Blush",
		"Monarch",
		"Murdoch",
		"Night Mode",
		"Oil Yellow",
		"Orion Blue",
		"Parachute Silk",
		"Pekin Chicken",
		"Pincushion",
		"Platonic Blue",
		"Pragmatic",
		"Purple Ragwort",
		"Rattan Palm",
		"Retro",
		"Rose Taupe",
		"Salina Springs",
		"Scrofulous Brown",
		"Shamanic Journey",
		"Silver Sage",
		"Snuggle Pie",
		"Spiced Nectarine",
		"Stone Harbour",
		"Sunday Niqab",
		"Tambua Bay",
		"Thought",
		"Transformer",
		"Ultraviolet Cryner",
		"Victorian Cottage",
		"Wasabi Nori",
		"White Mecca",
		"Wisteria Yellow",
		"Zelyony Green",
		"Zhēn Zhū Bái Pearl",
	}
	rand.Seed(time.Now().UTC().UnixNano())
	randColor := rand.Intn(len(colorNames))
	return colorNames[randColor]
}

func createJob() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	createJobRequest := &api.CreateJobRequest{
		Name:    fake.Company(),
		Street:  fake.StreetAddress(),
		Street2: "APT 5E",
		City:    fake.City(),
		State:   fake.State(),
		Zipcode: fake.Zip(),
		Notes:   fake.WordsN(20),
		Contact: &api.Contact{
			Name: fake.FullName(),
			Info: fake.EmailAddress(),
		},
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	response, err := basecoatClient.CreateJob(ctx, createJobRequest)
	if err != nil {
		log.Fatalf("could not create Job: %v", err)
	}

	JobIDs = append(JobIDs, response.Id)
}

func createFormula() {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "localhost", "8081"), opts...)
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	rand.Seed(time.Now().UTC().UnixNano())
	randJob := rand.Intn(5)

	createFormulaRequest := &api.CreateFormulaRequest{
		Name:   generateColor(),
		Number: fake.DigitsN(2) + "-" + fake.DigitsN(4),
		Notes:  fake.WordsN(20),
		Jobs:   []string{JobIDs[randJob]},
		Bases: []*api.Base{
			{
				Type:   fake.Company(),
				Name:   generateColor(),
				Amount: fake.DigitsN(2),
			},
		},
		Colorants: []*api.Colorant{
			{
				Type:   fake.Company(),
				Name:   generateColor(),
				Amount: fake.DigitsN(2),
			},
			{
				Type:   fake.Company(),
				Name:   generateColor(),
				Amount: fake.DigitsN(2),
			},
		},
	}

	md := metadata.Pairs("Authorization", "Bearer "+APIKey)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.CreateFormula(ctx, createFormulaRequest)
	if err != nil {
		log.Fatalf("could not create formula: %v", err)
	}
}

func populateDB(entries int) {
	for i := 0; i < 5; i++ {
		createJob()
	}

	for i := 0; i < entries; i++ {
		createFormula()
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough arguments")
	}

	entryString := os.Args[1]
	entryNum, _ := strconv.Atoi(entryString)
	populateDB(entryNum)
}
