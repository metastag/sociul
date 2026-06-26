### Sociul

A microservices-based social media platform.

---

### Auth Service

- User Accounts (sign up, login, logout, delete account)
- JWT Authentication with introspection endpoint (refresh token in redis cache + access tokens with short ttl)
- Otp flows - forgot password, verify email (otp stored in redis)
- Account lockout on failed login attempts
- Future Work -> Authorization/roles/scopes

