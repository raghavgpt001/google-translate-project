package cli

import (
	"log"
	"net/http"
	"sync"

	"github.com/Jeffail/gabs"
)

type RequestBody struct {
	SourceLang string
	TargetLang string
	SourceText string
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", translateUrl, nil)

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)
	req.URL.RawQuery = query.Encode()

	if err != nil {
		log.Fatalf("1. There was a problem: %s", err)
	}

	//this is the place where the request is actually made
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("2. There was a problem: %s", err)
	}

	defer res.Body.Close()

	//you may get blocked if there are too many requests because golang can make loads of requests in very less time
	if res.StatusCode == http.StatusTooManyRequests {
		str <- "You have been rate limited, Try again later."
		wg.Done()
		return
	}

	//parse the json using gabs package
	parsedJson, err := gabs.ParseJSONBuffer(res.Body)

	if err != nil {
		log.Fatalf("3. There was a problem:  %s", err)
	}
	//get the nested elements at the root of parsedJson variable
	nestOne, err := parsedJson.ArrayElement(0)
	if err != nil {
		log.Fatalf("4. There was a problem:  %s", err)
	}

	//get one level deeper in nested element
	nestTwo, err := nestOne.ArrayElement(0)
	if err != nil {
		log.Fatalf("5. There was a problem:  %s", err)
	}

	//get one level deeper in nested element
	translatedStr, err := nestTwo.ArrayElement(0)
	if err != nil {
		log.Fatalf("6. There was a problem:  %s", err)
	}

	str <- translatedStr.Data().(string)

	wg.Done()
}
