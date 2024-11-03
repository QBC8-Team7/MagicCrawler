# Git Strategy Guide

This document outlines our Git workflow and standards for contributing to the project. Following this strategy ensures consistency, traceability, and a smoother review process.

---

## Branching Strategy

We use a feature-branch workflow where all work is done in dedicated branches off of `develop`. **Direct commits to `develop` and `master` are not allowed** to maintain a clean and controlled history.

### Branch Types

- **Feature Branches**: For new features or enhancements, prefix with `feat/`.
  - Example: `feat/add-user-authentication`
- **Bug Fix Branches**: For bug fixes, prefix with `fix/`.
  - Example: `fix/login-button-not-responding`
- **Chore Branches**: For maintenance or non-feature updates, prefix with `chore/`.
  - Example: `chore/update-dependencies`

### Creating a Branch

To create a new branch, ensure you’re on `develop` and create your branch:

```bash
git checkout develop
git pull origin develop
git checkout -b feat/add-user-authentication
```

## Pull Request Workflow

1. **Push Changes**: Push your branch to the remote repository:
   ```bash
   git push origin feat/add-user-authentication
   ```

2. **Open a Pull Request**: Create a Pull Request (PR) from your branch to `develop`.

3. **Set Reviewer and Assignee**:
   - **Reviewer**: Choose a random team member for review.
   - **Assignee**: Set yourself as the assignee, indicating you’re responsible for the changes.

4. **PR Title**: Ensure that your PR has a meaningful, clear, and descriptive title. This is crucial because we squash commits, and the PR title will serve as the commit message summary.

5. **Rebase Before Merging**:
   - Before merging your branch to `develop`, rebase it with `develop` to maintain a clean Git history.
   ```bash
   git fetch origin
   git rebase origin/develop
   ```

6. **Squash Commits**: When merging the PR, squash all commits to create a single, concise commit message from the PR title.

7. **Approval and Merge**:
   - Once the PR is approved, **any team member** can merge the branch into `develop`.

## Release Plan: Merging `develop` to `master`

We follow a scheduled release plan where changes from `develop` are merged into `master` every 2 days. This helps maintain a stable main branch with frequent, manageable updates.

- **Release Frequency**: Every 2 days.
- **Release Process**:
  1. Review and test `develop` to ensure it’s ready for release.
  2. Create a release PR from `develop` to `master`.
  3. After approval, merge into `master`.

## Commit Message Convention

We follow the **Conventional Commits** specification to write clear and structured commit messages. This helps keep our Git history consistent and readable.

### Commit Message Format

Each commit message should follow this structure:

```
<type>(scope): <description>

[optional body]
[optional footer]
```

### Commit Types

- **feat**: A new feature.
- **fix**: A bug fix.
- **chore**: Minor changes that don’t add new functionality or fix bugs.

### Example Commit Messages

```plaintext
feat(auth): add user authentication with JWT

fix(login): correct button alignment on mobile view

chore(deps): update lodash to latest version
```

## Example Workflow

1. **Creating a Feature Branch**:
   ```bash
   git checkout develop
   git pull origin develop
   git checkout -b feat/user-profile-page
   ```

2. **Making Commits**:
   ```bash
   git commit -m "feat(profile): add user profile page layout"
   ```

3. **Pushing Changes**:
   ```bash
   git push origin feat/user-profile-page
   ```

4. **Creating and Setting Up a PR**:
   - Create a PR from `feat/user-profile-page` to `develop`.
   - Assign yourself as the assignee.
   - Choose a random reviewer.
   - Ensure the PR title is meaningful and descriptive.

5. **Rebase and Squash Commits**:
   - Rebase your branch with `develop` before merging:
     ```bash
     git fetch origin
     git rebase origin/develop
     ```
   - When merging, squash commits to keep a clean history.

6. **Review and Merge**:
   - After approval, any team member can merge the PR into `develop`.

## Example Release Process

Every 2 days, we merge `develop` into `master`:

1. **Prepare for Release**:
   - Ensure all PRs are merged and `develop` is stable.

2. **Create a Release PR**:
   - Create a PR from `develop` to `master`.

3. **Approval and Merge**:
   - Once the release PR is approved, merge it into `master`.
