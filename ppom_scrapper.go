package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)

type teleInfo struct {
	Token string `json:"Token"`
	Uid   int64  `json:"Uid"`
}

// 검색을 위해 Key - Struct 구조 사용

type products map[string]ProductInfo

type ProductInfo struct {
	ProdNm   string  `json:"ProdNm"`
	RgstDT   string  `json:"RgstDT"`
	RgstTime string  `json:"RgstTime"`
	ClickCnt int     `json:"ClickCnt"`
	Cpm      float32 `json:"Cpm"`
	NotiYN   bool    `json:"NotiYN"`
}

// Define products Methods
func (p *products) SetName(name string) {
	// fmt.Println("setname called")
	for k, v := range *p {
		v.ProdNm = name
		(*p)[k] = v
	}
}

func (p *products) SetRgstDTTM(time string) {
	// fmt.Println("SetRgstDTTM called")
	for k, v := range *p {
		v.RgstDT = timeString()[:10]
		v.RgstTime = time
		(*p)[k] = v
	}
}

func (p *products) SetClickCnt(cnt int) {
	// fmt.Println("SetClickCnt called")
	for k, v := range *p {
		v.ClickCnt = cnt
		(*p)[k] = v
	}
}

func (p *products) SetCpm() {
	// fmt.Println("SetCpm called")
	for k, v := range *p {
		// v.ClickCnt = cnt
		currentTime := timeString()[11:19]
		if strings.Contains(v.RgstTime, "/") {
			fmt.Println("found /")
			break
		}
		currentHH, _ := strconv.Atoi(currentTime[:2])
		currentMM, _ := strconv.Atoi(currentTime[3:5])
		// currentSS := strconv.Atoi(currentTime[7:9])
		RgstHH, _ := strconv.Atoi(v.RgstTime[:2])
		RgstMM, _ := strconv.Atoi(v.RgstTime[3:5])
		// RgstSS := strconv.Atoi(v.RgstTime[7:9])

		// 경과시간 (분) 계산
		min := (currentHH-RgstHH)*60 + (currentMM - RgstMM)
		cpm := float32(v.ClickCnt) / float32(min)
		fmt.Println("Cpm=", cpm)
		v.Cpm = float32(cpm)

		fmt.Println(currentTime, v.RgstTime, "minutes = ", min)
		(*p)[k] = v
	}
}

// 알람 대상 판단
func (p *products) JudgeAlarm() {
	fmt.Println("Judge Alarm called")
	for k, v := range *p {
		if v.ClickCnt > 500 && v.Cpm > 50 {
			v.NotiYN = true
		}
		(*p)[k] = v
	}
}

