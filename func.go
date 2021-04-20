package cdp

import (
	"time"

	"github.com/chromedp/cdproto/cdp"
)

func GetCookieExpireFromInt(ts int) cdp.TimeSinceEpoch {
	return cdp.TimeSinceEpoch(time.Now().Add(time.Duration(ts) * time.Second))
}
