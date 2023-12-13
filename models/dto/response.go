package dto

import "encoding/json"

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}

func (resp *Response) Wrap(result string, data json.RawMessage, err error) {
	if err == nil {
		resp.Error = ""
	} else {
		resp.Error = err.Error()
	}
	resp.Data = data
	resp.Result = result
}
