package applications

import (
	"net/http"

	"github.com/CelanMatjaz/job_application_tracker_api/pkg/db"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/service"
	"github.com/CelanMatjaz/job_application_tracker_api/pkg/types"
)

// =============== APLICATIONS =============================

type ApplicationsResponse struct {
	Pagination service.PaginationParams `json:"pagination"`
	Data       []types.Application      `json:"applications"`
}

func sendJsonApplications(w http.ResponseWriter, data []types.Application, pagination service.PaginationParams, statusCode int) {
	service.SendJsonResponse(w, ApplicationsResponse{pagination, data}, statusCode)
}

type ApplicationResponse struct {
	Data types.Application `json:"application"`
}

func sendJsonApplication(w http.ResponseWriter, data types.Application, statusCode int) {
	service.SendJsonResponse(w, ApplicationResponse{data}, statusCode)
}

type ApplicationHandler struct {
	store db.ApplicationStore
}

func (h ApplicationHandler) GetMultiple(userId int, pagination service.PaginationParams) ([]types.Application, error) {
	return h.store.GetUserApplications(userId, pagination)
}

func (h ApplicationHandler) GetSingle(userId int, applicationId int) (types.Application, error) {
	return h.store.GetUserApplication(userId, applicationId)
}

func (h ApplicationHandler) Create(userId int, data types.Application) (types.Application, error) {
	return h.store.CreateUserApplication(userId, data)
}

func (h ApplicationHandler) Update(userId int, data types.Application) (types.Application, error) {
	return h.store.UpdateUserApplication(userId, data)
}

func (h ApplicationHandler) Delete(userId int, applicationId int) (int, error) {
	return h.store.DeleteUserApplication(userId, applicationId)
}

// =============== SECTIONS =============================

type SectionsResponse struct {
	Pagination service.PaginationParams   `json:"pagination"`
	Data       []types.ApplicationSection `json:"sections"`
}

func sendJsonSections(w http.ResponseWriter, data []types.ApplicationSection, pagination service.PaginationParams, statusCode int) {
	service.SendJsonResponse(w, SectionsResponse{pagination, data}, statusCode)
}

type SectionResponse struct {
	Data types.ApplicationSection `json:"section"`
}

func sendJsonSection(w http.ResponseWriter, data types.ApplicationSection, statusCode int) {
	service.SendJsonResponse(w, SectionResponse{data}, statusCode)
}

type SectionHandler struct {
	store db.ApplicationStore
}

func (h SectionHandler) GetSingle(userId int, applicationId int) (types.ApplicationSection, error) {
	return types.ApplicationSection{}, nil
}

func (h SectionHandler) GetMultiple(userId int, pagination service.PaginationParams) ([]types.ApplicationSection, error) {
	return h.store.GetApplicationSections(userId, pagination)
}

func (h SectionHandler) Create(userId int, data types.ApplicationSection) (types.ApplicationSection, error) {
	return h.store.CreateApplicationSection(userId, data)
}

func (h SectionHandler) Update(userId int, data types.ApplicationSection) (types.ApplicationSection, error) {
	return h.store.UpdateApplicationSection(userId, data)
}

func (h SectionHandler) Delete(userId int, applicationId int) (int, error) {
	return h.store.DeleteUserApplication(userId, applicationId)
}
