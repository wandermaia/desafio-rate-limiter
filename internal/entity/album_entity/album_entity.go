package album_entity

import (
	"fmt"

	"github.com/twinj/uuid"
)

// album represents data about a record album.
type Album struct {
	Id     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Cria um novo album
func CreateAlbum(title, artist string, price float64) (*Album, error) {
	album := &Album{
		Id:     uuid.NewV4().String(),
		Title:  title,
		Artist: artist,
		Price:  price,
	}

	if err := album.Validate(); err != nil {
		return nil, err
	}

	return album, nil
}

// Valida se as propriedades do álbum estão definidas conrretamente.
func (al *Album) Validate() error {
	if len(al.Title) <= 2 || len(al.Artist) <= 2 || al.Price <= 0.1 {
		return fmt.Errorf("invalid object")
	}
	return nil
}
