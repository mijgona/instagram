package post

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
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


func (s Service) GetPost(ctx context.Context, postID int64, id int64) (*types.Post, error) {
	item := &types.Post{}
	err := s.pool.QueryRow(ctx, `SELECT p.id, u.username, p.content, p.photo, p.tags, p.active, u.photo, COALESCE(count(p.id),0) from users u
	JOIN posts p ON u.id=p.user_id and p.id=$1 and p.active
	JOIN (
		SELECT l.id, l.post_id from likes l, posts p 
		WHERE l.post_id=p.id and l.active
		group BY l.id
		) ss ON ss.post_id=p.id
	GROUP BY u.id, p.id;
	`, postID).Scan(&item.ID, &item.Author.Name, &item.Content, &item.Photo, &item.Tags, &item.Active, &item.Author.Avatar,  &item.Likes)
	if err != nil {
		log.Print("GetPost err:", err)
		return nil, types.ErrNoSuchPost
	}
	like := 0
	err = s.pool.QueryRow(ctx, `
		SELECT user_id from likes WHERE user_id=$1 AND post_id=$2
		`, id, item.ID).Scan(&like)
	if err != nil && err != pgx.ErrNoRows {
		log.Print("GetPost err:", err)
		return nil, types.ErrInternal
	}
	if like!=0 {item.LikedByMe = true}
	return item, nil
}

func (s Service) GetAllPost(ctx context.Context, userID int64) ([]types.Post, error) {
	items := []types.Post{}
	rows, err := s.pool.Query(ctx, `
	SELECT p.id, u.username, p.content, p.photo, p.tags, p.active, u.photo, COALESCE(count(p.id),0) likes from users u
	JOIN posts p ON p.active and p.user_id =$1 and p.user_id=u.id
	JOIN (
		SELECT l.id, l.post_id from likes l, posts p 
		WHERE l.post_id=p.id
		group BY l.id
		) ss ON ss.post_id=p.id
	GROUP BY u.id, p.id
	ORDER BY p.created DESC
	`, userID)
	if err != nil {
		log.Print("GetAllPosts err:", err)
		return nil, types.ErrNoSuchPost
	}

	for rows.Next() {
		item := types.Post{}
		err = rows.Scan(&item.ID, &item.Author.Name, &item.Content, &item.Photo, &item.Tags, &item.Active, &item.Author.Avatar,  &item.Likes)
		if err != nil {
			log.Print("GetAllPosts err:",err)
			return nil, err
		}
		like := 0
		err = s.pool.QueryRow(ctx, `
			SELECT user_id from likes WHERE user_id=$1 AND post_id=$2
			`, userID, item.ID).Scan(&like)
		if err != nil && err != pgx.ErrNoRows {
			log.Print("GetAllPosts err:", err)
			return nil, types.ErrInternal
		}
		if like !=0 {item.LikedByMe = true}
		items = append(items, item)
	}
	rows.Close()
	return items, nil
}

func (s Service) NewPost(ctx context.Context, userID int64, item *types.Post) (*types.Post, error) {
	ctx = context.Background()
	if item.ID==0 {
		err := s.pool.QueryRow(ctx, `
			INSERT INTO posts(user_id, content, photo, tags) VALUES ($1, $2, $3, $4) RETURNING id;
		`, userID, item.Content, item.Photo, item.Tags).Scan(&item.ID)
		if err != nil {
			log.Print("NewPost err:", err)
			return nil, types.ErrInternal
		}	
	} 
	return item, nil
}

func (s Service) DeletePost(ctx context.Context, postID int64, id int64) (error) {
	
	if postID !=0 {
		postUserID := int64(0)
		err := s.pool.QueryRow(ctx, `
		SELECT user_id FROM posts WHERE id=$1;
		`, postID).Scan(&postUserID)
		if postUserID != id {
			log.Print("NewPost err:", types.ErrNotAdmin)
			return types.ErrNotAdmin
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