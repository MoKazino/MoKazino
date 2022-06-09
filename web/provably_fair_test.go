package web_test

import (
	"log"
	"testing"

	"x.x/x/xmrdice/web"
)

func TestProbability(t *testing.T) {
	var client_seed = web.RandomString(16)
	var server_seed = web.RandomString(16)
	var answers = []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := uint(0); i < 1000000; i++ {
		answers[web.GenerateProvablyFairNumber(client_seed, server_seed, i, 9)]++
	}
	log.Println(answers)
	//t.Fail()
}
