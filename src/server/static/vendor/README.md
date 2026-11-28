# Vendored browser documentation tools

These pinned distributions are served locally by the Jokes binary. Runtime pages never load these assets from a CDN. Do not edit vendor distribution files; update them from the official package sources and preserve the included licenses, notices, and integrity records.

| Tool | Version | Source | License | Package integrity (SHA-512 SRI) |
|---|---:|---|---|---|
| Swagger UI (`swagger-ui-dist`) | 5.34.1 | <https://registry.npmjs.org/swagger-ui-dist/-/swagger-ui-dist-5.34.1.tgz> | Apache-2.0, with bundled dependency notices | `sha512-VSr8fXfI93obYaQicLJbLUKGRxDojj/SfRJDDLCioGTUBj4BK/Nj4eE9jCeCAnJxVLi7iS+08PbiZLyrQKBkJQ==` |
| GraphiQL (`graphiql`) | 4.0.0 | <https://registry.npmjs.org/graphiql/-/graphiql-4.0.0.tgz> | MIT | `sha512-zJYxj9W7Epzc7F43HjC8LQTv+HFk7LxDHV2COTgSskdq2YX0OvgngvDx+D5V3BXcoW+0R+GADYErVs5MPEoJ+g==` |
| React | 18.3.1 | <https://registry.npmjs.org/react/-/react-18.3.1.tgz> | MIT | `sha512-wS+hAgJShR0KhEvPJArfuPVN1+Hz1t0Y6n5jLrGQbkb4urgPE/0Rve+1kMB1v/oWgHgm4WIcV+i7F2pTVj+2iQ==` |
| ReactDOM | 18.3.1 | <https://registry.npmjs.org/react-dom/-/react-dom-18.3.1.tgz> | MIT | `sha512-5m4nQKp+rZRb09LNH59GM4BxTh9251/ylbKIbpe7TpGxfJ+9kv6BLkLBXIjjspbgbnIBNqlI23tRnTWT0snUIw==` |

The upstream license and notice files remain beside each distribution. Swagger UI's JavaScript also retains its upstream `.LICENSE.txt` dependency notices. GraphiQL 4.0.0 is pinned because it includes the standalone UMD distribution needed for local, CSP-compatible serving; GraphiQL's current major uses ESM and requires a separate bundle pipeline.
