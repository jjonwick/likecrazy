package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

var searchText = "Jimin - Like Crazy (English Version)"

func main() {

	startUrl := "https://linktr.ee/requestbtsradioalerts"

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", true),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	var nodes []*cdp.Node
	selector := "a.sc-pFZIQ"

	// navigate to a page, wait for an element, click
	err := chromedp.Run(ctx,
		chromedp.Navigate(startUrl),
		chromedp.WaitReady(selector),
		chromedp.Nodes(selector, &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		log.Fatal(err)
	}

	for i, node := range nodes {
		u := node.AttributeValue("href")
		log.Printf("%d: %s", i, u)
	}

	startIdx := 0

	// map string to bool
	var processed map[string]bool = make(map[string]bool)

	for i, node := range nodes {

		u := node.AttributeValue("href")

		// if u is already processed, skip

		if processed[u] {
			log.Printf("already processed %s", u)
			continue
		}
		processed[u] = true

		log.Printf("%d: %s", i, u)

		if i < startIdx {
			continue
		}

		if strings.Contains(u, "twitter.com") {
			continue
		}

		if strings.Contains(u, "d1gm7n6w0pishx.cloudfront.net") {
			ActionOnPage(ctx, u)
		}
	}
}

func ActionOnPage(ctx context.Context, url string) {

	clone, cancel := chromedp.NewContext(ctx)
	defer cancel()
	err := chromedp.Run(clone,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.SetAttributes("//*[@id='sub-app-Alerts']/section",
			map[string]string{"style": "display: block;"}, chromedp.BySearch),
		chromedp.Sleep(2*time.Second),
		chromedp.SendKeys(`#sub-app-Alerts > section > div.song-search > div.sub-app-color-border.search-dropdown > div.input-group > div.input > input[type=text]`, searchText, chromedp.BySearch),
		chromedp.WaitVisible("div[data-action='request']"),
		chromedp.Click("div[data-action='request']"),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}

func ActionOnPage2(ctx context.Context, url string) {
	log.Printf("ActionOnPage2 href: %s", url)

	selector := "#voting-section > div.song-search > div.sub-app-color-border.search-dropdown > div.input-group > div.input > input[type=text]"

	selector = `Artist and Song Search`

	clone, cancel := chromedp.NewContext(ctx)
	defer cancel()
	err := chromedp.Run(clone,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.Sleep(5*time.Second),
		chromedp.SendKeys(selector, searchText, chromedp.BySearch),
		chromedp.WaitVisible("div[data-action='request']"),
		chromedp.Click("div[data-action='request']"),
		chromedp.Sleep(3*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}
}
