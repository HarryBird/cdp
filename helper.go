package cdp

import (
	"github.com/chromedp/chromedp"
)

func InnerHTML(url string, sel string, buf *string, opts ...chromedp.QueryOption) error {
	return NewCDP().
		WithBrowserLog().
		WithAction(chromedp.Navigate(url)).
		WithAction(chromedp.InnerHTML(sel, buf, opts...)).
		Run()
}
