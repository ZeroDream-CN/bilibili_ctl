# Bilibili Comments Tool
Bilibili 评论控制工具，帮助你自动清理不想要的评论内容，支持关键字、黑名单以及正则表达式匹配。

## 项目简介
如果你遇到过 B 站的黑名单数量不够用，屏蔽功能不够智能，或者你想屏蔽某些评论但是又不想屏蔽某个人，那么这个工具将可以很方便的帮助你完成这个工作。它通过不断读取评论列表，当发现有评论符合你设定的规则时，就会自动删除该评论，对于拥有大量粉丝的 UP 主来说将是一个非常方便的功能，你不再需要手动去清理那些让人血压飙升的评论，一切都交由这个工具来完成。

## 如何使用
首先在 [Release](https://github.com/ZeroDream-CN/bilibili_ctl/releases) 页面下载最新版本的 bilibili_ctl，根据你的操作系统类型来选择即可，一般用户下载 Windows 版本即可，压缩包内含 32 位和 64 位版本，可以根据你的需要运行。

如果你拥有一台 Linux 的服务器或者群晖 NAS，你可以将 bilibili_ctl 部署到上面，通过 Service 或者 Screen、Nohup 等软件让 bilibili_ctl 在后台静默运行。

<details>
    <summary>Cookie 获取方法</summary>
    <hr>
  
第一次启动软件会提示你输入 cookie，这里推荐使用 Chrome 谷歌浏览器或者其他 Chromium 系浏览器。

打开 [Bilibili 创作中心](https://member.bilibili.com/platform/comment/article)，打开之后按下 F12 打开浏览器控制台，然后转到 “网络” 或者 “Network”，接着刷新一下网页。

![image](https://user-images.githubusercontent.com/34357771/137756642-19f9a28e-0e5c-4820-9327-b6577e128d51.png)

然后点击第一个请求 article，此时右侧会出现请求的详细信息，找到 “请求标头” 或者 “Request Header”，将 “cookie:” 后面的内容复制（也就是截图中红框的部分）

![image](https://user-images.githubusercontent.com/34357771/137757549-273a9b9b-8859-4581-a34f-b8372e9f859a.png)

复制完之后返回到软件，按下右键粘贴即可。

</details>

## 软件配置
下面就开始介绍如何配置 bilibili_ctl，以下是一份示例配置文件

```
{
    "Interval": 30,        // 监控刷新间隔频率，不宜过高，否则可能触发 B 站验证码
    "Block": {
        "Users": [
            "HK赤霊",      // 要屏蔽的用户名或者 uid，每行一条，结尾要加英文逗号 <,>
            "LONEDR"       // 最后一条的结尾不要加逗号
        ],
        "Video": [],       // 要监控的视频 BV 号，格式同上，如果设置为空 [] 则表示监控所有视频
        "Texts": [
            "测试评论内容", // 要屏蔽的关键字内容，格式同上
            "傻逼"
        ],
        "Regex": ""        // 正则表达式匹配关键字模式
    },
    "WhiteList": [],       // 白名单用户名或者 UID 列表，格式同上
    "Request": {
        "Order": "ctime",  // 一些 API 请求参数，通常情况下无需修改
        "Filter": "-1",
        "Type": "1",
        "Page": "1",
        "Size": "20"       // 如果你的视频异常火爆，在一个刷新循环内就会突破 20 条的话，请适当增加此处数值
    }
}
```
__请不要在配置文件中加 `//` 符号以及后面的内容，这里只是为了方便各位理解配置文件的内容而已。__

你可以使用 JSON 在线验证工具来验证你的配置文件是否正确：https://www.bejson.com/

## 开源协议
本软件使用 GPL v3 协议开放源代码，任何人可以在遵循开源协议的情况下对本软件进行修改和使用。
