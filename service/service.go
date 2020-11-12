package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/luccas-eng/tweetstorm_v2/model"
	"github.com/spf13/cast"
)

// constant with tweet lenght
const maxLenght = 45

//Service ...
type Service struct {
	Context        context.Context
	ContextTimeout time.Duration
}

//NewService ...
func NewService(ctx context.Context) *Service {
	return &Service{
		Context: ctx,
	}
}

//GetData ...
func (s *Service) GetData(ctx context.Context, token string) (texto string, err error) {

	req, err := http.NewRequest("GET", "https://n8e480hh63o547stou3ycz5lwz0958.herokuapp.com/1.1/statuses/home_timeline.json", nil)
	if err != nil {
		return "", fmt.Errorf("http.Get(): %w", err)
	}

	bearer := "Bearer " + token

	req.Header.Add("Authorization", bearer)
	client := http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client.Do(): %w", err)
	}
	defer resp.Body.Close()

	var result []map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("json.NewDecoder().Decode(): %w", err)
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(len(result))

	var p = &model.Payload{}
	if result[r]["text"] != nil {
		p.Texto = cast.ToString(result[r]["text"])
	}
	// fmt.Println(p.Texto)

	return p.Texto, nil
}

//RefreshToken ...
func (s *Service) RefreshToken(ctx context.Context) (token string, err error) {

	req, err := http.NewRequest("POST", "https://n8e480hh63o547stou3ycz5lwz0958.herokuapp.com/1.1/auth", nil)
	if err != nil {
		return "", fmt.Errorf("http.Get(): %w", err)
	}

	client := http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client.Do(): %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("json.NewDecoder().Decode(): %w", err)
	}

	var p = &model.Payload{}
	if result["token"] != nil {
		p.Token = cast.ToString(result["token"])
	}
	// fmt.Println(p.Token)

	return p.Token, nil
}

//GenerateTweets split text into an 45 characters string array and returns a string slice with tweets and index
func (s *Service) GenerateTweets(input string) (tweets []string, err error) {

	var (
		inputLenght, i int
	)

	inputLenght = len(input)

	if inputLenght > maxLenght {

		//string reader
		reader := strings.NewReader(input)

		tweetSize, maxIndexSize := s.MapInput(inputLenght)

		readerOffSet := maxLenght - maxIndexSize

		for i = 0; i < tweetSize; i++ {

			var (
				textPart []byte
			)

			// creates a prefix with index of text
			index := "Tweet #" + strconv.Itoa(i+1) + ":" + " "

			// creates a reader offset
			offset, err := reader.Seek(int64(i*readerOffSet), 0)
			if err != nil {
				return nil, fmt.Errorf("reader.Seek(): %w", err)
			}

			// validates the index to set text
			if int(offset)+readerOffSet > inputLenght {
				textPart = make([]byte, int64(inputLenght)-offset)
			} else {
				textPart = make([]byte, readerOffSet)
			}

			// set io to read at least 1 char
			if read, err := io.ReadAtLeast(reader, textPart, 1); err != nil && read == 0 {
				return nil, fmt.Errorf("io.ReadAtLeast(): %w", err)
			}

			// concatenates index and build the final tweet
			tweet := index + string(textPart)

			// append the tweet into an array of string
			tweets = append(tweets, tweet)
		}
	} else {
		tweets = append(tweets, input)
	}

	return
}

//MapInput used to separate in quantity of tweets with chars limit and then returns infos for prefixing each tweet
func (s *Service) MapInput(inputLenght int) (tweetSize int, maxIndexSize int) {

	//calculate total tweets from inputLenght
	tweetSize = inputLenght / maxLenght

	//multiply total of chars by 2 considering prefixed text with index
	maxIndexSize = (len(strconv.Itoa(tweetSize)) * 2) + 2

	//recalculate tweetSize considering prefixed text for each tweet
	tweetSize = (inputLenght + maxIndexSize*tweetSize) / maxLenght

	if inputLenght%maxLenght != 0 {
		tweetSize++
	}

	return tweetSize, maxIndexSize
}

//PrintTweets print tweets on a tweet slice, from the last to the first part
func (s *Service) PrintTweets(tweets []string) {
	for _, tweet := range tweets {
		fmt.Println(tweet)
	}
}

//Instructions ...
func (s *Service) Instructions() {
	fmt.Println("Use the app by executing the program followed by your input text --> go run main.go \"the text\"")
}
