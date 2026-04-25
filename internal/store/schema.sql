CREATE TABLE IF NOT EXISTS worktrees (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  repo_name   TEXT    NOT NULL,
  branch      TEXT    NOT NULL,
  path        TEXT    NOT NULL UNIQUE,
  created_at  INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS ports (
  worktree_id INTEGER NOT NULL REFERENCES worktrees(id) ON DELETE CASCADE,
  service     TEXT    NOT NULL,
  host_port   INTEGER NOT NULL UNIQUE,
  PRIMARY KEY (worktree_id, service)
);
