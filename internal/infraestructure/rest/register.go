package rest

import (
	"encoding/json"
	"gcpaccesss-svc-accesss/internal/application"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
	"net/http"
)

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
	var dto accessDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	saveUC := application.NewRegisterAccessUseCase(uc.repository, uc.hasher, uc.tokener)
	id, err := saveUC.Register(dto.Username, dto.Password, dto.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto)
}
