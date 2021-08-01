workspace(name = "com_seankhliao_go_mono")

licenses(["notice"]) # MIT

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

git_repository(
    name = "bazel_gazelle",
    commit = "0709f360ae257d660816eac7ac5b8000d278bb7e",
    remote = "https://github.com/bazelbuild/bazel-gazelle",
    # master ~ 2021-08-01
    shallow_since = "1627660378 -0600",
)

http_archive(
    name = "com_google_protobuf",
    sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
    strip_prefix = "protobuf-3.14.0",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
    ],
)

git_repository(
    name = "io_bazel_rules_docker",
    commit = "3b0187852774ef0e424c92ec7b78a2a199a0503b",
    remote = "https://github.com/bazelbuild/rules_docker",
    # master ~ 2021-08-01
    shallow_since = "1627591998 -0400",
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "8e968b5fcea1d2d64071872b12737bbb5514524ee5f0a4f54f5920266c261acb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.28.0/rules_go-v0.28.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.28.0/rules_go-v0.28.0.zip",
    ],
)

http_archive(
    name = "rules_proto",
    sha256 = "602e7161d9195e50246177e7c55b2f39950a9cf7366f74ed5f22fd45750cd208",
    strip_prefix = "rules_proto-97d8af4dc474595af3900dd85cb3a29ad28cc313",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
        "https://github.com/bazelbuild/rules_proto/archive/97d8af4dc474595af3900dd85cb3a29ad28cc313.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_rules_dependencies()

go_register_toolchains(version = "1.16.6")

gazelle_dependencies()

go_repository(
    name = "co_honnef_go_tools",
    importpath = "honnef.co/go/tools",
    sum = "h1:/hemPrYIhOhy8zYrNj+069zDB68us2sMGsfkFJO0iZs=",
    version = "v0.0.0-20190523083050-ea95bdfd59fc",
)

go_repository(
    name = "com_github_andybalholm_cascadia",
    importpath = "github.com/andybalholm/cascadia",
    sum = "h1:BuuO6sSfQNFRu1LppgbD25Hr2vLYW25JvxHs5zzsLTo=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_antihax_optional",
    importpath = "github.com/antihax/optional",
    sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_benbjohnson_clock",
    importpath = "github.com/benbjohnson/clock",
    sum = "h1:Q92kusRqC1XV2MjkWETPvjJVqKetz1OzxZB7mHJLju8=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_bradleyfalzon_ghinstallation",
    importpath = "github.com/bradleyfalzon/ghinstallation",
    sum = "h1:HNMrvJcBerI5bz9RyVA5Fx6xQ5HfH9b0CAjcCaLRFts=",
    version = "v1.1.2-0.20210508092528-3386a2062511",
)

go_repository(
    name = "com_github_burntsushi_toml",
    importpath = "github.com/BurntSushi/toml",
    sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_cenkalti_backoff_v4",
    importpath = "github.com/cenkalti/backoff/v4",
    sum = "h1:G2HAfAmvm/GcKan2oOQpBXOd2tT2G57ZnZGWa1PxPBQ=",
    version = "v4.1.1",
)

go_repository(
    name = "com_github_census_instrumentation_opencensus_proto",
    importpath = "github.com/census-instrumentation/opencensus-proto",
    sum = "h1:glEXhBS5PSLLv4IXzLA5yPRVX4bilULVyxxbrfOtDAk=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_cheekybits_is",
    importpath = "github.com/cheekybits/is",
    sum = "h1:SKI1/fuSdodxmNNyVBR8d7X/HuLnRpvvFO0AgyQk764=",
    version = "v0.0.0-20150225183255-68e9c0620927",
)

go_repository(
    name = "com_github_client9_misspell",
    importpath = "github.com/client9/misspell",
    sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
    version = "v0.3.4",
)

go_repository(
    name = "com_github_cncf_udpa_go",
    importpath = "github.com/cncf/udpa/go",
    sum = "h1:cqQfy1jclcSy/FwLjemeg3SR1yaINm74aQyupQ0Bl8M=",
    version = "v0.0.0-20201120205902-5459f2c99403",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man_v2",
    importpath = "github.com/cpuguy83/go-md2man/v2",
    sum = "h1:U+s90UTSYgptZMwQh2aRr3LuazLJIa+Pg3Kc1ylSYVY=",
    version = "v2.0.0-20190314233015-f79a8a8ca69d",
)

go_repository(
    name = "com_github_davecgh_go_spew",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_dgrijalva_jwt_go_v4",
    importpath = "github.com/dgrijalva/jwt-go/v4",
    sum = "h1:CaO/zOnF8VvUfEbhRatPcwKVWamvbYd8tQGRWacE9kU=",
    version = "v4.0.0-preview1",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    importpath = "github.com/dustin/go-humanize",
    sum = "h1:VSnTsYCnlFHaM2/igO1h6X3HA71jcobQuxemgkq4zYo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_envoyproxy_go_control_plane",
    importpath = "github.com/envoyproxy/go-control-plane",
    sum = "h1:QyzYnTnPE15SQyUeqU6qLbWxMkwyAyu+vGksa0b7j00=",
    version = "v0.9.9-0.20210217033140-668b12f5399d",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    sum = "h1:EQciDnbrYxy13PgWoY8AqoxGiPrpgBZ1R8UNe3ddc+A=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_felixge_httpsnoop",
    importpath = "github.com/felixge/httpsnoop",
    sum = "h1:+nS9g82KMXccJ/wp0zyRW9ZBHFETmMGtkk+2CTTrW4o=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    importpath = "github.com/fsnotify/fsnotify",
    sum = "h1:hsms1Qyu0jgnwNXIxa+/V/PDsU6CfLf6CNO8H7IWoS4=",
    version = "v1.4.9",
)

go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_go_logr_logr",
    importpath = "github.com/go-logr/logr",
    sum = "h1:K7/B1jt6fIBQVd4Owv2MqGQClcgf0R266+7C/QjRcLc=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_golang_glog",
    importpath = "github.com/golang/glog",
    sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
    version = "v0.0.0-20160126235308-23def4e6c14b",
)

