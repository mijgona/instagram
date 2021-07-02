package types

import (
	"errors"
	"time"
)

//ErrNotFound возвращается, когда пользователь не найден
var ErrNoSuchUser = errors.New("no such user")
//ErrInternal возвращается, когда произошла внутренная ошибка
var ErrInternal = errors.New("internal error")
//ErrPhoneUsed возвращается, когда телефон уже зарегистрирован
var ErrUserExist = errors.New("user already registered")
//ErrInvalidPassword возвращается, когда пароль не введён
var ErrInvalidPassword= errors.New("invalid password")
//ErrTokenNotFound возвращается, когда токен не найден
var ErrTokenNotFound= errors.New("token not found")
//ErrTokenExpired возвращается, когда у токена вышло время
var ErrTokenExpired= errors.New("token expired")
var ErrNotAdmin = errors.New("not admin")




//User представляет информацию о пользователе.
type User struct {
	ID			int64		`json:"id"`
	Name		string		`json:"name"`
	Password	string		`json:"password"`
	Photo		string		`json:"photo"`
	Phone		string		`json:"phone"`
	Bio			string		`json:"bio"`
	Roles		[]string	`json:"roles"`
	Active		bool		`json:"active"`
	Created		time.Time	`json:"created"`
}

//Follow представляет информацию подписках.
type Follow struct {
	ID			int64		`json:"id"`	
	Avatar		string		`json:"avatar"`
	Name		string		`json:"name"`
	Active		bool		`json:"active"`
	Created		time.Time	`json:"created"`
	
}


//Post представляет информацию о поcте.
type Post struct {
	ID			int64		`json:"id"`
	Author struct{
		Avatar	string		`json:"avatar"`
		Name	string		`json:"name"`
	}
	Content		string		`json:"content"`
	Photo		string		`json:"photo"`
	Likes		string		`json:"likes"`
	LikedByMe	bool		`json:"liked_by_me"`
	Tags		[]string	`json:"tags"`
	Active		bool		`json:"active"`
	Created		time.Time	`json:"created"`
}

type Comment struct {
	ID		int64			`json:"id"`	
	Author struct{
		Avatar	string		`json:"avatar"`
		Name	string		`json:"name"`
	}
	PostID	int64			`json:"post_id"`
	Comment	int64			`json:"comment"`
	Active		bool		`json:"active"`
	Created		time.Time	`json:"created"`
}

//Token описывает токен для покупателя
type Token struct {
	Token	string		`json:"token"`
}