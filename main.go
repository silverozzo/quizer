// main
package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"

	"github.com/pkg/errors"

	"quizer/config"
	"quizer/html"
	"quizer/model"
)

var (
	errWrongStatus = errors.New("wrong response status from quiz host")
)

func main() {
	cfg := config.NewConfig()

	var wg sync.WaitGroup

	for i := 0; i < cfg.GetGoroutinesCount(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := procQuiz(cfg.GetQuizHost()); err != nil {
				log.Fatalln(err)
			}

			log.Println("прошли тест")
		}()
	}

	wg.Wait()
}

func procQuiz(quizHost string) error {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	cli := http.Client{
		Jar: jar,
	}

	rsp, err := cli.Get(quizHost)
	if err != nil {
		return errors.Wrap(err, "ошибка запроса на главную страницу")
	}
	if rsp.StatusCode != http.StatusOK {
		return errors.Wrap(errWrongStatus, "неверный статус от ответа на главной сранице")
	}

	q1 := quizHost + "/question/1"
	rsp, err = cli.Get(q1)
	if err != nil {
		return errors.Wrap(err, "ошибка запроса на страницу первого вопроса")
	}
	if rsp.StatusCode != http.StatusOK {
		return errors.Wrap(errWrongStatus, "неверный статус от ответа на сранице первого вопроса")
	}

	for {
		var inputs []model.Input
		if err := html.Parse(rsp.Body, &inputs); err != nil {
			return errors.Wrap(err, "ошибка распознования html-страницы")
		}
		if len(inputs) == 0 {
			return nil
		}

		inp := model.Fields(&inputs)
		frm := make(url.Values)
		for key, vl := range inp {
			frm.Add(key, vl)
		}

		rsp, err = cli.PostForm(q1, frm)
		if err != nil {
			return errors.Wrap(err, "ошибка загрузки формы")
		}
		if rsp.StatusCode != http.StatusOK {
			return errors.Wrap(errWrongStatus, "неверный статус от ответа на сранице очередного вопроса")
		}

		q1 = rsp.Request.URL.String()
	}
}
