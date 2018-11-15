package models

// Job describes a company that a formula might be associated with
type Job struct {
	TableName   struct{} `sql:"jobs" json:"-"`
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	ContactName string   `json:"contact_name"`
	ContactInfo string   `json:"contact_info"`
	Street      string   `json:"street"`
	Street2     string   `json:"street2"`
	City        string   `json:"city"`
	State       string   `json:"state"`
	Zipcode     string   `json:"zipcode"`
	Notes       string   `json:"notes"`
	Created     int64    `json:"created"`
	Modified    int64    `json:"modified"`
}
