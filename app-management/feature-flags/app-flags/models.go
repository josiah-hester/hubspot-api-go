package appflags

type State string

const (
	Absent State = "ABSENT"
	Off    State = "OFF"
	On     State = "ON"
)

type FlagInfo struct {
	AppID         int    `json:"appId"`
	DefaultState  State  `json:"defaultState"`
	FlagName      string `json:"flagName"`
	OverrideState State  `json:"overrideState,omitempty"`
}

type FlagState struct {
	AppID     int    `json:"appId"`
	FlagName  string `json:"flagName"`
	FlagState State  `json:"flagState"`
	PortalID  int    `json:"PortalId"`
}
