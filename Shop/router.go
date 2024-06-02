package shop

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	// объявляем mux.Router
	router := mux.NewRouter()
	// вывод списка из Items с фильтрацией
	router.HandleFunc("/items", ListItemsHandler).Methods(http.MethodGet)
	// добавление нового Item
	router.HandleFunc("/items", CreateItemHandler).Methods(http.MethodPost)
	// получение Item по ID
	router.HandleFunc("/items/{id}", GetItemHandler).Methods(http.MethodGet)
	// изменение Item по ID
	router.HandleFunc("/items/{id}", UpdateItemHandler).Methods(http.MethodPut)
	// удаление Item по ID
	router.HandleFunc("/items/{id}", DeleteItemHandler).Methods(http.MethodDelete)
	// загрузка файла изображения
	router.HandleFunc("/items/upload_image", UploadItemImageHandler).Methods(http.MethodPost)
	// авторизация
	router.HandleFunc("/user/login", LoginHandler).Methods(http.MethodPost)
	// завершение сессии
	router.HandleFunc("/user/logout", LogoutHandler).Methods(http.MethodPost)
	return router
}
