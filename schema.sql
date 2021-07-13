
-- ТАБЛИЦА Пользователей
CREATE TABLE users (    
    id          BIGSERIAL   PRIMARY KEY,
    username    TEXT        NOT NULL UNIQUE,
    password    TEXT        NOT NULL,
    name    TEXT,    
    photo       TEXT,
    phone       TEXT,
    bio         TEXT,
    roles       TEXT[]      NOT NULL DEFAULT '{}',
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);


--Таблица подписок
CREATE TABLE follows
(     
    id              BIGSERIAL   PRIMARY KEY,    
    user_id         BIGINT      NOT NULL REFERENCES users,
    followed_id     BIGINT      NOT NULL REFERENCES users,
    active          BOOLEAN     NOT NULL    DEFAULT TRUE,
    created         TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);

-- ТАБЛИЦА Токенов Пользователей
CREATE TABLE tokens (    
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users,
    token       TEXT        NOT NULL,
    roles       TEXT[]      NOT NULL DEFAULT '{}',
    expire      TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 hour',
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);
-- ТАБЛИЦА Постов
CREATE TABLE posts
(    
    id          BIGSERIAL   PRIMARY KEY,    
    user_id     BIGINT      NOT NULL REFERENCES users,
    content     TEXT        NOT NULL    DEFAULT ' ',
    photo       TEXT        NOT NULL,
    tags        TEXT[]      NOT NULL DEFAULT '{}',
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);
-- Таблица лайков
CREATE TABLE likes
(    
    id          BIGSERIAL   PRIMARY KEY,    
    post_id     BIGINT      NOT NULL REFERENCES posts,
    user_id     BIGINT      NOT NULL REFERENCES users,
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);

-- Таблица комментавриев
CREATE TABLE comments
(    
    id          BIGSERIAL   PRIMARY KEY,    
    post_id     BIGINT      NOT NULL REFERENCES posts,
    user_id     BIGINT      NOT NULL REFERENCES users,
    comment     TEXT        NOT NULL,    
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);



-- DROP TABLE users;
-- DROP TABLE tokens;
-- drop table follows;
-- drop table comments;
-- drop table likes;
-- drop table posts;