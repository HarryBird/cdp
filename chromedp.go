// Package cdp 对 github.com/chromedp/chromedp 做一层常用封装
package cdp

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

type CDP struct {
	// timeout 整体超时时长，默认为10秒
	timeout time.Duration
	// tasks chromedp task列表
	tasks []chromedp.Action
	// execOpts chrome 启动选项列表
	execOpts []chromedp.ExecAllocatorOption
	// ctxOpts context 选项列表
	ctxOpts []chromedp.ContextOption
}

// Run 启动并执行chrome
func (self *CDP) Run(actions ...chromedp.Action) error {

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), self.getExecOptions()...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx, self.ctxOpts...)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, self.timeout)
	defer cancel()

	tasks := append(self.tasks, actions...)
	return chromedp.Run(ctx, tasks...)
}

// getExecOptions 获取启动chrome的选项列表
// 这里做了特殊处理，从默认的启动选项里摘除了默认的Headless选项，是否Headless由对象的self.headless来控制
// 最终会把自定义的option列表尾缀在默认列表后面
func (self *CDP) getExecOptions() []chromedp.ExecAllocatorOption {
	return append(
		append(chromedp.DefaultExecAllocatorOptions[0:2], chromedp.DefaultExecAllocatorOptions[3:]...),
		self.execOpts...,
	)
}

// WithDevice 设置设备类型
func (self *CDP) WithDevice(dev chromedp.Device) *CDP {
	self.tasks = append(self.tasks, chromedp.Emulate(dev))
	return self
}

// WithViewport 设置展示区域的尺寸
func (self *CDP) WithViewport(width, height int64) *CDP {
	self.tasks = append(self.tasks, chromedp.EmulateViewport(width, height))
	return self
}

// WithSleep 设置chromedp在渲染过程中的睡眠时间
// 通常用于等待静态资源加载的过程中
func (self *CDP) WithSleep(sl time.Duration) *CDP {
	self.tasks = append(self.tasks, chromedp.Sleep(sl))
	return self
}

// WithAction 添加自定义Action
func (self *CDP) WithAction(act chromedp.Action) *CDP {
	self.tasks = append(self.tasks, act)
	return self
}

// WithActionFunc 接受一个func，用户自定义丰富的Action过程
func (self *CDP) WithActionFunc(f chromedp.ActionFunc) *CDP {
	self.tasks = append(self.tasks, f)
	return self
}

// WithoutHeadless 设置关闭headless
// 设置了这个选项后，chromedp在执行过程中会启动chrome的GUI，用于观察chrome的行为
func (self *CDP) WithoutHeadless() *CDP {
	self.execOpts = append(self.execOpts, chromedp.Headless)
	return self
}

// WithWindowSize 设置浏览器窗口的宽高
func (self *CDP) WithWindowSize(width, height int) *CDP {
	self.execOpts = append(self.execOpts, chromedp.WindowSize(width, height))
	return self
}

// WithChromePath 设置chrome的执行路径
func (self *CDP) WithChromePath(path string) *CDP {
	self.execOpts = append(self.execOpts, chromedp.ExecPath(path))
	return self
}

// WithUserAgent 设置UserAgent
func (self *CDP) WithUserAgent(ua string) *CDP {
	self.execOpts = append(self.execOpts, chromedp.UserAgent(ua))
	return self
}

// WithEnv 设置启动时的环境变量
func (self *CDP) WithEnv(vars ...string) *CDP {
	self.execOpts = append(self.execOpts, chromedp.Env(vars...))
	return self
}

// WithFlag 设置启动时的flag
func (self *CDP) WithFlag(name string, value interface{}) *CDP {
	self.execOpts = append(self.execOpts, chromedp.Flag(name, value))
	return self
}

// WithTimeout 设置整体执行时长
func (self *CDP) WithTimeout(t time.Duration) *CDP {
	self.timeout = t
	return self
}

// WithBrowserLog 设置启动浏览器日志
func (self *CDP) WithBrowserDebugLog(f func(string, ...interface{})) *CDP {
	self.ctxOpts = append(self.ctxOpts, chromedp.WithDebugf(f))
	return self
}

func (self *CDP) WithBrowserErrorLog(f func(string, ...interface{})) *CDP {
	self.ctxOpts = append(self.ctxOpts, chromedp.WithErrorf(f))
	return self
}

func (self *CDP) WithBrowserInfoLog(f func(string, ...interface{})) *CDP {
	self.ctxOpts = append(self.ctxOpts, chromedp.WithLogf(f))
	return self
}

func NewCDP() *CDP {
	return &CDP{
		tasks:    []chromedp.Action{},
		execOpts: []chromedp.ExecAllocatorOption{},
		ctxOpts:  []chromedp.ContextOption{},
		timeout:  10 * time.Second,
	}
}
