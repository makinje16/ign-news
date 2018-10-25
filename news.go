package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// CHARACTER_LIMIT is the twilio api character limit
const CHARACTER_LIMIT = 1600

type headlines struct {
	Author      string
	Title       string
	Description string
	Url         string
	UrlToImage  string
	PublishedAt string
	Content     string
}

type IgnResponse struct {
	Status       string
	TotalResults int
	Articles     []headlines
}

func main() {
	if len(os.Args) == 2 {
		setTimer()
	} else {
		print("Usage:\n")
		print("      ./news <cellNumber>\n")
	}
}

func setTimer() {
	t := time.Now()
	n := time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, t.Location())
	d := n.Sub(t)

	if d < 0 {
		n = n.Add(24 * time.Hour)
		d = n.Sub(t)
	}
	for {
		time.Sleep(d)
		d = 24 * time.Hour
		parseHeadlines()
	}
}

func parseHeadlines() {
	count := 0
	var message string
	newsAPIKey := os.Getenv("NEWS_API_KEY")
	twilioNumber := os.Getenv("TWILIO_NUMBER")
	receivingNumber := os.Args[1]

	resp, err := http.Get("https://newsapi.org/v2/top-headlines?sources=ign&apiKey=" + newsAPIKey)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		var ignResponse IgnResponse
		err := json.Unmarshal(bodyBytes, &ignResponse)
		if err == nil {
			for _, articles := range ignResponse.Articles {
				if len(message)+len(articles.Description+"\n"+articles.Url+" \n\n") > CHARACTER_LIMIT {
					sendMessage(message, twilioNumber, receivingNumber)
					message = articles.Description + "\n" + articles.Url + " \n\n"
					count = 0
				} else {
					message += articles.Description + "\n" + articles.Url + " \n\n"
					count = len(message)
				}
			}
		}
	}
	if count > 0 {
		sendMessage(message, twilioNumber, receivingNumber)
	}
}

func sendMessage(body string, from string, to string) {

	accountSid := os.Getenv("TWILIO_SID")
	authToken := os.Getenv("TWILIO_API_KEY")

	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"

	msg := url.Values{}
	msg.Set("To", to)
	msg.Set("From", from)
	msg.Set("Body", body)
	msgDataReader := *strings.NewReader(msg.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)

	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		fmt.Println(resp.Status)
		reader, _ := ioutil.ReadAll(resp.Body)
		print(string(reader))
	}

}
