package main

import (
	"testing"
	"github.com/ChimeraCoder/anaconda"
	"github.com/joho/godotenv"
)

func TestGetTweetFromTimeLine(t *testing.T) {
	godotenv.Load(".env")
	api := connectTwitterAPi()
	tweets := getTweetFromTimeLine(api, "30")

	result := len(tweets)
	expect := 30

  if result != expect {
    t.Error("\n結果： ", result, "\n期待： ", expect)
  }

  t.Log("TestGetTweetFromTimeLine終了")
}

func TestRequestTalkAPI(t *testing.T) {
	var tweet anaconda.Tweet
	tweet.FullText = "あなたは誰ですか？"

	talkResponse := requestTalkAPI(tweet)

	result := talkResponse.Status
	expect := 0

  if result != expect {
    t.Error("\n結果： ", result, "\n期待： ", expect)
  }

  t.Log("TestRequestTalkAPI終了")
}

func TestHasSuffix(t *testing.T) {
	var tweet anaconda.Tweet
	tweet.FullText = "あなたは誰ですか？"

	result := HasSuffix(tweet)
	expect := true

	if result != expect {
    t.Error("\n結果： ", result, "\n期待： ", expect)
  }

  t.Log("TestHasSuffix終了")
}
