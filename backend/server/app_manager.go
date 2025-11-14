
// ============================
// backend/server/app_manager.go
// ============================
package server

import (
	"encoding/json"
	"time"
)

type AppManager struct {
	bus     *BusServer
	vfs     *VFSService
	sess    *SessionManager
	running map[string]*AppInstance
}

type AppInstance struct {
	AppID     string
	PID       string
	Windows   []WindowState
	Minimized bool
	Dirty     bool
}

type WindowState struct {
	ID    string
	X, Y  int
	W, H  int
	Z     int64
	State string // normal/minimized/maximized
}

func NewAppManager(bus *BusServer, vfs *VFSService, sess *SessionManager) *AppManager {
	am := &AppManager{bus: bus, vfs: vfs, sess: sess, running: make(map[string]*AppInstance)}
	bus.Subscribe("app:launch", am.handleLaunch)
	bus.Subscribe("app:close", am.handleClose)
	bus.Subscribe("app:minimize", am.handleMinimize)
	bus.Subscribe("app:focus", am.handleFocus)
	bus.Subscribe("app:save-response", am.handleAppSaveResponse)
	return am
}

func (am *AppManager) handleLaunch(msg *Message) {
	manifest := msg.Payload["manifest"].(map[string]interface{})
	//TODO: Auth checks

	pid := "pid-placeholder" // In a real scenario, we would start a process and get a real PID

	appId := manifest["id"].(string)
	win := WindowState{ID: "win-" + genUUID(), X: 100, Y: 80, W: 800, H: 600, Z: time.Now().UnixMilli(), State: "normal"}
	am.running[appId] = &AppInstance{AppID: appId, PID: pid, Windows: []WindowState{win}}
	am.bus.Publish(&Message{Topic: "app:launched", Payload: map[string]interface{}{"appId": appId, "appWindowId": win.ID, "pid": pid, "geometry": win}})
}

func (am *AppManager) handleClose(msg *Message) {
	appId := msg.Payload["appId"].(string)
	inst, ok := am.running[appId]
	if !ok {
		return
	}

	if inst.Dirty {
		reply := am.bus.PublishSync("app:save-request", map[string]interface{}{"appId": appId, "reason": "close"}, 5*time.Second)
		if reply != nil {
			am.handleAppSaveResponse(reply)
		} else {
			// If the app does not respond, force close it
			am.forceClose(appId)
		}
		return
	}

	am.forceClose(appId)
}

func (am *AppManager) forceClose(appId string) {
	inst, ok := am.running[appId]
	if !ok {
		return
	}
	// am.stopProcess(inst.PID)
	delete(am.running, appId)
	am.bus.Publish(&Message{Topic: "app:closed", Payload: map[string]interface{}{"appId": appId}})
}

func (am *AppManager) handleAppSaveResponse(msg *Message) {
	appId := msg.Payload["appId"].(string)
	snapshot, hasSnapshot := msg.Payload["snapshot"]
	if hasSnapshot && snapshot != nil {
		b, _ := json.Marshal(snapshot)
		path := "/sessions/latest/apps/" + appId + "/snapshot-" + time.Now().Format(time.RFC3339) + ".json"
		am.vfs.Write(path, b)
	}
	am.forceClose(appId)
}

func (am *AppManager) handleMinimize(msg *Message) {
	appId := msg.Payload["appId"].(string)
	if inst, ok := am.running[appId]; ok {
		inst.Minimized = !inst.Minimized
		am.bus.Publish(&Message{Topic: "app:minimized", Payload: map[string]interface{}{"appId": appId, "minimized": inst.Minimized}})
	}
}

func (am *AppManager) handleFocus(msg *Message) {
	appId := msg.Payload["appId"].(string)
	if inst, ok := am.running[appId]; ok {
		for _, win := range inst.Windows {
			win.Z = time.Now().UnixMilli()
		}
		am.bus.Publish(&Message{Topic: "app:focused", Payload: map[string]interface{}{"appId": appId}})
	}
}
