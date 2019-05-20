package tests

// import (
// 	"bytes"
// 	"encoding/json"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"os"
// 	"strconv"
// 	"time"

// 	"github.com/clintjedwards/basecoat/models"
// 	"github.com/icrowley/fake"
// )

// func generateColor() string {

// 	colorNames := []string{
// 		"old soul",
// 		"mountain peak white",
// 		"white cloud",
// 		"plaster of paris",
// 		"pacific ocean blue",
// 		"true green",
// 		"gulf shores",
// 		"cayman lagoon",
// 		"paradise peach",
// 		"antique copper",
// 	}

// 	rand.Seed(time.Now().UTC().UnixNano())
// 	randColor := rand.Intn(len(colorNames))
// 	return colorNames[randColor]
// }

// func populateDB(entries int) {
// 	client := &http.Client{}

// 	Jobs := []models.Job{}

// 	rand.Seed(time.Now().UTC().UnixNano())
// 	randJob := rand.Intn(5)

// 	for i := 0; i < 5; i++ {
// 		id, _ := strconv.Atoi(fake.Digits())
// 		job := models.Job{
// 			ID:      id,
// 			Name:    fake.Company(),
// 			Street:  fake.StreetAddress(),
// 			City:    fake.City(),
// 			State:   fake.State(),
// 			Zipcode: fake.Zip(),
// 		}

// 		Jobs = append(Jobs, job)

// 		requestBodyBytes, err := json.Marshal(job)
// 		if err != nil {
// 			log.Fatalf("could not marshal json data: %v", err)
// 		}

// 		request, err := http.NewRequest("POST", "http://localhost:8080/jobs", bytes.NewReader(requestBodyBytes))
// 		if err != nil {
// 			log.Fatalf("could not create request: %v", err)
// 		}

// 		request.Header.Set("Authorization", "aG9sYnJvOnRlc3R0b2tlbg==")

// 		response, err := client.Do(request)
// 		if err != nil {
// 			log.Fatalf("could not get response: %v", err)
// 		}
// 		defer response.Body.Close()

// 		if response.StatusCode != http.StatusCreated {
// 			log.Printf("expected status Created; got %v", response.Status)
// 		}
// 	}

// 	for i := 0; i < entries; i++ {

// 		formula := models.Formula{
// 			Name:   generateColor(),
// 			Number: fake.CharactersN(2) + "-" + fake.DigitsN(3),
// 			Notes:  fake.SentencesN(3),
// 			Base: map[string]string{
// 				generateColor(): fake.DigitsN(1),
// 			},
// 			Colorants: map[string]string{
// 				generateColor(): fake.DigitsN(1),
// 				generateColor(): fake.DigitsN(1),
// 			},
// 		}

// 		if i%2 == 0 {
// 			formula.Jobs = append(formula.Jobs, Jobs[randJob].ID)
// 		}

// 		requestBodyBytes, err := json.Marshal(formula)
// 		if err != nil {
// 			log.Fatalf("could not marshal json data: %v", err)
// 		}

// 		request, err := http.NewRequest("POST", "http://localhost:8080/formulas", bytes.NewReader(requestBodyBytes))
// 		if err != nil {
// 			log.Fatalf("could not create request: %v", err)
// 		}

// 		request.Header.Set("Authorization", "aG9sYnJvOnRlc3R0b2tlbg==")

// 		response, err := client.Do(request)
// 		if err != nil {
// 			log.Fatalf("could not get response: %v", err)
// 		}
// 		defer response.Body.Close()

// 		if response.StatusCode != http.StatusCreated {
// 			log.Printf("expected status Created; got %v", response.Status)
// 		}
// 	}
// }

func main() {
	// if len(os.Args) < 2 {
	// 	log.Fatal("not enough arguments")
	// }
	// entryString := os.Args[1]
	// entryNum, _ := strconv.Atoi(entryString)
	// populateDB(entryNum)
}
