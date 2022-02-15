## wgx

## 介绍
wgx 是对sync.WaitGroup的增强版，解决sync.WaitGroup的以下问题:

- sync.WaitGroup只确保子任务被执行，无法传递任务最终的执行结果。wgx支持三个执行结果[全失败，全成功，半失败半成功]。
- sync.WaitGroup不关注子任务的上下文细节，也不关注失败和成功的计数。wgx支持对失败的任务，进行上下文输出(失败个数, 失败原因，失败场景敲定, 失败链路补回)。

## 使用
```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/fwhezfwhez/wgx"
	"runtime"
)

func main() {
	wg := wgx.NewWaitGroup()

	for i := 0; i < 10000; i ++ {
		go func(i int) {
			wg.Add(1)
			defer func() {
				switch i % 4 {
				case 0:
					wg.DoneSuccess()
				case 1:
					wg.DoneFail()
				case 2:
					wg.DoneFailWithSceneArgs("login_log_insert", map[string]interface{}{
						"loc":   here(),
						"uname": "ft",
					})

				case 3:
					wg.DoneFailWithErr("login_log_update", fmt.Errorf("redis err"))
				}
			}()

			fmt.Println(1)
		}(i)
	}

	rs := wg.Wait()

	fmt.Println(JSON(rs))
}

func here() string {
	_, f, l, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", f, l)
}

func JSON(i interface{}) string {
	r, _ := json.MarshalIndent(i, "  ", "  ")
	return string(r)
}

```