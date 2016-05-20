package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"sync"
)

type Item struct {
	By, Title, Url string
}

func (i Item) String() string {
	return fmt.Sprintf("%v : %v", i.Title, i.Url)
}

func (i Item) HtmlLink() string {
	return fmt.Sprintf("<a href=\"%v\">%v</a>", i.Url, i.Title)
}

type Items []*Item

func (is Items) ToStrings() []string {
	ss := make([]string, len(is))
	for i := 0; i < len(is); i++ {
		ss[i] = fmt.Sprintf("%v", (*is[i]).HtmlLink())
	}
	return ss
}

func getItem(id int, s *Item, wg *sync.WaitGroup) {
	res, _ := http.Get(fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%v.json", id))
	item, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	dec := json.NewDecoder(strings.NewReader(string(item)))
	for {
		if err := dec.Decode(s); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	(*wg).Done()
}

func sendMail(body string) {
	auth := smtp.PlainAuth("", "", "", "smtp.mailgun.org")
	to := []string{""}
	msg := []byte("To: \r\n" +
		"Subject: Top 20 Hacker news\r\n" +
		"MIME-version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		fmt.Sprintf("<html>%v</html>", body))
	err := smtp.SendMail("smtp.mailgun.org:25", auth, "", to, msg)
	if err != nil {
		log.Fatal(err)
	}
}

func getItems(s string) Items {
	dec := json.NewDecoder(strings.NewReader(s))
	var m []int
	for {
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}
	if len(m) > 20 {
		m = m[:20]
	}
	var wg sync.WaitGroup
	wg.Add(len(m))
	items := make(Items, len(m))
	for i := 0; i < len(m); i++ {
		var s Item
		items[i] = &s
		go getItem(m[i], items[i], &wg)
	}
	wg.Wait()
	return items
}

func getUrl(url string) string {
	res, _ := http.Get(url)
	top, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	return string(top)
}

func main() {

	topStories := getItems(getUrl("https://hacker-news.firebaseio.com/v0/topstories.json"))
	jobs := getItems(getUrl("https://hacker-news.firebaseio.com/v0/jobstories.json"))

	body := "<h2>News</h2><br/>" + strings.Join(topStories.ToStrings(), "<br/>") + "<br/><br/>" + "<h2>Jobs</h2><br/>" + strings.Join(jobs.ToStrings(), "<br/>")
	fmt.Println(body)
	sendMail(body)
}
