package wgx

import "sync"

type sceneContext struct {
	l *sync.RWMutex
	m map[string]map[string]interface{}
}

func newSceneContext() *sceneContext {
	return &sceneContext{
		l: &sync.RWMutex{},
		m: make(map[string]map[string]interface{}),
	}
}

func (sc *sceneContext) addArgs(scene string, args map[string]interface{}) {
	sc.l.Lock()
	defer sc.l.Unlock()

	sc.m[scene] = args
}

func (sc sceneContext) ReasonPanel() []ReasonPanel {
	sc.l.RLock()
	defer sc.l.RUnlock()

	var rs = make([]ReasonPanel, 0, 10)

	for k, v := range sc.m {
		rs = append(rs, ReasonPanel{
			Scene: k,
			Args:  v,
		})
	}
	return rs
}
