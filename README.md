# go-logfilter

## if
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
