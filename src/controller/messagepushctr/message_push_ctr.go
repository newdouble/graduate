package messagepushctr

import (
	"encoding/json"
	"graduate/src/event"
	"graduate/src/session"
	"io/ioutil"
	"net/http"
)

type MessagePushController struct {
}

func (c *MessagePushController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s := session.New()
	var err error
	var req *MessagePushReq
	var eventType EventType
	var res []byte
	eventParam := make(event.Param)
	reqBody, err := c.PreHandle(r,s)
	err = json.Unmarshal(reqBody, req)
	if err != nil {
		goto RETURN
	}
	eventParam.Add("UID", req.UID)
	if eventType == EventIDCongestionPopup {
		eventParam.Add("ShowType", 1)
	}
	event.Event[event.EventNewSelf].Notify(event.EventNewSelf, eventParam)

RETURN:
	if err != nil {
		w.WriteHeader(http.StatusOK)
		errRes := map[string]interface{} {
			"ret": 1001,
			"msg": "ddd",
		}
		res, err = json.Marshal(errRes)
		if err != nil {

		}
		w.Write(res)
	}
}

func (c *MessagePushController) PreHandle(r *http.Request, s *session.Session) ([]byte, error) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	r.ParseForm()
	return requestBody,nil
}