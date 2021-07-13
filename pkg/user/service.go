package user

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mijgona/instagram/cmd/app/middleware"
	"github.com/mijgona/instagram/types"
	"golang.org/x/crypto/bcrypt"
)

//Service описывает сервис работы с покльзователями.
type Service struct {
	pool *pgxpool.Pool
}

//NewService создаёт сервис
func NewService(pool *pgxpool.Pool) *Service  {
	return &Service{pool: pool}
}

func (s Service) EditUser(ctx context.Context, item *types.User, auth middleware.Auth) (*types.User, error)  {
	ctx = context.Background()
	if item.ID==0 {
		hash, err := bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Print(err)
			return nil, err
		} 
		err = s.pool.QueryRow(ctx, `
			INSERT INTO users(username, name, password, phone, bio, photo) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;
		`, item.UserName, item.Name, hash, item.Phone, item.Bio, item.Photo).Scan(&item.ID)
		if err != nil {
			log.Print("EditUser err:", err)
			return nil, types.ErrInternal
		}
	} else {
		if item.ID != auth.ID{
			log.Print("can edit other users")
			return nil, types.ErrNotAdmin			
		}
		var hash []byte
		var err error
		if item.Password !="" {		
			hash, err = bcrypt.GenerateFromPassword([]byte(item.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Print("EditUser err:",err)
				return nil, err
			}
			_, err = s.pool.Query(ctx, `
				UPDATE users SET password=$1 WHERE id=$2 ;
			`, hash, item.ID)
			if err != nil {
				log.Print("EditUser err:",err)
				return nil, types.ErrInternal
			}
			return item, nil	
		}
		if item.Photo != "" {
			_, err = s.pool.Query(ctx, `
				UPDATE users SET photo=$1 WHERE id=$2 ;
			`, item.Photo, item.ID)
			if err != nil {
				log.Print("EditUser err:",err)
				return nil, types.ErrInternal
			}
		}
		if item.Name != "" && item.Phone != "" && item.Bio !=""{		
			_, err = s.pool.Query(ctx, `
				UPDATE users SET name = $2, phone = $3, bio =$4 WHERE id=$1 ;
			`, item.ID, item.Name, item.Phone, item.Bio)
			if err != nil {
				log.Print("EditUser err:",err)
				return nil, types.ErrInternal
			}	
		}	
	}
		return item, nil
}


//удаляет пользователя
func (s Service) DeleteUser(ctx context.Context, auth middleware.Auth) (error)  {		
		_, err := s.pool.Query(ctx, `
			DELETE FROM users_tokens WHERE user_id=$1;
			`, auth.ID)
		if err != nil {
			log.Print("DeleteUser err:",err)
			return types.ErrNoSuchUser
		}	
		_, err = s.pool.Query(ctx, `
			DELETE FROM follows WHERE user_id=$1;
			`, auth.ID)
		if err != nil {
			log.Print("DeleteUser err:",err)
			return types.ErrNoSuchUser
		}	
		_, err = s.pool.Query(ctx, `
			DELETE FROM follows WHERE followed_id=$1;
			`, auth.ID)
		if err != nil {
			log.Print("DeleteUser err:",err)
			return types.ErrNoSuchUser
		}
		_, err = s.pool.Query(ctx, `
			DELETE FROM users WHERE id=$1;
			`, auth.ID)
		if err != nil {
			log.Print("DeleteUser err:",err)
			return types.ErrNoSuchUser
		}		
		return nil
}



