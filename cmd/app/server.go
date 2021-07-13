package app

import (
	"log"
	"net/http"
	
	"github.com/gorilla/mux"
	"github.com/mijgona/instagram/cmd/app/middleware"
	"github.com/mijgona/instagram/pkg/admin"
	"github.com/mijgona/instagram/pkg/post"
	"github.com/mijgona/instagram/pkg/user"

)

//Представляет собой логический сервер нашего приложения
type Server struct {
	mux 			*mux.Router
	userSvc			*user.Service
	adminSvc		*admin.Service
	postSvc			*post.Service
}

type Token struct {
	Token		string		`json:"token"`	
}

//NewServer - фунция конструктор для создания сервера
func NewServer(mux *mux.Router, userSvc	*user.Service, adminSvc	*admin.Service, postSvc *post.Service) *Server {
	log.Print("server.NewServer(): start")
	return &Server{
		mux: mux, 
		userSvc: 	userSvc,
		adminSvc: 	adminSvc,
		postSvc: 	postSvc,}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request)  {
	log.Print("server.ServeHTTP(): start")
	s.mux.ServeHTTP(writer, request)
}

const (
	GET = "GET"
	POST = "POST"
	DELETE = "DELETE"
)

//Init - инициализирует сервер (регистрирует все handler-ы)
func (s *Server) Init(){
	usersAuthenticateMd := middleware.Authenticate(s.userSvc.IDByToken)
	usersSubrouter := s.mux.PathPrefix("/api/user").Subrouter()
	usersSubrouter.Use(usersAuthenticateMd)
	usersSubrouter.HandleFunc("", s.handleGetUser).Methods(GET)
	usersSubrouter.HandleFunc("", s.handleUserEdit).Methods(POST)
	usersSubrouter.HandleFunc("", s.handleUserDelete).Methods(DELETE)	
	usersSubrouter.HandleFunc("/auth", s.handleUserGetToken).Methods(POST)
	usersSubrouter.HandleFunc("/{username}", s.handleGetUserByUsername).Methods(GET)
	usersSubrouter.HandleFunc("/follow", s.handleUserFollow).Methods(POST)
	usersSubrouter.HandleFunc("/img", s.handleUserEditImg).Methods(POST)


	postsAuthenticateMd := middleware.Authenticate(s.userSvc.IDByToken)
	postsSubrouter := s.mux.PathPrefix("/api/post").Subrouter()
	postsSubrouter.Use(postsAuthenticateMd)	
	postsSubrouter.HandleFunc("", s.handleNewPost).Methods(POST)
	postsSubrouter.HandleFunc("", s.handleGetAllPosts).Methods(GET)
	postsSubrouter.HandleFunc("/user/{username}", s.handleGetUserAllPosts).Methods(GET)
	postsSubrouter.HandleFunc("/{postid}", s.handleGetPostById).Methods(GET)
	postsSubrouter.HandleFunc("/{postid}/like", s.handlePostLike).Methods(POST)
	postsSubrouter.HandleFunc("/{postid}/delete", s.handlePostDelete).Methods(DELETE)
	postsSubrouter.HandleFunc("/{postid}/comment", s.handleGetComments).Methods(GET)
	postsSubrouter.HandleFunc("/comment", s.handleNewComment).Methods(POST)
	postsSubrouter.HandleFunc("/comment/{commentid}", s.handleCommentDelete).Methods(DELETE)

	adminAuthenticateMd := middleware.Authenticate(s.userSvc.IDByToken)
	adminSubrouter := s.mux.PathPrefix("/api/admin").Subrouter()
	adminSubrouter.Use(adminAuthenticateMd)
	adminSubrouter.HandleFunc("", s.handleGetAdmin).Methods(GET)
	adminSubrouter.HandleFunc("", s.handleAdminRegister).Methods(POST)
	adminSubrouter.HandleFunc("/auth", s.handleAdminGetToken).Methods(POST)
	adminSubrouter.HandleFunc("/active/user/{username}", s.handleAdminActiveUser).Methods(POST)
	adminSubrouter.HandleFunc("/active/users", s.handleAdminGetActiveUsers).Methods(GET)
	adminSubrouter.HandleFunc("/block/user/{username}", s.handleAdminBlockUser).Methods(POST)
	adminSubrouter.HandleFunc("/block/users", s.handleAdminGetBlockedUsers).Methods(GET)
	adminSubrouter.HandleFunc("/user/{username}", s.handleAdminDeleteUser).Methods(DELETE)	
	adminSubrouter.HandleFunc("/active/post/{postid}", s.handleAdminActivePost).Methods(POST)
	adminSubrouter.HandleFunc("/block/post/{postid}", s.handleAdminBlockPost).Methods(POST)	
	adminSubrouter.HandleFunc("/post/{postid}", s.handleAdminDeletePost).Methods(DELETE)

}