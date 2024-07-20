package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wandermaia/desafio-rate-limiter/internal/usecase/album_usecase"
)

// Handler de albuns
type AlbumHandler struct {
	albumUseCase album_usecase.AlbumUseCaseInterface
}

// Construtora do handler
func NewAlbumHandler(album_UseCase album_usecase.AlbumUseCaseInterface) *AlbumHandler {
	return &AlbumHandler{
		albumUseCase: album_UseCase,
	}
}

// Função responsável pela criação de um novo Álbum. Necessita receber um AlbumInputDTO
func (ah AlbumHandler) CreateAlbum(c *gin.Context) {

	// DTO que vai receber os
	var album album_usecase.AlbumInputDTO

	if err := c.ShouldBindJSON(&album); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid album data"})
		return
	}

	// Caso de uso de criação de album
	if err := ah.albumUseCase.CreateNewAlbum(context.Background(), &album); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create album"})
		return
	}
	c.Status(http.StatusCreated)
}

// Busca album pelo ID
func (ah AlbumHandler) FindAlbumById(c *gin.Context) {

	// ID para remoção
	albumId := c.Param("albumId")

	// Verficando se o uuid é válido
	if err := uuid.Validate(albumId); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}

	// Verificando se o album existe
	albumOutputDTO, err := ah.albumUseCase.GetAlbumByID(context.Background(), albumId)
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// Retornando o status de ok e os dados do album
	c.IndentedJSON(http.StatusOK, albumOutputDTO)

}

// Busca album pelo ID
func (ah AlbumHandler) FindAllAlbums(c *gin.Context) {

	// Recuperando os dados de todos os álbuns
	albuns, err := ah.albumUseCase.GetAllAlbums(context.Background())
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get albums"})
		return
	}

	// Retornando o status e a lista de albums
	c.IndentedJSON(http.StatusOK, albuns)

}

// Criação de um novo Álbum. Necessita receber um AlbumInputDTO
func (ah AlbumHandler) DeleteAlbumByID(c *gin.Context) {

	// ID para remoção
	albumId := c.Param("albumId")

	// Verficando se o uuid é válido
	if err := uuid.Validate(albumId); err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid uuid"})
		return
	}

	// Verificando se o album existe
	_, err := ah.albumUseCase.GetAlbumByID(context.Background(), albumId)
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	// Deletando o álbum
	err = ah.albumUseCase.DeleteAlbumByID(context.Background(), albumId)
	if err != nil {
		log.Print(err)
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete album"})
		return
	}

	// Retornando status do delete executado com sucesso (204)
	c.Status(http.StatusNoContent)

}
