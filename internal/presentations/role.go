// Package presentations
// Automatic generated
package presentations

type (
	// RoleQuery parameter
	RoleQuery struct {
		Code     string `db:"code,omitempty" json:"code" url:"code,omitempty"`
		Email    string `db:"email,omitempty" json:"email" url:"email,omitempty"`
		Password string `db:"password,omitempty" json:"password" url:"password,omitempty"`
		Paging
		PeriodRange
	}

	// RoleParam input param
	RoleParam struct {
		Email    string `db:"email,omitempty" json:"email"`
		Password string `db:"password,omitempty" json:"password"`
		Code     string `db:"code,omitempty" json:"code"`
	}

	// RoleDetail detail response
	RoleDetail struct {
		Code      string `json:"code"`
		Name      string `json:"name"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		DeletedAt string `json:"deleted_at"`
	}
)
