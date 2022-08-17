package dependencies

import (
	"apart-deal-api/pkg/store/user"

	"go.uber.org/fx"
)

var RepositoryModule = fx.Provide(
	user.NewUserRepository,
)
