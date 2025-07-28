package handlers

import (
	"net/http"

	"github.com/gmcc94/attendance-go/db"
	"github.com/gmcc94/attendance-go/helpers"
	"github.com/gmcc94/attendance-go/types"
)

func CreateAccountsHandler(accountsStore db.AccountsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateAccountsRequest
		if !helpers.DecodeJSON(r, w, &req) {
			return
		}

		if err := accountsStore.AddAccountEntries(req.Income, "income"); err != nil {
			helpers.JSONError(w, http.StatusInternalServerError, "Failed to insert incomes")
			return
		}

		if err := accountsStore.AddAccountEntries(req.Expenditure, "expenditure"); err != nil {
			helpers.JSONError(w, http.StatusInternalServerError, "Failed to insert expenditure")
			return
		}

		helpers.WriteJSON(w, http.StatusCreated, types.CreateAccountsResponse{
			Message: "Accounts saved successfully",
		})
	}
}

func GetGroupedAccountsHandler(accountsStore db.AccountsStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		entries, err := accountsStore.GetAccounts()
		if err != nil {
			helpers.JSONError(w, http.StatusInternalServerError, "Failed to fetch accounts")
			return
		}

		grouped := helpers.GroupedAccounts(entries)
		resp := helpers.BuildGroupedResponse(grouped)

		helpers.WriteJSON(w, http.StatusOK, resp)
	}
}
