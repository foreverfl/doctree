# agent-worktree (`aw`)

`git worktree` + per-worktree `docker compose` 스택을 한 바이너리로 묶은 CLI.
zshrc에 흩어져 있던 `worktree-add` / `worktree-build` / `worktree-remove` 함수를 대체.

여러 worktree에서 동시에 컨테이너를 띄울 때 발생하는 **호스트 포트 race**를
SQLite 백엔드 데몬에서 직렬화해서 막는 게 핵심.

## Commands

| 명령어 | 동작 |
| --- | --- |
| `aw on` | 데몬 기동 (`~/.aw/aw.sock`, `~/.aw/aw.db`). 다른 모든 명령은 데몬이 떠있어야 동작 |
| `aw off` | 데몬 종료 |
| `aw add <branch>` | `<repo>/../.worktrees/<repo>/<branch>` 에 worktree 생성, 포트 할당, `.env.worktree` 작성 |
| `aw remove` | 현재 worktree의 `docker compose down` + worktree 폴더 삭제 |

데몬이 안 떠있으면 `aw add`/`aw remove`는 즉시 에러로 끊고 `aw on` 안내. (auto-start 안 함)

## Install

```bash
git clone <this-repo> ~/code/agent-worktree
cd ~/code/agent-worktree
go build -o ~/.local/bin/aw .   # PATH 안의 어디든 OK
```

## zshrc

```bash
# 셸 시작 시 데몬 한 번 띄움 (이미 떠있으면 noop)
aw on >/dev/null 2>&1

alias aw-add='aw add'
alias aw-on='aw on'
alias aw-off='aw off'
alias aw-remove='aw remove'
```

## Repo convention

`aw add` 가 호출되는 repo는 다음을 가정:

- `infra/docker/compose.local.yml` — compose 파일
- `infra/docker/.env.local` (선택) — 공용 기본값
- `infra/docker/.env.worktree` — `aw add` 가 worktree별 포트 매핑을 써넣는 파일

`compose.local.yml` 쪽 변경 사항:

- `container_name` 라인 제거 → compose가 `<project>-<service>-<n>` 으로 자동 부여
- 호스트 포트는 `"${POSTGRES_HOST_PORT:-5432}:5432"` 형태로 변수화
- 컨테이너 내부 포트와 서비스 간 통신(`postgres:5432`, `redis:6379`)은 그대로 둠

## Status

골격만 잡힌 상태. 실제 RPC / SQLite open / git / docker 호출은 TODO.
