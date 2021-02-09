package channel

import (
	"encoding/json"
	"fmt"
)

// VimWrap wraps the envelope in a required Vim
// specific syntax.
// see: https://vimhelp.org/channel.txt.html#channel-use
type VimWrap [2]json.RawMessage

// Envelope wraps a json-encoded rpc message.
//
// When Vim receives an Envelope it will route
// the message to the appropriate vim-script handler.
//
// When the Channel receives an Envelope it
// will be delivered to the correct mailbox.
//
// The Envelope's Err field must be checked for
// any underlying tcp errors before assuming
// the Body is a valid json response.
type Envelope struct {
	Mailbox uint32          `json:"mailbox"`
	RPC     string          `json:"rpc"`
	Body    json.RawMessage `json:"body"`
	ReqNum  int             `jso:"request_number"`
	Err     error           `json:"error"`
	In      bool            `json:"-"`
}

func (e Envelope) ToVim() (VimWrap, error) {
	vw := VimWrap{}
	var err error
	vw[0], err = json.Marshal(0)
	if err != nil {
		return vw, err
	}
	vw[1], err = json.Marshal(e)
	if err != nil {
		return vw, err
	}
	return vw, nil
}

func (e *Envelope) FromVim(vw VimWrap) error {
	err := json.Unmarshal(vw[1], e)
	if err != nil {
		return fmt.Errorf("failed to unmarshal envelope: %v", err)
	}
	var reqNum int
	err = json.Unmarshal(vw[0], &reqNum)
	if err != nil {
		return fmt.Errorf("failed to unmarshal req number: %v", err)
	}
	e.ReqNum = reqNum
	return nil
}
