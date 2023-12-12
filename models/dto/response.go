package dto

import "encoding/json"

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}

// FIXME: зачем две функции для врапа? можно одну сделать. Даже сейчас функция Wrap имеет частичный фунционал функции ErrorWrap.
func (resp *Response) Wrap(result string, data json.RawMessage, err error) {
	if err == nil {
		resp.Error = ""
	} else {
		resp.Error = err.Error()
	}
	resp.Data = data
	resp.Result = result
}

func (resp *Response) ErrorWrap(err error) {
	resp.Result = "error"
	resp.Error = err.Error()
	resp.Data = nil
}
