package app

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mijgona/instagram/cmd/app/middleware"
	"github.com/mijgona/instagram/types"
)

func (s *Server) handleAdminGetToken(writer http.ResponseWriter, request *http.Request) {
	var item *types.User
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := s.adminSvc.Token(request.Context(), item.UserName, item.Password)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&types.Token{Token: token})
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}


func (s *Server) handleAdminRegister(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	var item *types.User
	err = json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = s.adminSvc.Register(request.Context(), item, auth)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := s.adminSvc.Token(request.Context(), item.UserName, item.Password)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(&types.Token{Token: token})
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (s *Server) handleGetAdmin(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.adminSvc.GetAdmin(request.Context(), auth)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}

func (s *Server) handleAdminActiveUser(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.adminSvc.ActiveUser(request.Context(), auth, username)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}

func (s *Server) handleAdminBlockUser(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.adminSvc.BlockUser(request.Context(), auth, username)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}

func (s *Server) handleAdminDeleteUser(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.adminSvc.DeleteUser(request.Context(), auth, username)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write([]byte("deleted"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}
func (s *Server) handleAdminGetBlockedUsers(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.adminSvc.BlockUser(request.Context(), auth, "")
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}

func (s *Server) handleAdminGetActiveUsers(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item, err := s.adminSvc.ActiveUser(request.Context(), auth, "")
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}


func (s *Server) handleAdminActivePost(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	postIdString, ok := mux.Vars(request)["postid"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	postID, err := strconv.ParseInt(postIdString, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	item, err := s.adminSvc.ActivePost(request.Context(), auth, postID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}

func (s *Server) handleAdminBlockPost(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	postIdString, ok := mux.Vars(request)["postid"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	postID, err := strconv.ParseInt(postIdString, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	item, err := s.adminSvc.BlockPost(request.Context(), auth, postID)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
}

func (s *Server) handleAdminDeletePost(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 || !auth.IsAdmin {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	postIdString, ok := mux.Vars(request)["postid"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	
	postID, err := strconv.ParseInt(postIdString, 10, 64)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	err = s.postSvc.DeletePost(request.Context(), postID, auth)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write([]byte("success"))
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

}