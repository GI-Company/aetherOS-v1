package server

import (
	"encoding/json"
	"log"
	"time"
)

type AppManager struct {
	bus     *BusServer
	sess    *SessionManager
	running map[string]*AppSnapshot
}

func NewAppManager(bus *BusServer, sess *SessionManager) *AppManager {
	am := &AppManager{bus: bus, sess: sess, running: make(map[string]*AppSnapshot)}
	bus.SubscribeServer("app:launch", func(env *Envelope) { am.handleLaunch(env) })
	bus.SubscribeServer("app:close", func(env *Envelope) { am.handleClose(env) })
	bus.SubscribeServer("app:save-response", func(env *Envelope) { am.handleSaveResponse(env) })
	return am
}

func (am *AppManager) handleLaunch(env *Envelope) {
	manifest, _ := env.Payload["manifest"].(map[string]interface{})
	appId, _ := manifest["id"].(string)
	log.Println("launching app:", appId)
	snap := &AppSnapshot{
		AppID:     appId,
		Windows:   []WindowState{{ID: genUUID(), X: 100, Y: 60, W: 900, H: 600, Z: time.Now().UnixMilli(), State: "normal"}},
		AppState:  map[string]interface{}{},
		Dirty:     false,
		LastSaved: time.Now(),
		Version:   1,
	}
	am.running[appId] = snap
	am.bus.Publish(&Envelope{Topic: "app:launched", From: "kernel", Payload: map[string]interface{}{"appId": appId, "appWindowId": snap.Windows[0].ID, "pid": genUUID(), "geometry": snap.Windows[0]}, Time: time.Now()})
}

func (am *AppManager) handleClose(env *Envelope) {
	appId, _ := env.Payload["appId"].(string)
	if snap, ok := am.running[appId]; ok && snap.Dirty {
		am.bus.Publish(&Envelope{Topic: "app:save-request", From: "kernel", Payload: map[string]interface{}{"appId": appId, "reason": "close"}, Time: time.Now()})
		go func(a string) { time.Sleep(5 * time.Second); am.forceClose(a) }(appId)
		return
	}
	am.forceClose(appId)
}

func (am *AppManager) handleSaveResponse(env *Envelope) {
	appId, _ := env.Payload["appId"].(string)
	payload := env.Payload["snapshot"]
	if payload == nil {
		am.forceClose(appId)
		return
	}
	b, _ := json.Marshal(payload)
	am.sess.cache.Set("app:"+appId, string(b))
	am.sess.cache.Snapshot("app:" + appId)
	am.forceClose(appId)
}

func (am *AppManager) forceClose(appId string) {
	delete(am.running, appId)
	am.bus.Publish(&Envelope{Topic: "app:closed", From: "kernel", Payload: map[string]interface{}{"appId": appId}, Time: time.Now()})
}
