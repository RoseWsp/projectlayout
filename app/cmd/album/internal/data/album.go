package data

import (
	"database/sql"

	"github.com/pkg/errors"
)

func AlbumByID(id int64) (*Album, error) {
	querySql := "SELECT * FROM album WHERE id=?"
	var album Album
	row := DB.QueryRow(querySql, id)
	if err := row.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Wrapf(err, "albumByID(id=%d) no found album", id)
		}
		return nil, errors.Wrapf(err, "albumByID(id=%d) has error", id)
	}
	return &album, nil
}
