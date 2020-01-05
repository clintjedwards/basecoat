package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/icrowley/fake"
	"go.uber.org/zap"
)

var info = struct {
	account        string
	storage        storage.BoltDB
	contractorList []string
	jobList        []string
	formulaList    []string
}{}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	storage, err := storage.NewBoltDB(config.Database.Path, config.Database.IDLength)
	if err != nil {
		log.Fatal(err)
	}

	info.storage = storage
	info.storage.CreateAccount("test", "test")
	info.account = "test"
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

func createContractors(num int) {
	for i := 0; i < num; i++ {
		newContractor := &api.Contractor{
			Company: fake.Company(),
			Contact: &api.Contact{
				Name:  fake.FullName(),
				Email: fake.EmailAddress(),
				Phone: fake.Phone(),
			},
		}

		key, err := info.storage.AddContractor(info.account, newContractor)
		if err != nil {
			zap.S().Fatal(err)
		}

		info.contractorList = append(info.contractorList, key)
	}
}

func createJobs(num int) {
	for i := 0; i < num; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		contractorid := info.contractorList[rand.Intn(len(info.contractorList)-1)]

		newJob := &api.Job{
			Name: fake.Company(),
			Address: &api.Address{
				Street:  fake.StreetAddress(),
				Street2: "APT " + fake.DigitsN(2),
				City:    fake.City(),
				State:   fake.State(),
				Zipcode: fake.Zip(),
			},
			Contact: &api.Contact{
				Name:  fake.FullName(),
				Email: fake.EmailAddress(),
				Phone: fake.Phone(),
			},
			Notes:        fake.WordsN(30),
			ContractorId: contractorid,
		}

		key, err := info.storage.AddJob(info.account, newJob)
		if err != nil {
			zap.S().Fatal(err)
		}

		info.jobList = append(info.jobList, key)
	}
}

func createFormulas(num int) {
	for i := 0; i < num; i++ {
		rand.Seed(time.Now().UTC().UnixNano())
		jobs := []string{}

		for j := 0; j < rand.Intn(3); j++ {
			jobid := info.jobList[rand.Intn(len(info.jobList)-1)]
			jobs = append(jobs, jobid)
		}

		newFormula := &api.Formula{
			Name:   generateColor(),
			Number: fake.DigitsN(2) + "-" + fake.DigitsN(4),
			Notes:  fake.WordsN(30),
			Jobs:   jobs,
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

		key, err := info.storage.AddFormula(info.account, newFormula)
		if err != nil {
			zap.S().Fatal(err)
		}

		info.formulaList = append(info.formulaList, key)
	}
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run populateDB.go <numContractors> <numJobs> <numFormulas>")
		os.Exit(1)
	}

	numContractors, _ := strconv.Atoi(os.Args[1])
	createContractors(numContractors)
	numJobs, _ := strconv.Atoi(os.Args[2])
	createJobs(numJobs)
	numFormulas, _ := strconv.Atoi(os.Args[1])
	createFormulas(numFormulas)
}
