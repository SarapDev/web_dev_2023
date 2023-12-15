create table comments
(
    id         int auto_increment
        primary key,
    post_id    int          not null,
    author     varchar(255) not null,
    content    text         not null,
    created_at timestamp    not null
)
    collate = utf8mb4_general_ci;

create index comments_post_id_index
    on comments (post_id);

