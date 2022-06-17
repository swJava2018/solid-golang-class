# 이벤트 데이터 파이프라인 프로젝트
## 빌드 가이드 
### 로컬 환경
* 개발자 PC 환경의 OS 플랫폼에서만 빌드하고 실행합니다.
#### 빌드
제공된 `setup.sh` 쉘 스크립트 실행
```sh
$ ./scripts/setup.sh
```

#### 실행
빌드된 `event-data-pipeline` 바이너리 파일 실행

Debug 모드로 실행하기 예시
```sh
$ ./bin/event-data-pipeline --debug --config configs/config.json
```
--- 
