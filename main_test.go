package main

import (
	"context"
	"testing"

	"github.com/luccas-eng/tweetstorm_v2/service"
)

func TestTweetstormGenerator(t *testing.T) {

	// testData
	testData := "Mussum Ipsum, cacilds vidis litro abertis. Paisis, filhis, espiritis santis. Viva Forevis aptent taciti sociosqu ad litora torquent. Não sou faixa preta cumpadi, sou preto inteiris, inteiris. Praesent malesuada urna nisi, quis volutpat erat hendrerit non. Nam vulputate dapibus."

	s := service.NewService(context.Background())

	_, err := s.GenerateTweets(testData)
	if err != nil {
		t.Errorf("generateTweets(): %w", err)
	}

}
