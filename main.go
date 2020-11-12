package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/luccas-eng/tweetstorm_v2/service"
)

func main() {
	var (
		//input  string
		tweets []string
		err    error
	)

	// checking args
	// args := os.Args[1:]
	// if len(args) >= 1 {
	// 	input = strings.Join(os.Args[1:], " ")
	// } else {
	// 	//s.Instructions()
	// 	return
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	s := service.NewService(ctx)

	token, err := s.RefreshToken(ctx)
	if err != nil {
		log.Println(fmt.Errorf("s.RefreshToken(): %w", err))
	}

	data, err := s.GetData(ctx, token)
	if err != nil {
		log.Println(fmt.Errorf("s.GetData(): %w", err))
	}

	tweets, err = s.GenerateTweets(data)
	if err != nil {
		log.Panic(fmt.Errorf("generateTweets(): %w", err))
	}

	s.PrintTweets(tweets)
}
