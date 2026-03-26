# portal

独立部署的 Web 门户，认证与权限唯一事实源为 Keycloak，MongoDB 只保存门户投影与 portal session。

## 技术栈

- 后端: Go 1.23+, Gin, MongoDB 官方驱动
- 前端: Vue 3, TypeScript, Vite, Pinia, Vue Router, Element Plus
- 认证: Keycloak OIDC
- 权限与同步: Keycloak Admin REST API
- 部署: Docker Compose + Kubernetes manifests

## 目录

```text
portal/
  apps/
    portal-api/
    portal-web/
  internal/
    auth/
    config/
    handler/
    kcadmin/
    middleware/
    model/
    permission/
    repository/
    service/
    session/
    sync/
  deployments/
    docker-compose/
    k8s/
  docs/
  scripts/
  tests/
```

## 核心实现

- `portal-web` 和 `portal-api` 独立部署
- 前端只调用 `portal-api`
- 未登录时由 `portal-web` 引导到 `portal-api /api/v1/auth/login`
- Keycloak 回调到 `portal-api /api/v1/auth/callback`
- `portal-api` 用 OIDC code 换 token 后，立刻调用 Keycloak Admin API 同步:
  - 当前 realm 基础信息
  - 当前 realm client 列表
  - 当前用户基础资料
  - 当前用户 realm roles
  - 当前用户 client roles
- 同步写入 MongoDB 投影后，才创建 `portal_sessions`
- 默认空闲超时 15 分钟，由 portal 自己控制
- 退出时先删除 portal session，再跳转 Keycloak logout

## 本地启动

1. 复制环境变量:

```powershell
Copy-Item .env.example .env
```

2. 启动整套依赖:

```powershell
./scripts/dev-up.ps1
```

3. 访问:

- Web: [http://localhost:5173](http://localhost:5173)
- API: [http://localhost:8080](http://localhost:8080)
- Keycloak: [http://localhost:8081](http://localhost:8081)
- OpenAPI: [http://localhost:8080/openapi.yaml](http://localhost:8080/openapi.yaml)

4. 样例 Keycloak 用户:

- `portal-admin / Admin123!`
- `alice / Alice123!`

## 主要 API

- `GET /healthz`
- `GET /readyz`
- `GET /api/v1/auth/login`
- `GET /api/v1/auth/callback`
- `GET /api/v1/auth/logout`
- `GET /api/v1/me`
- `GET /api/v1/apps`
- `GET /api/v1/admin/client-metas`
- `PUT /api/v1/admin/client-metas/:clientId`
- `GET /api/v1/admin/settings`
- `PUT /api/v1/admin/settings`

## Mongo 集合与索引

- `kc_realms`
- `kc_clients`
- `portal_client_meta`
- `kc_users`
- `portal_sessions`
- `portal_settings`

初始化脚本: [`scripts/init-mongo.js`](scripts/init-mongo.js)

## 文档

- 架构说明: [`docs/architecture.md`](docs/architecture.md)
- OpenAPI: [`docs/openapi.yaml`](docs/openapi.yaml)

## 测试

- 单元测试:

```powershell
go test ./...
```

- 集成测试骨架:

```powershell
go test -tags integration ./tests/integration/...
```
