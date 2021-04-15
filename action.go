package cdp

import (
	"context"
	"math"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Action struct{}

func NewAction() *Action {
	return &Action{}
}

func (self *Action) SetCookies(cookies []*network.CookieParam) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		if err := network.SetCookies(cookies).Do(ctx); err != nil {
			return err
		}

		return nil
	})
}

func (self *Action) FullScreen(quality int64, buf *[]byte) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		// get layout metrics
		_, _, cssContentSize, err := page.GetLayoutMetrics().Do(ctx)
		if err != nil {
			return err
		}

		width, height := int64(math.Ceil(cssContentSize.Width)), int64(math.Ceil(cssContentSize.Height))

		// force viewport emulation
		err = emulation.SetDeviceMetricsOverride(width, height, 1, false).
			WithScreenOrientation(&emulation.ScreenOrientation{
				Type:  emulation.OrientationTypePortraitPrimary,
				Angle: 0,
			}).
			Do(ctx)
		if err != nil {
			return err
		}

		// capture screenshot
		*buf, err = page.CaptureScreenshot().
			WithQuality(quality).
			WithClip(&page.Viewport{
				X:      cssContentSize.X,
				Y:      cssContentSize.Y,
				Width:  cssContentSize.Width,
				Height: cssContentSize.Height,
				Scale:  1,
			}).Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
