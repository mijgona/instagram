
-- ТАБЛИЦА Пользователей
CREATE TABLE user
(    
    id          BIGSERIAL   PRIMARY KEY,
    name        TEXT        NOT NULL,
    password    TEXT    NOT NULL,
    photo       TEXT        NOT NULL,
    phone       TEXT,
    bio         TEXT,
    roles       []TEXT      NOT NULL,
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);

-- ТАБЛИЦА Постов
CREATE TABLE post
(    
    id          BIGSERIAL   PRIMARY KEY,    
    user_id     BIGINT      NOT NULL REFERENCES user,
    content     TEXT        NOT NULL,
    photo       TEXT        NOT NULL,
    tags        []TEXT,
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);
-- Таблица лайков
CREATE TABLE likes
(    
    id          BIGSERIAL   PRIMARY KEY,    
    post_id     BIGINT      NOT NULL REFERENCES post,
    user_id     BIGINT      NOT NULL REFERENCES user,
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);

-- Таблица комментавриев
CREATE TABLE comments
(    
    id          BIGSERIAL   PRIMARY KEY,    
    post_id     BIGINT      NOT NULL REFERENCES post,
    user_id     BIGINT      NOT NULL REFERENCES user,
    comment     TEXT        NOT NULL,    
    active      BOOLEAN     NOT NULL    DEFAULT TRUE,
    created     TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);


--Таблица подписок
CREATE TABLE follows
(     
    id              BIGSERIAL   PRIMARY KEY,    
    user_id     BIGINT      NOT NULL REFERENCES user,
    followed_id     BIGINT      NOT NULL REFERENCES user,
    active          BOOLEAN     NOT NULL    DEFAULT TRUE,
    created         TIMESTAMP   NOT NULL    DEFAULT CURRENT_TIMESTAMP
);
