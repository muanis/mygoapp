package main

import (
	"bytes"
	"net/http"
	"log"
	"errors"
	"io/ioutil"
	"strconv"
	"sync"
	"bufio"
	"os"
)

func main() {
	
	// we want X goroutines reading from channel and weighing cats
	cats := make(chan []byte, 10000)
	weights := make(chan int64, 10000)

	var wg sync.WaitGroup
	for i:=0; i < 10; i++ {
		wg.Add(1)
		go weighCats(cats, weights, &wg)
	}



	// we need another gorouting listening to results and gathering them

	// read each line of cats.json

	// hand the data of that len to c achannel
	file, err := os.Open("cats.json")
	if err !=nil {
		log.Fatal(err)
	}
	s := bufio.NewScanner(file)
	for s.Scan() {
		catBytes := s.Bytes()
		cat := append([]byte{}, catBytes...)
		cats <- cat
	}
	close(cats)
	file.Close()
	
	if err := s.Err(); err !=nil {
		log.Fatal(err)
	}
	// wait for the weighers to complete
	wg.Wait()
}

func weighCats(cats chan []byte, weights chan int64, wg *sync.WaitGroup) {
	defer wg.Done()

	for cat := range cats {
		weight, err := weighCat(cat)
		if err != nil {
			log.Print(err)
			continue
		}
		log.Print(weight)
		weights <- weight
	}

}

func weighCat( cat []byte) (int64, error) {
	resp, err := http.Post("http://localhost:8000/weighcats", "application/json", bytes.NewReader(cat))
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != 200 {
		return 0, errors.New("bad status code")
	}

	var bod []byte
	bod, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(string(bod), 10, 64)
}