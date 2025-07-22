# RPC Timeline Main Service

This service is a backend component for generating and serving timeline feeds in a scalable and efficient manner. It is designed using the gRPC architecture.

## 项目介绍     README

The Timeline Feed Optimization Service (`timeline_main_service_v2`) is built upon gRPC architecture, designed to aggregate real-time posts from numerous creators and optimize the user experience through dynamic, time-based ranking. It emphasizes low-latency (<100ms) responses even under peak traffic conditions by leveraging Redis as a high-speed key-value store and cloud-based auto-scaling.

## 快速上手     Getting Started

This section describes how to quickly set up, configure, and run the service.

### 环境要求

- Go 1.18+
- Redis 6.x+
- Protobuf compiler (`protoc`)
- Docker (optional)

### 安装步骤

```bash
git clone <repo-url>
cd trpc_timeline_main_service_v2

# Install dependencies
go mod tidy

# Setup configuration
cp config/config.example.yaml config/config.yaml
```

Edit `config.yaml` as per your Redis setup.

### 启动服务

```bash
go run main.go
```

### API

API definitions are provided in GRPC format (`proto/timeline.proto`). Compile them using:

```bash
protoc --go_out=. --go-grpc_out=. proto/timeline.proto
```

Example method call:

```proto
rpc GetTimelineList(TimelineListRequest) returns (TimelineListResponse) {}
```

Example JSON input:

```json
{
    "page_params": {
        "key": "user-feed",
        "scene": "main-feed"
    },
    "business_info": {
        "app_key": "example-app-key"
    }
}
```

### 运维

If the service encounters problems, follow these operational guidelines:

- **服务降级**: Temporarily disable vote counting to reduce Redis load if experiencing latency.
- **问题定位**: Check Redis status, error logs, and tRPC logs for errors.
- **监控查看**: View metrics on Tencent Cloud monitoring dashboards.
- **日志查看**: Logs available at standard path (`/var/log/timeline_service/`) or configured via `config.yaml`.

## 常见问题     FAQ

- **Q: What happens if Redis fails?**  
  A: The timeline data will become unavailable temporarily until Redis recovers; no persistent fallback is currently implemented.

- **Q: How is the content ranked?**  
  A: The current ranking algorithm is a simple weighted model influenced by posted time and content freshness.

## 行为准则    Code Of Conduct

Developers must adhere to internal Tencent guidelines for code collaboration, responsibility scope, software licensing, and conflict resolution.

## 团队介绍    Members

- **Maintainers**: Infrastructure Team  
- **Communication Channel**: DevOps Platform, video_app_short_video Group
