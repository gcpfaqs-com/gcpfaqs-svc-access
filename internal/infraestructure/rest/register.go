package rest

import (
	"encoding/json"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
	"gcpaccesss-svc-accesss/internal/domain/services"
	"net/http"
)

type registerDTO struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	Name     string `json:"name"`
}

type registerResponseDTO struct {
	ID string `json:"id"`
}

type RESTRegisterUseCase struct {
	repository persistence.AccessRepository
	hasher     secure.Hasher
	tokener    secure.Tokener
}

func NewRESTRegisterUseCase(
	repository persistence.AccessRepository,
	hasher secure.Hasher,
	tokener secure.Tokener) *RESTRegisterUseCase {
	return &RESTRegisterUseCase{
		repository: repository,
		hasher:     hasher,
		tokener:    tokener,
	}
}

func (uc *RESTRegisterUseCase) Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var dto registerDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: "invalid request body"})
		return
	}

	saveUC := services.NewRegisterAccessUseCase(uc.repository, uc.hasher, uc.tokener)
	id, err := saveUC.Do(dto.Username, dto.Password, dto.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: err.Error()})
		return
	}

	dto.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerResponseDTO{ID: dto.ID})
}
