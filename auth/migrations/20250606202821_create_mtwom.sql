-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_groups
(
    user_id  BIGINT REFERENCES users (id) ON DELETE CASCADE,
    group_id BIGINT REFERENCES pass_groups (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, group_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_groups;
-- +goose StatementEnd
