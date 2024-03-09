# cognito-authflow-handler

A small we application that implements the client part of authorization code flow (see https://tools.ietf.org/html/rfc6749#section-4.1) with AWS Cognito as an Authorization Server.


## Proxy Set Up

Should be run behind a http server like nginx to protect various resources using `auth_request`

Exanple nginx configuration where the cognito-authflow-handler is listening on `localhost:4000`

```yaml

http{
  ...
  server {
        ## retrieve the token after login from code or validate token and attempt refresh if expired
        location = /auth {
            proxy_pass http://localhost:4000/auth; 
            proxy_pass_request_body off;
            proxy_set_header Content-Length "";
            proxy_set_header X-Original-URI $request_uri;
        }

        ## a protected endpoint
        location /secure {
            auth_request /auth; 
            auth_request_set $auth_cookie $upstream_http_set_cookie;
            add_header Set-Cookie $auth_cookie;

            try_files /secure.html /404.html;
        }

        location = /login {
            return 301 https://<resourceserver>.auth.<region>.amazoncognito.com/oauth2/authorize?...;
        }

        ## remove tokens and invalidate the refresh token
        location = /logout {
            proxy_pass http://localhost:4000/logout;
            proxy_pass_request_body off;
            proxy_set_header Content-Length "";
            proxy_set_header X-Original-URI $request_uri;
        }
  }
}

```

## Configuration File