go_repository(
    name = "com_github_golang_mock",
    importpath = "github.com/golang/mock",
    sum = "h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    sum = "h1:ROPKBNFfQgOUMifHyP+KYbvpjbdoFNs+aK7DXlji0Tw=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_golang_snappy",
    importpath = "github.com/golang/snappy",
    sum = "h1:fHPg5GQYlCeLIPB9BZqMVR5nR9A+IM5zcgeTdjMYmLA=",
    version = "v0.0.3",
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:BKbKCqvP6I+rmFHt06ZmyQtvB8xAkWdhFyr0ZUNZcxQ=",
    version = "v0.5.6",
)

go_repository(
    name = "com_github_google_go_github_v35",
    importpath = "github.com/google/go-github/v35",
    sum = "h1:fU+WBzuukn0VssbayTT+Zo3/ESKX9JYWjbZTLOTEyho=",
    version = "v35.3.0",
)

go_repository(
    name = "com_github_google_go_querystring",
    importpath = "github.com/google/go-querystring",
    sum = "h1:Xkwi/a1rcvNg1PPYe5vI8GbeBY/jrVuDX5ASuANWTrk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_gofuzz",
    importpath = "github.com/google/gofuzz",
    sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_uuid",
    importpath = "github.com/google/uuid",
    sum = "h1:EVhdT+1Kseyi1/pUmXKaFxYsDNy9RQYkMWRH68J/W7Y=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_gorilla_feeds",
    importpath = "github.com/gorilla/feeds",
    sum = "h1:HwKXxqzcRNg9to+BbvJog4+f3s/xzvtZXICcQGutYfY=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    build_file_generation = "on",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
    sum = "h1:gmcG1KaJ57LophUzW0Hy8NmPhnMZb4M0+kPpLofRdBo=",
    version = "v1.16.0",
)

