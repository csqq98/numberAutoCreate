# numberAutoCreate
Numbers are automatically generated from expressions
## 跟据表达式生成对应规则的不重复的编号
### 1.表达式规则：
eg:p:1-p:2-d:yyMMdd-i:2:3/p:1-p:2-d:yyMMdd-i:d:1
表达式分为三个部分： p类型-d时间格式-i末尾编号规则
#### p:类型
类型可以多个 eg p:1-p:2-p:3 表示有三种类型,数字为对应类型,起始为1（暂支持最多9个分类）,
#### d:日期格式（*yyMMdd/yyyyMMdd*两种）
暂时只有yyMMdd/yyyyMMdd 两种格式,生成后的编号部分分别为 210819/20210819
#### i:末尾规则
i 后紧跟的数字/d 代表生成规则跟据第几个类型或者根据日期部分的日期生成 eg i:2  则生成的编号与p2的相关（不同模块的表达式有个需要维护的map,详细看代码）  eg i:d 则生成的编号与yyMMdd/yyyyMMdd相关
最后数字的数字表示生成几位编号:eg p:1-p:2-d:yyMMdd-i:2:3  则末尾生成编号为"001".

### 2.需要做的事情
需要本地实现代码里面的save()和Data()两个方法
Data()用来获取存库的数据
Save()用来更新维护存库的数据

### 3.用到的方法
Gen() // 用来生成编号的方法,
IfAccordNumberRule()  // 判断表达式是否符合规则
