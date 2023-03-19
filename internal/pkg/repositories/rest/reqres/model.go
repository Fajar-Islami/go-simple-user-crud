package reqres

type (
	ReqListUser struct {
		PerPage int
		Page    int
	}

	ResListUsers struct {
		Page       int        `json:"page"`
		PerPage    int        `json:"per_page"`
		Total      int        `json:"total"`
		TotalPages int        `json:"total_pages"`
		Data       []ListUser `json:"data"`
		Support    struct {
			URL  string `json:"url"`
			Text string `json:"text"`
		} `json:"support"`
	}

	ListUser struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	}
)
