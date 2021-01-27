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

	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

const (
	BASEURL = "https://api.a3rt.recruit-tech.co.jp/talk/v1/smalltalk"
	COUNT   = "1000"
)

func main() {
	// local用
	godotenv.Load(".env")

	api := connectTwitterAPi()
	tweets := getTweetFromTimeLine(api, COUNT)

	for _, tweet := range tweets {
		if strings.HasSuffix(tweet.FullText, "？") {
			talkResponse := requestTalkAPI(tweet)
			postTweet(talkResponse, tweet, api)
			os.Exit(0)
		}
	}
}

func connectTwitterAPi() *anaconda.TwitterApi {
	return anaconda.NewTwitterApiWithCredentials(
		os.Getenv("ACCESS_TOKEN"),
		os.Getenv("ACCESS_TOKEN_SECRET"),
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"))
}

func getTweetFromTimeLine(api *anaconda.TwitterApi, count string) []anaconda.Tweet {
	v := url.Values{}
	v.Set("count", count)

	tweets, err := api.GetHomeTimeline(v)
	checkError(err)

	return tweets
}

func requestTalkAPI(tweet anaconda.Tweet) TalkResponse {
	req := buildRequest(tweet)
	client := buildClient()

	resp, err := client.Do(req)
	checkError(err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	var talkResponse TalkResponse
	err = json.Unmarshal(body, &talkResponse)
	checkError(err)

	// debug
	fmt.Println(tweet.FullText)
	fmt.Println(talkResponse)

	return talkResponse
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func buildRequest(tweet anaconda.Tweet) *http.Request {
	endpoint := fmt.Sprintf("%s?apikey=%s", BASEURL, os.Getenv("APIKEY"))

	values := make(url.Values)
	values.Set("query", tweet.FullText)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(values.Encode()))
	checkError(err)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req
}

func buildClient() *http.Client {
	client := &http.Client{}
	client.Timeout = time.Second * 15
	return client
}

func postTweet(talkResponse TalkResponse, tweet anaconda.Tweet, api *anaconda.TwitterApi) {
	for _, result := range talkResponse.Results {
		status := fmt.Sprintf("Q. %s\nA. %s", tweet.FullText, result.Reply)
		api.PostTweet(status, nil)
	}
}

type TalkResponse struct {
	Status  int      `json:"status"`
	Massage string   `json:"message"`
	Results []Result `json:"results"`
}

type Result struct {
	Perplexity float64 `json:"perplexity"`
	Reply      string  `json:"reply"`
}
