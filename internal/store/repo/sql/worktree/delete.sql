DELETE FROM worktrees
WHERE repo_id = (SELECT id FROM repos WHERE root_path = ?)
  AND branch_name = ?
RETURNING id;
