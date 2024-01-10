package routers

import (
	"net/http"
	"time"

	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/domain/services"
	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/http/handlers"
	"github.com/Lucasvmarangoni/financial-file-manager/internal/modules/user/infra/repositories"
	"github.com/go-chi/chi"
	"github.com/go-chi/httprate"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5"
)

type UserRouter struct {
	Db            pgx.Tx
	method        string
	userHandler   *handlers.UserHandler
	jwtExpiriesIn int
}

func NewUserRouter(db pgx.Tx, jwtExpiriesIn int, tokenAuth *jwtauth.JWTAuth) *UserRouter {
	u := &UserRouter{
		Db:            db,
		jwtExpiriesIn: jwtExpiriesIn,
	}
	u.userHandler = u.init()
	return u
}

func (u *UserRouter) init() *handlers.UserHandler {
	userRepository := repositories.NewUserRepository(u.Db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)
	return userHandler
}

func (u *UserRouter) InitializeUserRoutes(r chi.Router) {
	r.Route("/authn", func(r chi.Router) {
		u.Method("POST").InitializeRoute(r, "/create", u.userHandler.Create)
		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				5,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			u.Method("POST").InitializeRoute(r, "/", u.userHandler.Authentication)
		})
	})
}

func (u *UserRouter) UserRoutes(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				10,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			u.Method("GET").InitializeRoute(r, "/me", u.userHandler.Me)
			u.Method("PUT").InitializeRoute(r, "/update", u.userHandler.Update)
		})

		r.Group(func(r chi.Router) {
			r.Use(httprate.Limit(
				3,
				60*time.Minute,
				httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
				httprate.WithLimitHandler(func(w http.ResponseWriter, r *http.Request) {
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				}),
			))
			u.Method("DELETE").InitializeRoute(r, "/del", u.userHandler.Delete)
			u.Method("PATCH").InitializeRoute(r, "/authz/{id}", u.userHandler.AdminAuthz)
		})
	})
}