go_repository(
    name = "com_github_json_iterator_go",
    importpath = "github.com/json-iterator/go",
    sum = "h1:Kz6Cvnvv2wGdaG/V8yMvfkmNiXq9Ya2KUv4rouJJr68=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_kballard_go_shellquote",
    importpath = "github.com/kballard/go-shellquote",
    sum = "h1:Z9n2FFNUXsshfwJMBgNA0RU6/i7WVaAegv3PtuIHPMs=",
    version = "v0.0.0-20180428030007-95032a82bc51",
)

go_repository(
    name = "com_github_klauspost_compress",
    importpath = "github.com/klauspost/compress",
    sum = "h1:wXr2uRxZTJXHLly6qhJabee5JqIhTRoLBhDOA74hDEQ=",
    version = "v1.13.1",
)

go_repository(
    name = "com_github_kr_pretty",
    importpath = "github.com/kr/pretty",
    sum = "h1:Fmg33tUaq4/8ym9TJN1x7sLJnHVwhP33CNkpYV/7rwI=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_kr_pty",
    importpath = "github.com/kr/pty",
    sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_kr_text",
    importpath = "github.com/kr/text",
    sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_matryer_try",
    importpath = "github.com/matryer/try",
    sum = "h1:JAEbJn3j/FrhdWA9jW8B5ajsLIjeuEHLi8xE4fk997o=",
    version = "v0.0.0-20161228173917-9ac251b645a2",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    importpath = "github.com/mattn/go-isatty",
    sum = "h1:wuysRhFDzyxgEmMf5xjvJ2M9dZoWAXNNr5LSBS7uHXY=",
    version = "v0.0.12",
)

go_repository(
    name = "com_github_mattn_go_sqlite3",
    importpath = "github.com/mattn/go-sqlite3",
    sum = "h1:dNPt6NO46WmLVt2DLNpwczCmdV5boIZ6g/tlDrlRUbg=",
    version = "v1.14.6",
)

go_repository(
    name = "com_github_mmcdole_gofeed",
    importpath = "github.com/mmcdole/gofeed",
    sum = "h1:pdrvMb18jMSLidGp8j0pLvc9IGziX4vbmvVqmLH6z8o=",
    version = "v1.1.3",
)

go_repository(
    name = "com_github_mmcdole_goxpp",
    importpath = "github.com/mmcdole/goxpp",
    sum = "h1:sWGE2v+hO0Nd4yFU/S/mDBM5plIU8v/Qhfz41hkDIAI=",
    version = "v0.0.0-20181012175147-0068e33feabf",
)

go_repository(
    name = "com_github_modern_go_concurrent",
    importpath = "github.com/modern-go/concurrent",
    sum = "h1:ZqeYNhU3OHLH3mGKHDcjJRFFRrJa6eAM5H+CtDdOsPc=",
    version = "v0.0.0-20180228061459-e0a39a4cb421",
)

go_repository(
    name = "com_github_modern_go_reflect2",
    importpath = "github.com/modern-go/reflect2",
    sum = "h1:Esafd1046DLDQ0W1YjYsBW+p8U2u7vzgW2SQVmlNazg=",
    version = "v0.0.0-20180701023420-4b7aa43c6742",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_prometheus_client_model",
    importpath = "github.com/prometheus/client_model",
    sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
    version = "v0.0.0-20190812154241-14fe0d1b01d4",
)

go_repository(
    name = "com_github_puerkitobio_goquery",
    importpath = "github.com/PuerkitoBio/goquery",
    sum = "h1:PSPBGne8NIUWw+/7vFBV+kG2J/5MOjbzc7154OaKCSE=",
    version = "v1.5.1",
)

