create table post_item
(
    id             text primary key,
    content        text not null,
    isItemFinished bool default false
)