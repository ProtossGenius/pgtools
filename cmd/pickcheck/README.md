# pick check 工具

检查分支pick 是否有遗漏

原理： 分析git log，比较主分支和子分支的差异

工具依赖一个假定：git log日志中包含`Maniphest Tasks:`(绑定了task)

- 注
其实可以创建分支时把分支名命名为 T0001\_XXX 的格式，这样arc会自动关联task
