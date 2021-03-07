# go-cas

[英文](https://github.com/xczh/go-cas/blob/main/README.md)

Golang实现的CAS协议Server和Client。

## CAS Client

`cas.Client`允许Go应用程序通过中央身份认证服务（CAS）服务器轻松地对用户进行身份认证。

其他类似项目：

  - [go-cas](https://github.com/go-cas/cas)

### 特点

  - 严格按照[CAS协议](https://apereo.github.io/cas/development/protocol/CAS-Protocol-Specification.html)文本的标准实现
  - 支持CAS协议`v1`、`v2`、`v3`版本
  - 支持代理认证
  - 不支持`SAML`协议
  - 高性能
  - 良好的可扩展性

已支持的客户端特性：

  - [ ] CAS 1.0
  - [ ] CAS 2.0
  - [ ] CAS 3.0
  - [ ] Gateway Authentication
  - [ ] Single Sign-Out
  - [ ] Local Application Configuration

## CAS Server

TODO
