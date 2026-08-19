# Changelog

## [0.1.1](https://github.com/yuyudeqiu/chronicle/compare/v0.1.0...v0.1.1) (2026-08-19)


### Bug Fixes

* publish standard semantic version tags ([#45](https://github.com/yuyudeqiu/chronicle/issues/45)) ([70e6976](https://github.com/yuyudeqiu/chronicle/commit/70e69762155ccae5e9b0d0952b36acf3199f248d))
* reject invalid task statuses ([#47](https://github.com/yuyudeqiu/chronicle/issues/47)) ([9224eb3](https://github.com/yuyudeqiu/chronicle/commit/9224eb3e1c5987bd356f25d1f3a3818372635ff0))

## 0.1.0 (2026-07-21)


### Features

* add archive feature ([ea9b7e9](https://github.com/yuyudeqiu/chronicle/commit/ea9b7e9f1eac33cca6f5c696a175058efc5f7d97))
* add archive feature with archived_at field ([ca03d40](https://github.com/yuyudeqiu/chronicle/commit/ca03d4048635e8b39a24b300babe15cd01b58627))
* add blocked task status ([#38](https://github.com/yuyudeqiu/chronicle/issues/38)) ([afccc3b](https://github.com/yuyudeqiu/chronicle/commit/afccc3b5cdf74d410596c644cd3b28860d9370ed))
* add chronicle version command with dual source support ([#36](https://github.com/yuyudeqiu/chronicle/issues/36)) ([cbe0183](https://github.com/yuyudeqiu/chronicle/commit/cbe0183a94e9e14d772516ec04a4d6ab2b0c9ac4))
* add export/import commands, remove unused exporter ([#35](https://github.com/yuyudeqiu/chronicle/issues/35)) ([5b67fd2](https://github.com/yuyudeqiu/chronicle/commit/5b67fd245294f8a8e28f0fe73f90fbe213df27f2))
* add start task button to move task to in-progress state ([ed05b09](https://github.com/yuyudeqiu/chronicle/commit/ed05b098e03f846d5fd3f512843da12697a499fd))
* add task deadline support and update obsidian export to zip template ([b4955ee](https://github.com/yuyudeqiu/chronicle/commit/b4955eec6038d78f6afc2d525cad2fa87dbc11ad))
* add weekly-summary command and API endpoint ([#34](https://github.com/yuyudeqiu/chronicle/issues/34)) ([af42f06](https://github.com/yuyudeqiu/chronicle/commit/af42f0634a89f4838a0fdcc34168956d19bd8b51))
* automate versioned releases ([#42](https://github.com/yuyudeqiu/chronicle/issues/42)) ([445740b](https://github.com/yuyudeqiu/chronicle/commit/445740bbbfa0ce453a8c3aa01cc3a6c2c4ef5562))
* CLI工具 + 补充必要字段 ([57eab13](https://github.com/yuyudeqiu/chronicle/commit/57eab136e1b442a0528e04eb219ebfb308cdc790))
* CLI工具 + 补充必要字段 ([5fc9925](https://github.com/yuyudeqiu/chronicle/commit/5fc9925ae7090046315e337d75cf0d73444bc10a))
* daily-summary API增加category字段 ([32791fa](https://github.com/yuyudeqiu/chronicle/commit/32791fa29abe1c323bd5bdd2bc09c8875aa6b67c))
* **exporter:** bind Links field in obsidian task template ([f9f9731](https://github.com/yuyudeqiu/chronicle/commit/f9f9731c653519c83d3ecd26d5db1ae138ada4ae))
* sort tasks by deadline ([47f29c2](https://github.com/yuyudeqiu/chronicle/commit/47f29c23d4239fa4476ea7f626bd135bd2b62857))
* support updating deadline when saving task progress & fix timezone ([6600352](https://github.com/yuyudeqiu/chronicle/commit/660035280a7e570c5a229c30504e4983ea1b0e34))
* support updating task title ([#41](https://github.com/yuyudeqiu/chronicle/issues/41)) ([a22e69f](https://github.com/yuyudeqiu/chronicle/commit/a22e69fc0898ebdef43917b06cc248b85ca050f7))
* **ui:** add custom confirm modal for delete confirmation ([c2c256a](https://github.com/yuyudeqiu/chronicle/commit/c2c256a071e748c1d74d991620a94edd582d7134))
* **ui:** add custom confirm modal for delete confirmation ([434f25b](https://github.com/yuyudeqiu/chronicle/commit/434f25b01e1e586f4a8ea655c1c92cdc6bebb727))
* **ui:** add stats bar to main page ([f73f7ff](https://github.com/yuyudeqiu/chronicle/commit/f73f7ff3597dfb9cd14dda40f15d6d0739d3a6d5))
* **ui:** add stats bar to main page ([df47ccc](https://github.com/yuyudeqiu/chronicle/commit/df47ccc09b72e843435241b23c080bebcef9e507))
* **ui:** change weekly stats to line chart for better trend visualization ([64305b2](https://github.com/yuyudeqiu/chronicle/commit/64305b2aad21f2f66e12233e2308a6d8be88c869))
* 任务支持关联文档链接 ([4088d62](https://github.com/yuyudeqiu/chronicle/commit/4088d62adede07fc936e42b558c2e97684cbcc03))
* 创建任务支持 deadline 参数 ([1ba8359](https://github.com/yuyudeqiu/chronicle/commit/1ba8359cc310a710a5d7bffd95e57b5cc3414c82))
* 创建任务支持 deadline 参数 ([acc85b2](https://github.com/yuyudeqiu/chronicle/commit/acc85b2c98c51e1e61c83dbc13fc1dcead9e1ac7))
* 支持 go install 安装和数据路径配置化 ([#29](https://github.com/yuyudeqiu/chronicle/issues/29)) ([a4becff](https://github.com/yuyudeqiu/chronicle/commit/a4becff995614bccfe510f3236413cdec48bd9c4))
* 更新任务支持修改 deadline ([caf4972](https://github.com/yuyudeqiu/chronicle/commit/caf49729f9827afcd69b9bdaa9429756c517102c))
* 更新任务支持修改 deadline ([9f2c4bf](https://github.com/yuyudeqiu/chronicle/commit/9f2c4bf50813a20d627497148454400701e4eed5))
* 替换为纯Go SQLite驱动 (glebarez/sqlite)，移除CGO依赖 ([7f644c7](https://github.com/yuyudeqiu/chronicle/commit/7f644c7e171e190e28ae8d9ff464776ad2eb6bed))
* 替换为纯Go SQLite驱动，移除CGO依赖 ([84d3499](https://github.com/yuyudeqiu/chronicle/commit/84d34991f30aa55e6c1d5548708f7ddbecf51750))
* 构建时注入git信息,新增version API ([#32](https://github.com/yuyudeqiu/chronicle/issues/32)) ([86faac7](https://github.com/yuyudeqiu/chronicle/commit/86faac78fe26b67d2e007052c451e17bad0ac8c3))
* 添加Docker支持 ([c471272](https://github.com/yuyudeqiu/chronicle/commit/c471272c57d49e0b0b2c5d2b63978ca2ca83d6ce))
* 添加Docker支持 ([e408642](https://github.com/yuyudeqiu/chronicle/commit/e408642dfd5b9d7b34de064df3c2ab66580c3b7f))
* 添加PATCH /tasks/:id编辑任务接口 ([1f498fb](https://github.com/yuyudeqiu/chronicle/commit/1f498fb13aba83e6aec8f16032a236f8bc7a6b1f))
* 添加删除worklog接口及前端功能 ([3c1545f](https://github.com/yuyudeqiu/chronicle/commit/3c1545f06e030f69f4492c7d0df0a0ba98bbaf3c))
* 添加删除worklog接口及前端功能 ([ddffc06](https://github.com/yuyudeqiu/chronicle/commit/ddffc069cafad3d68890fccca498077485d098ba))
* 添加统计API GET /stats/summary ([6d16e5b](https://github.com/yuyudeqiu/chronicle/commit/6d16e5bd39e5bb66ad8e132fbd474f045fb7f555))
* 页面增加编辑任务功能 ([5416cf8](https://github.com/yuyudeqiu/chronicle/commit/5416cf8e2d117bf1976dce3b0948336b5f6b36b5))
* 页面增加编辑任务功能 ([2efb954](https://github.com/yuyudeqiu/chronicle/commit/2efb9543dd2267800984b954c2fde4b7d0962771))


### Bug Fixes

* add worklog delete confirmation modal ([4b10941](https://github.com/yuyudeqiu/chronicle/commit/4b109411afdf70c41dce2775eb62d264f8414579))
* check error return values to satisfy errcheck ([#24](https://github.com/yuyudeqiu/chronicle/issues/24)) ([f3a4f76](https://github.com/yuyudeqiu/chronicle/commit/f3a4f766927313bd9e47b0bbdc4ad3d884b762d5))
* only create worklog when LogText is not empty ([#39](https://github.com/yuyudeqiu/chronicle/issues/39)) ([2b94a47](https://github.com/yuyudeqiu/chronicle/commit/2b94a4747c50c12104a9ac1f1f16342c9c3fac30))
* start releases at v0.1.0 ([#44](https://github.com/yuyudeqiu/chronicle/issues/44)) ([ac90ab4](https://github.com/yuyudeqiu/chronicle/commit/ac90ab4f1ad5ac2b127308d1745f74b14ae52eef))
* 修复 ProgressModal 更新任务后未刷新列表的问题 ([#28](https://github.com/yuyudeqiu/chronicle/issues/28)) ([ca43b99](https://github.com/yuyudeqiu/chronicle/commit/ca43b999d45b269208f22fd0b8dc2b74cbbc53f4))
* 移除 replace 指令以支持 go install ([#30](https://github.com/yuyudeqiu/chronicle/issues/30)) ([440f2ba](https://github.com/yuyudeqiu/chronicle/commit/440f2ba7477e1cee874dfc6a10904d640007b6a4))


### Reverts

* 移除不相关的 main.go 改动 ([ece82ac](https://github.com/yuyudeqiu/chronicle/commit/ece82acbe21f9ba276e84ee0153e472846aff55b))
