package main

type User struct {
	Email  string `json:"email"`
	Pubkey string `json:"pubkey"`
}

type Users struct {
	Items []User
}

func (u *Users) AddItem(item User) {
	u.Items = append(u.Items, item)
}

type KeyRequest struct {
	Email string `json:"email"`
}
