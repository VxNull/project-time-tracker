# 项目工时统计系统

[English](README.md)

## 项目简介

项目工时统计系统是一个基于 Go 语言开发的 Web 应用,用于帮助公司或团队管理和跟踪员工的工时情况。该系统提供了直观的界面,允许员工提交工时,管理员查看统计数据,并支持导出详细的 Excel 报表。

## 主要功能

- 员工工时提交:员工可以轻松记录每日工作时间和项目信息
- 管理员仪表盘:管理员可以查看整体工时统计和项目进展
- 项目管理:添加、编辑和删除项目
- 员工管理:添加、编辑和删除员工账户
- 数据导出:生成详细的 Excel 报表,包括项目和员工的工时统计

## 技术栈

- 后端:Go
- 数据库:SQLite
- 前端:HTML, CSS (Tailwind CSS), JavaScript
- 依赖管理:Go Modules

## 安装说明

1. 克隆仓库:
   ```
   git clone https://github.com/VxNull/project-time-tracker.git
   ```

2. 进入项目目录:
   ```
   cd project-time-tracker
   ```

3. 安装依赖:
   ```
   go mod tidy
   ```

4. 在项目根目录创建 `config.yaml` 文件,结构如下:
   ```yaml
   database:
     path: "./timetracker.db"
   
   admin:
     default_username: "admin"
     default_password: "password123"
   
   server:
     port: 8080
   
   session:
     secret_key: "your-secret-key"
   ```

5. 运行应用:
   ```
   go run main.go
   ```
   
   或者指定自定义配置文件路径:
   ```
   go run main.go -c /path/to/your/config.yaml
   ```

6. 在浏览器中访问 `http://localhost:8080` (或配置文件中指定的端口)

## 使用说明

### 管理员

1. 使用默认管理员账号登录 (用户名: admin, 密码: password123,或配置文件中指定的值)
2. 在仪表盘中查看整体工时统计
3. 管理项目和员工账户
4. 导出工时报表

### 员工

1. 使用分配的账号登录
2. 提交每日工时
3. 查看个人工时统计

## 贡献指南

我们欢迎任何形式的贡献,包括但不限于:

- 报告 bug
- 提出新功能建议
- 改进文档
- 提交代码修复或新功能

请遵循以下步骤:

1. Fork 本仓库
2. 创建您的特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交您的更改 (`git commit -m 'Add some AmazingFeature'`)
4. 将您的更改推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详细信息

## 联系方式

如果您有任何问题或建议,请通过以下方式联系我们:

- 项目链接: [https://github.com/VxNull/project-time-tracker](https://github.com/VxNull/project-time-tracker)
- 问题反馈: [https://github.com/VxNull/project-time-tracker/issues](https://github.com/VxNull/project-time-tracker/issues)
