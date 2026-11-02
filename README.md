# 漂流瓶项目 - 启动与部署操作清单

## 项目结构

```
piaoliu-v2-all/
├── piaoliu-v2-server/        # 后端服务 (Go)
│   ├── cmd/server/main.go    # 启动入口
│   ├── internal/            # 内部代码
│   ├── docker-compose.yml   # Docker 配置
│   └── .env                 # 环境变量
└── piaoliu-v2-miniprogram/ # 前端小程序 (Vue3 + UniApp)
    ├── src/                 # 源码目录
    └── package.json         # 依赖配置
```

---

## 环境要求

| 依赖 | 版本 |
|------|------|
| Docker | 20.10+ |
| Go | 1.21+ |
| Node.js | 18+ |
| npm | 9+ |

---

## 开发环境启动

### 步骤 1：启动 Docker 服务

```bash
cd piaoliu-v2-all/piaoliu-v2-server
docker-compose up -d mysql redis
```

验证服务：
```bash
docker-compose ps
```

预期输出：
```
Name                     Command                 State                 Ports
-------------------------------------------------------------------------------------
piaoliu-v2-server-mysql-1   docker-entrypoint.sh --cha...   Up             0.0.0.0:3306->3306/tcp
piaoliu-v2-server-redis-1   docker-entrypoint.sh redis...   Up             0.0.0.0:6379->6379/tcp
```

### 步骤 2：启动后端服务

```bash
cd piaoliu-v2-all/piaoliu-v2-server
go run cmd/server/main.go
```

验证服务：
```bash
curl http://localhost:8080/api/users/profile
```

预期输出：后端监听在 `:8080` 端口

### 步骤 3：启动前端开发服务

```bash
cd piaoliu-v2-all/piaoliu-v2-miniprogram
npm install
npm run dev:h5
```

访问地址：`http://localhost:5173`（或自动分配的其他端口）

---

## 微信小程序开发

### 编译小程序代码

```bash
cd piaoliu-v2-all/piaoliu-v2-miniprogram
npm run dev:mp-weixin
```

### 微信开发者工具导入

1. 打开微信开发者工具
2. 选择「导入项目」
3. 目录选择：`piaoliu-v2-all/piaoliu-v2-miniprogram/dist/build/mp-weixin`
4. AppID 填写测试号或正式号

---

## 生产部署

### 后端构建

```bash
cd piaoliu-v2-all/piaoliu-v2-server
go build -o drift-bottle cmd/server/main.go
```

### 前端构建

```bash
cd piaoliu-v2-all/piaoliu-v2-miniprogram
npm run build:h5      # H5 版本
npm run build:mp-weixin  # 微信小程序版本
```

### Docker 一键部署

```bash
cd piaoliu-v2-all/piaoliu-v2-server
docker-compose up -d
```

包含服务：
- **app**: 后端 API 服务 (端口 8080)
- **mysql**: 数据库 (端口 3306)
- **redis**: 缓存服务 (端口 6379)

---

## 环境变量配置

编辑 `piaoliu-v2-server/.env` 文件：

```env
# 数据库配置
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root_password
DB_NAME=drift_bottle

# Redis 配置
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=

# JWT 密钥（生产环境请修改）
JWT_SECRET=your-jwt-secret-key

# 腾讯云 COS 配置（上传图片用）
COS_BUCKET_URL=https://your-bucket.cos.ap-guangzhou.myqcloud.com
COS_SECRET_ID=your-secret-id
COS_SECRET_KEY=your-secret-key

# 微信小程序配置
WX_APP_ID=your-app-id
WX_APP_SECRET=your-app-secret
```

---

## 正式上线步骤（生产环境部署清单）

### 一、微信公众平台配置

