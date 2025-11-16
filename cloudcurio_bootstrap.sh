#!/usr/bin/env bash
# ======================================================================
# Script Name  : cloudcurio_bootstrap.sh
# Author       : CBW + ChatGPT
# Date         : 2025-11-15
#
# Summary:
#   Bootstrap the CloudCurio project structure:
#   - Creates core repo folders under a root directory.
#   - Populates each repo with standard documentation templates:
#       PROJECT_SUMMARY.md, RULES.md, AGENTS.md, INSTRUCTIONS.md,
#       JOURNAL.md, SRS.md, TASKS.md, TESTING.md
#   - Initializes local git repositories.
#   - Optionally creates remote GitHub repos (via gh) and pushes.
#   - Optionally adds repos as git submodules to cloudcurio-map.
#
# Inputs / Parameters:
#   Environment variables (optional):
#     CC_ROOT       : Root directory for all CloudCurio repos
#                     (default: "$HOME/dev/cloudcurio")
#     CC_GIT_HOST   : Git host type, currently supports: "github" or "none"
#                     (default: "none")
#     CC_GH_USER    : GitHub username/namespace (default: "cbwinslow")
#     CC_MAKE_SUBMODULES : "1" to add submodules under cloudcurio-map,
#                          "0" to skip (default: "1")
#
# Outputs:
#   - Directory tree with initialized git repos.
#   - Standard doc files in each repo.
#   - Optional remote GitHub repos created and pushed.
#   - Optional submodules registered in cloudcurio-map.
#
# Dependencies:
#   - git
#   - bash
#   - gh (GitHub CLI) if CC_GIT_HOST=github and you want remote creation
#
# Modification Log:
#   2025-11-15 - Initial version.
# ======================================================================

set -euo pipefail

# ------------- CONFIGURABLE SETTINGS ----------------------------------

CC_ROOT="${CC_ROOT:-$HOME/dev/cloudcurio}"
CC_GIT_HOST="${CC_GIT_HOST:-none}"      # "github" or "none"
CC_GH_USER="${CC_GH_USER:-cbwinslow}"   # GitHub username/namespace
CC_MAKE_SUBMODULES="${CC_MAKE_SUBMODULES:-1}"

# Repos in your empire
REPOS=(
  "dotfiles"
  "cloudcurio-knowledge"
  "cloudcurio-tools"
  "cloudcurio-agents"
  "cloudcurio-infra"
  "cloudcurio-apps"
  "cloudcurio-map"
)

# ------------- UTILITY FUNCTIONS --------------------------------------

log() {
  printf '[cloudcurio] %s\n' "$*" >&2
}

die() {
  printf '[cloudcurio:ERROR] %s\n' "$*" >&2
  exit 1
}

check_dependencies() {
  command -v git >/dev/null 2>&1 || die "git is required but not installed."

  if [[ "$CC_GIT_HOST" == "github" ]]; then
    if ! command -v gh >/dev/null 2>&1; then
      die "CC_GIT_HOST=github but 'gh' CLI not found. Install GitHub CLI or set CC_GIT_HOST=none."
    fi
  fi
}

ensure_root_dir() {
  if [[ ! -d "$CC_ROOT" ]]; then
    log "Creating root directory: $CC_ROOT"
    mkdir -p "$CC_ROOT"
  else
    log "Root directory already exists: $CC_ROOT"
  fi
}

# ------------- TEMPLATE GENERATORS ------------------------------------