go_repository(
    name = "com_github_remyoudompheng_bigfft",
    importpath = "github.com/remyoudompheng/bigfft",
    sum = "h1:OdAsTTz6OkFY5QxjkYwrChwuRruF69c169dPK26NUlk=",
    version = "v0.0.0-20200410134404-eec4a21b6bb0",
)

go_repository(
    name = "com_github_rogpeppe_fastuuid",
    importpath = "github.com/rogpeppe/fastuuid",
    sum = "h1:Ppwyp6VYCF1nvBTXL3trRso7mXMlRrw9ooo375wvi2s=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_russross_blackfriday_v2",
    importpath = "github.com/russross/blackfriday/v2",
    sum = "h1:lPqVAte+HuHNfhJ/0LC98ESWRz8afy9tM/0RK8m9o+Q=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_shurcool_sanitized_anchor_name",
    importpath = "github.com/shurcooL/sanitized_anchor_name",
    sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_spf13_pflag",
    importpath = "github.com/spf13/pflag",
    sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
    version = "v1.0.5",
)

go_repository(
    name = "com_github_stretchr_objx",
    importpath = "github.com/stretchr/objx",
    sum = "h1:4G4v2dO3VZwixGIRoQ5Lfboy6nUhCyYzaqnIAPPhYs4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    importpath = "github.com/stretchr/testify",
    sum = "h1:nwc3DEeHmmLAfoZucVR881uASk0Mfjw8xYJ99tb5CcY=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_tdewolff_minify_v2",
    importpath = "github.com/tdewolff/minify/v2",
    sum = "h1:QdCQhIqCeBnSWU/4e0nJ8QJeQUNO5CC7jrH3xwETkVI=",
    version = "v2.9.19",
)

go_repository(
    name = "com_github_tdewolff_parse_v2",
    importpath = "github.com/tdewolff/parse/v2",
    sum = "h1:Kjaj3KQOx/4elIxlBSglus4E2oMfdROphvbq2b+OBZ0=",
    version = "v2.5.19",
)

go_repository(
    name = "com_github_tdewolff_test",
    importpath = "github.com/tdewolff/test",
    sum = "h1:76mzYJQ83Op284kMT+63iCNCI7NEERsIN8dLM+RiKr4=",
    version = "v1.0.6",
)

go_repository(
    name = "com_github_urfave_cli",
    importpath = "github.com/urfave/cli",
    sum = "h1:FpNT6zq26xNpHZy08emi755QwzLPs6Pukqjlc7RfOMU=",
    version = "v1.22.3",
)

go_repository(
    name = "com_github_yuin_goldmark",
    importpath = "github.com/yuin/goldmark",
    sum = "h1:OtISOGfH6sOWa1/qXqqAiOIAO6Z5J3AEAE18WAq6BiQ=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_yuin_goldmark_meta",
    importpath = "github.com/yuin/goldmark-meta",
    sum = "h1:ScsatUIT2gFS6azqzLGUjgOnELsBOxMXerM3ogdJhAM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_google_cloud_go",
    importpath = "cloud.google.com/go",
    sum = "h1:eOI3/cP2VTU6uZLDYAoic+eyzzB9YyGmJ7eIjl8rOPg=",
    version = "v0.34.0",
)

go_repository(
    name = "com_lukechampine_uint128",
    importpath = "lukechampine.com/uint128",
    sum = "h1:pnxCASz787iMf+02ssImqk6OLt+Z5QHMoZyUXR4z6JU=",
    version = "v1.1.1",
)

go_repository(
    name = "in_gopkg_check_v1",
    importpath = "gopkg.in/check.v1",
    sum = "h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=",
    version = "v0.0.0-20161208181325-20d25e280405",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:clyUAQHOM3G0M3f5vQj7LuJrETvjVot3Z5el9nffUtU=",
    version = "v2.3.0",
)

