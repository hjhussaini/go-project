package main

import (
	"encoding/json"
	"net/http"

	"github.com/hjhussaini/go-project/wallet/models"
	"github.com/go-chi/chi/v5"
)

type UserList struct {
	Users	[]string	`json:"users"`
}

type Response struct {
	OK	bool		`json:"ok"`
	Message	string		`json:"message,omitempty"`
	Content	interface{}	`json:"content,omitempty"`
}

func (response *Response) write(writer http.ResponseWriter, statusCode int) {
	data, err := json.Marshal(response)
	if err != nil {
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(statusCode)
	writer.Write(data)
}

func (app *application) createWallet(writer http.ResponseWriter, request *http.Request) {
	var wallet models.Wallet

	err := json.NewDecoder(request.Body).Decode(&wallet)
	if err != nil {
		response := Response{
			OK:		false,
			Message:	err.Error(),
		}
		response.write(writer, http.StatusUnprocessableEntity)

		return
	}

	if err = app.database.InsertWallet(wallet); err != nil {
		response := Response{
			OK:		false,
			Message:	err.Error(),
		}
		response.write(writer, http.StatusInternalServerError)

		return
	}

	response := Response{
		OK:		true,
		Content:	wallet,
	}
	response.write(writer, http.StatusOK)
}

func (app *application) getWalletByID(writer http.ResponseWriter, request *http.Request) {
	phoneNumber := chi.URLParam(request, "id")
	wallet := app.database.GetWallet(phoneNumber)
	if wallet == nil {
		response := Response{
			OK:		false,
			Message:	"Not found",
		}
		response.write(writer, http.StatusNotFound)

		return
	}

	response := Response{
		OK:		true,
		Content:	wallet,
	}
	response.write(writer, http.StatusOK)
}
