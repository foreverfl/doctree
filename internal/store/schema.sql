CREATE TABLE IF NOT EXISTS worktrees (
  id               INTEGER PRIMARY KEY AUTOINCREMENT,

  repo_root        TEXT NOT NULL,
  repo_name        TEXT NOT NULL,
  branch_name      TEXT NOT NULL,
  safe_branch_name TEXT NOT NULL,
  worktree_path    TEXT NOT NULL,

  status           TEXT NOT NULL DEFAULT 'created',

  created_at       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at       TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,

  UNIQUE(repo_root, branch_name),
  UNIQUE(worktree_path)
);

CREATE TABLE IF NOT EXISTS ports (
  worktree_id INTEGER NOT NULL REFERENCES worktrees(id) ON DELETE CASCADE,
  service     TEXT    NOT NULL,
  host_port   INTEGER NOT NULL UNIQUE,
  PRIMARY KEY (worktree_id, service)
);