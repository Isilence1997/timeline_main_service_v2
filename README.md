# RPC Timeline Main Service

This service is a backend component for generating and serving timeline feeds in a scalable and efficient manner. It is designed using the gRPC architecture.

## 项目介绍  Overview

The Timeline Feed Service (v2) is a Go-based microservice designed to retrieve and optimize a user's timeline feed (e.g. video or social feed). It provides RPC endpoints for fetching personalized feed items, checking unread feed markers, and retrieving specialized feed data (such as profile picture update feeds) and related info. This service uses gRPC framework and integrates with internal systems (like **Xvkj** for feed configuration and ranking, and Redis for caching) to deliver a tailored feed experience.

The core responsibilities include:

- **Fetching timeline items** for a given user and scene (feed/avatar_list/refresh) on demand (using a _read diffusion_ model – i.e. content is pulled when the user reads the feed, rather than pushed on write).
    
- **Integrating with Xvkj (internal recommendation/search service)** to determine feed configuration, ranking algorithms, and to retrieve content IDs for the feed.
    
- **Caching and sorting feed items** using Redis to improve performance on subsequent requests or pagination.
    
- **Providing additional feed utilities** such as unread indicators and specialized feeds (e.g. a feed of friends' profile picture updates).

It emphasizes low-latency (<100ms) responses even under peak traffic conditions by leveraging Redis as a high-speed key-value store and cloud-based auto-scaling.

## 快速上手     Getting Started

This section describes how to quickly set up, configure, and run the service.

### 环境要求 Requirements

- Go 1.18+
- Redis 6.x+
- gRPC framework
- Protobuf compiler (`protoc`)
- Docker (optional)

### 安装步骤 install

```bash
git clone <repo-url>
cd trpc_timeline_main_service_v2

# Install dependencies
go mod tidy

# Setup configuration
mkdir -p /usr/local/trpc/conf
cp config/config.example.yaml config/config.yaml

# Run the service
go run main.go
```

Edit `config.yaml` as per your Redis setup.

### 服务架构 Project Structure

timeline_main_service_v2/
├── logic/                # Core business logic divided into sub-packages  
│   ├── config/       # Higher-level config handlers that work with the data from the config file and Xvkj  
│   ├── idlist/           # obtaining lists of content IDs that should appear in the timeline feed for a user  
│   ├── sortcache/        # Handles sorting of feed items and caching the results  
│   ├── feed/             # high-level feed composition logic  
│   ├── avatarlist/       #   
│   ├── common/           #   
│   └── backsource/       # handle _back-source_ IDs.  
├── dao/                  # Data access layer  
│   ├── Xvkj/             # Integration with the Xvkj platform  
│   └── redis/            # Encapsulates interactions with Redis for caching  
├── config/               # Configuration loading logic  
├── test/                 # Unit and integration tests  
├── main.go               # Service entrypoint  
├── common/               # Common utilities and definitions  
│   ├── constant/         # Integration with the Xvkj platform  
│   ├── retry/            # Ge tNetwork Error Retry  
│   ├── errcode/          # Encapsulates interactions with Redis for caching  
│   └── utility/          # Common utilities  
└── model/                # Data models and structures used across the service               

### 启动服务 Start and Running

```bash
go build -o timeline_service .
./timeline_service
```
To run integration tests:
```bash
sh unit_test.sh
cd test/
go test -v
```
Defined in `proto/` as GRPC services. Compile them using:
```bash
protoc --go_out=. --go-grpc_out=. proto/timeline.proto
```

### API

- `GetTimelineList(req: TimelineListReq) -> TimelineListRsp` – Retrieves a list of timeline feed items for the specified feed **scene** (context). The request typically includes pagination parameters (`PageParams`) and AppKey for authorization.  The response contains a list of feed items (content IDs and metadata) along with a `SeqNum`  as a unique token for the feed fetch.

- `GetUnReadMark(req: UnReadMarkRequest) -> UnReadMarkResponse` – Checks the timeline feed for unread content since the last read.

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


