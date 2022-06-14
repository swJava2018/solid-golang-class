# 이벤트 데이터 파이프라인 프로젝트
## 개발 가이드 
### 로컬 환경에서 빌드하기
제공된 `setup.sh` 쉘 스크립트 실행
```sh
$ ./scripts/setup.sh
```

### 로컬 환경에서 실행하기
빌드된 `event-data-pipeline` 바이너리 파일 실행

Debug 모드로 실행하기 예시
```sh
$ ./bin/event-data-pipeline --debug --config configs/config.json
```