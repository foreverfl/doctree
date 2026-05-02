-- Idempotently ensure a repos row exists for the given root_path and
-- return its id. On first call we INSERT with derived defaults: bare and
-- worktrees paths follow gitt's <root>/.bare and <root>/.worktrees layout
-- convention, while the repo metadata fields (default_branch, language,
-- framework, compose_monorepo) start blank because the daemon doesn't
-- detect them yet — a later `gitt repos set` flow fills them in. On
-- subsequent calls the ON CONFLICT branch matches via the UNIQUE
-- root_path index and bumps updated_at, so callers always see a fresh
-- "last referenced" timestamp and the existing id is returned.
INSERT INTO repos (
  root_path, bare_path, worktrees_path, default_branch,
  language, framework, compose_monorepo,
  created_at, updated_at
) VALUES (?, ?, ?, '', '', '', 0, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (root_path) DO UPDATE SET updated_at = CURRENT_TIMESTAMP
RETURNING id;
