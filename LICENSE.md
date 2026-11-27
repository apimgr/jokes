# MIT License

Copyright (c) 2025 APIMGR

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

---

## Third-Party Licenses

## Original joke corpus

The project owner states that the joke corpus in `src/data/jokes.json` is their own work and is in the public domain. This statement applies to the corpus and does not change the MIT license for the application code.

This project includes the following third-party software:

### Gin Web Framework
- License: MIT
- Copyright: Copyright (c) 2014 Manuel Martínez-Almeida
- URL: https://github.com/gin-gonic/gin

### YAML for Go
- License: MIT and Apache 2.0
- Copyright: Copyright 2011-2016 Canonical Ltd.
- URL: https://github.com/go-yaml/yaml

### Dracula Theme (UI Colors)
- License: MIT
- Copyright: Dracula Theme
- URL: https://draculatheme.com/

### Runtime GeoIP datasets
- DB-IP City databases: CC BY 4.0; source: https://db-ip.com/.
- IPtoASN ASN and Country databases: PDDL; source: https://iptoasn.com/.
- Required City database attribution: [IP Geolocation by DB-IP](https://db-ip.com/).

### Vendored browser documentation tools

| Library | Version | License | Copyright |
|---|---:|---|---|
| Swagger UI (`swagger-ui-dist`) | 5.34.1 | Apache-2.0 | 2020-2021 SmartBear Software Inc. |
| GraphiQL | 4.0.0 | MIT | 2021 GraphQL Contributors |
| React | 18.3.1 | MIT | Facebook, Inc. and its affiliates |
| ReactDOM | 18.3.1 | MIT | Facebook, Inc. and its affiliates |
| classnames (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | 2018 Jed Watson |
| extend (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | 2013-2018 Viacheslav Lotsmanov |
| buffer and safe-buffer (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | Feross Aboukhadijeh |
| JSON-Patch (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | 2017-2022 Joachim Wester |
| repeat-string (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | 2014-2015 Jon Schlinkert |
| DOMPurify (bundled with Swagger UI) | 3.4.17 | Apache-2.0 and MPL-2.0 | Cure53 and contributors |
| ieee754 (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | BSD-3-Clause | Feross Aboukhadijeh |
| regenerator-runtime (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | Facebook, Inc. |
| graphql-js (bundled with Swagger UI) | As distributed by Swagger UI 5.34.1 | MIT | 2014-present Lee Byron and contributors |
| scheduler and use-sync-external-store (bundled with React) | 18.3.1 distribution | MIT | Facebook, Inc. and Meta Platforms, Inc. |

The upstream license and notice files are retained beside the browser distributions in `src/web/static/vendor/`. Swagger UI's bundled dependency notices are also retained in its `.LICENSE.txt` files. The Swagger UI `NOTICE` content required by Apache-2.0 is:

> swagger-ui  
> Copyright 2020-2021 SmartBear Software Inc.

BSD-3-Clause non-endorsement condition: neither the names of the copyright holders nor the names of their contributors may be used to endorse or promote products derived from the software without specific prior written permission.

All embedded third-party software retains its original license and copyright notices.
