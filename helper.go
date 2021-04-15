package cdp

import (
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/mitchellh/mapstructure"
)

type Helper struct {
	url     string
	cookies []*network.CookieParam
	logger  func(string, ...interface{})
	debuger func(string, ...interface{})
	errorer func(string, ...interface{})
}

func NewHelper(url string) *Helper {
	return &Helper{
		url: url,
	}
}

func (self *Helper) WithInfoLogger(f func(string, ...interface{})) *Helper {
	self.logger = f
	return self
}

func (self *Helper) WithDebugLogger(f func(string, ...interface{})) *Helper {
	self.debuger = f
	return self
}

func (self *Helper) WithErrorLogger(f func(string, ...interface{})) *Helper {
	self.errorer = f
	return self
}

func (self *Helper) WithCookie(cookie map[string]interface{}) error {
	var cp network.CookieParam

	if err := mapstructure.Decode(cookie, &cp); err != nil {
		return err
	}

	self.cookies = append(self.cookies, &cp)
	return nil
}

func (self *Helper) WithCookies(cookies []map[string]interface{}) error {
	for _, cookie := range cookies {
		if err := self.WithCookie(cookie); err != nil {
			return err
		}
	}

	return nil
}

func (self *Helper) InnerHTML(sel string, buf *string, opts ...chromedp.QueryOption) error {
	return self.Init().
		WithAction(chromedp.InnerHTML(sel, buf, opts...)).
		Run()
}

func (self *Helper) Init() *CDP {
	chrome := NewCDP()

	if self.logger != nil {
		chrome.WithBrowserInfoLog(self.logger)
	}

	if self.debuger != nil {
		chrome.WithBrowserDebugLog(self.debuger)
	}

	if self.errorer != nil {
		chrome.WithBrowserErrorLog(self.errorer)
	}

	if len(self.cookies) > 0 {
		chrome.WithCookies(self.cookies)
	}

	return chrome.WithAction(chromedp.Navigate(self.url))
}
