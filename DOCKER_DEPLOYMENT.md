# Docker Deployment Guide

## 🐳 Current Docker Environment

### Active Containers:

```bash
docker ps
# CONTAINER ID   IMAGE     COMMAND   CREATED   STATUS    PORTS     NAMES
# agitated_napier (agent0ai)
# zealous_jang (codex)
# vtcode (vtcode)
```

### Docker Desktop Status:

- ✅ Docker Desktop running
- ✅ 3 active containers detected
- ✅ WSL2 integration functional

## 🚀 Multi-Agent System Deployment Options

### Option 1: Replace Existing Containers

```bash
# Stop current containers
docker stop agitated_napier zealous_jang vtcode

# Remove containers (optional)
docker rm agitated_napier zealous_jang vtcode

# Build multi-agent container
cd /home/cbwinslow/apps/qwen-code
docker build -t openrouter-multi-agent .

# Run new container
docker run -d \
  --name multi-agent-system \
  -p 8080:8080 \
  -e OPENROUTER_API_KEY=your_key_here \
  openrouter-multi-agent
```

### Option 2: Build Custom Docker Image

```dockerfile
# Dockerfile for Multi-Agent System
FROM python:3.12-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy multi-agent system
COPY openrouter_multi_agent.py .
COPY test_multi_agent.py .

# Create non-root user
RUN useradd -m -u 1000 agent
USER agent

# Expose port
EXPOSE 8080

# Default command
CMD ["python3", "openrouter_multi_agent.py", "--demo"]
```

### Option 3: Docker Compose Deployment

```yaml
# docker-compose.yml
version: '3.8'

services:
  multi-agent:
    build: .
    container_name: openrouter-multi-agent
    ports:
      - '8080:8080'
    environment:
      - OPENROUTER_API_KEY=${OPENROUTER_API_KEY}
      - LOG_LEVEL=INFO
    volumes:
      - ./results:/app/results
      - ./logs:/app/logs
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    container_name: multi-agent-redis
    ports:
      - '6379:6379'
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  redis_data:
```

## 📦 Container Management

### Build Commands:

```bash
# Build from Dockerfile
docker build -t openrouter-multi-agent:latest .

# Build with specific tag
docker build -t openrouter-multi-agent:v1.0 .

# Build with build args
docker build \
  --build-arg OPENROUTER_API_KEY=your_key \
  -t openrouter-multi-agent:custom \
  .
```

### Run Commands:

```bash
# Interactive mode
docker run -it --rm openrouter-multi-agent:latest

# Daemon mode
docker run -d --name multi-agent openrouter-multi-agent:latest

# With environment variables
docker run -d \
  --name multi-agent \
  -e OPENROUTER_API_KEY=your_key \
  -e LOG_LEVEL=DEBUG \
  -p 8080:8080 \
  openrouter-multi-agent:latest

# With volume mounts
docker run -d \
  --name multi-agent \
  -v $(pwd)/results:/app/results \
  -v $(pwd)/logs:/app/logs \
  openrouter-multi-agent:latest
```

### Management Commands:

```bash
# List containers
docker ps -a

# View logs
docker logs multi-agent-system

# Execute commands in container
docker exec -it multi-agent-system bash

# Stop container
docker stop multi-agent-system

# Start container
docker start multi-agent-system

# Remove container
docker rm multi-agent-system
```

## 🔧 Production Configuration

### Environment Variables:

```bash
OPENROUTER_API_KEY=your_api_key_here
LOG_LEVEL=INFO
MAX_AGENTS=10
CONSENSUS_THRESHOLD=0.6
MEMORY_LIMIT=1000
TIMEOUT_SECONDS=300
```

### Resource Limits:

```yaml
deploy:
  resources:
    limits:
      cpus: '2.0'
      memory: 2G
    reservations:
      cpus: '0.5'
      memory: 512M
```

### Health Check:

```dockerfile
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD python3 -c "import requests; requests.get('http://localhost:8080/health')" || exit 1
```

## 🌐 Network Configuration

### Port Mapping:

- **8080**: Main API endpoint
- **6379**: Redis (if using)
- **8000**: Web interface (future)

### Service Discovery:

```bash
# Register with service discovery
docker run -d \
  --name multi-agent \
  --network my-network \
  -e SERVICE_NAME=multi-agent \
  -e SERVICE_TAGS=ai,agents \
  openrouter-multi-agent:latest
```

## 📊 Monitoring & Logging

### Log Management:

```bash
# View real-time logs
docker logs -f multi-agent-system

# Export logs
docker logs multi-agent-system > multi-agent.log

# Log rotation
docker run --rm \
  -v /var/lib/docker/containers:/var/lib/docker/containers \
  alpine sh -c "find /var/lib/docker/containers -name '*.log' -exec truncate -s 0 {} \;"
```

### Monitoring Setup:

```yaml
# docker-compose.monitoring.yml
version: '3.8'

services:
  prometheus:
    image: prom/prometheus
    ports:
      - '9090:9090'
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - '3000:3000'
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

## 🚀 Deployment Script

```bash
#!/bin/bash
# deploy-multi-agent.sh

set -e

echo "🚀 Deploying OpenRouter Multi-Agent System..."

# Build image
echo "📦 Building Docker image..."
docker build -t openrouter-multi-agent:latest .

# Stop existing container
echo "🛑 Stopping existing container..."
docker stop multi-agent-system 2>/dev/null || true
docker rm multi-agent-system 2>/dev/null || true

# Run new container
echo "🏃 Starting new container..."
docker run -d \
  --name multi-agent-system \
  --restart unless-stopped \
  -p 8080:8080 \
  -e OPENROUTER_API_KEY="$OPENROUTER_API_KEY" \
  -v $(pwd)/results:/app/results \
  openrouter-multi-agent:latest

echo "✅ Deployment complete!"
echo "🌐 Container available at: http://localhost:8080"
echo "📊 Logs: docker logs -f multi-agent-system"
```

## 🔒 Security Considerations

### Best Practices:

1. **Use non-root users** in containers
2. **Limit resource usage** with constraints
3. **Scan images** for vulnerabilities
4. **Use secrets management** for API keys
5. **Network segmentation** for multi-container setups

### Security Commands:

```bash
# Scan for vulnerabilities
docker scan openrouter-multi-agent:latest

# Run as non-root
docker run --user 1000:1000 openrouter-multi-agent:latest

# Read-only filesystem
docker run --read-only --tmpfs /tmp openrouter-multi-agent:latest
```

---

**Created**: 2025-11-15  
**Environment**: WSL2 + Docker Desktop  
**Status**: Ready for Production
