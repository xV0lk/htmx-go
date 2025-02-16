package models

type Client struct {
	ID            int    `db:"id,omitempty"`
	Name          string `db:"name"`
	Phone         string `db:"phone"`
	Email         string `db:"email,unique" json:"studio_id"`
	Notifications bool   `db:"notifications" json:"notifications"`
	StudioID      int    `db:"studio_id" json:"studio_id"`
}
