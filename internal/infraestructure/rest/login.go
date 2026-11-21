package rest

import (
	"encoding/json"
	"gcpaccesss-svc-accesss/internal/application"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
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
	var dto loginDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	saveUC := application.NewLoginAccessUseCase(uc.repository, uc.hasher, uc.tokener)
	token, err := saveUC.Login(dto.Username, dto.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(loginResponseDTO{Token: *token})
}
