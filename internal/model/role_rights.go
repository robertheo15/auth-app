package model

type RoleRight struct {
	ID      string `json:"id"`
	Section string `json:"section"`
	Route   string `json:"route"`
	RCreate bool   `json:"r_create"`
	RRead   bool   `json:"r_read"`
	RUpdate bool   `json:"r_update"`
	RDelete bool   `json:"r_delete"`
}
