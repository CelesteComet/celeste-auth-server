package postgres

import (
	"github.com/CelesteComet/celeste-auth-service/app"
)

var _ app.UserService = UserService{}
