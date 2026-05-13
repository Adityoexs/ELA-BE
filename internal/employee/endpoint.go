package employee

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	Create  endpoint.Endpoint
	GetByID endpoint.Endpoint
	List    endpoint.Endpoint
	Update  endpoint.Endpoint
	Delete  endpoint.Endpoint
}

func NewEndpoints(svc Service) Endpoints {
	return Endpoints{
		Create:  MakeCreateEndpoint(svc),
		GetByID: MakeGetByIDEndpoint(svc),
		List:    MakeListEndpoint(svc),
		Update:  MakeUpdateEndpoint(svc),
		Delete:  MakeDeleteEndpoint(svc),
	}
}

type CreateRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Position string `json:"position"`
}

type CreateResponse struct {
	Employee *Employee `json:"employee,omitempty"`
	Error    string    `json:"error,omitempty"`
}

func MakeCreateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		employee, err := svc.Create(ctx, CreateEmployeeInput(req))
		if err != nil {
			return CreateResponse{Error: err.Error()}, nil
		}
		return CreateResponse{Employee: employee}, nil
	}
}

type GetByIDRequest struct {
	ID uint
}

type GetByIDResponse struct {
	Employee *Employee `json:"employee,omitempty"`
	Error    string    `json:"error,omitempty"`
}

func MakeGetByIDEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		employee, err := svc.GetByID(ctx, req.ID)
		if err != nil {
			return GetByIDResponse{Error: err.Error()}, nil
		}
		return GetByIDResponse{Employee: employee}, nil
	}
}

type ListResponse struct {
	Employees []Employee `json:"employees,omitempty"`
	Error     string     `json:"error,omitempty"`
}

func MakeListEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		employees, err := svc.List(ctx)
		if err != nil {
			return ListResponse{Error: err.Error()}, nil
		}
		return ListResponse{Employees: employees}, nil
	}
}

type UpdateRequest struct {
	ID       uint   `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Position string `json:"position"`
}

type UpdateResponse struct {
	Employee *Employee `json:"employee,omitempty"`
	Error    string    `json:"error,omitempty"`
}

func MakeUpdateEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateRequest)
		employee, err := svc.Update(ctx, req.ID, UpdateEmployeeInput{
			Name:     req.Name,
			Email:    req.Email,
			Position: req.Position,
		})
		if err != nil {
			return UpdateResponse{Error: err.Error()}, nil
		}
		return UpdateResponse{Employee: employee}, nil
	}
}

type DeleteRequest struct {
	ID uint
}

type DeleteResponse struct {
	Error string `json:"error,omitempty"`
}

func MakeDeleteEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteRequest)
		if err := svc.Delete(ctx, req.ID); err != nil {
			return DeleteResponse{Error: err.Error()}, nil
		}
		return DeleteResponse{}, nil
	}
}
