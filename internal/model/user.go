package model

type User struct {
	Id        string `json:"id"`
	IpAddress string `json:"ipaddress"`
	Refresh   string `json:"refresh"`
}
