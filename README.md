#ReSearch--考据级搜索引擎微服务系统<br>  
使用技术golang+goleveldb<br> 

##简介<br>  
ReSearch，从原“乾隆大藏经搜索引擎”和“四库全书搜索引擎”逐渐完善改进而成。<br> 
可用于整理大量的资料并且具备考据级别的搜索功能的系统。<br> <br> 
##天下没有搜索不到的词<br> 
自主研发的遍历分词技术，搜索引擎突破性技术，无需分词库，搜索成功率是100%。<br> 
其他分词技术如果出现新词，第一时间没有办法搜索到结果。<br> 
即使将新词加入分词库，要搜索到该词，又要经历对原数据大量的运算才能搜索到该词的结果。<br> 
遍历分词技术，实时添加实时搜索，没有搜索不到的词。<br>
其他分词技术，都需要维护分词库。<br>分词库就是搜索引擎的眼睛。<br>没有眼睛什么都看不见，什么都搜索不到。<br> 
ReSearch，则摒弃分词库。<br> 
如果分词库是其他搜索引擎的眼睛，这个眼睛是肉眼。<br> 
而遍历分词的眼睛则是天眼。<br> 
其他分词技术会因为分词库的不完善，会导致搜索不到结果的概率。<br>也就是做不到搜索成功率是100%。<br> <br> 
举个极端的例子：<br> 
将一篇文章全部倒过来，然后也用倒过来的词搜索。<br> 
其他分词技术估计什么都搜索不到。<br> 
遍历分词技术，无论如何，搜索成功率都是100%。<br> 

##高精准<br> 
可以自定义搜索粒度。<br> 
通用的搜索引擎如Google，百度，搜索粒度是整篇文章。<br> 
ReSearch，可以自定义到段落，句子等等。<br> 
通常定义精准粒度为句子。<br> 

##高性能<br> 
百亿级数据，毫秒级响应。<br> 
因为无需分词库，不需要解析分词，比其他分词技术的搜索引擎性能更高。<br> 
实时添加实时搜索。<br> 

##部署成本低<br> 
10G级别的文本数据，只需几十M内存。<br> 

##极简部署<br> 
打开对应系统的可执行文件即可运作。<br> 

