package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const findUserByLogin = `
	select id, login, hashed_password, created_at, updated_at
	from loyalty."user"
	where login = $1;
`

// FindUserByLogin finds a user by specified login.
func (r *Repo) FindUserByLogin(ctx context.Context, login string) (model.User, error) {
	var user model.User
	err := r.db.QueryRow(ctx, &user, findUserByLogin, login)
	return user, errs.Wrap(err, "query row")
}
