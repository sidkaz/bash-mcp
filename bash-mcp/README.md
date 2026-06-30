# bash-mcp

A minimal [Model Context Protocol](https://modelcontextprotocol.io) (MCP) server that exposes a single `bash.run` tool for executing bash commands inside a Docker container.

## Usage

### Build the image

```bash
docker build -t bash-mcp .
```

### Add to Claude Code

```bash
claude mcp add bash-mcp --transport stdio -- \
  /usr/local/bin/docker run --rm -i --init --log-driver=none \
  -v /path/to/your/workspace:/workspace bash-mcp
```

Replace `/path/to/your/workspace` with the directory you want mounted inside the container.

### Verify the connection

```bash
claude mcp list
```

## Tool

| Tool | Description |
|------|-------------|
| `bash.run` | Execute a bash command in `/workspace` inside the container |

### Example

```json
{
  "name": "bash.run",
  "arguments": {
    "command": "ls /workspace"
  }
}
```

Response:

```json
{
  "stdout": "...",
  "stderr": "",
  "exitCode": 0
}
```
