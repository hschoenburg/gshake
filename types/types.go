package types

type NameNotif struct {
  Name string   `json:"name"`
  Contact string  `json:"contact"`
  Week int      `json:"week"`
  Notified bool `json:"notified"`
  Verified bool `json:"verified"`
}

