package handler

import (
	domain2 "github.com/aifedorov/gophermart/internal/user/repository"
)

func ToBalanceResponse(balance domain2.Balance) BalanceResponse {
	return BalanceResponse{
		balance.Current,
		balance.Withdrawn,
	}
}
