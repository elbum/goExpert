package scrapper

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractedJob struct {
	id       string
	title    string
	location string
	salary   string
	summary  string
}

// Scrape Indeed by the term
func Scrape(term string) {
	var baseURL string = "https://kr.indeed.com/취업?q=" + term + "&limit=50"
	var jobs []extractedJob
	c := make(chan []extractedJob)
	pages := getPages(baseURL)
	fmt.Println("Total Pages = ", pages)

	for i := 0; i < pages; i++ {
		go getPage(i, baseURL, c)
		// jobs = append(jobs, extractedJobs...)

	}

	for i := 0; i < pages; i++ {
		extractedJobs := <-c
		jobs = append(jobs, extractedJobs...)
	}

	writeJobs(jobs)
	fmt.Println("Job Extracting Done : ", len(jobs))
}

func writeJobs(jobs []extractedJob) {
	file, err := os.Create("jobs.csv")
	checkErr(err)

	w := csv.NewWriter(file)
	defer w.Flush()

	headers := []string{"Link", "Title", "Location", "Salary", "Summary"}
	wErr := w.Write(headers)
	checkErr(wErr)

	for _, job := range jobs {
		jobSlice := []string{"https://kr.indeed.com/viewjob?jk=" + job.id, job.title, job.location, job.salary, job.summary}
		wErr := w.Write(jobSlice)
		checkErr(wErr)

	}

}

func getPage(page int, url string, mainC chan<- []extractedJob) {
	var jobs []extractedJob
	c := make(chan extractedJob)

	pageUrl := url + "&start=" + strconv.Itoa(page*50)
	fmt.Println("Requesting : ", pageUrl)

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	fmt.Println(doc)

	searchCards := doc.Find(".jobsearch-SerpJobCard")
	searchCards.Each(func(i int, s *goquery.Selection) {
		go extractJob(s, c)
	})

	for i := 0; i < searchCards.Length(); i++ {
		// job :=  <- c
		jobs = append(jobs, <-c)
	}
	mainC <- jobs
}

func extractJob(s *goquery.Selection, c chan<- extractedJob) {
	fmt.Println("JOB CARDS = ", s.Find("a").Length())
	id, _ := s.Attr("data-jk")
	title := CleanString(s.Find(".title>a").Text())
	location := CleanString(s.Find(".sjcl").Text())
	salary := CleanString(s.Find(".salaryText").Text())
	summary := CleanString(s.Find(".summary").Text())
	fmt.Println(id, title, location, salary, summary)
	c <- extractedJob{
		id:       id,
		title:    title,
		location: location,
		salary:   salary,
		summary:  summary,
	}
}

func getPages(url string) int {

	pages := 0

	res, err := http.Get(url)
	checkErr(err)
	checkCode(res)

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)
	fmt.Println(doc)

	doc.Find(".pagination").Each(func(i int, s *goquery.Selection) {
		pages = s.Find("a").Length()
	})
	return pages

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
