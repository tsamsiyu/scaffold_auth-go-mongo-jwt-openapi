package dependencies

import (
	authDomain "apart-deal-api/pkg/domain/auth"

	"go.uber.org/fx"
)

var AuthServicesModule = fx.Provide(
	authDomain.NewSignUpService,
	authDomain.NewConfirmSignUpService,
)
