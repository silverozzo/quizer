// main
package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"quizer/html"
	"quizer/model"
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
	if rsp.StatusCode != http.StatusOK {
		log.Fatalln("получили ошибочный ответ с начала опроса", rsp.StatusCode)
	}

	q1 := quizHost + "/question/1"
	rsp, err = cli.Get(q1)
	if err != nil {
		log.Fatalln(err)
	}
	if rsp.StatusCode != http.StatusOK {
		log.Fatalln("получили ошибочный ответ с вопроса 1", rsp.StatusCode)
	}

	for {
		var inputs []model.Input
		html.Parse(rsp.Body, &inputs)
		if len(inputs) == 0 {
			log.Println("внезапно прошли тест!!!")

			return
		}

		inp := model.Fields(&inputs)
		frm := make(url.Values)
		for key, vl := range inp {
			frm.Add(key, vl)
		}

		rsp, err = cli.PostForm(q1, frm)
		if err != nil {
			log.Fatalln(err)
		}
		if rsp.StatusCode != http.StatusOK {
			log.Fatalln("получили ошибочный ответ с ответа на очередной вопрос", rsp.StatusCode)
		}

		q1 = rsp.Request.URL.String()
		log.Println(q1)
	}
}
