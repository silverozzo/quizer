// main
package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
)

const (
	quizHost = "http://185.204.3.165"
)

func main() {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	cli := http.Client{
		Jar: jar,
	}

	rsp, err := cli.Get(quizHost)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("статус начала опроса", rsp.StatusCode)

	q1 := quizHost + "/question/1"
	rsp, err = cli.Get(q1)

	log.Println("статус первого вопроса", rsp.StatusCode)
}
