## CORS (Cross Origin Requests)
- 보내는 것은 OK
- 받을때 웹브라우저에서 Origin이 다르면 차단
- 허용하는 방법
  - `Access-Control-Allow-Origin: *`: allow any other origin
  - `Access-Control-Allow-Origin: https://www.example.com`
  - 단 한개의 origin만을 설정이 가능하다는 단점. (리스트 불가. 표준에는 있으나 구현을 안함)
  - API에서 req의 Origin 헤더를 보고 해당 내용을 반영해주는식으로 해결한다.
	- 중요: `Vary: Origin`을 꼭 포함할 것
### Auth and CORS
- API가 credential을 요구할 경우 다음 헤더가 필요
  - `Access-Control-Allow-Credentials: true`
### CORS 종류
- Simple
  - HEAD/GET/POST 메소드 (CORS-safe)
  - Request Header가 `Forbidden header` or 다음 CORS-safe 헤더일 것
    - Accept
	- Accept-Language
	- Content-Language
	- Content-Type
  - Content-Type의 값이 다음중 하나
    - application/x-www-form-urlencoded
	- multipart/form-data
	- text/plain
- Preflight
  - Simple CORS가 아닌 경우 실제 요청 전에 요청하는 Preflight 요청
    - 실 CORS가 허용될 것인지를 확인하기 위함
