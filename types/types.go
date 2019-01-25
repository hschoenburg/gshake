package types

type NameNotif struct {
  Name string   `json:"name"`
  Contact string  `json:"contact"`
  Week int      `json:"week"`
  Notified bool `json:"notified"`
}

type Contact struct {
  Email string `json:"email"`
  Phone string `json:"phone"`
  Subscribed bool `json:subscribed`
  Verified bool `json:verified`
}

  
type NameInfoData struct {
  Name string
	Result struct {
	  Info struct {} 		`json:"info"`
		Start StartData   `json:"start"`
  }
  Error struct {
    Message string  `json:"message"`
    Code int        `json:"code"`
  }                 `json:"error"`
}

type StartData struct {
		Reserved bool `json:"reserved"`
		Week     int  `json:"week"`
		Start    int  `json:"start"`
}
