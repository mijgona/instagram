package app

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/mijgona/instagram/cmd/app/middleware"
	"github.com/mijgona/instagram/types"
)

func (s *Server) handleGetUserByUsername(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	username, ok := mux.Vars(request)["username"]
	if !ok {
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	user, err := s.userSvc.GetUser(request.Context(), auth, username)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	posts, err := s.postSvc.GetAllPost(request.Context(),auth, username)
	if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
	}

	wall := types.Wall{
		User: *user,
		Posts: posts,
	}

	data, err := json.Marshal(wall)
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




func (s *Server) handleGetUser(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	user, err := s.userSvc.GetUser(request.Context(), auth, "")
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	posts := []*types.Post{}
	for _, users := range user.Follows {
		post, err :=s.postSvc.GetAllPost(request.Context(), auth, users.UserName)
		if err != nil {
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			log.Print(err)
			return
		}
		posts = append(posts, post...)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID>posts[j].ID
	})

	wall := types.Wall{
		User:  *user,
		Posts: posts,
	}

	data, err := json.Marshal(wall)
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



func (s *Server) handleUserFollow(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	id := auth.ID
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	item := &types.Follow{}
	err = json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	
	item, err = s.userSvc.Follow(request.Context(), item, id)
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
func (s *Server) handleUserEditImg(writer http.ResponseWriter, request *http.Request) {
	item := &types.User{}
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if auth.ID == 0 {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item.ID=auth.ID
	//?????????????????? ??????????????????????   
    item.Photo, err = saveImg(request, item.Photo)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	_, err = s.userSvc.EditUser(request.Context(), item, auth)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}


	data, err := json.Marshal([]byte("success"))
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

func (s *Server) handleUserEdit(writer http.ResponseWriter, request *http.Request) {
	item := &types.User{}
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	auth, err := middleware.Authentication(request.Context())
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}
	item.Photo, err = saveImg(request, item.Photo)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	item, err = s.userSvc.EditUser(request.Context(), item, auth)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Print(err)
		return
	}

	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	token := ""
	if item.UserName !="" && item.Password !=""{
		token, err = s.userSvc.Token(request.Context(), item.UserName, item.Password)
		if err != nil {
			log.Print(err)
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	data, err := json.Marshal(&types.Token{Token: token})
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

func saveImg(request *http.Request, name string) (string, error) {
	//???????????????????? ?????????????????? ???????????? ?????? ?????????? ??????????
	outStr := uuid.New()
	//???????????? ????????
	in, header, err := request.FormFile("image")
    if err !=nil && err!=http.ErrNotMultipart{
		log.Print("SaveImg err:", err)
		return "", types.ErrInternal
    }	
	
	//???????????????????? ????????????
	if err!=http.ErrNotMultipart{			
		i := strings.Index(header.Filename, ".")
		name ="../user_img/" + outStr.String()+"."+header.Filename[i:]
	}
	// ???????? ???????? ??????, ?????????????????????? ???????? ???? ??????????????????
	if err==http.ErrNotMultipart{
		name = "../user_img/"+ outStr.String()+ ".png"
		in, err = os.Open("../user_img/default.png")		
		if err !=nil {
			log.Print("SaveImg err:", err)
			return "", types.ErrInternal
		}
	}
    defer in.Close()

	//??????????????????
	out, err := os.Create(name) 
    if err !=nil {
		log.Print("SaveImg err:", err)
		return "", types.ErrInternal
    }
    defer out.Close()
    io.Copy(out, in)
	return name, nil		
}


func (s *Server) handleUserDelete(writer http.ResponseWriter, request *http.Request) {
	auth, err := middleware.Authentication(request.Context())
	id := auth.ID
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if id == 0 {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = s.userSvc.DeleteUser(request.Context(), auth)
	if err != nil {
		log.Print(err)	
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "text/plain")
	_, err = writer.Write([]byte("success"))
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}


	
}

func (s *Server) handleUserGetToken(writer http.ResponseWriter, request *http.Request) {
	log.Print("server.handleUserGetToken(): start")
	var item *types.User
	err := json.NewDecoder(request.Body).Decode(&item)
	if err != nil {
		log.Print(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	token, err := s.userSvc.Token(request.Context(), item.UserName, item.Password)
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

