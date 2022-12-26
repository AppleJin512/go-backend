package schemas

type Error struct {
	Message          interface{} `json:"message,omitempty"`
	ValidationErrors interface{} `json:"validationErrors,omitempty"`
}
