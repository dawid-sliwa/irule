-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';


alter table tags
add column created_by uuid,
add constraint tags_user_fkey foreign key (created_by) references users(id) on delete cascade;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

alter table tags
drop column created_by;

-- +goose StatementEnd
