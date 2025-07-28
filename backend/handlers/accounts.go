package handlers

import (
	"log"
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

		if err := helpers.ValidateEntries(append(req.Income, req.Expenditure...)); err != nil {
			log.Printf("Validation error :%v", err)
			helpers.JSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		entries := map[string][]types.AccountEntry{
			"income":      req.Income,
			"expenditure": req.Expenditure,
		}

		for category, slice := range entries {
			if len(slice) == 0 {
				continue
			}

			if err := accountsStore.AddAccountEntries(slice, category); err != nil {
				log.Printf("DB insert error (%s): %v", category, err)
				helpers.JSONError(w, http.StatusInternalServerError, "Failed to insert "+category)
				return
			}
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
