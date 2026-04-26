package daemon

import "fmt"

// handleRegisterWorktree persists a worktree row from the request args.
// Required args: repo_root, repo_name, branch_name, safe_branch_name,
// worktree_path. The unique constraint on (repo_root, branch_name) and
// worktree_path is enforced by the store; conflicts surface as the error.
func (server *server) handleRegisterWorktree(req Request) Response {
	repoRoot, err := stringArg(req, "repo_root")
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	repoName, err := stringArg(req, "repo_name")
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	branchName, err := stringArg(req, "branch_name")
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	safeBranchName, err := stringArg(req, "safe_branch_name")
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	worktreePath, err := stringArg(req, "worktree_path")
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}

	worktree, err := server.store.InsertWorktree(repoRoot, repoName, branchName, safeBranchName, worktreePath)
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Data: map[string]any{"worktree": worktree}}
}

// handleListWorktrees returns every persisted worktree row.
func (server *server) handleListWorktrees(_ Request) Response {
	worktrees, err := server.store.ListWorktrees()
	if err != nil {
		return Response{OK: false, Error: err.Error()}
	}
	return Response{OK: true, Data: map[string]any{"worktrees": worktrees}}
}

func stringArg(req Request, name string) (string, error) {
	raw, ok := req.Args[name]
	if !ok {
		return "", fmt.Errorf("missing arg: %s", name)
	}
	value, ok := raw.(string)
	if !ok {
		return "", fmt.Errorf("arg %s must be a string", name)
	}
	if value == "" {
		return "", fmt.Errorf("arg %s must not be empty", name)
	}
	return value, nil
}