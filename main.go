package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/chromedp/chromedp"
)

var (
	concurrency int
	urls        bool
	outputFile  string
)

const (
	//InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor = "\033[1;36m%s\033[0m"
	//WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor = "\033[1;31m%s\033[0m"
	//DebugColor   = "\033[0;36m%s\033[0m"
)

func banner() {
	fmt.Println(`
	 _____           _        _____               
	|  __ \         | |      |   __ \         
	| |__) _ __ ___ | |_ ___ |  |__) |
	|  ___| '__/ _ \| __/ _ \|  ____ /
	| |   | | | (_) | || (_) | |
	|_|   |_|  \___/ \__\___||_|
													
				-@androidh4ppi3r							
	`)
}

func main() {
	banner()
	flag.IntVar(&concurrency, "c", 10, "Set Concurrency ")
	flag.BoolVar(&urls, "u", false, "Scan Urls ")
	flag.StringVar(&outputFile, "o", "", "Save Result to OutputFile")
	flag.Parse()

	if outputFile != "" {
		emptyFile, err := os.Create(outputFile)
		if err != nil {
			log.Fatal(err)
		}
		//log.Println(emptyFile)
		emptyFile.Close()
		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				ProtoP()
				wg.Done()
			}()
			wg.Wait()
		}

	} else {
		var wg sync.WaitGroup
		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				ProtoP()
				wg.Done()
			}()
			wg.Wait()
		}
	}
}

// ProtoP scans
func ProtoP() {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		// create context
		url := sc.Text()
		//fmt.Println(url)
		ctx, cancel := chromedp.NewContext(context.Background())

		// run task list
		var res string
		if urls == true {
			err := chromedp.Run(ctx,
				chromedp.Navigate(url+"&__proto__[test]=test"),
				chromedp.Evaluate(`window.test`, &res),
			)
			cancel()
			if err != nil {
				log.Printf(ErrorColor, url+" [Not Vulnerable]")
				continue
			}
		} else {
			err := chromedp.Run(ctx,
				chromedp.Navigate(url+"/"+"?__proto__[test]=test"),
				chromedp.Evaluate(`window.test`, &res),
			)
			cancel()
			if err != nil {
				log.Printf(ErrorColor, url+" [Not Vulnerable]")
				continue
			}
		}
		if outputFile != "" {
			f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_WRONLY, 0644)
			if err != nil {
				log.Println(err)
			}
			if _, err := f.WriteString(url + "\n"); err != nil {
				log.Fatal(err)
			}
			f.Close()
		}
		log.Printf(NoticeColor, url+" [Vulnerable]")
	}
}
