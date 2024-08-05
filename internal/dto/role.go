package dto

import (
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/entity"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/presentations"
	"gitlab.privy.id/privypass/privypass-boilerplate/pkg/util"
)

func RoleToResponse(src entity.Role) presentations.RoleDetail {
	x := presentations.RoleDetail{
		Code:      src.Code,
		Name:      src.Name,
		CreatedAt: util.DateToString(src.CreatedAt),
		UpdatedAt: util.DateToString(src.UpdatedAt),
	}

	if !src.CreatedAt.IsZero() {
		x.CreatedAt = src.CreatedAt.Format(consts.LayoutDateTimeFormat)
	}
	//
	//if !src.UpdatedAt.IsZero() {
	//	x.UpdatedAt = src.UpdatedAt.Format(consts.LayoutDateTimeFormat)
	//}
	//
	//if !src.DeletedAt.IsZero() {
	//	x.DeletedAt = src.DeletedAt.Format(consts.LayoutDateTimeFormat)
	//}

	return x
}

func RolesToResponse(inputs []entity.Role) []presentations.RoleDetail {
	var (
		result = []presentations.RoleDetail{}
	)

	for i := 0; i < len(inputs); i++ {
		result = append(result, RoleToResponse(inputs[i]))
	}

	return result
}
