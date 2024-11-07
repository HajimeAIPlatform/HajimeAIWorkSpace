load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "com_github_armon_go_metrics",
        importpath = "github.com/armon/go-metrics",
        sum = "h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_bytedance_sonic",
        importpath = "github.com/bytedance/sonic",
        sum = "h1:oUp34TzMlL+OY1OUWxHqsdkgC/Zfc85zGqw9siXjrc0=",
        version = "v1.11.6",
    )
    go_repository(
        name = "com_github_bytedance_sonic_loader",
        importpath = "github.com/bytedance/sonic/loader",
        sum = "h1:c+e5Pt1k/cy5wMveRDyk2X4B9hF4g7an8N3zCYjJFNM=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_github_cloudwego_base64x",
        importpath = "github.com/cloudwego/base64x",
        sum = "h1:jwCgWpFanWmN8xoIUHa2rtzmkd5J2plF/dnLS6Xd/0Y=",
        version = "v0.1.4",
    )
    go_repository(
        name = "com_github_cloudwego_iasm",
        importpath = "github.com/cloudwego/iasm",
        sum = "h1:1KNIy1I1H9hNNFEEH3DVnI4UujN+1zjpuk6gwHLTssg=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_coreos_go_semver",
        importpath = "github.com/coreos/go-semver",
        sum = "h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_coreos_go_systemd_v22",
        importpath = "github.com/coreos/go-systemd/v22",
        sum = "h1:D9/bQk5vlXQFZ6Kwuu6zaiXJ9oTPe68++AzAJc1DzSI=",
        version = "v22.3.2",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=",
        version = "v1.1.2-0.20180830191138-d8f796af33cc",
    )
    go_repository(
        name = "com_github_dustin_go_humanize",
        importpath = "github.com/dustin/go-humanize",
        sum = "h1:GzkhY7T5VNhEkwH0PVJgjz+fX1rhBrR7pRT3mDkpeCY=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_fatih_color",
        importpath = "github.com/fatih/color",
        sum = "h1:qfhVLaG5s+nCROl1zJsZRxFeYrHLqWroPOQ8BWiNb4w=",
        version = "v1.14.1",
    )
    go_repository(
        name = "com_github_felixge_httpsnoop",
        importpath = "github.com/felixge/httpsnoop",
        sum = "h1:NFTV2Zj1bL4mc9sqWACXbQFVBBg2W3GPvqp8/ESS2Wg=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_frankban_quicktest",
        importpath = "github.com/frankban/quicktest",
        sum = "h1:7Xjx+VpznH+oBnejlPUj8oUpdxnVs4f8XU8WnHkI4W8=",
        version = "v1.14.6",
    )
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:8JEhPFa5W2WU7YfeZzPNqzMP6Lwt7L2715Ggo0nosvA=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_gabriel_vasile_mimetype",
        importpath = "github.com/gabriel-vasile/mimetype",
        sum = "h1:in2uUcidCuFcDKtdcBxlR0rJ1+fsokWf+uqxgUFjbI0=",
        version = "v1.4.3",
    )
    go_repository(
        name = "com_github_gin_contrib_cors",
        importpath = "github.com/gin-contrib/cors",
        sum = "h1:oLDHxdg8W/XDoN/8zamqk/Drgt4oVZDvaV0YmvVICQw=",
        version = "v1.7.2",
    )
    go_repository(
        name = "com_github_gin_contrib_sse",
        importpath = "github.com/gin-contrib/sse",
        sum = "h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gin_gonic_gin",
        importpath = "github.com/gin-gonic/gin",
        sum = "h1:nTuyha1TYqgedzytsKYqna+DfLos46nTv2ygFy86HFU=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_go_ini_ini",
        importpath = "github.com/go-ini/ini",
        sum = "h1:z6ZrTEZqSWOTyH2FlglNbNgARyHG8oLW9gMELqKr06A=",
        version = "v1.67.0",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:pKouT5E8xu9zeFC39JXRDukb6JFQPXM5p5I91188VAQ=",
        version = "v1.4.1",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_go_playground_assert_v2",
        importpath = "github.com/go-playground/assert/v2",
        sum = "h1:JvknZsQTYeFEAhQwI4qEt9cyV5ONwRHC+lYKSsYSR8s=",
        version = "v2.2.0",
    )
    go_repository(
        name = "com_github_go_playground_locales",
        importpath = "github.com/go-playground/locales",
        sum = "h1:EWaQ/wswjilfKLTECiXz7Rh+3BjFhfDFKv/oXslEjJA=",
        version = "v0.14.1",
    )
    go_repository(
        name = "com_github_go_playground_universal_translator",
        importpath = "github.com/go-playground/universal-translator",
        sum = "h1:Bcnm0ZwsGyWbCzImXv+pAJnYK9S473LQFuzCbDbfSFY=",
        version = "v0.18.1",
    )
    go_repository(
        name = "com_github_go_playground_validator_v10",
        importpath = "github.com/go-playground/validator/v10",
        sum = "h1:K9ISHbSaI0lyB2eWMPJo+kOS/FBExVwjEviJTixqxL8=",
        version = "v10.20.0",
    )
    go_repository(
        name = "com_github_go_sql_driver_mysql",
        importpath = "github.com/go-sql-driver/mysql",
        sum = "h1:LedoTUt/eveggdHS9qUFC1EFSa8bU2+1pZjSRpvNJ1Y=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_goccy_go_json",
        importpath = "github.com/goccy/go-json",
        sum = "h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        importpath = "github.com/gogo/protobuf",
        sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt",
        importpath = "github.com/golang-jwt/jwt",
        sum = "h1:IfV12K8xAKAnZqdXVzCZ+TOjboZ2keLg81eXfW3O+oY=",
        version = "v3.2.2+incompatible",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_github_golang_sql_civil",
        importpath = "github.com/golang-sql/civil",
        sum = "h1:au07oEsX2xN0ktxqI+Sida1w446QrXBRJ0nee3SNZlA=",
        version = "v0.0.0-20220223132316-b832511892a9",
    )
    go_repository(
        name = "com_github_golang_sql_sqlexp",
        importpath = "github.com/golang-sql/sqlexp",
        sum = "h1:ZCD6MBpcuOVfGVqsEmY5/4FtYiKz6tSyUv9LPEDei6A=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=",
        version = "v0.5.9",
    )
    go_repository(
        name = "com_github_google_gofuzz",
        importpath = "github.com/google/gofuzz",
        sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_google_s2a_go",
        importpath = "github.com/google/s2a-go",
        sum = "h1:60BLSyTrOV4/haCDW4zb1guZItoSq8foHCXrAnjBo/o=",
        version = "v0.1.7",
    )
    go_repository(
        name = "com_github_google_uuid",
        importpath = "github.com/google/uuid",
        sum = "h1:NIvaJDMOsjHA8n1jAhLSgzrAzy1Hgr+hNrb57e+94F0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:Vie5ybvEvT75RniqhfFxPRy3Bf7vr3h0cechB90XaQs=",
        version = "v0.3.2",
    )
    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:5/zPPDvw8Q1SuXjrqrZslrqT7dL/uJT2CQii/cLCKqA=",
        version = "v2.12.3",
    )
    go_repository(
        name = "com_github_googleapis_google_cloud_go_testing",
        importpath = "github.com/googleapis/google-cloud-go-testing",
        sum = "h1:zC34cGQu69FG7qzJ3WiKW244WfhDC3xxYMeNOX2gtUQ=",
        version = "v0.0.0-20210719221736-1c9a4c676720",
    )
    go_repository(
        name = "com_github_gopherjs_gopherjs",
        importpath = "github.com/gopherjs/gopherjs",
        sum = "h1:EGx4pi6eqNxGaHF6qqu48+N2wcFQ5qg5FXgOdqsJ5d8=",
        version = "v0.0.0-20181017120253-0766667cb4d1",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        importpath = "github.com/gorilla/mux",
        sum = "h1:TuBL49tXwgrFYWhqrNgrUNEY92u81SPhu7sTdzQEiWY=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_hashicorp_consul_api",
        importpath = "github.com/hashicorp/consul/api",
        sum = "h1:mXfkRHrpHN4YY3RqL09nXU1eHKLNiuAN4kHvDQ16k/8=",
        version = "v1.28.2",
    )
    go_repository(
        name = "com_github_hashicorp_errwrap",
        importpath = "github.com/hashicorp/errwrap",
        sum = "h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_cleanhttp",
        importpath = "github.com/hashicorp/go-cleanhttp",
        sum = "h1:035FKYIWjmULyFRBKPs8TBQoi0x6d9G4xc9neXJWAZQ=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_hclog",
        importpath = "github.com/hashicorp/go-hclog",
        sum = "h1:bI2ocEMgcVlz55Oj1xZNBsVi900c7II+fWDyV9o+13c=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_immutable_radix",
        importpath = "github.com/hashicorp/go-immutable-radix",
        sum = "h1:DKHmCUm2hRBK510BaiZlwvpD40f8bJFeZnpfm2KLowc=",
        version = "v1.3.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_multierror",
        importpath = "github.com/hashicorp/go-multierror",
        sum = "h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_rootcerts",
        importpath = "github.com/hashicorp/go-rootcerts",
        sum = "h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru",
        importpath = "github.com/hashicorp/golang-lru",
        sum = "h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=",
        version = "v0.5.4",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_serf",
        importpath = "github.com/hashicorp/serf",
        sum = "h1:Z1H2J60yRKvfDYAOZLd2MU0ND4AH/WDz7xYHDWQsIPY=",
        version = "v0.10.1",
    )
    go_repository(
        name = "com_github_jackc_pgpassfile",
        importpath = "github.com/jackc/pgpassfile",
        sum = "h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_jackc_pgservicefile",
        importpath = "github.com/jackc/pgservicefile",
        sum = "h1:L0QtFUgDarD7Fpv9jeVMgy/+Ec0mtnmYuImjTz6dtDA=",
        version = "v0.0.0-20231201235250-de7065d80cb9",
    )
    go_repository(
        name = "com_github_jackc_pgx_v5",
        importpath = "github.com/jackc/pgx/v5",
        sum = "h1:amBjrZVmksIdNjxGW/IiIMzxMKZFelXbUoPNb+8sjQw=",
        version = "v5.5.5",
    )
    go_repository(
        name = "com_github_jackc_puddle_v2",
        importpath = "github.com/jackc/puddle/v2",
        sum = "h1:RhxXJtFG022u4ibrCSMSiu5aOq1i77R3OHKNJj77OAk=",
        version = "v2.2.1",
    )
    go_repository(
        name = "com_github_jinzhu_inflection",
        importpath = "github.com/jinzhu/inflection",
        sum = "h1:K317FqzuhWc8YvSVlFMCCUb36O/S9MCKRDI7QkRKD/E=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_jinzhu_now",
        importpath = "github.com/jinzhu/now",
        sum = "h1:/o9tlHleP7gOFmsnYNz3RGnqzefHA47wQpKrrdTIwXQ=",
        version = "v1.1.5",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_jtolds_gls",
        importpath = "github.com/jtolds/gls",
        sum = "h1:xdiiI2gbIgH/gLH7ADydsJ1uDOEzR8yvV7C0MuV77Wo=",
        version = "v4.20.0+incompatible",
    )
    go_repository(
        name = "com_github_k3a_html2text",
        importpath = "github.com/k3a/html2text",
        sum = "h1:nvnKgBvBR/myqrwfLuiqecUtaK1lB9hGziIJKatNFVY=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_klauspost_compress",
        importpath = "github.com/klauspost/compress",
        sum = "h1:In6xLpyWOi1+C7tXUUWv2ot1QvBjxevKAaI6IXrJmUc=",
        version = "v1.17.11",
    )
    go_repository(
        name = "com_github_klauspost_cpuid_v2",
        importpath = "github.com/klauspost/cpuid/v2",
        sum = "h1:+StwCXwm9PdpiEkPyzBXIy+M9KUb4ODm0Zarf1kS5BM=",
        version = "v2.2.8",
    )
    go_repository(
        name = "com_github_knz_go_libedit",
        importpath = "github.com/knz/go-libedit",
        sum = "h1:0pHpWtx9vcvC0xGZqEQlQdfSQs7WRlAjuPvk3fOZDCo=",
        version = "v1.10.1",
    )
    go_repository(
        name = "com_github_kr_fs",
        importpath = "github.com/kr/fs",
        sum = "h1:Jskdu9ieNAYnjxsi0LbQp1ulIKZV1LAFgK1tWhpZgl8=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_kr_pretty",
        importpath = "github.com/kr/pretty",
        sum = "h1:flRD4NNwYAUpkphVc1HcthR4KEIFJ65n8Mw5qdRn3LE=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_kr_text",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_leodido_go_urn",
        importpath = "github.com/leodido/go-urn",
        sum = "h1:WT9HwE9SGECu3lg4d/dIA+jxlljEa1/ffXKmRjqdmIQ=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_magiconair_properties",
        importpath = "github.com/magiconair/properties",
        sum = "h1:IeQXZAiQcpL9mgcAe1Nu6cX9LLw6ExEHKjN0VQdvPDY=",
        version = "v1.8.7",
    )
    go_repository(
        name = "com_github_mattn_go_colorable",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:fFA4WZxdEF4tXPZVKMLwD8oUnCTTo08duU7wxecdEvA=",
        version = "v0.1.13",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:xfD0iDuEKnDkl03q4limB+vH+GxLEtL/jb4xVJSWWEY=",
        version = "v0.0.20",
    )
    go_repository(
        name = "com_github_mattn_go_sqlite3",
        importpath = "github.com/mattn/go-sqlite3",
        sum = "h1:vfoHhTN1af61xCRSWzFIWzx2YskyMTwHLrExkBOjvxI=",
        version = "v1.14.15",
    )
    go_repository(
        name = "com_github_microsoft_go_mssqldb",
        importpath = "github.com/microsoft/go-mssqldb",
        sum = "h1:Fto83dMZPnYv1Zwx5vHHxpNraeEaUlQ/hhHLgZiaenE=",
        version = "v0.17.0",
    )
    go_repository(
        name = "com_github_minio_md5_simd",
        importpath = "github.com/minio/md5-simd",
        sum = "h1:Gdi1DZK69+ZVMoNHRXJyNcxrMA4dSxoYHZSQbirFg34=",
        version = "v1.1.2",
    )
    go_repository(
        name = "com_github_minio_minio_go_v7",
        importpath = "github.com/minio/minio-go/v7",
        sum = "h1:2mdUHXEykRdY/BigLt3Iuu1otL0JTogT0Nmltg0wujk=",
        version = "v7.0.80",
    )
    go_repository(
        name = "com_github_mitchellh_go_homedir",
        importpath = "github.com/mitchellh/go-homedir",
        sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_mitchellh_mapstructure",
        importpath = "github.com/mitchellh/mapstructure",
        sum = "h1:jeMsZIYE/09sWLaz43PL7Gy6RuMjD2eJVyuac5Z2hdY=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_nats_io_nats_go",
        importpath = "github.com/nats-io/nats.go",
        sum = "h1:fnxnPCNiwIG5w08rlMcEKTUw4AV/nKyGCOJE8TdhSPk=",
        version = "v1.34.0",
    )
    go_repository(
        name = "com_github_nats_io_nkeys",
        importpath = "github.com/nats-io/nkeys",
        sum = "h1:RwNJbbIdYCoClSDNY7QVKZlyb/wfT6ugvFCiKy6vDvI=",
        version = "v0.4.7",
    )
    go_repository(
        name = "com_github_nats_io_nuid",
        importpath = "github.com/nats-io/nuid",
        sum = "h1:5iA8DT8V7q8WK2EScv2padNa/rTESc1KdnPw4TC2paw=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:aYUidT7k73Pcl9nb2gScu7NSrKCSHIDE89b3+6Wq+LM=",
        version = "v2.2.2",
    )
    go_repository(
        name = "com_github_pkg_errors",
        importpath = "github.com/pkg/errors",
        sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
        version = "v0.9.1",
    )
    go_repository(
        name = "com_github_pkg_sftp",
        importpath = "github.com/pkg/sftp",
        sum = "h1:JFZT4XbOU7l77xGSpOdW+pwIMqP044IyjXX6FGyEKFo=",
        version = "v1.13.6",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:Jamvg5psRIccs7FGNTlIRMkT8wgtp5eCXdBlqhYGL6U=",
        version = "v1.0.1-0.20181226105442-5d4384ee4fb2",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:73kH8U+JUqXU8lRuOHeVHaa/SZPifC7BkcraZVejAe8=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_rs_xid",
        importpath = "github.com/rs/xid",
        sum = "h1:fV591PaemRlL6JfRxGDEPl69wICngIQ3shQtzfy2gxU=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_crypt",
        importpath = "github.com/sagikazarmark/crypt",
        sum = "h1:WMyLTjHBo64UvNcWqpzY3pbZTYgnemZU8FBZigKc42E=",
        version = "v0.19.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_locafero",
        importpath = "github.com/sagikazarmark/locafero",
        sum = "h1:HApY1R9zGo4DBgr7dqsTH/JJxLTTsOt7u6keLGt6kNQ=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_slog_shim",
        importpath = "github.com/sagikazarmark/slog-shim",
        sum = "h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_sashabaranov_go_openai",
        importpath = "github.com/sashabaranov/go-openai",
        sum = "h1:/eNVa8KzlE7mJdKPZDj6886MUzZQjoVHyn0sLvIt5qA=",
        version = "v1.32.5",
    )
    go_repository(
        name = "com_github_smartwalle_alipay_v3",
        importpath = "github.com/smartwalle/alipay/v3",
        sum = "h1:i1VwJeu70EmwpsXXz6GZZnMAtRx5MTfn2dPoql/L3zE=",
        version = "v3.2.23",
    )
    go_repository(
        name = "com_github_smartwalle_ncrypto",
        importpath = "github.com/smartwalle/ncrypto",
        sum = "h1:P2rqQxDepJwgeO5ShoC+wGcK2wNJDmcdBOWAksuIgx8=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_smartwalle_ngx",
        importpath = "github.com/smartwalle/ngx",
        sum = "h1:pUXDvWRZJIHVrCKA1uZ15YwNti+5P4GuJGbpJ4WvpMw=",
        version = "v1.0.9",
    )
    go_repository(
        name = "com_github_smartwalle_nsign",
        importpath = "github.com/smartwalle/nsign",
        sum = "h1:8poAgG7zBd8HkZy9RQDwasC6XZvJpDGQWSjzL2FZL6E=",
        version = "v1.0.9",
    )
    go_repository(
        name = "com_github_smartystreets_assertions",
        importpath = "github.com/smartystreets/assertions",
        sum = "h1:zE9ykElWQ6/NYmHa3jpm/yHnI4xSofP+UP6SpjHcSeM=",
        version = "v0.0.0-20180927180507-b2de0cb4f26d",
    )
    go_repository(
        name = "com_github_smartystreets_goconvey",
        importpath = "github.com/smartystreets/goconvey",
        sum = "h1:fv0U8FUIMPNf1L9lnHLvLhgicrIVChEkdzIKYqbNC9s=",
        version = "v1.6.4",
    )
    go_repository(
        name = "com_github_sourcegraph_conc",
        importpath = "github.com/sourcegraph/conc",
        sum = "h1:OQTbbt6P72L20UqAkXXuLOj79LfEanQ+YQFNpLA9ySo=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_spf13_afero",
        importpath = "github.com/spf13/afero",
        sum = "h1:WJQKhtpdm3v2IzqG8VMqrr6Rf3UYpEF239Jy9wNepM8=",
        version = "v1.11.0",
    )
    go_repository(
        name = "com_github_spf13_cast",
        importpath = "github.com/spf13/cast",
        sum = "h1:GEiTHELF+vaR5dhz3VqZfFSzZjYbgeKDpBxQVS4GYJ0=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_spf13_pflag",
        importpath = "github.com/spf13/pflag",
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_spf13_viper",
        importpath = "github.com/spf13/viper",
        sum = "h1:RWq5SEjt8o25SROyN3z2OrDB9l7RPd3lwTWU8EcEdcI=",
        version = "v1.19.0",
    )
    go_repository(
        name = "com_github_stretchr_objx",
        importpath = "github.com/stretchr/objx",
        sum = "h1:xuMeJ0Sdp5ZMRXx/aWO6RZxdr3beISkG5/G/aIRr3pY=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        importpath = "github.com/stretchr/testify",
        sum = "h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_thanhpk_randstr",
        importpath = "github.com/thanhpk/randstr",
        sum = "h1:psAOktJFD4vV9NEVb3qkhRSMvYh4ORRaj1+w/hn4B+o=",
        version = "v1.0.6",
    )
    go_repository(
        name = "com_github_twitchyliquid64_golang_asm",
        importpath = "github.com/twitchyliquid64/golang-asm",
        sum = "h1:SU5vSMR7hnwNxj24w34ZyCi/FmDZTkS4MhqMhdFk5YI=",
        version = "v0.15.1",
    )
    go_repository(
        name = "com_github_ugorji_go_codec",
        importpath = "github.com/ugorji/go/codec",
        sum = "h1:9LC83zGrHhuUA9l16C9AHXAqEV/2wBQ4nkvumAE65EE=",
        version = "v1.2.12",
    )
    go_repository(
        name = "com_google_cloud_go",
        importpath = "cloud.google.com/go",
        sum = "h1:uJSeirPke5UNZHIb4SxfZklVSiWWVqW4oXlETwZziwM=",
        version = "v0.112.1",
    )
    go_repository(
        name = "com_google_cloud_go_compute",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:phWcR2eWzRJaL/kOiJwfFsPs4BaKq1j6vnpZrc1YlVg=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:mg4jlk7mCAj6xXp9UJ4fjI9VUI5rubuGBW5aJ7UnBMY=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_google_cloud_go_firestore",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:/k8ppuWOtNuDHt2tsRV42yI21uaGnKDEQnRFeBpbFF8=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_iam",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:1jTsCu4bcsNsE4iiqNT5SHwrDRCfRmIaaaVFhRveTJI=",
        version = "v1.1.5",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:GOE6pZFdSrTb4KAiKnXsJBtlE6mEyaW44oKyMILWnOg=",
        version = "v0.5.5",
    )
    go_repository(
        name = "com_google_cloud_go_storage",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:B59ahL//eDfx2IIKFBeT5Atm9wnNmj3+8xG/W4WB//w=",
        version = "v1.35.1",
    )
    go_repository(
        name = "com_nullprogram_x_optparse",
        importpath = "nullprogram.com/x/optparse",
        sum = "h1:xGFgVi5ZaWOnYdac2foDT3vg0ZZC9ErXFV57mr4OHrI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "in_gopkg_alexcesaro_quotedprintable_v3",
        importpath = "gopkg.in/alexcesaro/quotedprintable.v3",
        sum = "h1:2gGKlE2+asNV9m7xrywl36YYNnBG5ZQ0r/BOOxqPpmk=",
        version = "v3.0.0-20150716171945-2caba252f4dc",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_gomail_v2",
        importpath = "gopkg.in/gomail.v2",
        sum = "h1:n7WqCuqOuCbNr617RXOY0AWRXxgwEyPp2z+p0+hgMuE=",
        version = "v2.0.0-20160411212932-81ebce5c23df",
    )
    go_repository(
        name = "in_gopkg_ini_v1",
        importpath = "gopkg.in/ini.v1",
        sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
        version = "v1.67.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_etcd_go_etcd_api_v3",
        importpath = "go.etcd.io/etcd/api/v3",
        sum = "h1:W4sw5ZoU2Juc9gBWuLk5U6fHfNVyY1WC5g9uiXZio/c=",
        version = "v3.5.12",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_pkg_v3",
        importpath = "go.etcd.io/etcd/client/pkg/v3",
        sum = "h1:EYDL6pWwyOsylrQyLp2w+HkQ46ATiOvoEdMarindU2A=",
        version = "v3.5.12",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v2",
        importpath = "go.etcd.io/etcd/client/v2",
        sum = "h1:0m4ovXYo1CHaA/Mp3X/Fak5sRNIWf01wk/X1/G3sGKI=",
        version = "v2.305.12",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v3",
        importpath = "go.etcd.io/etcd/client/v3",
        sum = "h1:v5lCPXn1pf1Uu3M4laUE2hp/geOTc5uPcYYsNe1lDxg=",
        version = "v3.5.12",
    )
    go_repository(
        name = "io_filippo_edwards25519",
        importpath = "filippo.io/edwards25519",
        sum = "h1:FNf4tywRC1HmFuKW5xopWpigGjJKiJSV0Cqo0cJWDaA=",
        version = "v1.1.0",
    )
    go_repository(
        name = "io_gorm_datatypes",
        importpath = "gorm.io/datatypes",
        sum = "h1:uZmGAcK/QZ0uyfCuVg0VQY1ZmV9h1fuG0tMwKByO1z4=",
        version = "v1.2.4",
    )
    go_repository(
        name = "io_gorm_driver_mysql",
        importpath = "gorm.io/driver/mysql",
        sum = "h1:Ld4mkIickM+EliaQZQx3uOJDJHtrd70MxAUqWqlx3Y8=",
        version = "v1.5.6",
    )
    go_repository(
        name = "io_gorm_driver_postgres",
        importpath = "gorm.io/driver/postgres",
        sum = "h1:DkegyItji119OlcaLjqN11kHoUgZ/j13E0jkJZgD6A8=",
        version = "v1.5.9",
    )
    go_repository(
        name = "io_gorm_driver_sqlite",
        importpath = "gorm.io/driver/sqlite",
        sum = "h1:HBBcZSDnWi5BW3B3rwvVTc510KGkBkexlOg0QrmLUuU=",
        version = "v1.4.3",
    )
    go_repository(
        name = "io_gorm_driver_sqlserver",
        importpath = "gorm.io/driver/sqlserver",
        sum = "h1:t4r4r6Jam5E6ejqP7N82qAJIJAht27EGT41HyPfXRw0=",
        version = "v1.4.1",
    )
    go_repository(
        name = "io_gorm_gorm",
        importpath = "gorm.io/gorm",
        sum = "h1:I0u8i2hWQItBq1WfE0o2+WuL9+8L21K9e2HHSTE/0f8=",
        version = "v1.25.12",
    )
    go_repository(
        name = "io_opencensus_go",
        importpath = "go.opencensus.io",
        sum = "h1:y73uSU6J157QMP2kn2r30vwW1A2W2WFwSCGnAVxeaD0=",
        version = "v0.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
        importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
        sum = "h1:4Pp6oUg3+e/6M4C0A/3kJ2VYa++dsWVTtGgLVj5xtHg=",
        version = "v0.49.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
        importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
        sum = "h1:jq9TW8u3so/bN+JPT166wjOI6/vQPF6Xe7nMNIltagk=",
        version = "v0.49.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel",
        importpath = "go.opentelemetry.io/otel",
        sum = "h1:0LAOdjNmQeSTzGBzduGe/rU4tZhMwL5rWgtp9Ku5Jfo=",
        version = "v1.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:6EhoGWWK28x1fbpA4tYTOWBkPefTDQnb8WSGXlc88kI=",
        version = "v1.24.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:CsKnnL4dUAr/0llH9FKuc698G04IrpWV0MQA/Y1YELI=",
        version = "v1.24.0",
    )
    go_repository(
        name = "io_rsc_pdf",
        importpath = "rsc.io/pdf",
        sum = "h1:k1MczvYDUvJBe93bYd7wrZLLUEcLZAuF824/I4e5Xr4=",
        version = "v0.1.1",
    )
    go_repository(
        name = "org_golang_google_api",
        importpath = "google.golang.org/api",
        sum = "h1:w174hnBPqut76FzW5Qaupt7zY8Kql6fiVjgys4f58sU=",
        version = "v0.171.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        importpath = "google.golang.org/appengine",
        sum = "h1:IhEN5q69dyKagZPYMSdIjS2HqprW324FRQZJcGqPAsM=",
        version = "v1.6.8",
    )
    go_repository(
        name = "org_golang_google_genproto",
        importpath = "google.golang.org/genproto",
        sum = "h1:9+tzLLstTlPTRyJTh+ah5wIMsBW5c4tQwGTN3thOW9Y=",
        version = "v0.0.0-20240213162025-012b6fc9bca9",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_api",
        importpath = "google.golang.org/genproto/googleapis/api",
        sum = "h1:rIo7ocm2roD9DcFIX67Ym8icoGCKSARAiPljFhh5suQ=",
        version = "v0.0.0-20240311132316-a219d84964c2",
    )
    go_repository(
        name = "org_golang_google_genproto_googleapis_rpc",
        importpath = "google.golang.org/genproto/googleapis/rpc",
        sum = "h1:lfpJ/2rWPa/kJgxyyXM8PrNnfCzcmxJ265mADgwmvLI=",
        version = "v0.0.0-20240314234333-6e1732d8331c",
    )
    go_repository(
        name = "org_golang_google_grpc",
        importpath = "google.golang.org/grpc",
        sum = "h1:B4n+nfKzOICUXMgyrNd19h/I9oH0L1pizfk1d4zSgTk=",
        version = "v1.62.1",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:9ddQBjfCyZPOHPUiPxpYESBLc+T8P3E+Vo4IbKZgFWg=",
        version = "v1.34.1",
    )
    go_repository(
        name = "org_golang_x_arch",
        importpath = "golang.org/x/arch",
        sum = "h1:3wRIsP3pM4yUptoR96otTUOXI367OS0+c9eeRi9doIc=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:GBDwsMXVQi34v5CCYUm2jkJvu4cbtru2U4TN2PSyQnw=",
        version = "v0.28.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:GoHiUyI/Tp2nVkLI2mCxVkOjsbSXD66ic0XW0js0R9g=",
        version = "v0.0.0-20230905200255-921286631fa9",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:zY54UmvipHiNd+pm+m0x9KhZ9hl1/7QNMyxXbc6ICqA=",
        version = "v0.17.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:AcW1SDZMkb8IpzCdQUaIq2sP4sZ4zw+55h6ynffypl4=",
        version = "v0.30.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        importpath = "golang.org/x/oauth2",
        sum = "h1:09qnuIAgzdx1XplqJvW6CQqMCtGZykZWcXzPMPUusvI=",
        version = "v0.18.0",
    )
    go_repository(
        name = "org_golang_x_sync",
        importpath = "golang.org/x/sync",
        sum = "h1:3NFvSEYkUoMifnESzZl15y791HH1qU2xm6eCJU5ZPXQ=",
        version = "v0.8.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:KHjCJyddX0LoSTb3J+vWpupP9p0oznkqVk/IfjymZbo=",
        version = "v0.26.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:WtHI/ltw4NvSUig5KARz9h521QvRC8RmF/cuYqifU24=",
        version = "v0.25.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:kTxAhCbGbxhK0IwgSKiMO5awPoDQ0RpfiVYBfK860YM=",
        version = "v0.19.0",
    )
    go_repository(
        name = "org_golang_x_time",
        importpath = "golang.org/x/time",
        sum = "h1:o7cqy6amK/52YcAKIPlM3a+Fpj35zvRj2TP+e1xFSfk=",
        version = "v0.5.0",
    )
    go_repository(
        name = "org_golang_x_tools",
        importpath = "golang.org/x/tools",
        sum = "h1:vU5i/LfpvrRCpgM/VPfJLg5KjxD3E+hfT1SH+d9zLwg=",
        version = "v0.21.1-0.20240508182429-e35e4ccd0d2d",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:H2TDz8ibqkAF6YGhCdN3jS9O0/s90v0rJh3X/OLHEUk=",
        version = "v0.0.0-20220907171357-04be3eba64a2",
    )
    go_repository(
        name = "org_uber_go_atomic",
        importpath = "go.uber.org/atomic",
        sum = "h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=",
        version = "v1.9.0",
    )
    go_repository(
        name = "org_uber_go_multierr",
        importpath = "go.uber.org/multierr",
        sum = "h1:7fIwc/ZtS0q++VgcfqFDxSBZVv/Xo49/SYnDFupUwlI=",
        version = "v1.9.0",
    )
    go_repository(
        name = "org_uber_go_zap",
        importpath = "go.uber.org/zap",
        sum = "h1:WefMeulhovoZ2sYXz7st6K0sLj7bBhpiFaud4r4zST8=",
        version = "v1.21.0",
    )
