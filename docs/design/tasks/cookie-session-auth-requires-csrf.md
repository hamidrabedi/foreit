Your workspace is /home/hamid/Other/projects/foreit-wt/wave3. Always use absolute paths. Never read or edit outside it. Do NOT run any git command.

# Cookie session authentication is only accepted on unsafe methods when CSRF protection ran

Module root: /home/hamid/Other/projects/foreit-wt/wave3/forge (run go commands from there).
Only edit forge/identity/middleware/auth.go and its tests.

## Principles (mandatory)
- Tests FIRST. gofmt; functions < 40 lines; naming after behaviour.

## Facts (verified on master)
forge/identity/middleware/auth.go ~102-121 authenticates a user from the `X-Session-Key` header or, if absent, from the `session_key` cookie. A cookie is sent by the browser automatically, so a cross-site form POST is authenticated unless CSRF middleware is mounted on the same routes; nothing checks that (B29). Django's `SessionAuthentication` enforces CSRF for session-cookie auth; header-based credentials do not need it.
forge/server has CSRF middleware based on `github.com/gorilla/csrf` (already in go.mod); find it (grep `csrf.Protect` in forge/server).

## Change
- Record whether the session key came from the cookie.
- For a cookie-sourced key on an unsafe method (anything except GET, HEAD, OPTIONS, TRACE), accept it only if CSRF validation ran for this request. With gorilla/csrf that means `csrf.Token(r) != ""` (the token is only set in the request context by the middleware, which already rejected invalid tokens). Otherwise do not authenticate from the cookie (fall through as unauthenticated, same as an invalid session) — do not return a new error type.
- Header-sourced keys and safe methods are unchanged.
- Update the function's doc comment to state the rule.

## Tests
- cookie + POST without CSRF middleware → not authenticated.
- cookie + POST wrapped in `csrf.Protect(key, csrf.Secure(false))` with a valid token (GET first to obtain token and cookie, then POST with `X-CSRF-Token`) → authenticated.
- cookie + GET → authenticated; `X-Session-Key` header + POST without CSRF → authenticated.
Use fakes for the session and user repositories the way existing tests in forge/identity/middleware do.

## Acceptance (from the module root)
    gofmt -l identity                                 # empty
    go vet ./identity/...
    staticcheck ./identity/...                        # no output
    go test -race ./identity/... -count=1

Report (max 15 lines): failing output before, change, each command result.
