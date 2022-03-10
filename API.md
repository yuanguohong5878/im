

#### user：
|type|API|data|result|备注|
|:--:|:--:|:--:|:--:|:--------:|
|POST|xxx.xx.xx.xx/api/user/register|{"name":昵称,"email":邮箱，"password":密码}|{Code:200,Type:success}||
|POST|xxx.xx.xx.xx/api/user/login|{"name":邮箱，"password":密码}|{Code:200,Type:success，Message:"login success"}||
|GET|xxx.xx.xx.xx/api/user/id/:id|null|{Code:200,Type:success，Message:user.Info}||
|GET|xxx.xx.xx.xx/api/user/connect|null|null||
|GET|xxx.xx.xx.xx/api/user/self|null|{Code:200,Type:success，Message:user}||
|PUT|xxx.xx.xx.xx/api/user/message||{Code:200,Type:success}||
|POST|xxx.xx.xx.xx/api/user/addfriend|{"email":"好友的邮箱"}|{Code:200,Type:success}||
|GET|xxx.xx.xx.xx/api/user/friends|null|{Code:200,Type:success,MEssage:friends}||

### moment
|type|API|data|result|
|:--:|:--:|:--:|:--:|
|POST| xxx.xx.xx.xx/api/moment/publish |{"image":"朋友圈配图","content":"朋友圈文案"}|{Code:200,Type:success}|
|POST|xxx.xx.xx.xx/api/moment/comment|{"momentid":"朋友圈id","content":"评论内容"}|{Code:200,Type:success}|
|GET|xxx.xx.xx.xx/api/moment/all|null|{Code:200,Type:success，Message:momentsres}|
|GET|xxx.xx.xx.xx/api/moment/all/:fid|null|{Code:200,Type:success，Message:momentsres}|
|POST|xxx.xx.xx.xx/api/moment/like|{“momentid”：“朋友圈id”}|{Code:200,Type:success}|

### group

|  type  |                  API                   |                    data                     |                  result                   |      |
| :----: | :------------------------------------: | :-----------------------------------------: | :---------------------------------------: | ---- |
|  POST  |     xxx.xx.xx.xx/api/group/create      |              {"name":"群名称"}              |          {Code:200,Type:success}          |      |
|  GET   |   xxx.xx.xx.xx/api/group/add/:number   |                    null                     |          {Code:200,Type:success}          |      |
|  POST  |  xxx.xx.xx.xx/api/group/announcement   | {"number":"群号","announcement":"公告内容"} |          {Code:200,Type:success}          |      |
|  POST  |     xxx.xx.xx.xx/api/group/delete      |   {"number":"群号","memberid":"群成员id"}   |          {Code:200,Type:success}          |      |
| DELETE | xxx.xx.xx.xx/api/group/delete/:groupid |                    null                     |          {Code:200,Type:success}          |      |
|  GET   |       xxx.xx.xx.xx/api/group/all       |                    null                     | {Code:200,Type:success，Message:groupres} |      |
|  GET   |   xxx.xx.xx.xx/api/group/id/:number    |                    null                     | {Code:200,Type:success，Message:groupres} |      |

