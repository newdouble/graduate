package fetchactdetailctr

import (
	"encoding/json"
	"graduate/src/activity"
	"graduate/src/session"
	"io/ioutil"
	"net/http"
)

type FetchActDetailController struct {}

func (c *FetchActDetailController) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	var act *activity.Activity
	var res []byte
	var actID string
	var err error

	 s:= session.New()
	reqBody, _ := c.PreHandle(r, s)
	actID = r.Form.Get("actid")
	act = activity.Get(actID)
	res, err = act.QueryDetails(s, reqBody)
	if err != nil {
		goto RETURN
	}

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

func (c *FetchActDetailController) PreHandle(r *http.Request, s *session.Session) ([]byte, *error) {
	requestBody, _ := ioutil.ReadAll(r.Body)
	r.ParseForm()
	return requestBody,nil
}