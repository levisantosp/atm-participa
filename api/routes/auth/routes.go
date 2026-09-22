package auth

import (
	"github.com/danielgtaylor/huma/v2"
	"github.com/levisantosp/atm-participa/api/middlewares"
)

func Routes(api huma.API) {
	group := huma.NewGroup(api, "/auth")
	huma.Post(group, "/sign-in/email", SignInWithEmail)
	huma.Post(group, "/sign-up/email", SignUpWithEmail)
	huma.Get(group, "/me", GetMe)

	group.UseMiddleware(middlewares.Auth(group, false))
	huma.Post(group, "/sign-out", SignOut)
}
