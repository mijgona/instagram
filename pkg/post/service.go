package post

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mijgona/instagram/cmd/app/middleware"
	"github.com/mijgona/instagram/types"
)

//Service описывает сервис работы с покльзователями.
type Service struct {
	pool *pgxpool.Pool
}

//NewService создаёт сервис
func NewService(pool *pgxpool.Pool) *Service  {
	return &Service{pool: pool}
}


func (s Service) GetPost(ctx context.Context, postID int64, auth middleware.Auth) (*types.Post, error) {
	item := &types.Post{}
	//При наличии лайков выдет пост с количеством лайков
	//при отсуствии лайков выдаёт ошибку pgx.ErrNoRows
	err := s.pool.QueryRow(ctx, `SELECT p.id, u.username, p.content, p.photo, p.tags, p.active, u.photo, COALESCE(count(p.id),0) from users u
	JOIN posts p ON u.id=p.user_id and p.id=$1 and p.active
	JOIN (
		SELECT l.id, l.post_id from likes l, posts p 
		WHERE l.post_id=p.id and l.active
		group BY l.id
		) ss ON ss.post_id=p.id
	GROUP BY u.id, p.id;
	`, postID).Scan(&item.ID, &item.Author.Name, &item.Content, &item.Photo, &item.Tags, &item.Active, &item.Author.Avatar,  &item.Likes)
	// при отсуствии лайков меняем запрос
	if err == pgx.ErrNoRows {
		err := s.pool.QueryRow(ctx, `SELECT p.id, u.username, p.content, p.photo, p.tags, p.active, u.photo from users u
		JOIN posts p ON p.active and p.id =$1
		WHERE p.user_id=u.id
		GROUP BY u.id, p.id;`, postID).Scan(&item.ID, &item.Author.Name, &item.Content, &item.Photo, &item.Tags, &item.Active, &item.Author.Avatar)
		if err != nil {	
			log.Print("GetPost err:", err)
			return nil, types.ErrInternal
		}
		
		comments, err := s.GetComments(ctx, postID)
		if err != nil {	
			log.Print("GetPost err:", err)
			return nil, types.ErrInternal
		}
		item.Comments=comments	
		return item, nil
	}
	if err != nil{
		log.Print("GetPost err:", err)
		return nil, types.ErrNoSuchPost
	}
	//Получаем все комментарии	
	comments, err := s.GetComments(ctx, postID)
	if err != nil {	
		log.Print("GetPost err:", err)
		return nil, types.ErrInternal
	}
	item.Comments=comments
	// получаем likedByMe
	like := 0
	err = s.pool.QueryRow(ctx, `
		SELECT user_id from likes WHERE user_id=$1 AND post_id=$2 and active
		`, auth.ID, item.ID).Scan(&like)
	if err != nil && err != pgx.ErrNoRows {
		log.Print("GetPost err:", err)
		return nil, types.ErrInternal
	}
	if like!=0 {item.LikedByMe = true}
	return item, nil
}

func (s Service) GetAllPost(ctx context.Context, auth middleware.Auth, username string) ([]*types.Post, error) {
	id := auth.ID
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

	items := []*types.Post{}
	rows, err := s.pool.Query(ctx, `
	SELECT p.id FROM posts p WHERE p.user_id=$1;
	`, id)
	if err != nil {
		log.Print("GetAllPosts err:", err)
		return nil, types.ErrNoSuchPost
	}

	for rows.Next() {
		item := &types.Post{}
		err = rows.Scan(&item.ID)
		if err != nil {	
			log.Print("GetPost err:", err)
			return nil, types.ErrInternal
		}	
		item, err = s.GetPost(context.Background(), item.ID, auth)
		if err != nil {	
			log.Print("GetPost err:", err)
			return nil, types.ErrInternal
		}
		
		items = append(items, item)
	}
	rows.Close()
	return items, nil
}

