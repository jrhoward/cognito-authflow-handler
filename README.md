# cognito-authflow-handler

A small web application that implements the client part of OAUTH2 authorization code flow (see https://tools.ietf.org/html/rfc6749#section-4.1) with AWS Cognito as an Authorization Server.


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
        # Cognito hosted UI to login
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

A configuration file should be passed at start as the first argument.

eg:

```yaml
server:
  domain: localhost
  port: 4000
  authHandlerRedirect: localhost/redirect
  logoutHandlerRedirect: localhost
  idCookieName: aaaabbbbbb
  refreshCookieName: ccccdddd
cognito:
  oauthServer: https://<hosted endpoint>.auth.<region>.amazoncognito.com
  callBackUrl: http://localhost:80/auth
  clientId: <app client id> # Alternatively can be an environment variable: CLIENT_ID
  clientSecret: <app client secret> # Alternatively can be an environment variable: CLIENT_SECRET
  poolId: <user pool id>
  scope: <some scope>
  awsRegion: <aws region>

```

