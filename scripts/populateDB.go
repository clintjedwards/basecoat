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

var jobIDs []string

// dev only creds
const user string = "test"
const pass string = "test"
const certPath string = "./localhost.crt"
const key string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcnkiOjE1NzYzMjI4OTEsInVzZXJuYW1lIjoidGVzdCJ9.DB-U0U7KEXI5_c3uHN6H-1yBVv-W20YOOXP_f0lM2C0"

var opts []grpc.DialOption

func init() {
	creds, err := credentials.NewClientTLSFromFile(certPath, "")
	if err != nil {
		log.Fatalf("failed to get certificates: %v", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))

	hash, err := utils.HashPassword([]byte(pass))
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	storage, err := storage.InitStorage()
	if err != nil {
		log.Fatalf("could not connect to storage: %v", err)
	}

	err = storage.CreateUser(user, &api.User{
		Name: user,
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

	md := metadata.Pairs("Authorization", "Bearer "+key)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	response, err := basecoatClient.CreateJob(ctx, createJobRequest)
	if err != nil {
		log.Fatalf("could not create Job: %v", err)
	}

	jobIDs = append(jobIDs, response.Id)
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
		Jobs:   []string{jobIDs[randJob]},
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

	md := metadata.Pairs("Authorization", "Bearer "+key)
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
