package rest

import (
	"encoding/json"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
	"gcpaccesss-svc-accesss/internal/domain/services"
	"net/http"
)

type loginDTO struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

type loginResponseDTO struct {
	Token string `json:"token"`
}

type RESTLoginUseCase struct {
	repository persistence.AccessRepository
	hasher     secure.Hasher
	tokener    secure.Tokener
}

func NewRESTLoginUseCase(
	repository persistence.AccessRepository,
	hasher secure.Hasher,
	tokener secure.Tokener) *RESTLoginUseCase {
	return &RESTLoginUseCase{
		repository: repository,
		hasher:     hasher,
		tokener:    tokener,
	}
}

func (uc *RESTLoginUseCase) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto loginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: "invalid request body"})
		return
	}

	saveUC := services.NewLoginAccessUseCase(uc.repository, uc.hasher, uc.tokener)
	token, err := saveUC.Do(dto.Username, dto.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponseDTO{Token: *token})
}
