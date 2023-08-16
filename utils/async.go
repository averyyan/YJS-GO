package utils

import (
	"fmt"
	"runtime"
	"sync"
)

// PanicBufLen 是出异常后堆栈的打印缓冲长度
var PanicBufLen = 1024

func SafeAsyncExe(exeFunc func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("error:%+v", r)
		}
	}()

	exeFunc()
}

// GoAndWait 提供安全并发能力，每一个函数启动一个协程处理，如果其中有一个panic或者出现异常了，整体退出，并返回第最先到达的error,
// panic异常会打印异常堆栈的前`PanicBufLen`个字节
// 如果全部执行正常则只返回nil
func GoAndWait(handlers ...func() error) error {
	var (
		wg   sync.WaitGroup
		once sync.Once
		err  error
	)
	for _, f := range handlers {
		wg.Add(1)
		go func(handler func() error) {
			defer func() {
				if e := recover(); e != nil {
					buf := make([]byte, PanicBufLen)
					buf = buf[:runtime.Stack(buf, false)]
					fmt.Printf("[PANIC]%v\n%s\n", e, buf)
				}
				wg.Done()
			}()
			if e := handler(); e != nil {
				once.Do(func() {
					err = e
				})
			}
		}(f)
	}
	wg.Wait()
	return err
}
