package wgx

type WaitGroup struct {
	wr *wgResult
}

func NewWaitGroup() WaitGroup {
	wr := newWgReulst()
	return WaitGroup{
		wr: &wr,
	}
}

func (wg *WaitGroup) Add(delta int) {
	wg.wr.Add(delta)
}
func (wg *WaitGroup) Done() {
	wg.DoneSuccess()
}
func (wg *WaitGroup) DoneSuccess() {
	wg.wr.Success()

}

func (wg *WaitGroup) DoneFail() {
	wg.wr.Fail()
}

func (wg *WaitGroup) DoneFailWithSceneArgs(scene string, args map[string]interface{}) {
	wg.wr.FailWithSceneArgs(scene, args)
}

func (wg *WaitGroup) DoneFailWithErr(scene string, e error) {
	wg.wr.FailWithErr(scene, e)
}

func (wg *WaitGroup) Wait() Result {
	r := wg.wr.Wait()

	return r
}
