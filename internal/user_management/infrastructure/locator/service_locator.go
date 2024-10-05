package locator

import (
	"database/sql"
	"league-management/internal/user_management/application/service"
)

type ServiceLocator struct {
	db          *sql.DB
	userService *service.UserService
}

func (sl *ServiceLocator) GetUserService() *service.UserService {
	if sl.userService == nil {
		sl.userService = service.NewUserService()
	}
	return sl.userService
}
