SELECT p.id, u.username, p.content, p.photo, p.tags, p.active, u.photo, COALESCE(count(p.id),0) likes from users u
	JOIN posts p ON p.active and p.user_id =17 and p.user_id=u.id
	JOIN (
		SELECT l.id, l.post_id from likes l, posts p 
		WHERE l.post_id=p.id
		group BY l.id
		) ss ON ss.post_id=p.id
	GROUP BY u.id, p.id
	ORDER BY p.created DESC













SELECT (
    SELECT u.username FROM users u  WHERE u.id=fu.followed_id
    ), p.content, p.tags, count(ss.id) FROM users u
JOIN posts p on p.user_id=u.id 
JOIN (
    SELECT l.id, p.id post_id from likes l, posts p 
    WHERE l.post_id=p.id 
    group BY l.id, p.id
    ) ss ON ss.post_id=p.id
JOIN (
    SELECT u.username, f.user_id, f.followed_id, f.active from follows f, users u
    WHERE f.user_id=u.id and u.id=20
    GROUP BY f.id, u.id
    ) fu ON fu.followed_id=u.id 
GROUP BY fu.followed_id, p.content, p.tags, p.created
ORDER BY p.created DESC




JOIN (
    SELECT f.user_id, f.followed_id, u.username from follows f, users u 
    WHERE f.user_id=u.id
    group BY u.id, f.id
    ) fu ON fu.user_id=u.id and fu.followed_id=17
JOIN posts p ON p.active and p.user_id=fu.followed_id 

SELECT c.id, u.name, u.photo, c.comment, c.active FROM users u, comments c 
			WHERE c.user_id=u.id AND c.post_id=1 ORDER BY c.created DESC



JOIN (
    select l.post_id from likes l, posts p where l.post_id=p.id and l.active
    GROUP BY l.post_id
    )ss ON ss.post_id=p.id
group By p.id


--Считывает сумму всех сдeлок менеджеров и сортирует их по убыванию
SELECT m.id, m.name, m.salary*1000 salary, m.plan*1000 plan, COALESCE(ss.sum,0) total FROM managers m
LEFT JOIN (
    SELECT s.manager_id, sum(sp.price*sp.qty) sum FROM
    sales s, sale_positions sp 
    WHERE s.id=sp.sale_id
    GROUP BY  s.manager_id
) ss ON m.id=ss.manager_id
ORDER BY COALESCE(ss.sum,0) DESC

--Находит ТОП-3 продукта
SELECT p.id, p.name, ss.sum total from products p
JOIN (
    SELECT sp.product_id, SUM(sp.price*sp.qty) FROM
    sales s, sale_positions sp 
    WHERE s.id=sp.sale_id
    GROUP BY  sp.product_id
) ss ON p.id=ss.product_id 
ORDER BY ss.sum DESC LIMIT 3 


SELECT p.id, p.content, p.photo, p.tags, p.active, u.photo, u.username, COUNT(l.user_id) from posts p
right join users u on u.id=17
right join likes l on l.post_id=p.id and l.user_id=u.id



	FROM posts p, users u, likes l WHERE p.user_id=17 AND l.post_id=p.id
