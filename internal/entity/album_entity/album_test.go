package album_entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wandermaia/desafio-rate-limiter/internal/entity/album_entity"
)

func TestNewAlbum(t *testing.T) {
	album, err := album_entity.CreateAlbum("Led Zeppelin IV", "Led Zeppelin", 47.99)
	assert.Nil(t, err)
	assert.Equal(t, album.Title, "Led Zeppelin IV")
	assert.Equal(t, album.Artist, "Led Zeppelin")
	assert.Equal(t, album.Price, 47.99)
}
