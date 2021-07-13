package admin

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mijgona/instagram/cmd/app/middleware"
	"github.com/mijgona/instagram/types"
	"golang.org/x/crypto/bcrypt"
)

//Service описывает сервис работы с покупателями.
type Service struct {
	pool *pgxpool.Pool
}

//NewService создаёт сервис
func NewService(pool *pgxpool.Pool) *Service  {
	return &Service{pool: pool}
}

//Register регистрирует нового админа
func (s Service) Register(ctx context.Context, item *types.User, auth middleware.Auth) (*types.User, error)  {
	if item.ID==0 && auth.IsAdmin {
		hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Print("Register err:",err)
			return nil, err
		} 
		item.Roles=append(item.Roles, "ADMIN")
		err = s.pool.QueryRow(ctx, `
			INSERT INTO users(username, name, password, phone, bio, roles) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
		`, item.UserName, item.Name, hash, item.Phone, item.Bio, item.Roles).Scan(&item.ID)
		if err != nil {
			log.Print("Register err:",err)
			return nil, types.ErrInternal
		}	
		item.Photo = fmt.Sprint(item.ID)+".png"
		_, err = s.pool.Query(ctx, `
			UPDATE users SET photo=$1 WHERE id=$2 ;
		`, item.Photo, item.ID)
		if err != nil {
			log.Print("Register err:",err)
			return nil, types.ErrInternal
		}
	} 
		return item, nil
}

