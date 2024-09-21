package models

type (
	WsState struct {
		State    int8   `form:"state" json:"state"` // 0: stop; 1: running; 2: finish
		SentByte int64  `form:"sent_byte" json:"sent_byte"`
		Message  string `form:"message" json:"message"`
	}

	VueCmd struct {
		Command string `form:"command" json:"command"`
	}
)