go_repository(
    name = "in_gopkg_yaml_v3",
    importpath = "gopkg.in/yaml.v3",
    sum = "h1:dUUwHk2QECo/6vqA44rthZ8ie2QXMNeKRTHCNY2nXvo=",
    version = "v3.0.0-20200313102051-9f266ea9e77c",
)

go_repository(
    name = "io_k8s_klog_v2",
    importpath = "k8s.io/klog/v2",
    sum = "h1:D7HV+n1V57XeZ0m6tdRkfknthUaM06VFbWldOFh8kzM=",
    version = "v2.9.0",
)

go_repository(
    name = "io_k8s_sigs_yaml",
    importpath = "sigs.k8s.io/yaml",
    sum = "h1:kr/MCeFWJWTwyaHoR9c8EjH9OumOmoF9YGiZd7lFm/Q=",
    version = "v1.2.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib",
    importpath = "go.opentelemetry.io/contrib",
    sum = "h1:RMJ6GlUVzLYp/zmItxTTdAmr1gnpO/HHMFmvjAhvJQM=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
    importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
    sum = "h1:G1vNyNfKknFvrKVC8ga8EYIECy0s5D/QPW4QPRSMhwc=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_contrib_propagators",
    importpath = "go.opentelemetry.io/contrib/propagators",
    sum = "h1:Wnio4Ffi9MoLrUkN/J5yqtHf2F9a7wa2VClFkKcQcOk=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel",
    importpath = "go.opentelemetry.io/otel",
    sum = "h1:4CeoX93DNTWt8awGK9JmNXzF9j7TyOu9upscEdtcdXc=",
    version = "v1.0.0-RC1",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlpmetric",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlpmetric",
    sum = "h1:Yi+SrcMhA21eMWN7rf5agAyiwbTM65TdS0xo3otuzAo=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlpmetric_otlpmetricgrpc",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc",
    sum = "h1:OB+C8gNWhwtvzm5zRUX+HyeCb6Y1Fuzo+NgGKtqFxLg=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace",
    sum = "h1:GHKxjc4EDldz8ScMDpiNwX4BAub6wGFUUo5Axm2BimU=",
    version = "v1.0.0-RC1",
)

go_repository(
    name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracegrpc",
    importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
    sum = "h1:ZOQXuxKJ9evGspu3LvbZxx3KOOQvKAPBJVMOfGf1cOM=",
    version = "v1.0.0-RC1",
)

go_repository(
    name = "io_opentelemetry_go_otel_internal_metric",
    importpath = "go.opentelemetry.io/otel/internal/metric",
    sum = "h1:gZlIBo5O51hZOOZz8vEcuRx/l5dnADadKfpT70AELoo=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_metric",
    importpath = "go.opentelemetry.io/otel/metric",
    sum = "h1:ZtcJlHqVE4l8Su0WOLOd9fEPheJuYEiQ0wr9wv2p25I=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_oteltest",
    importpath = "go.opentelemetry.io/otel/oteltest",
    sum = "h1:G685iP3XiskCwk/z0eIabL55XUl2gk0cljhGk9sB0Yk=",
    version = "v1.0.0-RC1",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk",
    importpath = "go.opentelemetry.io/otel/sdk",
    sum = "h1:Sy2VLOOg24bipyC29PhuMXYNJrLsxkie8hyI7kUlG9Q=",
    version = "v1.0.0-RC1",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk_export_metric",
    importpath = "go.opentelemetry.io/otel/sdk/export/metric",
    sum = "h1:4tSMVkDbvrowOeP/6rOfGABEWv5n+0gCfhI/TWleUvc=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_sdk_metric",
    importpath = "go.opentelemetry.io/otel/sdk/metric",
    sum = "h1:LNLUj35NNdEpyJQwj/htiEsfnY6GeTIwYHweCJNV+nc=",
    version = "v0.21.0",
)

go_repository(
    name = "io_opentelemetry_go_otel_trace",
    importpath = "go.opentelemetry.io/otel/trace",
    sum = "h1:jrjqKJZEibFrDz+umEASeU3LvdVyWKlnTh7XEfwrT58=",
    version = "v1.0.0-RC1",
)

