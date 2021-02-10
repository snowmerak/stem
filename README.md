# stem

## launch

`multipass launch`와 같습니다.

`stem launch <dir>`를 실행하여 해당 경로 내 `template.json`을 읽어서 템플릿에 맞춰서 인스턴스를 만듭니다.

만약 경로 내에 `template.json`이 존재하지 않는다면 새로운 `template.json`을 만들고 종료합니다.

`template.json`은 다음과 같은 json 형식을 가지고 있습니다.

```json
{
  "name": "mongodb-community",
  "image": "20.04",
  "cpu": 1,
  "ram": 2048,
  "disk": 8192
}
```

이름은 인스턴스 이름으로 사용되며 같은 이름이 있을 경우 하이픈과 숫자를 추가하여 만듭니다.  
이미지는 우분투 버전을 의미합니다. 현재는 20.10이 최신이며 LTS로는 20.04가 최신입니다.  
CPU는 총 cpu 쓰레드 수를 의미하며 정수로 작성하여야합니다.  
RAM은 램 메모리를 의미하며 단위는 MB입니다. 최소 128MB를 확보하여야합니다.  
Disk는 디스크를 의미하며 단위는 MB입니다. 최소 512MB를 확보하여야합니다.

## insert

`stem insert <vmName> <dir>`를 실행하여 경로 내 `fertilizer` 폴더를 인스턴스에 복사합니다.

인스턴스에 복사할 때 홈 폴더의 경로 내 `<dir>`의 베이스 경로 이름의 폴더로 복사합니다.  
`~/datas/mongo`를 `<dir>`로 주어 mongo 폴더 내의 fertilizer를 복사하면 인스턴스 내 `~/mongo` 폴더 내로 복사합니다.  
부모 폴더 이름을 기준으로 하기 때문에 스크립트를 작성할 때 주의해야합니다.

만약 경로 내 `fertilizer` 폴더가 없을 경우 새로 만들고 종료합니다.

## install

`stem install <vmName> <dir>`를 실행하여 경로 내 `seed.json`을 읽어서 템플릿에 맞춰서 인스턴스를 만듭니다.

만약 경로 내 `seed.json` 파일이 없다면 새로운 파일을 만들고 종료합니다.

```json
{
  "using_python": false,
  "using_ruby": false,
  "using_julia": true,
  "scripts": ["scripts/init.jl"]
}
```

`using_python`, `using_ruby`, `using_julia`은 각각 파이썬, 루비, 줄리아를 사용하는 지 체크합니다.  
사용한다고 체크되어 있다면 스크립트를 실행하기 전 루비와 줄리아를 snap store에서 설치합니다.  
이후 `scripts` 항목에 기재된 스크립트들을 순서대로 인스턴스 내로 복사하여 실행하게 됩니다.

## list

`stem list <regexp>`를 실행하여 해당 정규표현식에 부합하는 이름을 가진 모든 인스턴스의 상태와 IPv4를 출력합니다.

정규표현식의 해당 부분은 `()`로 감싸서 작성해야합니다.

한 예로 모든 인스턴스의 상태를 출력하는 커맨드는 `stem list (.*)`입니다.

## remove

`stem remove <regexp>`를 실행하여 해당 정규표현식에 부합하는 이름을 가진 모든 인스턴스를 삭제합니다.

list와 같은 형태로 작성하면 됩니다.

한 예로 mongo로 시작하는 모든 인스턴스를 삭제하는 커맨드는 `stem remove mongo(.*)`입니다.