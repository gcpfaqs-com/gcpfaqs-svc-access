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
	idStr := r.Header.Get("Authorization")
	subject, err := uc.tokener.DecodeToken(idStr)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	var dto updateDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	updateUC := application.NewUpdateAccessUseCase(uc.repository)
	err = updateUC.Update(subject, dto.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
