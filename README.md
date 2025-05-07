# ArxBackend - Multi-Service Bazel Project

This project contains two main services:
1. Python API server (`//pythonp:api_service`)
2. Go application (`//golangp/apps/arx_center:arx_center`)

## Prerequisites

- Bazel 6.0+
- Go 1.23.4
- Python 3.12
- Git with SSH configured (for private repos)

## Configuration (`.bazelrc`)

The following environment configurations are set in `.bazelrc`:

### Go Configuration
```bash
# GOPROXY settings for China mainland users
common --repo_env=GOPROXY=https://goproxy.cn
common --action_env=GOPROXY=https://goproxy.cn

# Proxy settings (replace with your actual proxy)
common --repo_env=HTTP_PROXY=172.25.96.1:7890
common --repo_env=HTTPS_PROXY=172.25.96.1:7890

# Private repository access
common --repo_env=GOPRIVATE="github.com/Arxtect/*"
common --action_env=GOPRIVATE="github.com/Arxtect/*"
Python Configuration
```

# API server port configuration
common --action_env=PYTHON_API_SERVER_PORT=9002
Running the Services
## Option 1: Separate Terminals (Recommended for Development)
### 1、For the Python API server:

```bash
# Terminal 1
bazel run //pythonp:api_service
```
### For the Go application:

```bash
# Terminal 2
bazel run //golangp/apps/arx_center:arx_center
```

## Option 2: Using tmux/screen
```bash
# Create a new tmux session
tmux new-session -d -s arxbackend

# Split window vertically
tmux split-window -v

# Run Python API in top pane
tmux send-keys -t 0 "bazel run //pythonp:api_service" C-m

# Run Go app in bottom pane
tmux send-keys -t 1 "bazel run //golangp/apps/arx_center:arx_center" C-m

# Attach to session
tmux attach-session -t arxbackend
```

## Option 3: Log to Files
```bash
# Run Python API with log output
bazel run //pythonp:api_service > python_api.log 2>&1 &

# Run Go application with log output
bazel run //golangp/apps/arx_center:arx_center > go_app.log 2>&1 &

# Tail both logs (in separate terminal)
tail -f python_api.log go_app.log
```
# Environment Variables
The Python API server reads its port from the environment variable:

```python
port = int(os.getenv("PYTHON_API_SERVER_PORT", "9002"))
```

This is configured in .bazelrc and can be overridden at runtime:

```bash
bazel run //pythonp:api_service -- --action_env=PYTHON_API_SERVER_PORT=8080
```

### Troubleshooting
####  If you encounter dependency issues:

Ensure your Go and Python toolchains are properly configured

Verify your proxy settings in .bazelrc are correct

Check that your SSH keys are configured for GitHub access

Clean and rebuild if needed:

```bash
bazel clean --expunge
bazel build //...
```