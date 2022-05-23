# RabbitMQ 설치하고 구동하는 법
## Homebrew 
### 홈브루 업데이트
```
$ brew update
```
### RabbitMQ 설치
```
$ brew install rabbitmq
```

### RabbitMQ 구동
```
$ brew services start rabbitmq
==> Successfully started `rabbitmq` (label: homebrew.mxcl.rabbitmq)
```

### RabbitMQ GUI 접속
브라우저를 열고 http://localhost:15672 로 접속합니다.
기본 로그인 `ID/Password`는 `guest/guest` 입니다.

[공식 가이드 문서 원문 링크](https://www.rabbitmq.com/install-homebrew.html)