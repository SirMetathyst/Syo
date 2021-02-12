package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	StrHttp    = "https://"
	SyosetuURL = ".syosetu.com"
	MaxWorkers = flag.Int("w", 1, "number of workers")
	Ids        = flag.String("n", "", "ids of novels")
	Overwrite  = flag.Bool("o", false, "overwrite existing content")
	Limit      = flag.Duration("l", 200*time.Millisecond, "limit requests to syosetu")
	SyoType    = flag.String("s", "ncode", "which syosetu site")
)

// DownloadPage ...
func DownloadPage(id string) {
	log.Printf("Downloading: %s\n", id)
	client := &http.Client{}
	req, err := http.NewRequest("GET", StrHttp+*SyoType+SyosetuURL+id, nil)
	if err != nil {
		log.Println(err)
	}
	
	req.AddCookie(&http.Cookie{Name: "over18", Value: "yes", Domain: ".syosetu.com"})
	
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	title := doc.Find(".novel_subtitle").Text()
	content := doc.Find(".novel_view").Text()
	b := bytes.Buffer{}
	b.WriteString(title)
	b.WriteString("\n\n")
	b.WriteString(content)
	WriteFile("./wn"+id[:len(id)-1]+".txt", b.Bytes(), 0666)
}

// PageLink ...
type PageLink struct {
	URL        string
	Title      string
	LongUpdate string
}

// Download ...
func Download(id string, jobs chan Job, wg *sync.WaitGroup) {
	url := StrHttp+*SyoType+SyosetuURL + "/" + id
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	
	req.AddCookie(&http.Cookie{Name: "over18", Value: "yes", Domain: ".syosetu.com"})
	
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	title := doc.Find("title").Text()
	description := doc.Find("#novel_ex").Text()
	pageLinks := make([]PageLink, 0)
	doc.Find(".novel_sublist2").Each(func(i int, s *goquery.Selection) {
		link := s.Find("a")
		linkTitle := link.Text()
		url, _ := link.Attr("href")
		longUpdate := s.Find(".long_update").Text()
		pageLink := PageLink{
			URL:        url,
			Title:      linkTitle,
			LongUpdate: longUpdate,
		}
		pageLinks = append(pageLinks, pageLink)
	})
	nPath := "./wn/" + id
	if MkDirAll(nPath, 0666) {
		WriteFile(nPath+"/source.txt", []byte(url), 0666)
		WriteFile(nPath+"/title.txt", []byte(title), 0666)
		WriteFile(nPath+"/desc.txt", []byte(description), 0666)
		b := bytes.Buffer{}
		b.WriteString("Timestamp, Revised, URL, Title")
		for _, pageLink := range pageLinks {
			b.WriteRune('\n')
			longUpdate := strings.Replace(pageLink.LongUpdate, "\n", "", -1)
			pageURL := StrHttp+*SyoType+SyosetuURL + strings.Replace(pageLink.URL, "\n", "", -1)
			title := strings.Replace(pageLink.Title, "\n", "", -1)
			b.WriteString(longUpdate[:16])
			b.WriteString(", ")
			if len(longUpdate) > 17 {
				b.WriteString("yes")
			} else {
				b.WriteString("no")
			}
			b.WriteString(", ")
			b.WriteString(pageURL)
			b.WriteString(", ")
			b.WriteString(title)
		}
		WriteFile(nPath+"/toc.csv", b.Bytes(), 0666)
	}
	for _, pageLink := range pageLinks {
		log.Printf("Found: %s\n", pageLink.URL)
		url := pageLink.URL
		f := func() {
			DownloadPage(url)
			wg.Done()
		}
		wg.Add(1)
		jobs <- Job{ID: url, f: f}
	}
}

// WriteFile ...
func WriteFile(filename string, data []byte, perm os.FileMode) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) || *Overwrite {
		err = ioutil.WriteFile(filename, data, perm)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	log.Println(err)
	return false
}

// MkDirAll ...
func MkDirAll(path string, perm os.FileMode) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, perm)
		if err != nil {
			log.Println(err)
			return false
		}
		return true
	}
	return true
}

// Job ...
type Job struct {
	ID string
	f  func()
}

// Worker ...
func Worker(id int, jobs <-chan Job, limiter <-chan time.Time) {
	for job := range jobs {
		<-limiter
		jobID := job.ID
		jobFunc := job.f
		log.Printf("Worker(%d): Job: %q\n", id, jobID)
		jobFunc()
	}
}

// Run ...
func Run(wm int, ids ...string) {
	wg := sync.WaitGroup{}
	if wm <= 0 || len(ids) <= 0 {
		return
	}
	MkDirAll("./wn", 0666)
	limiter := time.Tick(*Limit)
	jobs := make(chan Job)
	for w := 1; w <= wm; w++ {
		log.Printf("Worker(%d): Created\n", w)
		go Worker(w, jobs, limiter)
	}
	for _, id := range ids {
		log.Printf("Downloading: /%s\n", id)
		Download(id, jobs, &wg)
	}
	wg.Wait()
	close(jobs)
}

func main() {
	flag.Parse()
	splitIds := strings.Split(*Ids, ",")
	Run(*MaxWorkers, splitIds...)
}
