package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type teleInfo struct {
	Token string `json:"Token"`
	Uid   int64  `json:"Uid"`
}

// Scrape 11st HappyMoney
func main() {
	var baseURL string = "http://www.11st.co.kr/products/2778024489"
	fn := "/docker/goapps/11st_happymoney.log"

	if judgeLogger(fn) == true {

		itemStatus := getPages(baseURL)
		fmt.Println("itemStatus = ", itemStatus)

		if itemStatus == "현재 판매중인 상품이 아닙니다." || itemStatus == "일시품절로 구매가 불가합니다." {
			fmt.Println("판매중단")

		} else {
			fmt.Println("구매가능")
			teleSend()
			writePurchaseLog(fn)

		}

	}

}

func judgeLogger(fn string) bool {
	fmt.Println("===read file====")
	fmt.Println(timeString())
	dat, err := ioutil.ReadFile(fn)
	checkErr(err)

	lastNotiLog := string(dat)[len(string(dat))-20 : len(string(dat))-7]
	fmt.Println("LastLog = ", lastNotiLog)

	nowTime := timeString()[:13]
	fmt.Println("NowTime = ", nowTime)

	if lastNotiLog != nowTime {
		// 이번 타임에 알람로그가 없으면 상품을 트래킹한다.
		fmt.Println("가격 추적대상")
		return true

	} else {
		// 이번 타임에 이미 알람로그가 있으면
		fmt.Println("가격 추적 비대상")
		return false
	}

}

func writePurchaseLog(fn string) {

	today := timeString() + "\n"

	// 파일 스트림 생성
	file, err := os.OpenFile(fn, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.FileMode(0655))
	checkErr(err)

	// 메소드 종료 시 파일 닫기
	defer file.Close()

	// 쓰기 버퍼 선언
	w := bufio.NewWriter(file)
	w.WriteString(today)

	// Flush
	w.Flush()
}

func teleSend() {
	// Read Tele info jsonfile
	b, err := ioutil.ReadFile("/docker/goapps/telegram.json")
	if err != nil {
		fmt.Println("Telegram infofile reading failure")
		fmt.Println(err)
		return
	}
	conf := new(teleInfo) // new 로 만들면 포인터
	json.Unmarshal(b, &conf)

	bot, err := tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	txt := "해피머니 GO\n\nhttp://www.11st.co.kr/products/2778024489"

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)

	fmt.Println(updates)
	msg := tgbotapi.NewMessage(conf.Uid, txt)
	bot.Send(msg)
}

func getPages(url string) string {

	itemStatus := ""

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	doc.Find(".b_product_info_price").Each(func(i int, s *goquery.Selection) {
		itemStatus = CleanString(s.Text())
	})
	return itemStatus

}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Result is not 200")
	}
}

func CleanString(s string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(s)), " ")
}

func timeString() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
	// return now.Format(time.RFC3339)
}