write_common_docs() {
  local repo_dir="$1"
  local project_name
  project_name="$(basename "$repo_dir")"

  # PROJECT_SUMMARY.md
  cat > "$repo_dir/PROJECT_SUMMARY.md" <<EOF
# Project Summary

## Name
$project_name

## Purpose
Describe what this project does and why it exists in the CloudCurio ecosystem.

## Scope
- In-scope:
  - TODO
- Out-of-scope:
  - TODO

## Relationships
- Related repos:
  - cloudcurio-map
  - cloudcurio-knowledge
  - cloudcurio-tools
  - cloudcurio-agents
  - cloudcurio-infra
  - cloudcurio-apps

## Primary Users
- Human: CBW (developer / operator)
- Agents:
  - See \`AGENTS.md\` in this project.

## Status
- Current state: active
- Last major change: $(date +%Y-%m-%d)

EOF

  # RULES.md
  cat > "$repo_dir/RULES.md" <<'EOF'
# Rules & Operating Constraints

> These rules are binding for both humans and AI agents working in this project.

## 1. Safety & Scope

1.1. Do **not** modify or delete files outside this repo.  
1.2. Do **not** remove any existing configuration or script without explicit justification in `JOURNAL.md`.  
1.3. Respect secrets: never hardcode tokens, API keys, or passwords.

## 2. Required Reading for Agents

Before performing **any** non-trivial change, agents MUST:

- Read:
  - `PROJECT_SUMMARY.md`
  - `RULES.md`
  - `AGENTS.md`
  - `INSTRUCTIONS.md` (if relevant)
- Summarize:
  - The task in their own words.
  - The key rules that apply.

## 3. Change Management

- Log reasoning for any non-trivial change in `JOURNAL.md`.
- Prefer small, atomic commits with clear messages.
- Avoid large refactors without explicit request.

## 4. Anti-Hallucination Rules

- If unsure, agents must explicitly state uncertainty and propose options.
- Agents must not invent:
  - Files that do not exist.
  - APIs, tools, or configuration not present in this repo or docs.
- When context is insufficient, agents should request clarification.

## 5. Second-Opinion and Parallel Agents

- Significant changes should be reviewed by a reviewer agent or human.
- If an agent appears stuck or drifting:
  - Stop that agent.
  - Summarize progress in `JOURNAL.md`.
  - Start a fresh agent using the summary + rules.

## 6. Logging and Traceability

- Every working session (human or agent) should append to `JOURNAL.md`:
  - Date
  - Actor
  - Task
  - Summary of changes
  - TODOs / follow-ups

EOF

  # AGENTS.md
  cat > "$repo_dir/AGENTS.md" <<EOF
# Agents Overview for $project_name

> This file defines how AI agents should behave in this project.

## 1. Agent Profiles

### 1.1. Primary Builder Agent

- Name: \`builder-agent\`
- Role: Implement tasks, write code/configs, follow instructions.
- Must:
  - Read \`PROJECT_SUMMARY.md\`, \`RULES.md\`, \`INSTRUCTIONS.md\` before editing.
  - Log decisions in \`JOURNAL.md\`.

### 1.2. Reviewer / Critic Agent

- Name: \`reviewer-agent\`
- Role: Review changes, check for hallucinations, enforce rules.
- Must:
  - Verify that changes align with rules and requirements.
  - Log review notes under a **REVIEW** heading in \`JOURNAL.md\`.

### 1.3. Archivist / Knowledge Agent

- Name: \`archivist-agent\`
- Role:
  - Turn raw output, logs, and conversations into:
    - Notes
    - Rules
    - Memories
    - Datasets
  - Coordinate with \`cloudcurio-knowledge\` for long-term storage.

## 2. Allowed Tools

- Git (status, diff, commit) with human confirmation for pushes.
- Local file edits within this repo.
- Project-specific tools documented in \`INSTRUCTIONS.md\`.
- External APIs only when explicitly configured and documented.

## 3. Required Workflow

1. Read:
   - \`PROJECT_SUMMARY.md\`
   - \`RULES.md\`
   - \`AGENTS.md\`
2. Plan:
   - Write a short plan in \`JOURNAL.md\` under "PLAN".
3. Execute:
   - Make small, incremental changes.
4. Review:
   - Have a reviewer agent or human confirm.
5. Document:
   - Log summary, decisions, and next steps in \`JOURNAL.md\`.

EOF

  # INSTRUCTIONS.md
  cat > "$repo_dir/INSTRUCTIONS.md" <<EOF
# Instructions & Task Recipes for $project_name

> Common workflows for humans and agents.

## 1. Getting Started

1. Ensure repo is cloned:
   \`\`\`bash
   git clone <REPO_URL>
   cd $project_name
   \`\`\`
2. Read:
   - \`PROJECT_SUMMARY.md\`
   - \`RULES.md\`
   - \`AGENTS.md\`

## 2. Common Tasks

### 2.1. Add or Modify a Script / Config

1. Identify appropriate directory.
2. Make minimal, safe edits.
3. Run any tests or validations (if defined in \`TESTING.md\`).
4. Log the change in \`JOURNAL.md\`.

### 2.2. Logging a Working Session

For each session (human or agent):

- Add an entry in \`JOURNAL.md\` with:
  - Date
  - Actor
  - Task
  - Steps taken
  - Result
  - TODOs

## 3. Project-Specific Recipes

- TODO: Add concrete workflows for this project.

EOF

  # JOURNAL.md
  cat > "$repo_dir/JOURNAL.md" <<EOF
# Project Journal for $project_name

> Log of reasoning, decisions, and changes by humans and agents.

---

## $(date +%Y-%m-%d) - INITIAL SETUP

### SESSION: Bootstrap documentation
- Actor: Human (CBW)
- Summary:
  - Initialized standard project docs.
  - Connected this project to the CloudCurio ecosystem.
- Follow-ups:
  - [ ] Customize \`PROJECT_SUMMARY.md\` details.
  - [ ] Add project-specific workflows to \`INSTRUCTIONS.md\`.
  - [ ] Add tasks to \`TASKS.md\`.

---

## TEMPLATE ENTRY

### SESSION: <short description>
- Date: <YYYY-MM-DD>
- Actor: <Human / Agent name>
- Task:
  - <task description>
- Steps:
  - <step 1>
  - <step 2>
- Result:
  - <what changed>
- Follow-ups:
  - [ ] TODO item 1
  - [ ] TODO item 2

EOF

  # SRS.md
  cat > "$repo_dir/SRS.md" <<EOF
# Software Requirements Specification (SRS) for $project_name

## 1. Overview

- Project name: $project_name
- Owner: CBW
- Purpose:
  - TODO: High-level description.

## 2. Functional Requirements

- FR-1: TODO
- FR-2: TODO

## 3. Non-Functional Requirements

- NFR-1: TODO
- NFR-2: TODO

## 4. Actors

- Human operators
- Builder agents
- Reviewer agents
- Archivist agents

## 5. Dependencies

- TODO: List services, repos, external APIs.

EOF

  # TASKS.md
  cat > "$repo_dir/TASKS.md" <<EOF
# Task Breakdown for $project_name

> High-level tasks, micro-goals, and completion criteria.

## 1. Task List

- [ ] TASK-1: TODO (short description)
  - Micro-goals:
    - [ ] Goal 1
    - [ ] Goal 2
  - Success criteria:
    - Clearly defined in \`TESTING.md\` or below.

- [ ] TASK-2: TODO

## 2. Backlog

- [ ] BACKLOG-1: TODO
- [ ] BACKLOG-2: TODO

EOF

  # TESTING.md
  cat > "$repo_dir/TESTING.md" <<EOF
# Testing & Success Criteria for $project_name

## 1. Testing Strategy

- Unit tests: TODO
- Integration tests: TODO
- Manual tests: TODO

## 2. Pass/Fail Criteria

- Task is considered **complete** when:
  - All defined tests pass.
  - Documentation (including \`JOURNAL.md\` and \`TASKS.md\`) is updated.
  - No major open TODOs remain for that task.

## 3. Test Cases

- TC-1: TODO
- TC-2: TODO

EOF
}

create_repo_structure() {
  local repo="$1"
  local dir="$CC_ROOT/$repo"

  if [[ -d "$dir/.git" ]]; then
    log "Repo already exists: $dir"
    return
  fi

  log "Creating repo directory: $dir"
  mkdir -p "$dir"

  # Minimal extra structure per repo (you can expand later)
  case "$repo" in
    dotfiles)
      mkdir -p "$dir"/{.config,bin,aliases,functions,themes,templates}
      ;;
    cloudcurio-knowledge)
      mkdir -p "$dir"/raw/{chat_logs,transcripts,bookmarks,pdfs,images,agent_dialogue,logs,conversations,scrape_results}
      mkdir -p "$dir"/processed/{notes,rules,memories,summaries,examples,playbooks}
      mkdir -p "$dir"/datasets/{jsonl,parquet}
      mkdir -p "$dir"/{catalog,schema}
      ;;
    cloudcurio-tools)
      mkdir -p "$dir"/tools/{cli,monitoring,install,linter,scanners}
      mkdir -p "$dir"/scripts
      ;;
    cloudcurio-agents)
      mkdir -p "$dir"/agents/{openai,claude,crewai,toolsets,mcp/{mcp-servers,mcp-gateway}}
      mkdir -p "$dir"/{rules,memories,configs,dialogue}
      ;;
    cloudcurio-infra)
      mkdir -p "$dir"/ansible/inventories/{homelab/{group_vars,host_vars},personal/{group_vars,host_vars}}
      mkdir -p "$dir"/ansible/{playbooks,roles}
      mkdir -p "$dir"/terraform/envs/{homelab,cloud,network}
      mkdir -p "$dir"/pulumi/stacks/{homelab,personal,cloud}
      mkdir -p "$dir"/docker/{compose,images}
      ;;
    cloudcurio-apps)
      mkdir -p "$dir"/{mcp-servers,nextjs-apps,crawlers/{crawl4ai,web},api-wrappers,devtools,dashboards}
      ;;
    cloudcurio-map)
      mkdir -p "$dir"/{catalog,diagrams,scripts}
      ;;
    *)
      ;;
  esac

  write_common_docs "$dir"

  # Initialize git
  log "Initializing git repo in $dir"
  (
    cd "$dir"
    git init -b main >/dev/null 2>&1
    git add .
    git commit -m "Initial project skeleton for $repo" >/dev/null 2>&1
  )
}

create_github_remote_and_push() {
  local repo="$1"
  local dir="$CC_ROOT/$repo"
  local remote="git@github.com:${CC_GH_USER}/${repo}.git"

  log "Creating GitHub repo: ${CC_GH_USER}/${repo}"
  (
    cd "$dir"
    if ! gh repo view "${CC_GH_USER}/${repo}" >/dev/null 2>&1; then
      gh repo create "${CC_GH_USER}/${repo}" --public --confirm >/dev/null 2>&1
    else
      log "GitHub repo already exists: ${CC_GH_USER}/${repo}"
    fi

    if ! git remote get-url origin >/dev/null 2>&1; then
      git remote add origin "$remote"
    fi

    git push -u origin main >/dev/null 2>&1 || log "Warning: push failed for $repo (check credentials/keys)."
  )
}

add_submodules_to_map() {
  local map_dir="$CC_ROOT/cloudcurio-map"

  if [[ ! -d "$map_dir/.git" ]]; then
    log "cloudcurio-map repo not found or not initialized; skipping submodule creation."
    return
  fi

  log "Adding submodules to cloudcurio-map (if not already present)..."

  (
    cd "$map_dir"
    mkdir -p repos

    for repo in "${REPOS[@]}"; do
      [[ "$repo" == "cloudcurio-map" ]] && continue

      local sub_path="repos/$repo"

      if [[ -d "$sub_path/.git" || -d "$sub_path" ]]; then
        log "Submodule path already exists for $repo, skipping."
        continue
      fi

      local url=""
      if [[ "$CC_GIT_HOST" == "github" ]]; then
        url="git@github.com:${CC_GH_USER}/${repo}.git"
      else
        log "No remote URL configured for $repo; skipping submodule add (set CC_GIT_HOST=github to use remote URLs)."
        continue
      fi

      log "Adding submodule: $repo -> $sub_path"
      git submodule add "$url" "$sub_path" >/dev/null 2>&1 || log "Warning: failed to add submodule for $repo"
    done

    git add .gitmodules repos || true
    git commit -m "Add core repos as submodules" >/dev/null 2>&1 || log "No changes to commit in cloudcurio-map."
  )
}

main() {
  log "Starting CloudCurio bootstrap..."
  check_dependencies
  ensure_root_dir

  # Create local repos + docs
  for repo in "${REPOS[@]}"; do
    create_repo_structure "$repo"
  done

  # Optional: create remotes and push
  if [[ "$CC_GIT_HOST" == "github" ]]; then
    log "Creating GitHub remotes and pushing initial commits..."
    for repo in "${REPOS[@]}"; do
      create_github_remote_and_push "$repo"
    done
  else
    log "CC_GIT_HOST is 'none'; skipping remote creation."
  fi

  # Optional: add submodules
  if [[ "$CC_MAKE_SUBMODULES" == "1" ]]; then
    add_submodules_to_map
  else
    log "CC_MAKE_SUBMODULES=0; skipping submodule setup."
  fi

  log "Bootstrap complete. Root directory: $CC_ROOT"
}

main "$@"
