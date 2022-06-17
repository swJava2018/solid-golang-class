# solid-golang-class
---
## Github Action 를 이용한 컨테이너 이미지 빌드
--- 
### Step1. Workflow 설정파일 편집
solid-golang-class/.github/workflows/deploy-app.yml

`<BRANCH_NAME>` 과 `<DOCKERHUB_USERNAME>`를 수정
```yaml
name: Build and Push Docker Image
on:
  push:
    branches:
      - <BRANCH_NAME>
jobs:
  build-and-push-image:
    runs-on: ubuntu-latest
    defaults:
      run:
        working-directory: ./app
    steps:
    - name: Checkout
      uses: actions/checkout@v2

    - name: Set up Docker Buildx
      uses: docker/setup-buildx-action@v1

    - name: Login to DockerHub
      uses: docker/login-action@v1
      with:
        username: ${{ secrets.DOCKERHUB_USERNAME }}
        password: ${{ secrets.DOCKERHUB_TOKEN }}

    - name: Build and push
      id: docker_build
      uses: docker/build-push-action@v2
      with:
        context: app
        dockerfile: ./app/Dockerfile
        push: true
        tags: <DOCKERHUB_USERNAME>/event-data-pipeline:latest
```
--- 
### Step2. Github Secrets 설정
#### Github Repository > Settings > Secrets > Actions 메뉴에서 New repository secret 클릭 후 각각 생성
`DOCKERHUB_USERNAME` : [Dockerhub](https://hub.docker.com/) 사용자 계정명
`DOCKERHUB_TOKEN` : [Dockerhub](https://hub.docker.com/) > Account Settings > Security > New Access Token 생성

--- 
### Step3. 해당 브랜치 원격지에 Push 반영
```sh
$ git push 
```

### Step4. Github Repository > Actions 에서 확인
