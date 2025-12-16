package info

type AccountType string

const (
	AppDeveloper  AccountType = "APP_DEVELOPER"
	DeveloperTest AccountType = "DEVELOPER_TEST"
	Sandbox       AccountType = "SANDBOX"
	Standard      AccountType = "STANDARD"
)

type AccountDetails struct {
	PortalID              int         `json:"portalId"`
	AccountType           AccountType `json:"accountType"`
	TimeZone              string      `json:"timeZone"`
	CompanyCurrency       string      `json:"companyCurrency"`
	AdditionalCurrencies  []string    `json:"additionalCurrencies"`
	UTCOffset             string      `json:"utcOffset"`
	UTCOffsetMilliseconds int         `json:"utcOffsetMilliseconds"`
	UIDomain              string      `json:"uiDomain"`
	DataHostingLocation   string      `json:"dataHostingLocation"`
}

type PrivateAppAPIUsage struct {
	Name         string `json:"name"`
	UsageLimit   int    `json:"usageLimit"`
	CurrentUsage int    `json:"currentUsage"`
	CollectedAt  string `json:"collectedAt"`
	FetchStatus  string `json:"fetchStatus"`
	ResetsAt     string `json:"resetsAt"`
}