func (s Service) NewPost(ctx context.Context, auth middleware.Auth, item *types.Post) (*types.Post, error) {
	if item.ID==0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO posts(user_id, content, photo, tags) VALUES ($1, $2, $3, $4) RETURNING id;
		`, auth.ID, item.Content, item.Photo, item.Tags).Scan(&item.ID)
		if err != nil {
			log.Print("NewPost err:", err)
			return nil, types.ErrInternal
		}
		_, err = s.pool.Query(ctx, `		
			INSERT INTO likes(user_id, post_id)VALUES ($1, $2);
		`, auth.ID, item.ID)
		if err != nil {
			log.Print("NewPost err:", err)
			return nil, types.ErrInternal
		}	
	} 
	return item, nil
}

func (s Service) DeletePost(ctx context.Context, postID int64, auth middleware.Auth) (error) {
	if postID !=0 {
		postUserID := int64(0)
		err := s.pool.QueryRow(ctx, `
		SELECT user_id FROM posts WHERE id=$1;
		`, postID).Scan(&postUserID)
		if postUserID != auth.ID{
			if !auth.IsAdmin{
				log.Print("NewPost err:", types.ErrNotAdmin)
				return types.ErrNotAdmin
			}
		}
		if err != nil {
			log.Print("NewPost err:", err)
			return types.ErrInternal
		}
		
		_, err = s.pool.Query(ctx, `
			DELETE FROM likes WHERE post_id=$1 ;
		`, postID)
		if err != nil {
			log.Print("NewPost err:", err)
			return types.ErrInternal
		}
		_, err = s.pool.Query(ctx, `
			DELETE FROM posts WHERE id=$1 ;
		`, postID)
		if err != nil {
			log.Print("NewPost err:", err)
			return types.ErrInternal
		}	
	} 
	return nil
}

func (s Service) LikePost(ctx context.Context, postID int64, userID int64) (error) {
	id := 0
	err := s.pool.QueryRow(ctx, `
	SELECT id FROM likes WHERE user_id=$1 and post_id=$2;
	`, userID, postID).Scan(&id)
	if err != nil && err!= pgx.ErrNoRows {
		log.Print("LikePost err:",err)
		return types.ErrInternal
	}
	if id == 0{
		err = s.pool.QueryRow(ctx, `
		INSERT INTO likes(user_id, post_id) VALUES ($1, $2) RETURNING id;
		`, userID, postID).Scan(&id)

		if err != nil {
			log.Print("LikePost err:",err)
			return types.ErrInternal
		}
	} else {
		_, err = s.pool.Query(ctx, `
		UPDATE likes SET active = NOT active WHERE id=$1;
		`, id)
		if err != nil {
			log.Print("LikePost err:",err)
			return types.ErrInternal
		}
	}
	return nil
}


func (s Service) NewComment(ctx context.Context, userID int64, item *types.Comment) (*types.Comment, error) {
	if item.ID==0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO comments(user_id, post_id, comment) VALUES ($1, $2, $3) RETURNING id;
		`, userID, item.PostID, item.Comment).Scan(&item.ID)
		if err != nil {
			log.Print("NewComment err:", err)
			return nil, types.ErrInternal
		}	
	} 
	return item, nil
}

func (s Service) GetComments(ctx context.Context, postID int64) ([]types.Comment, error) {
	var items []types.Comment
	if postID!=0 {
		rows, err := s.pool.Query(context.Background(), `
			SELECT c.id, u.name, u.photo, c.comment, c.active FROM users u, comments c 
			WHERE c.user_id=u.id AND c.post_id=$1 ORDER BY c.created DESC
		`, postID)
		if err != nil {
			log.Print("GetComment err:", err)
			return nil, types.ErrInternal
		}	

		for rows.Next() {
			item := types.Comment{
				PostID:      postID,
			}
			err = rows.Scan(&item.ID, &item.Author.Name, &item.Author.Avatar, &item.Comment, &item.Active)
			if err != nil {
				log.Print("GetComment err:",err)
				return nil, err
			}
			
			items = append(items, item)
		}
		rows.Close()
	} 
	return items, nil
}

func (s Service) DeleteComment(ctx context.Context, cmntID int64, id int64) (error) {
	
	if cmntID !=0 {
		commentUserID := int64(0)
		err := s.pool.QueryRow(ctx, `
		SELECT user_id FROM comments WHERE id=$1;
		`, cmntID).Scan(&commentUserID)
		if commentUserID != id {
			log.Print("DeleteComment err:", types.ErrNotAdmin)
			return types.ErrNotAdmin
		}
		if err != nil {
			log.Print("DeleteComment err:", err)
			return types.ErrInternal
		}

		
		_, err = s.pool.Query(ctx, `
			DELETE FROM comments WHERE id=$1;
		`, cmntID)
		if err != nil {
			log.Print("DeleteComment err:", err)
			return types.ErrInternal
		}	
	} 
	return nil
}