##案例<br> 
[四库全书搜索引擎](http://www.skqs12.com)
<br> 
<br> 
<br> 
API接口文档<br> 
RESTful 风格，返回信息为json格式。<br>添加、修改、删除仅支持单进程操作。<br> <br> 
目录接口<br> 
包括添加、修改、删除、打开一个目录信息，目录子集信息，目录路径。<br> 
如果项目仅作搜索中介，该目录则是可选，不一定需要。<br> <br> 
添加<br> 
POST /api/cata/?id=&title=&fid&psw=<br> 
id,可选。<br> 
不填写则系统按自动增值id。<br> 
如果是原有系统导入需要保持id一致，则填写。<br> 
title，目录名称。<br> 
fid，上级目录的id。<br> 顶级目录fid=0.<br> 
psw，管理员密码。<br> 
返回成功信息：{ Msg: "提交成功", Succ: true }<br> <br> 
修改<br> 
PUT /api/cata/?id=&title=&fid&psw=<br> 
id,修改的目录id。<br> 
title，目录名称。<br> 
fid，上级目录的id。<br> 顶级目录fid=0.<br> 
psw，管理员密码。<br> 
返回成功信息：{ Msg: "提交成功", Succ: true }<br> <br> 
删除<br> 
DELETE /api/cata/?id=&psw=<br> 
id,删除的目录id。<br> 
psw，管理员密码。<br> 
返回成功信息：{ Msg: "提交成功", Succ: true }<br><br>  
打开一个目录信息<br> 
GET /api/cata/?id=<br> 
id,目录id。<br> 
返回信息示例：{ id: 11, fid: 5, title: "佛经" }<br> 
打开一个目录的所有子目录信息<br> <br> 
GET /api/cata/?fid=<br> 
fid,目录id。<br> 
返回信息示例： [ { id: 3, fid: 1, title: "金刚经" }, { id: 4, fid: 1, title: "六祖坛经" } ]<br><br>  
打开目录路径<br> 
GET /api/cata/?dir=<br> 
dir,目录id。<br> 
返回信息示例： [{"id":21,"title":"机缘品"},{"id":4,"title":"六祖坛经"},{"id":1,"title":"佛经"}]<br><br>  
打开目录下文章<br> 
GET /api/artitem/?id=&p=&count=<br> 
id,目录id。<br> 
p,当前页的定位键值。<br> 第一页为空。<br>第二页开始根据返回结果获得p的定位键值。<br> 
支持无限分页。<br>但是不支持跳页。<br> 
count,每页返回条数。<br> 最大108<br> 
返回信息示例：<br> 
{"ArtItem":[{"id":104,"title":"《实相般若波罗蜜经》"},{"id":3740,"title":"《实相般若波罗蜜经》"}],<br> "p":"105,3740"}<br> 
"id"，文章id。<br> 
"title"，文章标题。<br> 
"p"，定位键值。<br> <br> 
文章接口<br> 
包括添加、修改、删除、打开一篇文章信息，文章某段落信息。<br> <br> 
添加<br> 
POST /api/art/?id=&title=&text=&fid&split=&url=&psw=<br> 
id，文章id,可选。<br> 
不填写则系统按自动增值id。<br> 
如果是原有系统导入需要保持id一致，则填写。<br> 
title，文章标题。<br> 
text，文章内容。<br> 
split，分隔符。<br> 
多个用“|”间隔。<br>
如果不填系统默认按分行符\n分隔。<br> 
中文一般以“。<br>”分隔，就是以每个句子为搜索精度。<br> 
fid，上级目录的id。<br>搜索时可以按目录范围搜索。<br> 可选。<br> 
url， 文章来源url，可选。<br> 
psw，管理员密码。<br> 
返回成功信息：{ Msg: "提交成功", Succ: true }<br> <br> 
修改<br> 
PUT /api/art/?id=&title=&text=&fid&split=&url=&psw=<br> 
id,文章id。<br> 
不填写则系统按自动增值id。<br> 
如果是原有系统导入需要保持id一致，则填写。<br> 
title，文章标题。<br> 
text，文章内容。<br> 
split，分隔符。<br> 
url， 可选。<br> 
fid，上级目录的id。<br>搜索时可以按目录范围搜索。<br> 可选。<br> 
psw，管理员密码。<br> 
返回成功信息：{ Msg: "提交成功", Succ: true }<br> <br> 
删除<br> 
DELETE /api/art/?id=&psw=<br> 
id,文章id。<br> 
psw，管理员密码。<br> 
返回成功信息：{ Msg: "提交成功", Succ: true }<br> <br> 
打开一个文章信息<br> 
GET /api/art/?id=<br> 
id,文章id。<br> 
返回信息，参考【添加】接口的字段说明。<br> 
打开一个文章某段落信息<br> 
GET /api/art/?id=&secid=<br> 
id,文章id。<br> 
secid,文章分割后的段id。<br>搜索时返回结果即是文章id和段落id。<br> 
返回信息：该段落的内容。<br> <br> 
搜索接口<br> 
/api/search/?kw=&p=&count=&caids=&order=<br> 
kw,搜索词。<br> 
搜索是按精准前缀匹配。<br> 
如果要模糊搜索，则使用空格将两个或多个词隔开。<br>空格等于是省略号。<br>
模糊搜索，所有的词同在一个段落内，则是符合条件。<br>
p,当前页的定位键值。<br> 第一页为空。<br>第二页开始根据返回结果获得p的定位键值。<br>
搜索结果支持无限分页。<br>但是不支持跳页。<br>
count,每页返回条数。<br>
caids,目录id，在某个目录下搜索。<br>支持多个，用"|"分隔。<br> 可选。<br>
order,排序。0，是升序；1，降序。比如最新的文章排在前面，则order=1<br>
返回信息:<br>
{"Rset":<br>
[{"CataDir":[{"id":5301,"fid":247,"title":"大般泥洹经卷第二"},{"id":247,"fid":8,"title":"大般泥洹经"},{"id":8,"fid":1,"title":"大乘涅槃部"},{"id":1,"fid":0,"title":"经"}], //目录路径<br>
"Artid":5301, //文章id<br>
"Title":"《哀叹品第四》", //文章标题<br>
"ArtUrl":"",//文章网址<br>
"Secid":94,//文章段落id<br>
"LastSecid":99,//文章最后段落id。<br>搜索结果节录默认最少49字。<br>
如果不足，会顺序向下读取段落，直到大于等于49。<br>最后一个段落id即是LastSecid。<br>
"Text":"时诸比丘白佛言：“世尊，我当云何如世尊教......\n"},//搜索结果节录<br>
......."}],<br>
"p":"如来？我等于瞿,3948,10", //下一页的定位键值。<br>
"SeTime":"3.6088ms", //搜索用时。<br>
"SetTime":"11.436ms"//数据集结用时。<br>
} 