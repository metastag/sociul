\### Sociul



A microservices-based social media platform.





\---



\### API Gateway



* Routes requests to microservices
* Future work -> auth token verification/extract claims, rate limiting





\### Auth Service



* User Accounts (sign up, login, logout, delete account)
* Authentication (refresh token in redis cache + access tokens with short ttl)
* Otp flows (forgot password, verify email, otp stored in redis)
* Future Work -> Authorization/roles, Account lockout on failed login attempts, Rate limiting

