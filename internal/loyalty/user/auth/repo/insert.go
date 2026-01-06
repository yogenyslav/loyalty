package repo

import (
	"context"

	"github.com/yogenyslav/loyalty/internal/loyalty/user/auth/model"
	"github.com/yogenyslav/loyalty/pkg/errs"
)

const insertUser = `
	insert into loyalty."user" (login, hashed_password)
	values ($1, $2)
	returning id;
`

// InsertUser inserts a new user into the database and returns the new user's ID.
func (r *Repo) InsertUser(ctx context.Context, u *model.User) (int64, error) {
	var userID int64
	err := r.db.TxQueryRow(ctx, &userID, insertUser, u.Login, u.HashedPassword)
	if err != nil {
		return 0, errs.Wrap(err, "query row")
	}
	return userID, nil
}
