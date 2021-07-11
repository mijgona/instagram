package comment

import (
	"context"
	"log"

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

func (s Service) GetComments(ctx context.Context, postID int64) ([]*types.Comment, error) {
	var items []*types.Comment
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
			
			items = append(items, &item)
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
			DELETE FROM comments WHERE id=$1 ;
		`, cmntID)
		if err != nil {
			log.Print("DeleteComment err:", err)
			return types.ErrInternal
		}	
	} 
	return nil
}