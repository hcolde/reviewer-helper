# reviewer-helper

该项目是[reviewer](https://github.com/hcolde/reviewer)的辅助程序，目的是减轻reviewer的压力，以下将此项目简称为helper



> 功能

- [x] 转账
- [ ] 中文分词切割
- [ ] 垃圾文本识别

* `转账`：[reviewer](https://github.com/hcolde/reviewer)将转账信息存放至队列，helper从中取出持久化至数据库；
* `中文分词切割`：将使用[HanLP](https://github.com/hankcs/HanLP)实现，由helper远程调用返回结果，供用户实时添加敏感词；
* `垃圾文本识别`：任重道远。



> 安装

wait...