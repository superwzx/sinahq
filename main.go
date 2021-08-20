package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

const sinahqURLPrefix = "https://hq.sinajs.cn/list="

type cfg struct {
	Stocks []struct {
		Name string `yaml:"name"`
		Code string `yaml:"code"`
	}
}

type stockModel struct {
	name  string
	code  string
	start float64
	now   float64
	drise float64
}

func (s *stockModel) hq() {
	URL := sinahqURLPrefix + s.code
	res, err := http.Get(URL)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
	}

	t := strings.Split(string(body), "\"")

	d := strings.Split(t[1], ",")

	s.start, _ = strconv.ParseFloat(d[2], 64)

	s.now, _ = strconv.ParseFloat(d[3], 64)

	s.drise, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", (s.now-s.start)/s.start*100), 64)

	log.Println(s)
}

func printStocks(c *cfg) {
	for _, value := range c.Stocks {
		s := stockModel{name: value.Name, code: value.Code}
		s.hq()
	}
}

func main() {
	yamlFile, err := ioutil.ReadFile("stock.yaml")

	if err != nil {
		log.Fatal(err)
	}

	conf := cfg{}

	err = yaml.Unmarshal(yamlFile, &conf)

	if err != nil {
		log.Fatal(err)
	}

	printStocks(&conf)

	ticker := time.NewTicker(30 * time.Second)

	defer ticker.Stop()

	for range ticker.C {
		printStocks(&conf)
	}

}
