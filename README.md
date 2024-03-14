# SQL-Generate
SQL语句生成器

部署文档：https://github.com/zhashut/zhashut.github.io

# 功能大全

- 可视化建表
- 快捷导入建表
  - 智能导入
  - 导入表
  - 导入配置
  - 导入建表 SQL
  - 导入 Excel
- 一键生成
  - SQL 建表、插入数据语句
  - 模拟数据
  - JSON 数据
  - Go 代码
  - Java 代码
- 多种模拟数据生成规则
  - 固定值
  - 随机值
  - 正则表达式
  - 递增
  - 定制词库
- 词库共享
  - 创建词库
  - 词库继承
  - 一键创建字典表
  - 根据词库生成模拟数据
- 表信息共享
  - 创建表信息
  - 一键复制建表语句
  - 一键导入表
- 字段共享
  - 创建字段
  - 一键复制创建字段语句
  - 一键导入字段
- 举报

# 技术栈

主要技术：

- Gin
- Gorm

依赖库：
- Viper : 配置管理
- Vitess : 解析 SQL 语句
- zap : 日志库
- Excelize : 导入导出
- Gofakeit ：模拟数据
- template ：模板引擎

# 快速启动

1. 安装依赖库：`go mod tidy`
2. 运行 docs/sql/ 目录下的 create_table.sql 建表
3. 修改 config-dev.yaml 中的数据库地址为自己的
4. 已经编写好了 Dockerfile，支持 Docker 镜像部署。
5. 注意，如果是 Win 系统，是启动不了的，因为用到了 vitess（不支持 Win），要启动需要在 WSL 才可，也可以选择注释掉相关代码
   - 相关代码位置：core/builder/tableSchemaBuilder.go 中的 BuildFromSQL 方法（不要把方法全部注释，因为有引用的地方，注释方法里面的代码即可，返回 nil）和 getExprVal（这个全部注释掉即可）

# 整体架构设计

核心设计理念：将各输入方式统一为明确的 Schema，并根据 Schema 生成各类内容。

系统分为以下几个核心模块，各模块职责分明：

1. Schema 构造器：将各种不同的输入源转为统一的 Table Schema 定义
2. 统一 Schema 定义：本质是一个 Java 类（JSON 配置），用于保存表和字段的信息
3. 生成器：负责根据 Schema 生成数据和代码
4. 共享服务：包括词库、表信息、字段信息共享

> 核心模块的代码都在后端 core 目录下

### Schema 构造器

核心类：TableSchemaBuilder，作用是将不同的参数统一收敛为 TableSchema 对象。

### Schema 定义

用于保存表和字段的信息，结构如下：

```json
{
  "dbName": "库名",
  "tableName": "test_table",
  "tableComment": "表注释",
  "mockNum": 20,
  "fieldList": [{
    "fieldName": "username",
    "comment": "用户名",
    "fieldType": "varchar(256)",
    "mockType": "随机",
    "mockParams": "人名",
    "notNull": true,
    "primaryKey": false,
    "autoIncrement": false
  }]
}
```



### 生成器

#### 多种生成类型

将每种生成类型定义为一个 Builder（core/builder 目录）：

其中，对于 SQL 代码生成器（ SqlBuilder），使用方言来支持不同的数据库类型（策略模式），并使用单例模式 + 工厂模式创建方言实例。

对于 Go、Java代码生成器，编写了 codeBuilderGenertor 模板方法，并通过入参来决定代码生成。使用内置 template 模板引擎。 

#### 多种模拟数据生成规则

每种生成规则定义为一个 Generator，使用 DataGeneratorFactory（工厂模式）对多个 Generator 实例进行统一的创建和管理。

使用 Gofakeit 库实现随机数据生成（RandomDataGenerator）。

#### 统一的生成入口

使用门面模式聚合各种生成类型，提供统一的生成调用和校验方法

### 共享服务

包括词库、表信息、字段信息共享
