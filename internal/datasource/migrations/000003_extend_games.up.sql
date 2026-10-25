ALTER TABLE games
ADD COLUMN state INT NOT NULL DEFAULT 0,
    ADD COLUMN turn_player_id UUID,
    ADD COLUMN winner_id UUID,
    ADD COLUMN vs_computer BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN cross_player_id UUID,
    ADD COLUMN nought_player_id UUID;
UPDATE games
SET vs_computer = true,
    state = 4;