package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/DanyJDuque/go_lib_response/response"
	"github.com/DanyJDuque/gocourse_course/internal/course"
	"github.com/gorilla/mux"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

func NewCourseHTTPServer(ctx context.Context, endpoints course.Endpoints) http.Handler {

	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateCourse, encodeResponse,
		opts...,
	)).Methods("POST")

	r.Handle("/courses", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllCourse, encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetCourse, encodeResponse,
		opts...,
	)).Methods("GET")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateCourse, encodeResponse,
		opts...,
	)).Methods("PATCH")

	r.Handle("/courses/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteCourse, encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

func decodeCreateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}
	return req, nil
}

func decodeGetCourse(_ context.Context, r *http.Request) (interface{}, error) {

	p := mux.Vars(r)
	req := course.GetReq{
		ID: p["id"],
	}
	return req, nil
}

func decodeGetAllCourse(_ context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := course.GetAllReq{
		Name:  v.Get("name"),
		Limit: limit,
		Page:  page,
	}
	return req, nil
}

func decodeUpdateCourse(_ context.Context, r *http.Request) (interface{}, error) {
	var req course.UpdateReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid request format: '%v'", err.Error()))
	}

	path := mux.Vars(r)
	req.ID = path["id"]

	return req, nil
}

func decodeDeleteCourse(_ context.Context, r *http.Request) (interface{}, error) {

	path := mux.Vars(r)
	req := course.DeleteReq{
		ID: path["id"],
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Contet-Type", "application/json; charset=utf8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Contect-type", "application/json; charset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
