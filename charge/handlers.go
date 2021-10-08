package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/hjhussaini/go-project/charge/models"
	"github.com/go-chi/chi/v5"
)

type Wallet struct {
	PhoneNo	string	`json:"phone_no"`
	Credit	int	`json:"credit"`
}

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

func (app *application) giftCharge(writer http.ResponseWriter, request *http.Request) {
	var charge models.Charge

	err := json.NewDecoder(request.Body).Decode(&charge)
	if err != nil {
		response := Response{
			OK:		false,
			Message:	err.Error(),
		}
		response.write(writer, http.StatusUnprocessableEntity)

		return
	}

	if app.database.CountCharge(charge.Code) >= app.config.charge.max {
		response := Response{
			OK:		false,
			Message:	"all charges have been used",
		}
		response.write(writer, http.StatusServiceUnavailable)

		return
	}

	wallet := Wallet{
		PhoneNo:	charge.PhoneNo,
		Credit:		app.config.charge.credit,
	}

	data, err := json.Marshal(wallet)
	if err != nil {
		response := Response{
			OK:		false,
			Message:	err.Error(),
		}
		response.write(writer, http.StatusInternalServerError)

		return
	}

	resp, err := http.Post(
				app.config.charge.walletAPI,
				"application/json",
				bytes.NewBuffer(data),
				)
	if err != nil {
		response := Response{
			OK:		false,
			Message:	err.Error(),
		}
		response.write(writer, http.StatusInternalServerError)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		response := Response{
			OK:		false,
			Message:	"Unable to create wallet",
		}
		response.write(writer, http.StatusInternalServerError)

		return
	}

	if err = app.database.InsertCharge(charge); err != nil {
		response := Response{
			OK:		false,
			Message:	err.Error(),
		}
		response.write(writer, http.StatusInternalServerError)

		return
	}

	response := Response{
		OK:		true,
		Content:	charge,
	}
	response.write(writer, http.StatusOK)
}

func (app *application) getUsersByCharge(writer http.ResponseWriter, request *http.Request) {
	var list UserList

	code, _ := strconv.Atoi(chi.URLParam(request, "code"))

	chargeList := app.database.GetChargeList(code)
	for _, charge := range chargeList {
		list.Users = append(list.Users, charge.PhoneNo)
	}

	response := Response{
		OK:		true,
		Content:	list,
	}

	response.write(writer, http.StatusOK)
}
