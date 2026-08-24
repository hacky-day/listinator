-- +goose Up
-- +goose StatementBegin

-- ClientUpdatedAt is when the user actually made the edit (captured on the
-- frontend at the moment of the optimistic local write), not when the
-- request reached the server. It drives edit-time last-write-wins so a
-- delayed offline write can't clobber a newer edit purely by arriving
-- later.
ALTER TABLE entries ADD COLUMN client_updated_at datetime;

-- +goose StatementEnd
