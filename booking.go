package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Hotel struct {
	Id      int
	Comment string
	Score   int
}

type Hotels []Hotel

func (hs Hotels) Len() int      { return len(hs) }
func (hs Hotels) Swap(i, j int) { hs[i], hs[j] = hs[j], hs[i] }
func (hs Hotels) Less(i, j int) bool {
	if hs[i].Score == hs[j].Score {
		return hs[i].Id < hs[j].Id
	}
	return hs[i].Score > hs[j].Score
}

var (
	M int
)

func catchError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getScore(h Hotel, ws []string) int {
	var score int
	for i := 0; i < len(ws); i++ {
		w := ws[i]
		score = score + strings.Count(h.Comment, w)
	}
	return score
}

func main() {

	reader := bufio.NewReader(os.Stdin)

	text, err := reader.ReadString('\n')
	catchError(err)

	fmt.Fscan(reader, &M)

	words := strings.Split(strings.Replace(text, "\n", "", 1), " ")

	hotels := make(Hotels, M)

	for i := 0; i < M; i++ {
		var Id int
		h := Hotel{}
		fmt.Fscan(reader, &Id)
		text, err := reader.ReadString('\n')
		catchError(err)
		h.Id = Id
		h.Comment = strings.Replace(text, "\n", "", 1)
		h.Score = getScore(h, words)
		hotels[i] = h
	}

	fmt.Println(hotels)

	mh := make(map[int]int)

	for i := 0; i < len(hotels); i++ {
		mh[hotels[i].Id] += hotels[i].Score
	}

	newHotels := make(Hotels, 0)
	for k, v := range mh {
		newHotels = append(newHotels, Hotel{
			Id:    k,
			Score: v,
		})
	}

	sort.Sort(newHotels)

	fmt.Println(newHotels)

	hotelIds := make([]string, M)
	for i := 0; i < len(newHotels); i++ {
		hotelIds[i] = fmt.Sprintf("%v", newHotels[i].Id)
	}

	fmt.Println(strings.Join(hotelIds, " "))
}
