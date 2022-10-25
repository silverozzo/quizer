package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
}

func NewConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalln("не загружен файл окружения:", err)
	}

	cfg := &Config{}

	return cfg
}

func (cfg *Config) GetQuizHost() string {
	return cfg.loadField("QUIZ_HOST")
}

func (cfg *Config) GetGoroutinesCount() int {
	str := cfg.loadField("GOROUTINES_COUNT")
	res, err := strconv.Atoi(str)
	if err != nil {
		log.Fatalln("не правильное числовое поле:", err)
	}

	return res
}

func (cfg *Config) loadField(fld string) string {
	val, ok := os.LookupEnv(fld)
	if !ok {
		log.Fatalln("в окружении нет поля:", fld)
	}

	return val
}