//GetAdmin выводит данные авторизованного админа
func (s Service) GetAdmin(ctx context.Context, auth middleware.Auth) (*types.User, error)  {
	if auth.IsAdmin{
	item := &types.User{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, username, name, phone, roles, photo, bio, active FROM users WHERE id=$1
		`, auth.ID).Scan(&item.ID, &item.UserName, &item.Name, &item.Phone, &item.Roles, &item.Photo, &item.Bio, &item.Active)
	if err != nil {
		log.Print("GetAdmin err:",err)
		return nil, types.ErrNoSuchUser
	}
	return item, nil
}
	return nil, types.ErrNotAdmin
}

//ActivePost активирует выбранный пост
func (s Service) ActivePost(ctx context.Context, auth middleware.Auth, postID int64) (*types.Post, error)  {
	item := &types.Post{}
	if !auth.IsAdmin {
			log.Print("ActivePost err:","not admin")
			return nil, types.ErrNotAdmin	
	}
	if postID != 0{
		err := s.pool.QueryRow(ctx, `
			UPDATE posts SET active = true WHERE id=$1 RETURNING id, content, photo;
			`, postID).Scan(&item.ID, &item.Content, &item.Photo)
		if err != nil {
			log.Print("ActiveUser err:",err)
			return nil, types.ErrNoSuchUser
		}
	} 
	return item, nil
}

//BlockPost блокирует выбранный пост
func (s Service) BlockPost(ctx context.Context, auth middleware.Auth, postID int64) (*types.Post, error)  {
	item := &types.Post{}
	if !auth.IsAdmin {
			log.Print("BlockPost err:","not admin")
			return nil, types.ErrNotAdmin	
	}
	if postID != 0{
		err := s.pool.QueryRow(ctx, `
			UPDATE posts SET active = false WHERE id=$1 RETURNING id, content, photo;
			`, postID).Scan(&item.ID, &item.Content, &item.Photo)
		if err != nil {
			log.Print("BlockUser err:",err)
			return nil, types.ErrNoSuchUser
		}
	} 
	return item, nil
}

//ActiveUser выводит список всех активных пользователей, если введен username активирует выбранного пользователя
func (s Service) ActiveUser(ctx context.Context, auth middleware.Auth, username string) ([]*types.User, error)  {
	items := []*types.User{}
	if !auth.IsAdmin {
			log.Print("ActiveUser err:","not admin")
			return nil, types.ErrNotAdmin	
	}
	if username != ""{
		item := types.User{}
		err := s.pool.QueryRow(ctx, `
			UPDATE users SET active = true WHERE username=$1 RETURNING id, username, name, phone, roles, photo, bio, active;
			`, username).Scan(&item.ID, &item.UserName, &item.Name, &item.Phone, &item.Roles, &item.Photo, &item.Bio, &item.Active)
		if err != nil {
			log.Print("ActiveUser err:",err)
			return nil, types.ErrNoSuchUser
		}
		items = append(items, &item)
	} else {
		rows, err := s.pool.Query(ctx, `
			SELECT id, username, name, phone, roles, photo, bio, active FROM users WHERE active  LIMIT 300;
			`)
		if err != nil {
			log.Print("ActiveUser err:",err)
			return nil, types.ErrNoSuchUser
		}
		for rows.Next() {
		item := types.User{}
			err = rows.Scan(&item.ID, &item.UserName, &item.Name, &item.Phone, &item.Roles, &item.Photo, &item.Bio, &item.Active)
			if err != nil {
				log.Print("ActiveUser err:",err)
				return nil, err
			}
			items = append(items, &item)
		}
		err = rows.Err()
		if err != nil {
			log.Print("ActiveUser err:",err)
			return nil, err
		}
		rows.Close()
	}
	return items, nil
}

//BlockUser выводит список всех заблокированных пользователей, если введен username блокирует выбранного пользователя
func (s Service) BlockUser(ctx context.Context, auth middleware.Auth, username string) ([]*types.User, error)  {
	items := []*types.User{}
	if !auth.IsAdmin {
		log.Print("BlockUser err:","not admin")
		return nil, types.ErrNotAdmin	
}
	if username != "" {
		item := types.User{}
		err := s.pool.QueryRow(ctx, `
			UPDATE users SET active = false WHERE username=$1 RETURNING id, username, name, phone, roles, photo, bio, active;
			`, username).Scan(&item.ID, &item.UserName, &item.Name, &item.Phone, &item.Roles, &item.Photo, &item.Bio, &item.Active)
		if err != nil {
			log.Print("BlockUser err:",err)
			return nil, types.ErrNoSuchUser
		}
		items = append(items, &item)
	} else {
		rows, err := s.pool.Query(ctx, `
			SELECT id, username, name, phone, roles, photo, bio, active FROM users WHERE NOT active
			`)
		if err != nil {
			log.Print("BlockUser err:",err)
			return nil, types.ErrNoSuchUser
		}
		for rows.Next() {
		item := types.User{}
			err = rows.Scan(&item.ID, &item.UserName, &item.Name, &item.Phone, &item.Roles, &item.Photo, &item.Bio, &item.Active)
			if err != nil {
				log.Print("BlockUser err:",err)
				return nil, err
			}
			items = append(items, &item)
		}
		err = rows.Err()
		if err != nil {
			log.Print("BlockUser err:",err)
			return nil, err
		}
		rows.Close()
	}
	return items, nil
}
//DeleteUser удаляет выбранного пользователя
func (s Service) DeleteUser(ctx context.Context, auth middleware.Auth, username string) ( error)  {
	if !auth.IsAdmin {
		log.Print("DeleteUser err:","not admin")
		return types.ErrNotAdmin	
	}
	id := int64(0)
	err := s.pool.QueryRow(ctx, `
	SELECT id FROM users WHERE username=$1;
	`, username).Scan(&id)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}
	
	_, err = s.pool.Query(ctx, `
	DELETE FROM likes WHERE user_id=$1;
	`, id)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}

	_, err = s.pool.Query(ctx, `
	DELETE FROM comments WHERE user_id=$1;
	`, id)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}

	_, err = s.pool.Query(ctx, `
	DELETE FROM follows WHERE user_id=$1;
	`, id)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}

	_, err = s.pool.Query(ctx, `
	DELETE FROM posts WHERE user_id=$1;
	`, id)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}
	_, err = s.pool.Query(ctx, `
	DELETE FROM tokens WHERE user_id=$1;
	`, id)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}
	_, err = s.pool.Query(ctx, `
	DELETE FROM users WHERE username=$1;
	`, username)
	if err != nil {
		log.Print("DeleteUser err:",err)
		return types.ErrNoSuchUser
	}
		return nil
}
	
//Token если пользователь ввел верные логин и пароль система присваевает ему токен для авторизации
func (s *Service) Token(ctx context.Context, name string, password string) (token string, err error){
	var hash string
	var id int64
	if err != nil {
		return "", types.ErrInvalidPassword
	}
	roles := []string{}
	err = s.pool.QueryRow(ctx, `SELECT id, password, roles FROM users where username = $1`, name).Scan(&id, &hash, &roles)
	if err == pgx.ErrNoRows {
		log.Print("Token err:",err)
		return "", types.ErrNoSuchUser
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return "", types.ErrInternal
	}
	buffer := make([]byte, 256)
	n,err := rand.Read(buffer)
	if n != len(buffer) || err!= nil {
		log.Print("Token err:",err)
		return "", types.ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO tokens(token, user_id, roles, expire) VALUES($1, $2, $3, $4);`, token, id, roles, time.Now().UTC().Add(time.Hour))
	if err != nil {
		log.Print("Token err:",err)
		return "", types.ErrInternal
	}

	return token, nil
}

// //IDByToken возвращает ИД для авторизации admin
// func (s *Service) IDByToken(ctx context.Context, token string) (middleware.Auth, error)  {
// 	var auth middleware.Auth
// 	var expire time.Time
// 	var roles []string
// 	err := s.pool.QueryRow(ctx, `
// 	SELECT user_id, expire, roles from tokens WHERE token = $1
// 	`, token).Scan(&auth.ID, &expire, &roles)
// 	if err == pgx.ErrNoRows {
// 		return  middleware.Auth{
// 			ID: 0,
// 			IsAdmin: false,
// 		}, nil
// 	}
// 	if err !=nil {
// 		return  middleware.Auth{
// 			ID: 0,
// 			IsAdmin: false,
// 		}, nil
// 	}
	
// 	timeNow := time.Now().UTC().Format("2006-01-02 15:04:05")
// 	timeEnd := expire.Format("2006-01-02 15:04:05")
// 	if timeNow > timeEnd {
// 		log.Print("token expired")
// 		return middleware.Auth{
// 			ID: 0,
// 			IsAdmin: false,
// 		}, types.ErrTokenExpired
// 	}

// 	for _, role := range roles {
// 		if role=="ADMIN"{auth.IsAdmin=true}
// 	}
// 	return auth, nil
// }