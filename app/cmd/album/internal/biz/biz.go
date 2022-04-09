package biz

import "projectlayout/app/cmd/album/internal/data"

func AlbumByID(id int64) (*data.Album, error) {
	return data.AlbumByID(id)
}
