package wgx

import (
	"fmt"
	"sync"
)

const (
	AllFail         = 1 // 执行结果全失败
	AllSuccess      = 2 // 执行结果全成功
	FailSuccessWith = 3 // 执行结果失败和成功都有
)

type wgResult struct {
	wg *sync.WaitGroup // 任务等待使用了官方的wg

	sceneContext *sceneContext // 场景上下文存储, 仅会对DoneWithSceneArgs调用时，将链路存储起来。scene隔离,只存放最新一条

	l            *sync.RWMutex
	resultState  int // 任务执行完后，会在Wait()调用时，计算结果
	total        int // 总计数
	successCount int // 每次执行Done/DoneSuccess时,会计算加1
	failCount    int // 每次调用doneFail时，会计数加1
	addTimes     int // 每次调用Add时, 会计数加1
	doneTimes    int // 每次调用Done时，计数会加1
}

func newWgReulst() wgResult {
	return wgResult{
		l:            &sync.RWMutex{},
		wg:           &sync.WaitGroup{},
		sceneContext: newSceneContext(),
	}
}

func (r *wgResult) countResultState() {
	r.l.Lock()
	defer r.l.Unlock()

	if r.total == r.successCount {
		r.resultState = AllSuccess
		return
	}

	if r.total == r.failCount {
		r.resultState = AllFail
		return
	}
	r.resultState = FailSuccessWith

	return
}

func (r *wgResult) Add(delta int) {
	if delta == 0 {
		return
	}

	if delta < 0 {
		panic(fmt.Errorf("wg.add only accept delta>0"))
	}

	r.l.Lock()
	defer r.l.Unlock()

	r.total += delta
	r.addTimes += delta

	r.wg.Add(delta)
}

func (r *wgResult) Success() {
	r.l.Lock()
	defer r.l.Unlock()
	r.successCount ++

	r.doneTimes ++

	r.wg.Done()
}

type Result struct {
	ResultState  int           `json:"result_state"`  // 执行结果. 1全部失败,2全部成功，3混合
	TotalCount   int           `json:"total_count"`   // 总个数, 调用wg.Add(delta)的累计值
	SuccessCount int           `json:"success_count"` // 成功数量, 调用wg.Done/wg.DoneSuccess时计数+1
	FailCount    int           `json:"fail_count"`    // 失败个数, 调用wg.Fail/wg.FailWithSceneArgs计数+1
	Panels       []ReasonPanel `json:"panels"`        // 结果面板, 会按照失败场景，对失败原因进行上下文结果.
}

func (r *wgResult) Wait() Result {
	r.wg.Wait()
	r.countResultState()
	return Result{
		ResultState:  r.resultState,
		Panels:       r.sceneContext.ReasonPanel(),
		TotalCount:   r.total,
		SuccessCount: r.successCount,
		FailCount:    r.failCount,
	}
}

func (r *wgResult) FailWithSceneArgs(scene string, m map[string]interface{}) {
	r.failWithSceneArgsDepth(scene, 4, m)
}

func (r *wgResult) FailWithErr(scene string, e error) {
	r.failWithSceneArgsDepth(scene, 4, map[string]interface{}{
		"error": e.Error(),
	})
}

func (r *wgResult) failWithSceneArgsDepth(scene string, depth int, m map[string]interface{}) {
	caller := hereWithDepth(depth)

	m["location"] = caller

	r.fail()
	r.sceneContext.addArgs(scene, m)
}

func (r *wgResult) Fail() {
	r.failWithSceneArgsDepth(hereWithDepth(3), 4, map[string]interface{}{
		"location": hereWithDepth(2),
	})
}

func (r *wgResult) fail() {
	r.l.Lock()
	defer r.l.Unlock()
	r.failCount ++

	r.doneTimes ++

	r.wg.Done()
}
