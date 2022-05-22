# Mock 데이터를 생성하고 테스트 하기 위한 서비스를 띄웁니다. 
```
docker-compose up -d
```
# 컨테이너 구동상태를 확인합니다.
```
docker-compose ps
```

# 컨트롤 센터 컨테이너 프로세스가 제대로 뜨지 않을 경우 재시작 합니다.
```
docker-compose restart control-center
```

## Mock 데이터 생성을 위한 pageviews topic 을 생성합니다. 

### 브라우저에서 http://localhost:9021 를 입력하여 Control Center 를 접속합니다.
![controlcenter.cluster 를 클릭합니다.](./controlcenter.cluster.png)

### Topics 를 클릭하여 목록을 엽니다. Add topic 을 클릭 후 pageviews 토픽을 생성합니다.
![pageviews 토픽을 생성합니다.](./create.pageviews.topic.png)
![add topic 클릭.](./new.topic.png)


### Topics 를 클릭하여 목록을 엽니다. Add topic 을 클릭 후 pageviews 토픽을 생성합니다.