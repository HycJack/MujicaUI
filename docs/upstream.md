# 上游归并与上游缺口

## 归并（PRD 4.11）

MyGo 已有能力纳入下列组件，不增加计数。

| MyGo 能力 | MujicaUI 组件 | 文件 |
| --- | --- | --- |
| Primary button、Menu button | Button 主按钮变体；DropdownMenu | `input/button.go`、`overlay/dropdown_menu.go` |
| Stepper、Number input | NumberInput 仅步进模式与输入模式 | `input/number_input.go` |
| Search field、Find bar | SearchInput 搜索模式与查找栏模式 | `input/search_input.go` |
| Autocomplete、Combobox | Combobox 自由输入与选项选择 | `input/combobox.go` |
| Token field、Editable text | TagsInput；Text 的 `EditableText` | `input/tags_input.go`、`display/text.go` |
| List、Grid view | List 列表模式与网格模式 | `data/list.go` |
| Table、Outline、Tree | DataTable；TreeTable；Tree | `data/data_table.go`、`data/tree_table.go`、`data/tree.go` |
| Date input、Time input | DatePicker；TimePicker | `datetime/date_picker.go`、`datetime/time_picker.go` |
| Back and forward buttons | NavigationMenu 历史导航 | `navigation/navigation_menu.go` |
| Meter、Progress bar | Progress 计量模式与进度模式 | `feedback/progress.go` |

## 上游缺口（待提交 MyGo）

PRD 第 2 节要求公开接口不足时向上游提交扩展。以下缺口在 MyGo v0.2.10 中核实，组件文档的“限制”节写明影响。本仓库不访问 MyGo internal，也未修改依赖。

| 缺口 | 位置 | 受影响组件 |
| --- | --- | --- |
| 自定义元素无法设置 expanded、mixed、数值范围与值 | `ui/element.go:282-295`、`ui/base.go:128,297` | Checkbox 半选、MultiSelect、TreeSelect、RangeSlider、Rating、全部日期时间弹出按钮 |
| 无 menu、menuitem 角色 | `ui/access.go:22-79` | Menubar、DropdownMenu、ContextMenu |
| 无 heading 角色 | `ui/access.go` | Heading |
| 输入法组合状态不公开 | `ui/editor.go:73`、`ui/input.go:689` | TextInput、TextArea、SearchInput、Dialog、AlertDialog |
| 打开链接失败不回传 | `ui/context.go:270`、`ui/window.go:583` | Link |
| 浮层焦点恢复与浮层快捷键不公开 | `ui/scope.go:109,123,137` | Dialog、Drawer、Overlay |
| `PopoverBase` 背板吞掉外部点击；只有上下两个位置 | `ui/base.go:373`、`ui/widgets.go:689` | Popover、Popconfirm、DropdownMenu、HoverCard |
| `activeDescendant` 不公开 | `ui/element.go:331` | Menubar、CommandPalette |
| `SliderBase` 方向键步长固定 1%，构建时处理指针 | `ui/base.go:111-117` | Slider |
| `SelectBase` 高亮无设置接口 | `ui/base.go:260` | Select |
| `TextInputBase` 无只读开关 | `ui/editor.go:139,811` | TextInput 只读模式 |
| `ui.Stepper` 禁用时仍可步进，固定 9 位取整 | `ui/indicators.go:216`、`ui/widgets.go` `snap` | NumberInput（已改为自绘） |
| 路径遮罩缓存按四分之一像素取键，进程内共享 | `ui/path.go:189-210` | 自绘路径以 `internal.Snap` 对齐；MyGo 自绘的箭头（Outline 等）无法对齐；快照测试因此在独立子进程中渲染，结果与测试顺序无关 |

## mygo v0.3 迁移中发现的缺口（待提交 MyGo）

在迁移到 mygo v0.3.4（值类型 Element、延迟输入）时核实：

| 缺口 | 位置 | 影响 |
| --- | --- | --- |
| 模态 overlay 打开期间下层视图不构建，opener 的状态被 prune 剪除；关闭时 `restoreFocus` 因 `rt.states[opener] != nil` 不成立而失效，即使 opener 的 id 仍在 `rt.focused` 上 | `ui/scope.go:141-152`、`ui/runtime.go:591-636` | AlertDialog、ToolApprovalDialog 关闭后焦点落在 OverlaySite 的隐形占位上，不回到触发按钮（焦点不丢，Tab 可继续，但焦点环消失）。引擎侧需让 openers 记录的 opener 在模态期间 keep-alive，或恢复时按构建内容重挂 |
| 输入分发期（HandleInput 回调）与绘制期（Draw/DrawOver 回调）元素句柄过期并 panic，剪贴板/焦点/查询 API 静默失效；组件必须在构建期捕获值或改用 `ui.Handle` 与 `Context.Services()` | `ui/frame.go:96-107`、`ui/api_gen.go` | 全部自定义控件的输入/绘制回调；已在 fork 内完成适配 |
