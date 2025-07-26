# Revly: AI-Powered Code Review CLI

Revly is an intelligent command-line interface (CLI) tool designed to streamline your code review process. Leveraging the power of Large Language Models (LLMs), Revly analyzes your Git changes, provides instant feedback on code quality, and even helps you craft meaningful commit messages. Say goodbye to waiting for PR reviews and get actionable insights directly in your terminal.

## Features

*   **AI-Powered Code Review:** Get immediate, AI-generated suggestions and feedback on your code changes.
*   **Flexible Review Scopes:** Review unstaged changes, staged changes, specific commits, or the latest commit (HEAD).
*   **AI-Generated Commit Messages:** Automatically generate descriptive and concise commit messages based on your staged changes.
*   **Custom Commit Messages:** Option to provide your own commit message, bypassing AI generation.
*   **Dry Run Mode:** Preview AI-generated commit messages without actually committing.
*   **Optional Push:** Automatically push your commits to the remote repository after committing.
*   **Severity Highlighting:** AI review outputs are highlighted with `[CRITICAL]`, `[WARNING]`, and `[INFO]` tags for quick identification of important feedback.

## Installation

Revly is built with Go. To install it, you need to have Go installed on your system.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/nareshkarthigeyan/revly.git
    cd revly
    ```

2.  **Install Revly:**
    ```bash
    go install .
    ```
    This command will install the `revly` executable in your `GOPATH/bin` directory, which should ideally be in your system's PATH.

### Configuration

Revly uses an LLM for its AI capabilities. You need to provide an API key for the LLM service. Currently, Revly is configured to use `OPENROUTER_KEY`.

You can set your `OPENROUTER_KEY` in one of the following ways:

1.  **Environment Variable:**
    ```bash
    export OPENROUTER_KEY="your-api-key-here"
    ```
    Add this line to your shell's profile file (e.g., `~/.bashrc`, `~/.zshrc`, `~/.profile`) to make it persistent.

2.  **`.env` file:**
    Create a file named `.env` in the root directory where you run `revly` and add your key:
    ```
    OPENROUTER_KEY="your-api-key-here"
    ```

On the first run, Revly will check for the `OPENROUTER_KEY` and provide guidance if it's not found.

## Usage

### `revly` (Root Command)

The main command provides a brief introduction to the tool.

```bash
revly
```

**Output:**
```
Revly CLI - AI Code Review Assistant

Try `revly review` to get started.
For more help, use `revly --help`.
```

### `revly review`

Run an AI-powered code review on your Git changes. By default, it reviews unstaged changes in your working directory.

**Usage:**
```bash
revly review [flags]
```

**Flags:**
*   `-s`, `--staged`: Review only staged changes (`git diff --cached`).
*   `-c`, `--commit <hash>`: Review a specific commit by its SHA hash. If `<hash>` is omitted, it defaults to `HEAD` (the latest commit).
*   `--head`: Review the latest commit (`HEAD`).
*   `--diff`: Display the Git diff before running the AI review.

**Examples:**

*   **Review all current unstaged changes:**
    ```bash
    revly review
    ```

*   **Review only staged changes:**
    ```bash
    revly review --staged
    # or
    revly review -s
    ```

*   **Review a specific commit by its SHA:**
    ```bash
    revly review --commit <commit-hash>
    # or
    revly review -c <commit-hash>
    ```

*   **Review the most recent commit (HEAD):**
    ```bash
    revly review --head
    # or
    revly review -c
    ```

*   **Display diff before review:**
    ```bash
    revly review --diff
    ```

### `revly commit`

Stage changes, generate a commit message via AI or custom input, commit, and optionally push.

**Usage:**
```bash
revly commit [file/folder] [flags]
```

**Arguments:**
*   `[file/folder]` (optional): Specifies a particular file or folder to stage and commit. If omitted, and `--all` is not used, it defaults to staging the current directory (`.`).

**Flags:**
*   `-a`, `--all`: Stage all changes recursively from the project root before committing.
*   `--push`: Push the commit to the remote repository after committing.
*   `--dry-run`: Show the AI-generated commit message without actually committing or pushing.
*   `-m`, `--message <message>`: Use a custom commit message instead of an AI-generated one. This flag bypasses AI generation and the confirmation prompt.

**Examples:**

*   **Stage all changes, generate AI commit message, and prompt for confirmation:**
    ```bash
    revly commit
    ```

*   **Explicitly stage all changes and proceed with AI commit message workflow:**
    ```bash
    revly commit --all
    ```

*   **Stage only the `src/utils/` folder, generate AI commit message, and commit:**
    ```bash
    revly commit src/utils/
    ```

*   **Show what the AI would generate as a commit message (dry run):**
    ```bash
    revly commit --dry-run
    ```

*   **Stage everything, commit with an AI-generated message, and push to remote:**
    ```bash
    revly commit --all --push
    ```

*   **Target a specific file, commit it with an AI message, and push:**
    ```bash
    revly commit src/index.ts --push
    ```

*   **Commit immediately with a custom message (no AI, no prompt):**
    ```bash
    revly commit -m "Fix typo in README"
    ```

### `revly version`

Print the version number of Revly.

**Usage:**
```bash
revly version
# or
revly -v
# or
revly --version
```

**Output:**
```
Revly v0.1.0
```

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.