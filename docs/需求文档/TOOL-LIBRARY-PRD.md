# AI 热点追踪平台 - 工具库模块产品需求文档

**版本**: v1.0
**创建日期**: 2026-04-11
**负责人**: 产品团队
**关联文档**: [需求文档](./AI%20热点追踪平台%20-%20需求文档.md)

---

## 文档说明

### 文档目的
本文档详细描述 AI 热点追踪平台中**工具库模块**的产品设计，包括功能需求、用户体验设计、数据模型、API 接口和开发计划。

### 适用范围
- 工具库模块的功能设计
- 用户体验和界面设计
- 数据库设计
- API 接口设计
- 开发和测试计划

---

## 目录

1. [概述](#1-概述)
2. [产品定位](#2-产品定位)
3. [用户分析](#3-用户分析)
4. [功能需求](#4-功能需求)
5. [用户体验设计](#5-用户体验设计)
6. [数据模型设计](#6-数据模型设计)
7. [API 接口设计](#7-api-接口设计)
8. [非功能需求](#8-非功能需求)
9. [开发计划](#9-开发计划)
10. [测试计划](#10-测试计划)
11. [风险与应对](#11-风险与应对)

---

## 1. 概述

### 1.1 模块背景

工具库模块是 AI 热点追踪平台的重要功能，旨在帮助用户：
- 发现和试用 AI 工具
- 了解各工具的特点和价格
- 获取使用技巧和提示词
- 参与社区评测和分享

### 1.2 模块目标

| 目标类别 | 具体目标 | 衡量指标 |
|---------|---------|---------|
| **内容质量** | 收录主流、实用的 AI 工具 | 200+ 工具收录 |
| **用户价值** | 帮助用户快速找到合适的工具 | 用户转化率 > 10% |
| **社区活跃** | 激励用户参与评测和分享 | 每日评测数 > 20 条 |
| **内容覆盖** | 涵盖主流应用场景 | 7 大类工具全覆盖 |

### 1.3 核心价值主张

> **"让 AI 工具选择更容易，让 AI 学习更高效"**

通过系统化的工具整理、客观的评测体系和活跃的社区交流，帮助用户快速找到适合自己的 AI 工具，并掌握使用技巧。

---

## 2. 产品定位

### 2.1 模块定位

| 维度 | 定位 |
|------|------|
| **核心角色** | 工具发现者 + 评测参与者 |
| **核心价值** | 工具整理 + 评测体系 + 社区交流 |
| **差异化** | 多维度评测 + 提示词模板库 + 工作流分享 |
| **内容来源** | 官方数据 + 用户评测 + 社区投稿 |

### 2.2 用户旅程

```
┌─────────────────────────────────────────────────────────────────┐
│                      工具库用户旅程                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  发现 → 浏览 → 试用 → 评测 → 分享 → 成就                        │
│   │        │         │         │         │         │           │
│   │        │         │         │         │         └─ 徽章/积分 │
│   │        │         │         │         └─ 收藏/分享          │
│   │        │         │         └─ 提交评测                  │
│   │        │         └─ 查看评测/提示词                      │
│   │        └─ 筛选/搜索定位                                    │
│   └─ 首页推荐/职业关联                                          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 2.3 功能架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      工具库功能架构                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 📚 工具浏览   │  │ 📝 工具评测   │  │ 💬 社区交流   │          │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤          │
│  │ • 分类列表   │  │ • 用户评分   │  │ • 评论讨论   │          │
│  │ • 搜索筛选   │  │ • 评论系统   │  │ • 工作流分享 │          │
│  │ • 排序功能   │  │ • 优点缺点   │  │ • 社区问答   │          │
│  │ • 提示词库   │  │ • 效果对比   │  │ • 徽章成就   │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                 │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │ 🎯 工具详情   │  │ ⚙️ 个性化推荐 │  │ 📊 数据统计   │          │
│  ├──────────────┤  ├──────────────┤  ├──────────────┤          │
│  │ • 工具介绍   │  │ • 根据职业推荐│  │ • 热门工具   │          │
│  │ • 参数对比   │  │ • 基于兴趣推荐│  │ • 评测统计   │          │
│  │ • 价格信息   │  │ • 相关工具推荐│  │ • 用户画像   │          │
│  │ • 官方链接   │  │               │  │               │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 3. 用户分析

### 3.1 用户画像

#### 用户类型 A：AI 新手用户

```
基本信息：
- 年龄：18-25 岁
- 身份：学生、刚入职场的职场新人
- 对 AI 的了解：零基础到入门

特征：
- 对 AI 很感兴趣但不知道从何开始
- 想学习如何使用 AI 工具提高效率
- 预算有限，希望找到免费工具
- 怕踩坑，需要可靠的推荐

核心需求：
- 找到适合新手的 AI 工具
- 获取简单易用的教程
- 了解工具的基本用法

典型场景：
"我听说 AI 很厉害，但不知道哪些工具好用，
想学习但不想花太多钱，有没有适合新手的推荐？"
```

#### 用户类型 B：进阶探索者

```
基本信息：
- 年龄：25-40 岁
- 职业：开发者、产品经理、运营等
- 对 AI 的了解：有一定了解

特征：
- 已经在使用一些 AI 工具
- 想探索更多工具和高级功能
- 关注工具的价格和性价比
- 愿意分享自己的使用经验

核心需求：
- 发现高质量的专业工具
- 对比不同工具的优缺点
- 获取高级使用技巧
- 参与社区交流

典型场景：
"我在用 ChatGPT，但听说还有其他更好的工具，
想了解一下 Midjourney、Claude 这些的区别，"
```

#### 用户类型 C：内容创作者

```
基本信息：
- 年龄：20-35 岁
- 职业：自媒体、设计师、作家
- 对 AI 的了解：主动学习者

特征：
- 创作内容需要辅助工具
- 经常尝试新的 AI 工具
- 愿意创作教程和评测
- 关注工具的实用性和效果

核心需求：
- 寻找提升效率的专业工具
- 学习高级使用技巧
- 撰写评测文章
- 参与社区内容创作

典型场景：
"我需要快速生成图片和文案，
想了解哪些工具效果最好，"
```

### 3.2 用户需求矩阵

| 需求类别 | 用户需求 | 优先级 |
|---------|---------|-------|
| **发现工具** | 发现新的 AI 工具 | P0 |
| | 按场景筛选工具 | P0 |
| | 搜索工具 | P0 |
| | 查看工具详情 | P0 |
| **了解工具** | 工具介绍 | P0 |
| | 价格信息 | P0 |
| | 官方链接 | P0 |
| | 参数对比 | P1 |
| **使用工具** | 提示词模板 | P1 |
| | 使用教程 | P1 |
| | 工作流分享 | P2 |
| **参与评测** | 提交评测 | P1 |
| | 查看评测 | P1 |
| | 评论交流 | P1 |
| | 评分系统 | P1 |
| **社区交流** | 收藏工具 | P1 |
| | 关注工具 | P1 |
| | 社区问答 | P2 |

---

## 4. 功能需求

### 4.1 核心功能列表

#### 4.1.1 工具浏览（P0）

| 功能点 | 描述 | 详情 |
|-------|------|------|
| **分类浏览** | 按功能类别浏览工具 | 7 大类工具分类 |
| **搜索功能** | 全文搜索 + 筛选 | 搜索、筛选、排序 |
| **排序功能** | 按评分、名称、更新时间排序 | 多种排序方式 |
| **标签过滤** | 按使用场景/行业/难度筛选 | 多维度标签 |

**分类体系**：
```
├── 写作类 (ChatGPT/Claude/Notion AI)
├── 图像类 (Midjourney/Stable Diffusion/DALL-E)
├── 视频类 (Runway/Pika/Sora)
├── 音频类 (ElevenLabs/Suno)
├── 代码类 (GitHub Copilot/Cursor)
├── 办公类 (Microsoft 365/Google Workspace)
└── 其他 (数据分析/研究/教育)
```

**筛选维度**：
- 按类别（7 个大类）
- 按价格（免费/付费/开源）
- 按难度（入门/进阶/高级）
- 按语言（中文/英文）
- 按更新时间（近 7 天/30 天/90 天）
- 按评分（> 4.0/4.5/4.8）

#### 4.1.2 工具详情（P0）

| 功能点 | 描述 | 详情 |
|-------|------|------|
| **工具介绍** | 详细介绍工具的功能特点 | 起源、功能、适用场景 |
| **参数对比** | 与同类工具对比 | 对比表格 |
| **价格信息** | 免费/付费价格 | 定价模式 |
| **官方链接** | 官网、文档、社区链接 | 多渠道跳转 |
| **使用场景** | 推荐使用场景 | 文字说明 |

**详情页包含信息**：
- 工具名称、图标、描述
- 官方网站、文档链接
- 定价模式（免费/订阅/按次）
- 类别、难度标签
- 评分和评测数量
- 最近更新时间
- 相关工具推荐

#### 4.1.3 工具评测（P1）

| 功能点 | 描述 | 详情 |
|-------|------|------|
| **用户评分** | 1-5 星评分系统 | 平均分、评分人数 |
| **用户评论** | 评论系统 | 文字评论、点赞 |
| **优点缺点** | 用户总结工具优缺点 | 标签化展示 |
| **效果对比** | 效果截图对比 | 评分维度对比 |

**评分维度**：
- 易用性（1-5 分）
- 效果质量（1-5 分）
- 价格性价比（1-5 分）
- 功能丰富度（1-5 分）
- 更新频率（1-5 分）
- 客服支持（1-5 分）

#### 4.1.4 提示词模板（P1）

| 功能点 | 描述 | 详情 |
|-------|------|------|
| **模板库** | 常用提示词模板 | 分类整理 |
| **分类浏览** | 按使用场景分类 | 写作/代码/设计等 |
| **收藏功能** | 收藏常用模板 | 本地保存 |
| **复制功能** | 一键复制到剪贴板 | 快速使用 |
| **用户投稿** | 用户分享模板 | 社区贡献 |

**模板分类**：
```
├── 写作类
│   ├── 博客写作
│   ├── 邮件撰写
│   ├── 脚本撰写
│   └── 内容创作
├── 代码类
│   ├── 代码生成
│   ├── 代码解释
│   ├── 代码优化
│   └── 调试帮助
├── 设计类
│   ├── 文案生成
│   ├── 标题生成
│   ├── 创意构思
│   └── 用户洞察
└── 其他
    ├── 学习辅导
    ├── 翻译润色
    └── 摘要总结
```

#### 4.1.5 社区交流（P1）

| 功能点 | 描述 | 详情 |
|-------|------|------|
| **评论系统** | 对工具的评论 | 回复、点赞 |
| **问答社区** | 用户提问解答 | 匿名提问、专家解答 |
| **工作流分享** | 用户分享使用工作流 | 文字描述 + 截图 |
| **徽章系统** | 鼓励用户贡献 | 评测徽章、创作徽章 |

**徽章体系**：
```
├── 评测徽章
│   ├── 首次评测（1 个）
│   ├── 精选评测（10 个）
│   ├── 深度评测（50 个）
│   └── 优质评测（100 个）
├── 创作徽章
│   ├── 提示词达人（5 个模板）
│   ├── 工作流专家（3 个分享）
│   └── 社区活跃（每日登录 30 天）
└── 贡献徽章
    ├── 优质评论（50 条）
    ├── 提出问题（10 个）
    └── 帮助解答（30 个）
```

#### 4.1.6 个性化推荐（P2）

| 功能点 | 描述 | 详情 |
|-------|------|------|
| **职业推荐** | 根据用户职业推荐工具 | 基于职业画像 |
| **兴趣推荐** | 基于浏览历史推荐 | 相关工具推荐 |
| **相关工具** | 同类或互补工具推荐 | 协同过滤 |

**推荐逻辑**：
- 职业推荐：用户选择职业 → 推荐该职业常用的工具
- 兴趣推荐：根据用户浏览、收藏的工具 → 推荐相似工具
- 相关工具：同一类别、使用场景相似的工具

### 4.2 功能优先级矩阵

| 功能点 | 用户价值 | 开发难度 | 优先级 |
|-------|---------|---------|-------|
| 分类浏览 | 高 | 低 | P0 |
| 搜索筛选 | 高 | 中 | P0 |
| 工具详情 | 高 | 中 | P0 |
| 参数对比 | 中 | 中 | P1 |
| 提示词模板 | 高 | 中 | P1 |
| 用户评测 | 高 | 高 | P1 |
| 评论系统 | 高 | 低 | P1 |
| 社区问答 | 中 | 高 | P2 |
| 工作流分享 | 中 | 中 | P2 |
| 徽章系统 | 低 | 中 | P2 |
| 个性化推荐 | 中 | 高 | P3 |

---

## 5. 用户体验设计

### 5.1 信息架构

```
┌─────────────────────────────────────────────────────────────────┐
│                      工具库信息架构                              │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  首页
│  ├─ 搜索框（核心搜索）
│  ├─ 分类标签（7 大类）
│  ├─ 热门工具排行榜
│  ├─ 新上线路由
│  └─ 推荐工具
│                                                                 │
│  分类页（以写作类为例）
│  ├─ 搜索框
│  ├─ 筛选器（价格/难度/评分/更新）
│  ├─ 工具列表
│  │   └─ 卡片（图标、名称、评分、免费/付费标签）
│  └─ 分页
│                                                                 │
│  工具详情页
│  ├─ 工具信息（图标、名称、描述）
│  ├─ 快速操作（评分、收藏、评论）
│  ├─ 参数对比
│  ├─ 提示词模板入口
│  ├─ 用户评测（评分、评论列表）
│  └─ 相关工具推荐
│                                                                 │
│  提示词模板页
│  ├─ 搜索筛选
│  ├─ 模板分类
│  └─ 模板列表（标题、用途、复制按钮）
│                                                                 │
│  评测提交页
│  ├─ 工具信息
│  ├─ 评分维度（6 个维度）
│  ├─ 优点/缺点输入
│  └─ 评论内容
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 5.2 页面布局设计

#### 5.2.1 首页布局

```
┌─────────────────────────────────────────────────────────────────┐
│                        工具库首页布局                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  顶部导航
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  AI 热点追踪      资讯    职业    学习路径    工具库        │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  Hero 区域
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  🚀 发现最适合你的 AI 工具                               │ │
│  │  按分类浏览、对比评测、获取使用技巧                      │ │
│  │                                                          │ │
│  │  [搜索工具]                                              │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  分类入口（7 大类）
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  写作  图像  视频  音频  代码  办公  其他                 │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  热门工具排行榜
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  🔥 热门工具 TOP 5                                       │ │
│  │  1. ChatGPT      ⭐ 4.9        ⭐ 2.3k 评测              │ │
│  │  2. Midjourney   ⭐ 4.8        ⭐ 1.8k 评测              │ │
│  │  3. Claude       ⭐ 4.8        ⭐ 1.5k 评测              │ │
│  │  4. Copilot      ⭐ 4.7        ⭐ 1.2k 评测              │ │
│  │  5. Runway       ⭐ 4.6        ⭐ 980   评测              │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  新上线路由
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  ✨ 本周新工具                                           │ │
│  │  [工具卡片]  [工具卡片]  [工具卡片]  [工具卡片]          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  推荐工具（基于职业）
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  👤 推荐给：产品经理                                     │ │
│  │  ChatGPT、Notion AI、Gamma                             │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

#### 5.2.2 工具详情页布局

```
┌─────────────────────────────────────────────────────────────────┐
│                      工具详情页布局                             │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  顶部信息
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  [返回]  ChatGPT  [评论] [收藏]                         │ │
│  │                                                          │ │
│  │  🤖 ChatGPT                                                │ │
│  │  OpenAI 开发的大型语言模型                                │ │
│  │  📊 评分: 4.8  🔥 2.3k 评测  💬 580 评论                 │ │
│  │  💰 定价: 免费版/付费版  🏷️ 写作类  ⭐ 入门               │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  工具介绍
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  ChatGPT 是由 OpenAI 开发的大型语言模型，能够理解并生成  │ │
│  │  人类语言，适用于写作、编程、学习等多种场景。             │ │
│  │  [查看更多]                                              │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  官方链接
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  🌐 官网: https://chat.openai.com                        │ │
│  │  📚 文档: https://help.openai.com                         │ │
│  │  💬 社区: https://chat.openai.com/chat                    │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  价格信息
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  💰 定价模式                                              │ │
│  │  • 免费版: 每日消息限额                                   │ │
│  │  • Plus: $20/月，无限额                                   │ │
│  │  • 企业版: 定制价格                                       │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  快速评分
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  👍 快速评分                                              │ │
│  │  [⭐⭐⭐⭐⭐]                                            │ │
│  │  (点击即可评分)                                          │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  参数对比
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  📊 与同类工具对比                                       │ │
│  │  [表格展示]                                              │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  提示词模板入口
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  📝 ChatGPT 提示词模板 (320+ 模板)                       │ │
│  │  [查看全部] [收藏模板]                                  │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  用户评测
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  💬 用户评测 (580 条)                                     │ │
│  │  [评分分布图]                                            │ │
│  │  [最新评论]                                              │ │
│  │  [查看全部]                                              │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
│  相关工具
│  ┌──────────────────────────────────────────────────────────┐ │
│  │  🔗 相关工具                                              │ │
│  │  [工具卡片]  [工具卡片]  [工具卡片]                      │ │
│  └──────────────────────────────────────────────────────────┘ │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 5.3 交互设计

#### 5.3.1 搜索交互

- **搜索框**：
  - 支持实时搜索（输入即过滤）
  - 支持模糊搜索
  - 搜索历史记录
  - 搜索建议（自动补全）

- **筛选器**：
  - 多维度筛选（类别、价格、难度等）
  - 筛选条件支持组合
  - 筛选结果实时更新
  - 清空筛选按钮

- **排序**：
  - 按评分排序（默认）
  - 按名称排序
  - 按更新时间排序
  - 按评测数量排序

#### 5.3.2 评分交互

- **快速评分**：
  - 1-5 星点击评分
  - 支持取消评分
  - 评分后显示感谢语

- **详细评分**：
  - 6 个维度评分
  - 每个维度可单独修改
  - 评分后弹出表单

#### 5.3.3 评测提交交互

- **表单设计**：
  - 预填工具信息
  - 6 个评分维度（滑块 + 数字）
  - 优点/缺点输入框
  - 详细评论输入框

- **交互流程**：
  - 点击"写评测"按钮
  - 填写评分和评论
  - 提交后立即显示
  - 获得评测徽章

---

## 6. 数据模型设计

### 6.1 核心数据表

#### 6.1.1 工具表 (tools)

```sql
CREATE TABLE tools (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(200) NOT NULL COMMENT '工具名称',
  slug VARCHAR(200) UNIQUE NOT NULL COMMENT 'URL友好的标识',
  icon VARCHAR(500) COMMENT '工具图标URL',
  description TEXT COMMENT '工具描述',
  official_url VARCHAR(500) COMMENT '官方网站',
  documentation_url VARCHAR(500) COMMENT '文档链接',
  pricing TEXT COMMENT '定价信息（JSON格式）',
  category_id INT COMMENT '所属类别ID',
  difficulty VARCHAR(20) COMMENT '难度等级：beginner/intermediate/advanced',
  rating DECIMAL(2,1) DEFAULT 0.00 COMMENT '平均评分',
  review_count INT DEFAULT 0 COMMENT '评测数量',
  view_count INT DEFAULT 0 COMMENT '浏览量',
  popularity INT DEFAULT 0 COMMENT '热度值',
  is_free BOOLEAN DEFAULT TRUE COMMENT '是否免费',
  is_active BOOLEAN DEFAULT TRUE COMMENT '是否上架',
  last_updated_at DATETIME COMMENT '最后更新时间',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_slug (slug),
  INDEX idx_category (category_id),
  INDEX idx_rating (rating),
  INDEX idx_popularity (popularity),
  INDEX idx_is_free (is_free),
  INDEX idx_is_active (is_active)
) COMMENT='AI工具表';
```

#### 6.1.2 工具类别表 (tool_categories)

```sql
CREATE TABLE tool_categories (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '类别名称',
  slug VARCHAR(100) UNIQUE NOT NULL COMMENT 'URL友好的标识',
  icon VARCHAR(500) COMMENT '图标',
  description TEXT COMMENT '类别描述',
  sort_order INT DEFAULT 0 COMMENT '排序',
  is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_slug (slug),
  INDEX idx_sort_order (sort_order)
) COMMENT='工具类别表';
```

#### 6.1.3 工具标签表 (tool_tags)

```sql
CREATE TABLE tool_tags (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '标签名称',
  slug VARCHAR(100) UNIQUE NOT NULL COMMENT 'URL友好的标识',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_slug (slug)
) COMMENT='工具标签表';
```

#### 6.1.4 工具-标签关联表 (tool_tag_relations)

```sql
CREATE TABLE tool_tag_relations (
  tool_id BIGINT NOT NULL,
  tag_id INT NOT NULL,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  PRIMARY KEY (tool_id, tag_id),
  INDEX idx_tag_id (tag_id),
  INDEX idx_tool_id (tool_id),
  FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE,
  FOREIGN KEY (tag_id) REFERENCES tool_tags(id) ON DELETE CASCADE
) COMMENT='工具标签关联表';
```

#### 6.1.5 用户评测表 (tool_reviews)

```sql
CREATE TABLE tool_reviews (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(100) COMMENT '用户ID（登录用户/访客）',
  tool_id BIGINT NOT NULL COMMENT '工具ID',
  rating TINYINT COMMENT '评分1-5',
  ease_of_use TINYINT COMMENT '易用性1-5',
  effectiveness TINYINT COMMENT '效果质量1-5',
  value_for_money TINYINT COMMENT '性价比1-5',
  features TINYINT COMMENT '功能丰富度1-5',
  update_frequency TINYINT COMMENT '更新频率1-5',
  support TINYINT COMMENT '客服支持1-5',
  pros TEXT COMMENT '优点',
  cons TEXT COMMENT '缺点',
  comment TEXT COMMENT '详细评论',
  is_anonymous BOOLEAN DEFAULT FALSE COMMENT '是否匿名',
  status TINYINT DEFAULT 1 COMMENT '状态：1-审核通过，0-待审核，2-已拒绝',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_user_id (user_id),
  INDEX idx_tool_id (tool_id),
  INDEX idx_rating (rating),
  INDEX idx_created_at (created_at),
  FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
) COMMENT='用户评测表';
```

#### 6.1.6 评论表 (comments)

```sql
CREATE TABLE comments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(100) COMMENT '用户ID',
  commentable_id BIGINT NOT NULL COMMENT '评论对象ID',
  commentable_type VARCHAR(50) NOT NULL COMMENT '评论对象类型：tool_review/tool',
  parent_id BIGINT COMMENT '父评论ID',
  content TEXT NOT NULL COMMENT '评论内容',
  is_anonymous BOOLEAN DEFAULT FALSE COMMENT '是否匿名',
  likes INT DEFAULT 0 COMMENT '点赞数',
  is_spam BOOLEAN DEFAULT FALSE COMMENT '是否为垃圾评论',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_commentable (commentable_id, commentable_type),
  INDEX idx_parent_id (parent_id),
  INDEX idx_created_at (created_at),
  FOREIGN KEY (parent_id) REFERENCES comments(id) ON DELETE CASCADE
) COMMENT='评论表';
```

#### 6.1.7 提示词模板表 (prompt_templates)

```sql
CREATE TABLE prompt_templates (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(200) NOT NULL COMMENT '模板名称',
  slug VARCHAR(200) UNIQUE NOT NULL COMMENT 'URL友好的标识',
  description TEXT COMMENT '模板描述',
  content TEXT NOT NULL COMMENT '提示词内容',
  tool_id BIGINT COMMENT '适用工具ID（可为空）',
  category_id INT COMMENT '模板类别ID',
  use_cases TEXT COMMENT '使用场景（JSON数组）',
  tags TEXT COMMENT '标签（JSON数组）',
  example_response TEXT COMMENT '示例回复',
  likes INT DEFAULT 0 COMMENT '点赞数',
  views INT DEFAULT 0 COMMENT '浏览量',
  is_featured BOOLEAN DEFAULT FALSE COMMENT '是否精选',
  status TINYINT DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
  created_by VARCHAR(100) COMMENT '创建者',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_slug (slug),
  INDEX idx_tool_id (tool_id),
  INDEX idx_category_id (category_id),
  INDEX idx_created_at (created_at),
  INDEX idx_status (status),
  FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE SET NULL,
  FOREIGN KEY (category_id) REFERENCES tool_categories(id) ON DELETE SET NULL
) COMMENT='提示词模板表';
```

#### 6.1.8 提示词模板分类表 (prompt_template_categories)

```sql
CREATE TABLE prompt_template_categories (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '分类名称',
  slug VARCHAR(100) UNIQUE NOT NULL COMMENT 'URL友好的标识',
  description TEXT COMMENT '分类描述',
  icon VARCHAR(500) COMMENT '图标',
  sort_order INT DEFAULT 0 COMMENT '排序',
  is_active BOOLEAN DEFAULT TRUE COMMENT '是否启用',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

  INDEX idx_slug (slug),
  INDEX idx_sort_order (sort_order)
) COMMENT='提示词模板分类表';
```

#### 6.1.9 用户收藏表 (user_favorites)

```sql
CREATE TABLE user_favorites (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(100) NOT NULL COMMENT '用户ID',
  tool_id BIGINT NOT NULL COMMENT '工具ID',
  note VARCHAR(500) COMMENT '收藏备注',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  UNIQUE KEY uk_user_tool (user_id, tool_id),
  INDEX idx_user_id (user_id),
  INDEX idx_tool_id (tool_id),
  FOREIGN KEY (tool_id) REFERENCES tools(id) ON DELETE CASCADE
) COMMENT='用户收藏表';
```

#### 6.1.10 徽章表 (badges)

```sql
CREATE TABLE badges (
  id INT PRIMARY KEY AUTO_INCREMENT,
  name VARCHAR(100) NOT NULL COMMENT '徽章名称',
  slug VARCHAR(100) UNIQUE NOT NULL COMMENT 'URL友好的标识',
  description TEXT COMMENT '徽章描述',
  icon VARCHAR(500) COMMENT '徽章图标',
  type VARCHAR(50) NOT NULL COMMENT '类型：review/contribution/social',
  condition_type VARCHAR(50) NOT NULL COMMENT '获取条件类型',
  condition_value INT NOT NULL COMMENT '获取条件值',
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP,

  INDEX idx_type (type)
) COMMENT='徽章表';
```

#### 6.1.11 用户徽章表 (user_badges)

```sql
CREATE TABLE user_badges (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  user_id VARCHAR(100) NOT NULL COMMENT '用户ID',
  badge_id INT NOT NULL COMMENT '徽章ID',
  issued_at DATETIME DEFAULT CURRENT_TIMESTAMP,
  unique_index (user_id, badge_id)
) COMMENT='用户徽章表';
```

### 6.2 关系说明

```
┌──────────────┐
│   users      │
└──────┬───────┘
       │
       │ 1:N
       ▼
┌──────────────┐      1:N      ┌──────────────┐
│  tool_reviews │──────────────>│   tools      │
└──────────────┘               └──────┬───────┘
                                      │ 1:N
                                      ▼
                               ┌──────────────┐
                               │ user_favorites│
                               └──────────────┘

┌──────────────┐
│   tools      │
└──────┬───────┘
       │ 1:N
       ▼
┌──────────────┐
│  tool_reviews│
└──────────────┘

┌──────────────┐
│   tools      │
└──────┬───────┘
       │ 1:N
       ▼
┌──────────────┐
│ prompt_templates
└──────────────┘

┌──────────────┐
│   users      │
└──────┬───────┘
       │ 1:N
       ▼
┌──────────────┐
│  comments    │
└──────────────┘
```

---

## 7. API 接口设计

### 7.1 API 基础信息

- **基础路径**: `/api/v1/tools`
- **认证方式**: 部分接口需要（评测提交、收藏等）
- **响应格式**: JSON
- **编码格式**: UTF-8

### 7.2 接口列表

#### 7.2.1 工具浏览接口

##### 获取工具列表

**接口**: `GET /api/v1/tools`

**请求参数**:
| 参数 | 类型 | 必填 | 说明 | 示例 |
|------|------|------|------|------|
| page | int | 否 | 页码，默认 1 | 1 |
| page_size | int | 否 | 每页数量，默认 20 | 20 |
| category_id | int | 否 | 类别ID | 1 |
| is_free | boolean | 否 | 是否免费 | true/false |
| difficulty | string | 否 | 难度等级 | beginner/intermediate/advanced |
| min_rating | decimal | 否 | 最低评分 | 4.0 |
| sort_by | string | 否 | 排序字段 | rating/update_time/popularity |
| order | string | 否 | 排序方向 | asc/desc |
| search | string | 否 | 搜索关键词 | ChatGPT |

**成功响应** (200 OK):
```json
{
  "list": [
    {
      "id": 1,
      "name": "ChatGPT",
      "slug": "chatgpt",
      "icon": "🤖",
      "description": "OpenAI 开发的大型语言模型",
      "official_url": "https://chat.openai.com",
      "documentation_url": "https://help.openai.com",
      "category": {
        "id": 1,
        "name": "写作类",
        "slug": "writing"
      },
      "difficulty": "beginner",
      "is_free": true,
      "rating": 4.8,
      "review_count": 2300,
      "view_count": 15000,
      "popularity": 95,
      "tags": ["大语言模型", "对话"],
      "last_updated_at": "2026-04-10T10:00:00Z"
    }
  ],
  "total": 200,
  "page": 1,
  "page_size": 20
}
```

##### 获取工具详情

**接口**: `GET /api/v1/tools/{slug}`

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 工具slug |

**成功响应** (200 OK):
```json
{
  "id": 1,
  "name": "ChatGPT",
  "slug": "chatgpt",
  "icon": "🤖",
  "description": "OpenAI 开发的大型语言模型，能够理解并生成人类语言",
  "official_url": "https://chat.openai.com",
  "documentation_url": "https://help.openai.com",
  "pricing": {
    "free": "每天有限的消息配额",
    "plus": "$20/月，无限制",
    "enterprise": "请联系我们获取报价"
  },
  "category": {
    "id": 1,
    "name": "写作类",
    "slug": "writing"
  },
  "difficulty": "beginner",
  "is_free": true,
  "rating": 4.8,
  "review_count": 2300,
  "view_count": 15000,
  "popularity": 95,
  "tags": ["大语言模型", "对话"],
  "feature_comparison": {
    "gpt4": true,
    "联网搜索": true,
    "插件支持": true
  },
  "created_at": "2026-01-01T00:00:00Z",
  "updated_at": "2026-04-10T10:00:00Z"
}
```

##### 获取工具分类列表

**接口**: `GET /api/v1/tools/categories`

**成功响应** (200 OK):
```json
{
  "list": [
    {
      "id": 1,
      "name": "写作类",
      "slug": "writing",
      "icon": "✍️",
      "description": "用于写作、文案创作的 AI 工具",
      "tool_count": 45
    },
    {
      "id": 2,
      "name": "图像类",
      "slug": "image",
      "icon": "🎨",
      "description": "用于图像生成、编辑的 AI 工具",
      "tool_count": 38
    }
  ]
}
```

#### 7.2.2 工具评测接口

##### 获取工具评测列表

**接口**: `GET /api/v1/tools/{slug}/reviews`

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 工具slug |

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |
| sort_by | string | 否 | 排序方式：newest/highest_rating |
| user_id | string | 否 | 用户ID（获取个人评测） |

**成功响应** (200 OK):
```json
{
  "list": [
    {
      "id": 1,
      "user": {
        "id": "user123",
        "nickname": "张三",
        "avatar": "https://example.com/avatar.jpg"
      },
      "rating": 5,
      "ease_of_use": 5,
      "effectiveness": 5,
      "value_for_money": 4,
      "features": 5,
      "update_frequency": 4,
      "support": 5,
      "pros": ["响应速度快", "功能强大", "易用性好"],
      "cons": ["免费版限额"],
      "comment": "非常好用的AI工具，强烈推荐！",
      "likes": 120,
      "is_liked": false,
      "created_at": "2026-04-10T10:00:00Z"
    }
  ],
  "total": 2300,
  "page": 1,
  "page_size": 20
}
```

##### 获取评测详情

**接口**: `GET /api/v1/reviews/{review_id}`

**成功响应** (200 OK):
```json
{
  "id": 1,
  "user": {
    "id": "user123",
    "nickname": "张三",
    "is_anonymous": false
  },
  "tool": {
    "id": 1,
    "name": "ChatGPT",
    "slug": "chatgpt",
    "icon": "🤖"
  },
  "rating": 5,
  "ease_of_use": 5,
  "effectiveness": 5,
  "value_for_money": 4,
  "features": 5,
  "update_frequency": 4,
  "support": 5,
  "pros": ["响应速度快", "功能强大", "易用性好"],
  "cons": ["免费版限额"],
  "comment": "非常好用的AI工具，强烈推荐！",
  "likes": 120,
  "is_liked": false,
  "is_anonymous": false,
  "created_at": "2026-04-10T10:00:00Z"
}
```

##### 提交工具评测

**接口**: `POST /api/v1/tools/{slug}/reviews`

**请求参数**:
```json
{
  "rating": 5,
  "ease_of_use": 5,
  "effectiveness": 5,
  "value_for_money": 4,
  "features": 5,
  "update_frequency": 4,
  "support": 5,
  "pros": ["响应速度快", "功能强大", "易用性好"],
  "cons": ["免费版限额"],
  "comment": "非常好用的AI工具，强烈推荐！",
  "is_anonymous": false
}
```

**成功响应** (200 OK):
```json
{
  "success": true,
  "message": "评测提交成功",
  "review_id": 1234
}
```

##### 获取评分分布

**接口**: `GET /api/v1/tools/{slug}/reviews/rating-distribution`

**成功响应** (200 OK):
```json
{
  "distribution": {
    "1": 50,
    "2": 120,
    "3": 250,
    "4": 450,
    "5": 1430
  },
  "total": 2300,
  "average": 4.8
}
```

#### 7.2.3 提示词模板接口

##### 获取模板列表

**接口**: `GET /api/v1/prompts`

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |
| category_id | int | 否 | 模板类别ID |
| tool_id | int | 否 | 工具ID |
| search | string | 否 | 搜索关键词 |

**成功响应** (200 OK):
```json
{
  "list": [
    {
      "id": 1,
      "name": "博客文章写作助手",
      "slug": "blog-article-writer",
      "description": "帮助快速撰写高质量博客文章",
      "content": "你是一名专业的内容作家，擅长撰写...",
      "category": {
        "id": 1,
        "name": "写作类",
        "slug": "writing"
      },
      "tool": {
        "id": 1,
        "name": "ChatGPT",
        "slug": "chatgpt",
        "icon": "🤖"
      },
      "use_cases": ["博客写作", "内容创作", "文章优化"],
      "tags": ["写作", "博客"],
      "likes": 150,
      "views": 1200,
      "is_featured": true
    }
  ],
  "total": 320,
  "page": 1,
  "page_size": 20
}
```

##### 获取模板详情

**接口**: `GET /api/v1/prompts/{slug}`

**成功响应** (200 OK):
```json
{
  "id": 1,
  "name": "博客文章写作助手",
  "slug": "blog-article-writer",
  "description": "帮助快速撰写高质量博客文章",
  "content": "你是一名专业的内容作家，擅长撰写博客文章。请按照以下步骤帮我写一篇博客文章：\n\n1. 确定文章主题\n2. 构建文章结构\n3. 撰写内容\n4. 审查优化\n\n主题是...",
  "category": {
    "id": 1,
    "name": "写作类",
    "slug": "writing"
  },
  "tool": {
    "id": 1,
    "name": "ChatGPT",
    "slug": "chatgpt",
    "icon": "🤖"
  },
  "use_cases": ["博客写作", "内容创作", "文章优化"],
  "tags": ["写作", "博客"],
  "example_response": "（示例回复）",
  "likes": 150,
  "views": 1200,
  "is_featured": true,
  "created_by": {
    "id": "user123",
    "nickname": "李四"
  },
  "created_at": "2026-03-01T10:00:00Z",
  "updated_at": "2026-04-05T15:30:00Z"
}
```

##### 提交提示词模板

**接口**: `POST /api/v1/prompts`

**请求参数**:
```json
{
  "name": "代码解释助手",
  "slug": "code-explainer",
  "description": "帮助解释代码逻辑",
  "content": "你是一名资深程序员，擅长解释代码逻辑。请帮我解释以下代码...",
  "category_id": 2,
  "tool_id": 3,
  "use_cases": ["代码学习", "调试帮助"],
  "tags": ["代码", "编程"]
}
```

**成功响应** (200 OK):
```json
{
  "success": true,
  "message": "模板提交成功",
  "template_id": 123
}
```

#### 7.2.4 收藏功能接口

##### 收藏工具

**接口**: `POST /api/v1/tools/{slug}/favorite`

**成功响应** (200 OK):
```json
{
  "success": true,
  "message": "已收藏工具"
}
```

##### 取消收藏

**接口**: `DELETE /api/v1/tools/{slug}/favorite`

**成功响应** (200 OK):
```json
{
  "success": true,
  "message": "已取消收藏"
}
```

##### 获取用户收藏列表

**接口**: `GET /api/v1/users/favorites`

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页数量，默认 20 |

**成功响应** (200 OK):
```json
{
  "list": [
    {
      "id": 1,
      "tool": {
        "id": 1,
        "name": "ChatGPT",
        "slug": "chatgpt",
        "icon": "🤖"
      },
      "note": "工作必备",
      "created_at": "2026-04-01T10:00:00Z"
    }
  ],
  "total": 15,
  "page": 1,
  "page_size": 20
}
```

#### 7.2.5 推荐功能接口

##### 获取职业推荐工具

**接口**: `GET /api/v1/tools/recommendations/profession/{profession_slug}`

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| profession_slug | string | 是 | 职业slug |

**成功响应** (200 OK):
```json
{
  "profession": {
    "id": 1,
    "name": "产品经理",
    "slug": "product-manager"
  },
  "recommendations": [
    {
      "tool": {
        "id": 1,
        "name": "ChatGPT",
        "slug": "chatgpt",
        "icon": "🤖"
      },
      "reason": "适合产品需求文档撰写和用户调研",
      "matched_features": ["写作", "调研"]
    },
    {
      "tool": {
        "id": 5,
        "name": "Notion AI",
        "slug": "notion-ai",
        "icon": "📝"
      },
      "reason": "适合知识管理和文档协作",
      "matched_features": ["写作", "协作"]
    }
  ]
}
```

##### 获取相关工具

**接口**: `GET /api/v1/tools/{slug}/related`

**路径参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| slug | string | 是 | 工具slug |

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| limit | int | 否 | 返回数量，默认 6 |

**成功响应** (200 OK):
```json
{
  "tool": {
    "id": 1,
    "name": "ChatGPT",
    "slug": "chatgpt",
    "icon": "🤖"
  },
  "related": [
    {
      "tool": {
        "id": 2,
        "name": "Claude",
        "slug": "claude",
        "icon": "🧠"
      },
      "similarity": 0.85,
      "category": "写作类"
    },
    {
      "tool": {
        "id": 3,
        "name": "Copilot",
        "slug": "copilot",
        "icon": "💻"
      },
      "similarity": 0.78,
      "category": "代码类"
    }
  ]
}
```

---

## 8. 非功能需求

### 8.1 性能要求

| 指标 | 目标值 | 说明 |
|------|-------|------|
| 页面加载时间 | < 2 秒 | 工具列表页 |
| API 响应时间 | < 200ms | P95 延迟 |
| 图片加载时间 | < 1 秒 | CDN 加速 |
| 并发支持 | 1000+ 并发 | 初期目标 |

### 8.2 可用性要求

| 指标 | 目标值 | 说明 |
|------|-------|------|
| 系统可用性 | > 99% | 月度 |
| 数据备份 | 每日自动备份 | 保留 30 天 |

### 8.3 安全要求

| 要求 | 说明 |
|------|------|
| HTTPS | 全站强制 HTTPS |
| 参数校验 | 所有输入参数校验 |
| SQL 注入防护 | 参数化查询 |
| XSS 防护 | 内容过滤、CSP |
| 速率限制 | API 限流防刷 |

### 8.4 SEO 要求

| 要求 | 说明 |
|------|------|
| Meta 标签 | 每页独立 title/description |
| Sitemap | 自动生成 |
| 结构化数据 | Schema.org |
| URL 规范 | 语义化、静态化 |

---

## 9. 开发计划

### 9.1 开发里程碑

```
┌─────────────────────────────────────────────────────────────────┐
│                      开发里程碑                                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  M1 (第 2 周)    │  工具列表页、详情页、分类浏览                │
│  M2 (第 4 周)    │  搜索、筛选、排序功能                        │
│  M3 (第 5 周)    │  评测系统、评分功能                          │
│  M4 (第 6 周)    │  提示词模板、收藏功能                        │
│  M5 (第 7 周)    │  评论系统、徽章系统                          │
│  M6 (第 8 周)    │  个性化推荐、数据统计                        │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

### 9.2 详细排期

#### Week 1-2: 基础功能开发

| 任务 | 工作量 | 负责人 | 产出 |
|------|-------|--------|------|
| 数据库设计 | 2天 | 后端 | 11张数据表 |
| API 接口开发 | 3天 | 后端 | 工具CRUD接口 |
| 工具列表页 | 2天 | 前端 | 列表+筛选+排序 |
| 工具详情页 | 2天 | 前端 | 详情展示 |
| 分类管理 | 1天 | 后端 | 7大类分类 |
| 基础样式 | 1天 | 前端 | 工具卡片样式 |

**验收标准**：
- 工具列表页可正常浏览
- 工具详情页信息完整
- 分类功能正常

#### Week 3-4: 搜索功能开发

| 任务 | 工作量 | 负责人 | 产出 |
|------|-------|--------|------|
| 搜索功能 | 2天 | 后端 | 搜索接口 |
| 筛选器开发 | 2天 | 后端 | 多维度筛选 |
| 排序功能 | 1天 | 后端 | 排序接口 |
| 搜索页开发 | 3天 | 前端 | 搜索+筛选+排序 |
| 性能优化 | 2天 | 全栈 | 搜索优化 |

**验收标准**：
- 搜索响应时间 < 200ms
- 多维度筛选正常
- 排序功能正确

#### Week 5: 评测系统开发

| 任务 | 工作量 | 负责人 | 产出 |
|------|-------|--------|------|
| 评测表开发 | 1天 | 后端 | 评测数据结构 |
| 评测接口 | 2天 | 后端 | 评测CRUD接口 |
| 评分功能 | 2天 | 前端 | 评分组件 |
| 评测列表 | 2天 | 前端 | 评测展示 |
| 评测提交 | 1天 | 前端 | 评测表单 |
| 数据统计 | 2天 | 后端 | 评分分布 |

**验收标准**：
- 评测可正常提交
- 评分功能正常
- 评分分布准确

#### Week 6: 提示词模板功能

| 任务 | 工作量 | 负责人 | 产出 |
|------|-------|--------|------|
| 模板表开发 | 1天 | 后端 | 模板数据结构 |
| 模板接口 | 2天 | 后端 | 模板CRUD接口 |
| 模板分类 | 1天 | 后端 | 模板分类管理 |
| 模板列表页 | 3天 | 前端 | 模板浏览 |
| 模板详情页 | 2天 | 前端 | 模板展示 |
| 复制功能 | 1天 | 前端 | 一键复制 |

**验收标准**：
- 模板可正常提交
- 模板分类清晰
- 复制功能可用

#### Week 7-8: 社区功能与推荐

| 任务 | 工作量 | 负责人 | 产出 |
|------|-------|--------|------|
| 评论系统 | 3天 | 后端 | 评论接口 |
| 徽章系统 | 2天 | 后端 | 徽章管理 |
| 收藏功能 | 2天 | 后端 | 收藏接口 |
| 评论模块 | 3天 | 前端 | 评论展示 |
| 徽章模块 | 2天 | 前端 | 徽章展示 |
| 收藏功能 | 2天 | 前端 | 收藏管理 |
| 推荐功能 | 3天 | 后端 | 推荐算法 |
| 个性化推荐 | 2天 | 前端 | 推荐展示 |

**验收标准**：
- 评论功能正常
- 徽章可正常获取
- 推荐准确

### 9.3 人力配置

| 角色 | 人数 | 职责 |
|------|------|------|
| 后端开发 | 1 | API接口、数据库、业务逻辑 |
| 前端开发 | 1 | 页面开发、交互实现 |
| UI 设计 | 0.5 | 页面设计 |
| 测试 | 0.5 | 测试用例、测试执行 |

---

## 10. 测试计划

### 10.1 测试策略

| 测试类型 | 覆盖范围 | 优先级 |
|---------|---------|-------|
| 单元测试 | 核心业务逻辑 | P0 |
| 集成测试 | API接口 | P0 |
| E2E测试 | 核心用户流程 | P1 |
| 性能测试 | 搜索、列表加载 | P1 |
| 兼容性测试 | 主流浏览器、移动端 | P1 |

### 10.2 测试用例

#### 功能测试

1. 工具列表页
   - 分页功能
   - 筛选功能
   - 排序功能
   - 搜索功能

2. 工具详情页
   - 信息展示
   - 官方链接跳转
   - 相关工具推荐

3. 评测系统
   - 评分提交
   - 评测列表
   - 评论功能

4. 提示词模板
   - 模板浏览
   - 模板详情
   - 模板复制

5. 社区功能
   - 收藏功能
   - 评论功能
   - 徽章获取

#### 性能测试

1. 搜索功能
   - 响应时间 < 200ms
   - 并发1000用户测试

2. 列表加载
   - 列表响应时间 < 500ms
   - 100条数据加载测试

### 10.3 测试环境

- 开发环境
- 测试环境
- 预生产环境

---

## 11. 风险与应对

### 11.1 风险矩阵

| 风险 | 概率 | 影响 | 应对策略 |
|------|------|------|---------|
| **内容风险** | | | |
| 工具数据不完整 | 中 | 中 | 建立内容生产流程，人工审核 |
| 工具链接失效 | 高 | 低 | 定期检查，自动提示失效 |
| **技术风险** | | | |
| 搜索性能问题 | 中 | 中 | 使用Elasticsearch，缓存优化 |
| 数据库性能 | 中 | 中 | 索引优化，读写分离 |
| **运营风险** | | | |
| 评测数据质量 | 中 | 中 | 建立审核机制，评分分布验证 |
| 社区氛围差 | 低 | 中 | 社区规则，举报机制 |

### 11.2 应急预案

```
┌─────────────────────────────────────────────────────────────────┐
│                      应急预案                                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  场景 1：搜索性能问题                                            │
│  → 检查索引是否正确                                            │
│  → 调整查询语句                                              │
│  → 启用Redis缓存                                            │
│  → 升级到Elasticsearch（后期）                                │
│                                                                 │
│  场景 2：数据库性能问题                                        │
│  → 分析慢查询日志                                            │
│  → 优化索引                                                    │
│  → 读写分离                                                  │
│  → 数据库扩容                                                │
│                                                                 │
│  场景 3：内容质量不佳                                          │
│  → 建立内容审核机制                                          │
│  → 增加人工审核                                              │
│  → 标记低质量内容                                            │
│  → 优化内容生产流程                                          │
│                                                                 │
│  场景 4：社区运营困难                                          │
│  → 制定社区规则                                              │
│  → 引导优质内容                                              │
│  → 增加社区活动                                              │
│  → 建立激励机制                                              │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

---

## 附录

### 附录 A：工具分类映射表

| 类别ID | 类别名称 | Slug | 排序 | Icon |
|-------|---------|------|------|------|
| 1 | 写作类 | writing | 1 | ✍️ |
| 2 | 图像类 | image | 2 | 🎨 |
| 3 | 视频类 | video | 3 | 🎬 |
| 4 | 音频类 | audio | 4 | 🔊 |
| 5 | 代码类 | code | 5 | 💻 |
| 6 | 办公类 | office | 6 | 📊 |
| 7 | 其他类 | other | 7 | 🔧 |

### 附录 B：推荐工具清单（初期）

#### 写作类
1. ChatGPT - OpenAI
2. Claude - Anthropic
3. Notion AI - Notion
4. Grammarly - Grammarly
5. Jasper - Jasper
6. Copy.ai - Copy.ai

#### 图像类
1. Midjourney - Midjourney
2. Stable Diffusion - Stability AI
3. DALL-E - OpenAI
4. Leonardo AI - Leonardo
5. Canva AI - Canva
6. Bing Image Creator - Microsoft

#### 视频类
1. Runway - Runway
2. Pika - Pika Labs
3. Sora - OpenAI（未发布）
4. HeyGen - HeyGen
5. D-ID - D-ID

#### 音频类
1. ElevenLabs - ElevenLabs
2. Suno - Suno
3. Soundraw - Soundraw
4. AIVA - AIVA

#### 代码类
1. GitHub Copilot - GitHub
2. Cursor - Cursor
3. Tabnine - Tabnine
4. Codeium - Codeium
5. Replit AI - Replit

#### 办公类
1. Microsoft 365 Copilot - Microsoft
2. Google Workspace AI - Google
3. Notion AI - Notion
4. Zapier AI - Zapier

#### 其他类
1. Perplexity AI - Perplexity
2. Otter.ai - Otter
3. Canva AI - Canva
4. Frase.io - Frase
5. Surfer SEO - Surfer

### 附录 C：术语表

| 术语 | 解释 |
|------|------|
| PRD | Product Requirements Document，产品需求文档 |
| P0 | 最高优先级 |
| P1 | 高优先级 |
| P2 | 中优先级 |
| P3 | 低优先级 |
| SSE | Server-Sent Events，服务端推送事件 |
| E2E Testing | End-to-End Testing，端到端测试 |
| RTO | Recovery Time Objective，恢复时间目标 |
| RPO | Recovery Point Objective，恢复点目标 |

### 附录 D：文档变更记录

| 版本 | 日期 | 变更内容 | 作者 |
|------|------|---------|------|
| v1.0 | 2026-04-11 | 初始版本 | 产品团队 |

---

**文档结束**

*本文档将作为工具库模块开发的基准文件，后续变更需同步更新此文档。*
