// Copyright (c) 2015 Mattermost, Inc. All Rights Reserved.
// See License.txt for license information.

package api

import (
	"net/http"

	l4g "github.com/alecthomas/log4go"
	"github.com/gorilla/websocket"
	"github.com/mattermost/platform/app"
	"github.com/mattermost/platform/model"
	"github.com/mattermost/platform/utils"
)

func InitWebSocket() {
	l4g.Debug(utils.T("api.web_socket.init.debug"))
	BaseRoutes.Users.Handle("/websocket", ApiAppHandlerTrustRequester(connect)).Methods("GET")
	app.HubStart()
}

func connect(c *Context, w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  model.SOCKET_MAX_MESSAGE_SIZE_KB,
		WriteBufferSize: model.SOCKET_MAX_MESSAGE_SIZE_KB,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		l4g.Error(utils.T("api.web_socket.connect.error"), err)
		c.Err = model.NewLocAppError("connect", "api.web_socket.connect.upgrade.app_error", nil, "")
		return
	}

	wc := app.NewWebConn(ws, c.Session, c.T, c.Locale)
	app.HubRegister(wc)
	go wc.WritePump()
	wc.ReadPump()
}
