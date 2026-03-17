## gin 开发模板
### 1. 项目结构
```
ai-pet-backend/
├── cmd/                    # 启动入口 (绝不在根目录放 main.go)
│   └── server/
│       └── main.go         # 初始化配置、日志、组装依赖、启动你封装的 httpsrv
├── configs/                # 配置文件存放地
│   ├── config.yaml         # 本地开发配置
│   └── config.prod.yaml    # 线上生产配置
├── internal/               # 核心私有代码 (Go 的特性：internal 外的包无法导入这里的代码)
│   ├── api/                # 1. 表现层 (Handler/Controller)：处理 Gin 路由、解析入参、校验参数、返回 JSON
│   ├── service/            # 2. 业务逻辑层 (Biz)：核心大脑，处理具体的业务规则，不关心 HTTP 还是 RPC
│   ├── repo/               # 3. 数据访问层 (Dao/Repository)：专职和 MySQL(GORM)、Redis、硬件消息队列打交道
│   ├── model/              # 4. 数据模型：存放数据库映射结构体 (Entity) 和业务实体
│   ├── middleware/         # 5. 业务中间件：比如 JWT 鉴权、限流中间件 (基础日志中间件在你的 infra 里)
│   └── router/             # 6. 路由中心：专门用来把 internal/api 里的 Handler 挂载到你的 Gin 上
├── pkg/                    # 业务无关的公共组件包 (如果你不把 infra 抽成独立 Git 仓库，就放这里)
│   └── ecode/              # 全局错误码定义
├── deployments/            # 运维部署相关
│   ├── Dockerfile
│   └── k3s-deployment.yaml # K3s/K8s 的部署编排文件
├── Makefile                # 编译、测试、代码生成的快捷指令集
├── go.mod
└── go.sum
```


```
aigc/
├── cmd/
│   └── server/          (main.go, wire.go)
├── configs/             (配置文件)
├── internal/
│   ├── api/             (HTTP 接口层，收发 JSON)
│   ├── worker/          (后台消费者层，收 Kafka 消息)
│   ├── service/         (核心业务逻辑层，大脑)
│   ├── repo/            (数据组装层，负责调度各种存储)
│   ├── data/            (底层存储层：MySQL, Redis, ClickHouse, Kafka Producer)
│   ├── conf/            (配置解析)
│   └── middleware/      (Gin 中间件：跨域、JWT、Prometheus 拦截器)
├── pkg/
│   └── metrics/         (Prometheus 业务自定义打点池)
└── go.mod
```