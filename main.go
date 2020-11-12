package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/luccas-eng/tweetstorm_v2/service"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	s := service.NewService(ctx)

	//checking args
	args := os.Args[1:]
	if len(args) >= 1 {
		if args[0] == "help" || args[0] == "--help" {
			s.Instructions()
			return
		}
		return
	}

	err := s.StartPrint(ctx)
	if err != nil {
		log.Println(err)
	}
}
