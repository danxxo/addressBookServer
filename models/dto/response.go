package dto

import "encoding/json"

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data"`
	Error  string          `json:"error"`
}
