package search

import "github.com/Latezly/nyaa_go/models"

// TorrentCache torrent cache struct
type TorrentCache struct {
	Torrents []models.Torrent
	Count    int
}
