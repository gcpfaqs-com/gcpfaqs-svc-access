package rest

import (
	"encoding/json"
	"gcpaccesss-svc-accesss/internal/application"
	"gcpaccesss-svc-accesss/internal/domain/persistence"
	"gcpaccesss-svc-accesss/internal/domain/secure"
	"net/http"
)

type updateDTO struct {
	Name string `json:"name"`
}

type RESTUpdateUseCase struct {
	repository persistence.AccessRepository
	tokener    secure.Tokener
}

func NewRESTUpdateUseCase(repository persistence.AccessRepository, tokener secure.Tokener) *RESTUpdateUseCase {
	return &RESTUpdateUseCase{
		repository: repository,
		tokener:    tokener,
	}
}

func (uc *RESTUpdateUseCase) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.Header.Get("Authorization")
	subject, err := uc.tokener.DecodeToken(idStr)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: "invalid token"})
		return
	}

	var dto updateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: "invalid request body"})
		return
	}

	updateUC := application.NewUpdateAccessUseCase(uc.repository)
	err = updateUC.Update(subject, dto.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResponseDTO{Error: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
}