//GetUser возвращает данные авторизованного пользователя, если введён username то его данные 
func (s Service) GetUser(ctx context.Context, auth middleware.Auth, username string) (*types.User, error)  {
	id := auth.ID
	//при наличии имени пользователя, найти его ИД и в дальнейшем использовать его
	if username != "" {
		err := s.pool.QueryRow(ctx, `
		SELECT id FROM users WHERE username=$1 and active
		`, username).Scan(&id)
		if err != nil {
			if err == pgx.ErrNoRows{			
				log.Print("GetUser err:",err)
				return nil, types.ErrNoSuchUser
			}
			log.Print("GetUser err:",err)
			return nil, types.ErrNoSuchUser
		}
	}
	item := &types.User{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, username, name, phone, photo, bio FROM users WHERE id=$1 and active
		`, id).Scan(&item.ID, &item.UserName, &item.Name, &item.Phone, &item.Photo, &item.Bio)
	if err != nil {
		if err == pgx.ErrNoRows{			
			log.Print("GetUser err:",err)
			return nil, types.ErrNoSuchUser
		}
		log.Print("GetUser err:",err)
		return nil, types.ErrNoSuchUser
	}

	rows, err := s.pool.Query(ctx, `
		SELECT f.id, u.photo, u.name, u.username, f.active, f.created FROM users u, follows f
		WHERE f.user_id=$1 AND f.followed_id=u.id AND f.active
	`, id)
	if err != nil {
		log.Print("GetUser err:",err)
		return nil, types.ErrNoSuchUser
	}

	for rows.Next() {
		follow := types.Follow{}
		err = rows.Scan(&follow.ID, &follow.Avatar, &follow.Name, &follow.UserName, &follow.Active, &follow.Created)
		if err != nil {
			log.Print("GetUser err:",err)
			return nil, err
		}
		item.Follows = append(item.Follows, follow)
	}
	err = rows.Err()
	if err != nil {
		log.Print("GetUser err:",err)
		return nil, types.ErrInternal
	}
	rows.Close()

	rows, err = s.pool.Query(ctx, `
		SELECT f.id, u.photo, u.name, u.username, f.active, f.created FROM users u, follows f
		WHERE u.id=f.user_id AND f.followed_id=$1 AND f.active
	`, id)
	if err != nil {
		log.Print("GetUser err:",err)
		return nil,  types.ErrNoSuchUser
	}

	for rows.Next() {
		follow := types.Follow{}
		err = rows.Scan(&follow.ID, &follow.Avatar, &follow.Name, &follow.UserName, &follow.Active, &follow.Created)
		if err != nil {
			log.Print("GetUser err:",err)
			return nil, types.ErrInternal
		}
		item.Followers = append(item.Followers, follow)
	}
	rows.Close()
	return item, nil
}


//Follow подписывает авторизованного пользователя на выбранного
func (s Service) Follow(ctx context.Context, item *types.Follow, id int64) (*types.Follow, error)  {
	err := s.pool.QueryRow(ctx, `
	SELECT id FROM follows WHERE user_id=$1 and followed_id=$2;
	`, id, item.UserID).Scan(&item.ID)

	if err != nil && err != pgx.ErrNoRows {
		log.Print("Follow err:",err)
		return nil, types.ErrInternal
	}
	if item.ID==0{
		err = s.pool.QueryRow(ctx, `
		INSERT INTO follows(user_id, followed_id) VALUES ($1, $2) RETURNING id;
		`, id, item.UserID).Scan(&item.ID)

		if err != nil {
			log.Print("Follow err:",err)
			return nil, types.ErrInternal
		}
	} else {
		_, err = s.pool.Query(ctx, `
		UPDATE follows SET active = NOT active WHERE id=$1;
		`, item.ID)
		if err != nil {
			log.Print("Follow err:",err)
			return nil, types.ErrInternal
		}
	}
	return item, nil
}


func (s *Service) Token(ctx context.Context, name string, password string) (token string, err error){
	var hash string
	var id int64
	if err != nil {
		log.Print("Token err:",err)
		return "", types.ErrInvalidPassword
	}	
	err = s.pool.QueryRow(ctx, `SELECT id, password FROM users where username = $1`, name).Scan(&id, &hash)
	if err == pgx.ErrNoRows {
		log.Print("Token err:",err)
		return "", types.ErrNoSuchUser
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		log.Print("Token err:",err)
		return "", types.ErrInternal
	}
	buffer := make([]byte, 256)
	n,err := rand.Read(buffer)
	if n != len(buffer) || err!= nil {
		log.Print("Token err:",err)
		return "", types.ErrInternal
	}

	token = hex.EncodeToString(buffer)
	_, err = s.pool.Exec(ctx, `INSERT INTO tokens(token, user_id, expire) VALUES($1, $2, $3);`, token, id, time.Now().UTC().Add(time.Hour))
	if err != nil {
		log.Print("Token err:",err)
		return "", types.ErrInternal
	}

	return token, nil
} 


func (s *Service) IDByToken(ctx context.Context, token string) (middleware.Auth, error)  {
	var auth middleware.Auth
	var expire time.Time
	var start time.Time
	var roles []string
	err := s.pool.QueryRow(ctx, `
	SELECT user_id, expire, created, roles from tokens WHERE token = $1
	`, token).Scan(&auth.ID, &expire, &start, &roles)
	if err == pgx.ErrNoRows {
		log.Print("token not found")
		return  middleware.Auth{
			ID: 0,
			IsAdmin: false,
		}, nil
	}
	if err !=nil {
		log.Print("token Err:", err)
		return  middleware.Auth{
			ID: 0,
			IsAdmin: false,
		}, nil
	}
	timeNow := time.Now().UTC().Format("2006-01-02 15:04:05")
	timeEnd := expire.Format("2006-01-02 15:04:05")
	if timeNow > timeEnd {
		log.Print("token expired")
		return middleware.Auth{
			ID: 0,
			IsAdmin: false,
		}, types.ErrTokenExpired
	}

	for _, role := range roles {
		if role=="ADMIN"{auth.IsAdmin=true}
	}
	return auth, nil
}
