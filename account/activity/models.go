package activity

type AuditLog struct {
	ActingUser     ActingUser `json:"actingUser"`
	Action         string     `json:"action"`
	Category       string     `json:"category"`
	ID             string     `json:"id"`
	OccurredAt     string     `json:"occurredAt"`
	SubCategory    string     `json:"subCategory"`
	TargetObjectID string     `json:"targetObjectId"`
}

type ActingUser struct {
	UserID    int    `json:"userId"`
	UserEmail string `json:"userEmail"`
}

type Paging struct {
	Next struct {
		After string `json:"after"`
		Link  string `json:"link"`
	} `json:"next"`
}

type LoginActivity struct {
	ID             string `json:"id"`
	LoginAt        string `json:"loginAt"`
	LoginSucceeded bool   `json:"loginSucceeded"`
	CountryCode    string `json:"countryCode"`
	Email          string `json:"email"`
	IPAddress      string `json:"ipAddress"`
	Location       string `json:"location"`
	RegionCode     string `json:"regionCode"`
	UserAgent      string `json:"userAgent"`
	UserID         int    `json:"userId"`
}

type SecurityHistory struct {
	CreatedAt   string `json:"createdAt"`
	ID          string `json:"id"`
	Type        string `json:"type"`
	UserID      int    `json:"userId"`
	ActingUser  string `json:"actingUser"`
	CountryCode string `json:"countryCode"`
	InfoURL     string `json:"infoUrl"`
	IPAddress   string `json:"ipAddress"`
	Location    string `json:"location"`
	ObjectID    string `json:"objectId"`
	RegeionCode string `json:"regionCode"`
}
