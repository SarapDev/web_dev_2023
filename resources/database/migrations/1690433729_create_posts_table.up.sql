create table posts
(
    id         int auto_increment
        primary key,
    title      varchar(255) not null,
    content    text         not null,
    created_at timestamp    null
)
    collate = utf8mb4_general_ci;

