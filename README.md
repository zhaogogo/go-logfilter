# go-logfilter

# 字段格式约定

## key

支持两种方式

1.   字面量显示 
     1.   `"a_key": "a_value"`
2.   多级添加
     1.   `"[x1][x2]": [1,2,3]`

```yaml
inputs:
  - stdin:
      prometheus_counter:
        namespace: ""
        subsystem: ""
        name: "logfilter_stdin"
      add_fields:  # 使用数组保证顺序
        # key添加
        - "[x1][x2]": [1,2,3]               # 多级添加,value为interface支持golang多种数据类型
          "a_key": "a_value"                # 字面量添加
#        # value
        - "jsonpath_test1": "$.message"
        - "jsonpath_test2": "$.messagexxx"   # 如果jsonpath返回值为空，不添加value
        - "[][]test1": "[x1][x2]"
        - "[][]test2": "[a1][b1]"           # 如果[a1][b1]返回值为空，不添加value
        - "[c1][c2]": "[a2][b2]"
        - "go_template1": "{{ timeFormat .timestamp }}"   # 使用自定义函数
        - "go_template2": "{{ timeFormat .timestamp \"2006-01-02\"}}"   # 使用自定义函数
        - "es_key": "aa%{x1}{x2}-%{+2006.01.02}"
      overwrite: true  # add_fields 是否覆盖
      failed_tag: true # 有执行失败的添加 failed_tag 字段
      codec: json
      timestamp: timestamp  # timestamp_key值
```

输出

```json
{
    "a_key": "a_value",
    "message": "a",
    "failed_tag":["$.messagexxx","[a1][b1]","[a2][b2]","{{ timeFormat .timestamp }}"]
    "timestamp": "2024-10-26T17:12:41.599093+08:00",
    "x1": {
        "x2": [
            1,
            2,
            3
        ]
    }
}
```

## value

1.   jsonpath渲染value
     1.   如果jsonpath返回值为空，不添加value
2.   \[a1\]\[b2\] 渲染value
     1.   如果a1.b2值不存在，不添加value
3.   {{XXX}} golang模版字符串
     1.   执行失败，不添加value
4.   %{XXX}{YYY}
     1.   执行失败，添加字面量值

# 添加字段

inputs、filters、outputs块都可配置add_fields字段

# 条件判断

if条件判断仅支持`filters plugin`块和`outputs plugin`块。如下示例
- template 条件 -->         {{if .name}}y{{end}}
- 自己实现的一套简单的DSL --> Exist(a) && (!Exist(b) || !Exist(c))
**注意**: if 数组中的条件是 AND 关系, 需要全部满足.
```yaml
filters:
  - Drop:
      if:
        - '{{if .name}}y{{end}}'
        - '{{if eq .name "childe"}}y{{end}}'
        - '{{if or (before . "-24h") (after . "24h")}}y{{end}}'
  - Drop:
      if:
        - 'EQ(name,"childe")'
        - 'Before(-24h) || After(24h)'
outputs:
  Stdout: 
    codec: json
    if:
      - '{{if .name}}y{{end}}'
      - '{{if eq .name "childe"}}y{{end}}'
      - '{{if or (before . "-24h") (after . "24h")}}y{{end}}'
```
## 自己实现的一套简单的DSL
目前 if 支持两种语法, 一种是 golang 自带的 template 语法, 一种是我自己实现的一套简单的DSL, 实现的常用的一些功能, 性能远超 template , 我把上面的语法按自己的DSL翻译一下.

```yaml
  Drop:
    if:
    - 'EQ(name,"childe")'
    - 'Before(-24h) || After(24h)'
```

也支持括号和逻辑运算符, 像 Exist(a) && (!Exist(b) || !Exist(c))

目前支持的函数如下:

注意: EQ/IN 函数需要使用双引号代表字符串, 因为他们也可能做数字的比较, 其他所有函数都不需要双引号, 因为他们肯定是字符串函数

EQ IN HasPrefix HasSuffix Contains Match , 这几个函数可以使用 jsonpath 表示, 除 EQ/IN 外需要使用双引号