go_repository(
    name = "io_opentelemetry_go_proto_otlp",
    importpath = "go.opentelemetry.io/proto/otlp",
    sum = "h1:C0g6TWmQYvjKRnljRULLWUVJGy8Uvu0NEL/5frY2/t4=",
    version = "v0.9.0",
)

go_repository(
    name = "org_golang_google_appengine",
    importpath = "google.golang.org/appengine",
    sum = "h1:/wp5JvzpHIxhs/dumFmF7BXTf3Z+dd4uXta4kVyO508=",
    version = "v1.4.0",
)

go_repository(
    name = "org_golang_google_genproto",
    importpath = "google.golang.org/genproto",
    sum = "h1:+kGHl1aib/qcwaRi1CbqBZ1rk19r85MNUf8HaBghugY=",
    version = "v0.0.0-20200526211855-cb27e3aa2013",
)

go_repository(
    name = "org_golang_google_grpc",
    importpath = "google.golang.org/grpc",
    sum = "h1:/9BgsAsa5nWe26HqOlvlgJnqBuktYOLCgjCPqsa56W0=",
    version = "v1.38.0",
)

go_repository(
    name = "org_golang_google_protobuf",
    importpath = "google.golang.org/protobuf",
    sum = "h1:bxAC2xTBsZGibn2RTntX0oH50xLsqy1OxA9tTL3p/lk=",
    version = "v1.26.0",
)

go_repository(
    name = "org_golang_x_crypto",
    importpath = "golang.org/x/crypto",
    sum = "h1:psW17arqaxU48Z5kZ0CQnkZWQJsqcURM6tKiBApRjXI=",
    version = "v0.0.0-20200622213623-75b288015ac9",
)

go_repository(
    name = "org_golang_x_exp",
    importpath = "golang.org/x/exp",
    sum = "h1:c2HOrn5iMezYjSlGPncknSEr/8x5LELb/ilJbXi9DEA=",
    version = "v0.0.0-20190121172915-509febef88a4",
)

go_repository(
    name = "org_golang_x_lint",
    importpath = "golang.org/x/lint",
    sum = "h1:XQyxROzUlZH+WIQwySDgnISgOivlhjIEwaQaJEJrrN0=",
    version = "v0.0.0-20190313153728-d0100b6bd8b3",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:RM4zey1++hCTbCVQfnWeKs9/IEsaBLA8vTkd0WVtmH4=",
    version = "v0.3.0",
)

go_repository(
    name = "org_golang_x_net",
    importpath = "golang.org/x/net",
    sum = "h1:IX6qOQeG5uLjB/hjjwjedwfjND0hgjPMMyO1RoIXQNI=",
    version = "v0.0.0-20201021035429-f5854403a974",
)

go_repository(
    name = "org_golang_x_oauth2",
    importpath = "golang.org/x/oauth2",
    sum = "h1:TzXSXBo42m9gQenoE3b9BGiEpg5IG2JkU5FkPIawgtw=",
    version = "v0.0.0-20200107190931-bf48bf16ab8d",
)

go_repository(
    name = "org_golang_x_sync",
    importpath = "golang.org/x/sync",
    sum = "h1:SQFwaSi55rU7vdNs9Yr0Z324VNlrF+0wMqRXT4St8ck=",
    version = "v0.0.0-20201020160332-67f06af15bc9",
)

go_repository(
    name = "org_golang_x_sys",
    importpath = "golang.org/x/sys",
    sum = "h1:gG67DSER+11cZvqIMb8S8bt0vZtiN6xWYARwirrOSfE=",
    version = "v0.0.0-20210510120138-977fb7262007",
)

go_repository(
    name = "org_golang_x_text",
    importpath = "golang.org/x/text",
    sum = "h1:cokOdA+Jmi5PJGXLlLllQSgYigAEfHXJAERHVMaCc2k=",
    version = "v0.3.3",
)

