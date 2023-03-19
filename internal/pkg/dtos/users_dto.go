package dtos

type (
	FilterUsers struct {
		Limit  int    `query:"limit" validate:"required,gt=0"`
		Page   int    `query:"page" validate:"required,gt=0"`
		Search string `query:"search"`
	}

	ReqCreateDataUser struct {
		Email     string `json:"email" validate:"required,email"`
		Firstname string `json:"first_name" validate:"required,min=2"`
		Lastname  string `json:"last_name" validate:"required,min=2"`
		Avatar    string `json:"avatar"`
	}

	ReqUpdateDataUser struct {
		Email     string `json:"email,omitempty" validate:"omitempty,email"`
		Firstname string `json:"first_name,omitempty" validate:"omitempty,min=2"`
		Lastname  string `json:"last_name,omitempty" validate:"omitempty,min=2"`
		Avatar    string `json:"avatar,omitempty"`
	}

	ResDataUserSingle struct {
		DtosModel
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	}

	ResDataUsers struct {
		Data []ResDataUserSingle `json:"data"`
		Pagination
	}
)
