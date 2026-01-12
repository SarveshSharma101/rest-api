package datamodels

type User struct {
	Age       int32  `json:"age,omitempty"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Name      string `json:"name"`
	Gender    string `json:"gender,omitempty"`
	UID       string `json:"uid"`
	IsDeleted bool   `json:"isDeleted"`
}

type UserUpdateReq struct {
	Age    int32  `json:"age,omitempty"`
	Email  string `json:"email,omitempty"`
	Name   string `json:"name,omitempty"`
	Gender string `json:"gender,omitempty"`
}