go_repository(
    name = "org_golang_x_tools",
    importpath = "golang.org/x/tools",
    sum = "h1:M8tBwCtWD/cZV9DZpFYRUgaymAYAr+aIUTWzDaM3uPs=",
    version = "v0.0.0-20201124115921-2c860bdd6e78",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:go1bK/D/BFZV2I8cIQd1NKEZ+0owSTG1fDTci4IqFcE=",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repository(
    name = "org_modernc_cc_v3",
    importpath = "modernc.org/cc/v3",
    sum = "h1:r63dgSzVzRxUpAJFPQWHy1QeZeY1ydNENUDaBx1GqYc=",
    version = "v3.33.6",
)

go_repository(
    name = "org_modernc_ccgo_v3",
    importpath = "modernc.org/ccgo/v3",
    sum = "h1:dEuUSf8WN51rDkprFuAqjfchKEzN0WttP/Py3enBwjk=",
    version = "v3.9.5",
)

go_repository(
    name = "org_modernc_httpfs",
    importpath = "modernc.org/httpfs",
    sum = "h1:AAgIpFZRXuYnkjftxTAZwMIiwEqAfk8aVB2/oA6nAeM=",
    version = "v1.0.6",
)

go_repository(
    name = "org_modernc_libc",
    importpath = "modernc.org/libc",
    sum = "h1:QUxZMs48Ahg2F7SN41aERvMfGLY2HU/ADnB9DC4Yts8=",
    version = "v1.9.11",
)

go_repository(
    name = "org_modernc_mathutil",
    importpath = "modernc.org/mathutil",
    sum = "h1:GCjoRaBew8ECCKINQA2nYjzvufFW9YiEuuB+rQ9bn2E=",
    version = "v1.4.0",
)

go_repository(
    name = "org_modernc_memory",
    importpath = "modernc.org/memory",
    sum = "h1:utMBrFcpnQDdNsmM6asmyH/FM9TqLPS7XF7otpJmrwM=",
    version = "v1.0.4",
)

go_repository(
    name = "org_modernc_opt",
    importpath = "modernc.org/opt",
    sum = "h1:/0RX92k9vwVeDXj+Xn23DKp2VJubL7k8qNffND6qn3A=",
    version = "v0.1.1",
)

go_repository(
    name = "org_modernc_sqlite",
    importpath = "modernc.org/sqlite",
    sum = "h1:ShWQpeD3ag/bmx6TqidBlIWonWmQaSQKls3aenCbt+w=",
    version = "v1.11.2",
)

go_repository(
    name = "org_modernc_strutil",
    importpath = "modernc.org/strutil",
    sum = "h1:xv+J1BXY3Opl2ALrBwyfEikFAj8pmqcpnfmuwUwcozs=",
    version = "v1.1.1",
)

go_repository(
    name = "org_modernc_tcl",
    importpath = "modernc.org/tcl",
    sum = "h1:N03RwthgTR/l/eQvz3UjfYnvVVj1G2sZqzFGfoD4HE4=",
    version = "v1.5.5",
)

go_repository(
    name = "org_modernc_token",
    importpath = "modernc.org/token",
    sum = "h1:a0jaWiNMDhDUtqOj09wvjWWAqd3q7WpBulmL9H2egsk=",
    version = "v1.0.0",
)

go_repository(
    name = "org_modernc_z",
    importpath = "modernc.org/z",
    sum = "h1:WyIDpEpAIx4Hel6q/Pcgj/VhaQV5XPJ2I6ryIYbjnpc=",
    version = "v1.0.1",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

rules_proto_dependencies()

rules_proto_toolchains()

load("@io_bazel_rules_docker//repositories:repositories.bzl", container_repositories = "repositories")

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//go:image.bzl", _go_image_repos = "repositories")

_go_image_repos()