1. **注册小程序** - 在 [微信公众平台](https://mp.weixin.qq.com/) 注册小程序账号
2. **获取凭证** - 在「开发」→「开发设置」中获取：
   - `AppID`（小程序唯一标识）
   - `AppSecret`（小程序密钥，需妥善保管）
3. **配置服务器域名** - 在「开发」→「开发设置」→「服务器域名」中配置：
   - `request 合法域名`：`https://你的后端域名`（如 `https://api.piaoliu.com`）
   - `uploadFile 合法域名`：同上
   - `downloadFile 合法域名`：同上（如果有文件下载）
   - `socket 合法域名`：同上（WebSocket 用，如 `wss://api.piaoliu.com`）
4. **配置业务域名** - 如果有 web-view 需求，在「开发」→「开发设置」→「业务域名」中配置

### 二、后端配置修改

#### 1. 修改 `.env` 文件

编辑 `piaoliu-v2-server/.env`：

```env
# 数据库配置（生产环境使用独立数据库）
DB_HOST=你的数据库地址
DB_PORT=3306
DB_USER=你的数据库用户名
DB_PASSWORD=你的数据库密码
DB_NAME=drift_bottle

# Redis 配置
REDIS_ADDR=你的Redis地址:6379
REDIS_PASSWORD=你的Redis密码

# JWT 密钥（生产环境必须修改为复杂随机字符串）
JWT_SECRET=你的复杂随机密钥_至少32位

# 腾讯云 COS 配置（上传图片用）
COS_BUCKET_URL=https://你的bucket.cos.ap-guangzhou.myqcloud.com
COS_SECRET_ID=你的secret-id
COS_SECRET_KEY=你的secret-key

# 微信小程序配置（必须填写真实值）
WX_APP_ID=你的小程序AppID
WX_APP_SECRET=你的小程序AppSecret
```

#### 2. 修改后端登录逻辑 - 接入真实微信 code2Session

编辑 `piaoliu-v2-server/internal/handler/auth.go`，将 mock 登录改为真实调用：

```go
// 替换 Login 函数中的 mock 部分为：
func Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Fail(c, http.StatusBadRequest, "参数错误")
        return
    }

    // 调用微信 code2Session 接口获取 openid
    wxURL := fmt.Sprintf(
        "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
        config.LoadConfig().WxAppID,
        config.LoadConfig().WxAppSecret,
        req.Code,
    )
    wxResp, err := http.Get(wxURL)
    if err != nil {
        response.Fail(c, http.StatusInternalServerError, "微信登录失败")
        return
    }
    defer wxResp.Body.Close()
    body, _ := io.ReadAll(wxResp.Body)
    var wxSession WxSessionResponse
    json.Unmarshal(body, &wxSession)
    if wxSession.ErrCode != 0 {
        response.Fail(c, http.StatusInternalServerError, fmt.Sprintf("微信登录失败: %s", wxSession.ErrMsg))
        return
    }
    openID := wxSession.OpenID
    unionID := wxSession.UnionID

    // ... 后续查找/创建用户逻辑不变
}
```

需要确保 `config` 包中有 `WxAppID` 和 `WxAppSecret` 的读取逻辑，并添加 `WxSessionResponse` 结构体：

```go
type WxSessionResponse struct {
    OpenID     string `json:"openid"`
    SessionKey string `json:"session_key"`
    UnionID    string `json:"unionid"`
    ErrCode    int    `json:"errcode"`
    ErrMsg     string `json:"errmsg"`
}
```

#### 3. 后端部署

```bash
cd piaoliu-v2-all/piaoliu-v2-server
docker-compose up -d
```

或使用服务器直接运行：
```bash
go build -o drift-bottle cmd/server/main.go
./drift-bottle
```

建议使用 `systemd` 或 `supervisor` 管理进程，配合 `nginx` 做反向代理和 HTTPS 终止。

### 三、前端配置修改

#### 1. 修改 API 地址

编辑 `piaoliu-v2-miniprogram/src/utils/request.ts`，将 `BASE_URL` 改为正式服务器地址：

```typescript
// 开发环境
const BASE_URL = 'http://localhost:8080/api'

// 生产环境改为：
const BASE_URL = 'https://你的后端域名/api'
```

或使用条件编译区分环境：

```typescript
// #ifdef MP-WEIXIN
const BASE_URL = 'https://你的后端域名/api'
// #endif
// #ifndef MP-WEIXIN
const BASE_URL = 'http://localhost:8080/api'
// #endif
```

#### 2. 去掉测试账号相关代码

编辑 `piaoliu-v2-miniprogram/src/pages/login/login.vue`：

- 已去掉测试账号切换 UI（用户A / 用户B 按钮）
- 已去掉 `selectedAccount` 变量和 `switchAccount` 方法
- `wechatLogin` 已改为直接使用 `uni.login({ provider: 'weixin' })` 获取真实 code

#### 3. 构建生产版本

```bash
cd piaoliu-v2-all/piaoliu-v2-miniprogram
npm run build:mp-weixin
```

### 四、微信开发者工具上传

1. 打开微信开发者工具
2. 导入项目目录：`piaoliu-v2-all/piaoliu-v2-miniprogram/dist/build/mp-weixin`
3. AppID 填写**正式小程序的 AppID**
4. 点击「上传」按钮，填写版本号和备注
5. 登录 [微信公众平台](https://mp.weixin.qq.com/) →「版本管理」→ 提交审核

### 五、上线前检查清单

| 检查项 | 状态 | 说明 |
|--------|------|------|
| 后端 `.env` 中 `WX_APP_ID` 和 `WX_APP_SECRET` 已填写真实值 |  | |
| 后端登录接口已接入真实 `code2Session` | ☐ | |
| 后端 `JWT_SECRET` 已改为复杂随机字符串 | ☐ | |
| 前端 `BASE_URL` 已改为正式 HTTPS 地址 | ☐ | |
| 前端测试账号代码已清理 | ☐ | |
| 微信公众平台已配置服务器域名（HTTPS） | ☐ | |
| 后端已部署到服务器并配置 HTTPS | ☐ | |
| 数据库已初始化（自动迁移） |  | |
| Redis 已部署并可用 | ☐ | |
| 文件上传 COS 配置已填写 | ☐ | |
| 微信开发者工具上传成功 | ☐ | |
| 小程序已提交审核 | ☐ | |

### 六、上线后注意事项

1. **HTTPS 必须** - 微信小程序要求所有请求必须使用 HTTPS
2. **域名备案** - 后端域名必须已完成 ICP 备案
3. **证书更新** - SSL 证书到期前及时续期
4. **数据备份** - 定期备份 MySQL 数据库
5. **日志监控** - 关注后端日志，及时发现异常
6. **版本管理** - 每次修改后通过微信开发者工具上传新版本，在公众平台提交审核

---

## 常用命令

### 停止所有服务

```bash
# 停止后端服务（MySQL + Redis）
cd piaoliu-v2-all/piaoliu-v2-server
docker-compose down

# 停止后端 Go 服务（如果在终端运行，按 Ctrl+C；如果在后台运行）
pkill -f "go run cmd/server/main.go"

# 停止前端服务（如果在终端运行，按 Ctrl+C；如果在后台运行）
pkill -f "npm run dev:h5"
# 或根据端口停止
lsof -ti:5173,5174 | xargs kill -9
```

### 查看日志

```bash
# Docker 日志
docker-compose logs -f mysql
docker-compose logs -f redis

# 后端日志（运行中）
tail -f /var/log/drift-bottle.log
```

### 数据库操作

```bash
# 进入 MySQL
mysql -h 127.0.0.1 -u root -proot_password drift_bottle

# 查看数据
SELECT * FROM users LIMIT 10;
SELECT * FROM bottles LIMIT 10;
```

### 代码检查

```bash
# 后端
cd piaoliu-v2-server
go vet ./...

# 前端
cd piaoliu-v2-miniprogram
npm run type-check
```

---

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| 后端 API | 8080 | RESTful API |
| MySQL | 3306 | 数据库 |
| Redis | 6379 | 缓存 |
| 前端 H5 | 5173/5174 | 开发服务器 |

---

## 目录说明

### 后端

```
piaoliu-v2-server/
├── cmd/server/          # 启动入口
├── internal/
│   ├── handler/         # HTTP 处理器
│   ├── models/          # 数据库模型
│   ├── middleware/      # 中间件
│   ├── router/          # 路由配置
│   ├── service/         # 业务逻辑
│   └── websocket/       # WebSocket
├── migrations/          # 数据库迁移
├── pkg/                 # 公共工具包
└── docker-compose.yml   # Docker Compose
```

### 前端

```
piaoliu-v2-miniprogram/
├── src/
│   ├── api/             # API 请求封装
│   ├── components/      # 公共组件
│   ├── pages/           # 页面
│   │   ├── ocean/       # 海洋主页
│   │   ├── pick/        # 捡瓶子
│   │   ├── chat/        # 聊天
│   │   └── mine/        # 个人中心
│   ├── stores/          # Pinia 状态管理
│   ├── types/           # TypeScript 类型
│   ├── utils/           # 工具函数
│   └── App.vue          # 根组件
└── package.json         # 依赖配置
```
