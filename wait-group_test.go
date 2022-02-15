package wgx

import (
	"fmt"
	"testing"
)

func TestWaitGroup(t *testing.T) {
	wg := NewWaitGroup()

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
						"err":   fmt.Errorf("db err"),
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
