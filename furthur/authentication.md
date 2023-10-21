## HTTP basic authentication
- client는 Authorization header에 credential을 넣어 요청
- credential format: `username:password`, and base64-encoded
  - `Authorization: Basic <base64-encoded credential>`
- 간단하나, 서버에서 요청마다 costly op를 수행하는 부담이 있다.

## Token(bearer) authentication
- client가 credential을 포함하여 요청전송
- API는 credential을 validation한후 bearer token을 생성하여 전달
  - bearer token: 사용자를 represent함. 보통 시간제한 있음
- 이후 client는 bearer token을 포함하여 요청전송
  - `Authorization: Bearer <token>`
- API는 bearer를 받은 후 만료유무를 체크한 후 사용자를 결정함.
- bearer token 만들때만 costly op를 수행한다.
- client측면에서 token caching/monitoring/expire&regen 하려면 성가심

### Stateful token authentication
- token값은 cryptographically secure random string.
- {token, userid, expiry}가 db에 저장됨
- 간단하며, API가 token을 control하는 장점.
- db op가 있으나 어차피 해야할 일이다.

### Stateless token authentication
- userid와 expiry를 token에 저장후 crypto sign하여 전달.
- ```JWT```, PASETO, Branca, nacl/secretbox 등이 있다.
- 장점: 토큰 op가 메모리에서 수행되며, token에 user정보가 존재. no db op.
- 단점: revoke하기 어려움.
  - secret을 바꿔 all-user revocation을 하거나, blacklist를 관리 (stateful해짐 -_-)
- stateless token에 application state를 넣지 마라 --> stale해지는 위험존재.
- JWT: highly configurable하나 보안위헙 많다
  - https://auth0.com/blog/critical-vulnerabilities-in-json-web-token-libraries/
  - https://curity.io/resources/learn/jwt-best-practices/
  - 이로 인해 JWT는 best auth scheme이 아니지만,
  - token을 발행하는 주체와 소비하는 주체가 분리되어 있을땐 매우 유용하다. (delegated authentication)
  
## API-key authentication
- 사용자는 자신의 계정에 강한 보안강도의 영구 secret를 보유
- 사용자는 요청마다 키를 전송
  - `Authorization: Key <key>`
- API는 해당 키에 대한 해시를 계산하여 user id를 db에서 검색
  - stateful token은 temporary key, API-key는 permanent key
- 클라이언트가 토큰관리를 하지않는 장점 vs 패스워드와 APIkey의 보안문제 단점

## OAuth 2.0 / OpenID Connect
- credential을 서드파티 ID provider가 관리
- OAuth 2.0은 auth proro가 아님.
  - https://oauth.net/articles/authentication/
- OpenID Connect: auth proto
  - req를 인증하기 위해 user를 IDP의 `auth & consent` form으로 redirect
  - 사용자가 consent(OK)하면 IDP는 API에 authorization code를 전송
  - API는 authorization code를 IDP의 다른 endpoint로 전송하고
    IDP는 코드가 확인되면 ID token(JWT)을 포함한 JSON 응답을 API에게 전송
  - 이후는 token auth와 동일
- 장점: user credential의 관리가 불필요
- 단점: 복잡하다

## 선택
- API가 실제 사용자 계정과 연관이 없다면: HTTP basic auth도 좋은 방안
- 웹사이트 백엔드환경이고 패스워드 관리가 부담이라면: OpenID Connect
- delegated auth환경, 즉 MSA구조에서 auth와 service가 분리되었다면: stateless auth
- 그밖이라면: API-key 혹은 stateful auth token.
  - Stateful auth token은 website/SPA의 backend에 적합
  - API-key는 general purpose