// 크롤링 메인 프로세스
func getPrices(url string, fn string, c chan bool) {

	// Read previous productInfo jsonfile
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println("Product infofile reading failure")
		fmt.Println(err)
		c <- false
	}
	fmt.Println(string(b), len(b))
	if len(b) < 5 {
		fmt.Println("Null file read")
		panic("errrrrrr")
	}

	js := new(products) // new 로 만들면 포인터
	// js := make([]products, 1)
	// js := []products{}
	// js := make(map[string]interface{})  // 이건 되는데 인덱싱이 안됨.

	fmt.Printf("%+v\n", js)
	json.Unmarshal(b, &js)
	fmt.Println("JSON File = ", js)

	// 기존 json 데이터 접근 테스트
	// v, found := (*js)["12451251"]
	// fmt.Println("Approaching = ", v, found)

	// v1 := (*js)["12451251"].ProdNm
	// fmt.Println("Approaching = ", v1)

	// if thisval, ok := (*js)["12451251"]; ok {
	// 	thisval.ProdNm = "UPDATED Product"
	// 	(*js)["12451251"] = thisval
	// }
	// v2, found2 := (*js)["12451251"]
	// fmt.Println("Approaching = ", v2, found2)

	// for k, v := range *js {
	// 	fmt.Println(k, v)
	// }

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	// Get data 확인
	fmt.Println("GET = ", res.Body)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	// // Array 타입의 맵을 만들고 인덱스로 접근
	newData := make([]products, 20)
	// newData[0] = products{"312": {"test2", "test", "", 123, 1.0, true}}
	// newData[0] = products{"312": {"test2", "test", "", 123, 1.0, true}}
	// newData[3] = products{"312": {"test2", "test", "", 123, 1.0, true}}
	// for k, v := range newData[0] {
	// 	fmt.Println(k, v)
	// }
	// newData[0].SetName("312")
	// fmt.Println(newData)

	doc.Find(".container .list_vspace").Each(func(i int, s *goquery.Selection) {
		itemStatus := CleanString(s.Text())
		item, err := encodeFromEUCKR(itemStatus)
		checkErr(err)

		if i >= 17 && (i-17)%9 == 0 {
			// 상품ID
			idx := (i - 17) / 9
			newData[idx] = products{item: {"NA", "NA", "NA", 0, 1.0, false}}

			fmt.Println("index=", i, "no=", idx, "       item=", item)
		}

		if i >= 17 && (i-21)%9 == 0 {
			// 상품명
			idx := (i - 21) / 9
			newData[idx].SetName(item)
			fmt.Println(newData[idx])
			fmt.Println("index=", i, "no=", (i-21)/9+1, "       item=", item)
		}

		if i >= 17 && (i-22)%9 == 0 {
			// 등록일시
			idx := (i - 22) / 9
			newData[idx].SetRgstDTTM(item)
			fmt.Println(newData[idx])
			fmt.Println("index=", i, "no=", (i-22)/9+1, "       item=", item)
		}

		if i >= 17 && (i-25)%9 == 0 {
			// 조회 수
			idx := (i - 25) / 9
			cnt, err := strconv.Atoi(item)
			checkErr(err)
			newData[idx].SetClickCnt(cnt)

			// CPM 계산
			newData[idx].SetCpm()

			// 알람 대상인지 판단.
			newData[idx].JudgeAlarm()

			fmt.Println(newData[idx])
			fmt.Println("index=", i, "no=", (i-25)/9+1, "       item=", item)
		}
		// fmt.Println("index=", i, "       item=", item)
	})

	// 새로 수집한 데이터 정리
	// 기존 json 에 집어넣기 위해,
	// array map 을 map 으로 변환해야 함.

	for i, p := range newData {

		for k, v := range p {
			fmt.Println("NewData [", i, "] = ", k, v)
			if _, ok := (*js)[k]; ok {
				fmt.Println("Found in previous json")

				// 기존 데이터인데,  이미 알람을 보냈으면,  알람 비대상으로 상태가 변경되도 변경하지 않음.
				// 1발송 원칙
				if (*js)[k].NotiYN == true && v.NotiYN == false {
					v.NotiYN = true
				}
				// 기존 데이터인데,  알람 비대상에서 -> 알람 대상으로 바뀐경우
				// 기존에 알람을 보내지 않았을 경우에만 알람
				if (*js)[k].NotiYN == false && v.NotiYN == true {
					// 알람처리
					teleSend(url, k, v)
				}
				(*js)[k] = v

			} else {
				fmt.Println("Newdata")

				// 새로운 데이터인데 알람을 보내지 않았다면 알람전송
				if v.NotiYN == true {
					// 알람처리
					teleSend(url, k, v)
				}
				(*js)[k] = v
			}

		}
	}

	// Case3. Bytes -> write File
	marshalbytes, _ := json.Marshal(js)

	err = ioutil.WriteFile(fn, marshalbytes, 0644)
	if err != nil {
		log.Print("Error is = ", err)
		os.Exit(1)
	}
	c <- true

}

// Scrape PPomppu Main Page
func main() {
	// Goroutine Channel
	c := make(chan bool)

	// 뽐게 핫딜 프로세싱
	var ppomURL string = "http://www.ppomppu.co.kr/zboard/zboard.php?id=ppomppu"
	go getPrices(ppomURL, "kor-product.json", c)
	// go getPrices(ppomURL, "/docker/goapps/kor-product.json", c) // for docker

	fmt.Println("Korea Deal Result = ", <-c)

	// 해외 핫딜 프로세싱
	var overseaURL string = "http://www.ppomppu.co.kr/zboard/zboard.php?id=ppomppu4"
	go getPrices(overseaURL, "oversea-product.json", c)
	// go getPrices(overseaURL, "/docker/goapps/oversea-product.json", c) // for docker
	fmt.Println("Oversea Deal Result = ", <-c)

	// 뽐게 특정상품 검출

	// 해외 특정상품 검출

}

func teleSend(url string, id string, v ProductInfo) {
	// Read Tele info jsonfile
	// b, err := ioutil.ReadFile("/docker/goapps/telegram.json")  // for docker
	b, err := ioutil.ReadFile("./telegram.json")
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

	cpm := fmt.Sprint(v.Cpm)
	txt := v.ProdNm + "\n\n" + "CPM (Click Per Minutes) : " + cpm + "\n\n" + url + "&no=" + id

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

func encodeToEUCKR(s string) (string, error) {
	var buf bytes.Buffer
	wr := transform.NewWriter(&buf, korean.EUCKR.NewEncoder())
	defer wr.Close()

	_, err := wr.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func encodeFromEUCKR(s string) (string, error) {
	var buf bytes.Buffer
	wr := transform.NewWriter(&buf, korean.EUCKR.NewDecoder())
	defer wr.Close()

	_, err := wr.Write([]byte(s))
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func timeString() string {
	now := time.Now()
	return now.Format("2006-01-02 15:04:05")
	// return now.Format(time.RFC3339)
}
