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

	//Саброутер пользователей
	usersAuthenticateMd := middleware.Authenticate(s.userSvc.IDByToken)
	usersSubrouter := s.mux.PathPrefix("/api/user").Subrouter()
	usersSubrouter.Use(usersAuthenticateMd)
	//Выдает посты всех пользователей на которых подписан авторизованный пользователя
	usersSubrouter.HandleFunc("", s.handleGetUser).Methods(GET)
	//Меняет данные пользователя
	usersSubrouter.HandleFunc("", s.handleUserEdit).Methods(POST)
	//Удаляет пользователя
	usersSubrouter.HandleFunc("", s.handleUserDelete).Methods(DELETE)
	//Авторизует пользователя	
	usersSubrouter.HandleFunc("/auth", s.handleUserGetToken).Methods(POST)
	//Выдает данные и посты пользователя с указанным username
	usersSubrouter.HandleFunc("/{username}", s.handleGetUserByUsername).Methods(GET)
	//Подписывается на пользователя
	usersSubrouter.HandleFunc("/{username}/follow", s.handleUserFollow).Methods(POST)
	//Меняет изображение пользователя
	usersSubrouter.HandleFunc("/img", s.handleUserEditImg).Methods(POST)

	//Саброутер Постов
	postsAuthenticateMd := middleware.Authenticate(s.userSvc.IDByToken)
	postsSubrouter := s.mux.PathPrefix("/api/post").Subrouter()
	postsSubrouter.Use(postsAuthenticateMd)	
	//Создаёт новый пост
	postsSubrouter.HandleFunc("", s.handleNewPost).Methods(POST)
	//Выдаёт все посты авторизованного пользователя
	postsSubrouter.HandleFunc("", s.handleGetAllPosts).Methods(GET)
	//Выдаёт все посты пользователя с указанным usename
	postsSubrouter.HandleFunc("/user/{username}", s.handleGetUserAllPosts).Methods(GET)
	//Выдаёт пост с указанным ID
	postsSubrouter.HandleFunc("/{postid}", s.handleGetPostById).Methods(GET)
	// Лайкает пост с указанным ID
	postsSubrouter.HandleFunc("/{postid}/like", s.handlePostLike).Methods(POST)
	//Удаляет пост с указанным ID
	postsSubrouter.HandleFunc("/{postid}/delete", s.handlePostDelete).Methods(DELETE)
	//Комментирует пост с указанним ID
	postsSubrouter.HandleFunc("/comment/{postid}", s.handleGetComments).Methods(GET)
	//Создаёт пост от имени авторизованного пользователя
	postsSubrouter.HandleFunc("/comment", s.handleNewComment).Methods(POST)
	//Удаляет комментарий с указанным ID
	postsSubrouter.HandleFunc("/comment/{commentid}", s.handleCommentDelete).Methods(DELETE)


	//Саброутер Админов
	adminAuthenticateMd := middleware.Authenticate(s.userSvc.IDByToken)
	adminSubrouter := s.mux.PathPrefix("/api/admin").Subrouter()
	adminSubrouter.Use(adminAuthenticateMd)
	//Выдаёт данные админа
	adminSubrouter.HandleFunc("", s.handleGetAdmin).Methods(GET)
	//Регистрирует нового админа
	adminSubrouter.HandleFunc("", s.handleAdminRegister).Methods(POST)
	//Авторизует админа
	adminSubrouter.HandleFunc("/auth", s.handleAdminGetToken).Methods(POST)
	//Активирует посльзователя с указанным username
	adminSubrouter.HandleFunc("/active/user/{username}", s.handleAdminActiveUser).Methods(POST)
	//Выдаёт всех активных пользователей
	adminSubrouter.HandleFunc("/active/users", s.handleAdminGetActiveUsers).Methods(GET)
	//Блокирует посльзователя с указанным username
	adminSubrouter.HandleFunc("/block/user/{username}", s.handleAdminBlockUser).Methods(POST)
	//выдаёт всех заблокированных пользователей
	adminSubrouter.HandleFunc("/block/users", s.handleAdminGetBlockedUsers).Methods(GET)
	//Активирует посльзователя с указанным username
	adminSubrouter.HandleFunc("/user/{username}", s.handleAdminDeleteUser).Methods(DELETE)	
	//Активирует пост с указанным ID
	adminSubrouter.HandleFunc("/active/post/{postid}", s.handleAdminActivePost).Methods(POST)	
	//Блокирует пост с указанным ID
	adminSubrouter.HandleFunc("/block/post/{postid}", s.handleAdminBlockPost).Methods(POST)		
	//Удаляет пост с указанным ID
	adminSubrouter.HandleFunc("/post/{postid}", s.handleAdminDeletePost).Methods(DELETE)
}