```text
Exist(user,name) [user][name]存在
EQ(user,age,20) EQ($.user.age,20) [user][age]存在并等于20
EQ(user,age,"20") EQ($.user.age,"20") [user][age]存在并等于"20" (字符串)
IN(tags,"app") IN($.tags,"app") "app"存在于 tags 数组中, tags 一定要是数组,否则认为条件不成立
HasPrefix(user,name,liu) HasPrefix($.user.name,"liu") [user][name]存在并以 liu 开头
HasSuffix(user,name,jia) HasSuffix($.user.name,"jia") [user][name]存在并以 jia 结尾
Contains(user,name,jia) Contains($.user.name,"jia") [user][name]存在并包含 jia
Match(user,name,^liu.*a$) Match($.user.name,"^liu.*a$") [user][name]存在并能匹配正则 ^liu.*a$
Random(20) 1/20 的概率返回 true
Before(24h) @timestamp 字段存在, 并且是 time.Time 类型, 并且在当前时间+24小时之前
After(-24h) @timestamp 字段存在, 并且是 time.Time 类型, 并且在当前时间-24小时之后
```

# Filter

## Filter插件通用配置

### if条件判断

详细见条件判断

```yaml
filters:
  Drop:
    if:
    - 'EQ(name,"childe")'
    - 'Before(-24h) || After(24h)'
```

### add_fields

当Filter执行成功时, 可以添加一些字段. 如果Filter失败, 则忽略. 下面具体的Filter说明中, 提到的"返回false", 就是指Filter失败

```yaml
filters:
- Grok:
    src: message
    match:
    - '^(?P<logtime>\S+) (?P<name>\w+) (?P<status>\d+)$'
    - '^(?P<logtime>\S+) (?P<status>\d+) (?P<loglevel>\w+)$'
    remove_fields: ['message']
    add_fields:
    - "es_key": "aa%{x1}{x2}-%{+2006.01.02}"
    overwrite: true  # add_fields 是否覆盖
    failed_tag: true # 有执行失败的添加 failed_tag 字段
```

### remove_fields

例子如上. 当Filter执行成功时, 可以删除一些字段. 如果Filter失败, 则忽略.

```yaml
      delete_fields:
        - "jsonpath_test1"
        - "[][]test1"
        - "a"
        - a_key
        - abc
        - "es_key"
        - go_template2
        - "x1"
```

## cover

```yaml
convert:
    fields:
        time_taken:
            remove_if_fail: false
            setto_if_nil: 0.0
            setto_if_fail: 0.0
            to: float
        sc_bytes:
	          failed_tag: true
            to: int
            remove_if_fail: true
        status:
						to: bool
            remove_if_fail: false
            setto_if_fail: true
        map_struct:
            to: string
            setto_if_fail: ""
```

-   remove_if_fail: 如果转换失败刚删除这个字段, 默认 false

-   setto_if_fail: XX: 如果转换失败, 刚将此字段的值设置为 XX . 优先级比 remove_if_fail 低. 如果 remove_if_fail 设置为 true, 则setto_if_fail 无效.

-   setto_if_nil: XX: 如果没有这个字段, 刚将此字段的值设置为 XX . 优先级最高

## grok

```yaml
  - grok:
      delete_fields: ['message']
      match:
        - '^(?P<logtime>\S+) (?P<name>\w+) (?P<status>\d+)$'
        - '^(?P<logtime>\S+) (?P<status>\d+) (?P<loglevel>\w+)$'
        - '(?m)^(?P<first>\w+) (?P<last>\w+)$'
      ignore_blank: true
      src: message
      overwrite: false
      pattern_paths:
        - 'https://raw.githubusercontent.com/vjeantet/grok/master/patterns/grok-patterns'
        - '/opt/gohangout/patterns/'
```

源字段不存在, 返回 false. 所有格式不匹配, 返回 false

-   src: 源字段, 默认 message

-   target: 目标字段, 默认为空, 直接写入根下. 如果不为空, 则创建target字段, 并把解析后的字段写到target下.

-   match: 依次匹配, 直到有一个成功.

-   pattern_paths: 会加载定义的 patterns 文件. 如果是目录会加载目录下的所有文件.
    -   这里推荐 https://github.com/vjeantet/grok 项目, 里面把 logstash 中使用的 pattern 都翻译成了 golang 的正则库可以使用的.

-   ignore_blank: 默认 true. 如果匹配到的字段为空字符串, 则忽略这个字段. 如果 ignore_blank: false , 则添加此字段, 其值为空字符串.

## filter

目的是为了一个 if 条件后跟多个Filter

```yaml
Filters:
    if:
        - '{{if eq .name "childe"}}y{{end}}'
    filters:
        - Add:
            fields:
                a: 'xyZ'
        - Lowercase:
            fields: ['url', 'domain']
```

