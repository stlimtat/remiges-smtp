load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_dependencies():
    go_repository(
        name = "com_github_armon_go_metrics",
        importpath = "github.com/armon/go-metrics",
        sum = "h1:hR91U9KYmb6bLBYLQjyM+3j+rcd/UhE+G78SFnF8gJA=",
        version = "v0.4.1",
    )
    go_repository(
        name = "com_github_beorn7_perks",
        importpath = "github.com/beorn7/perks",
        sum = "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_bsm_ginkgo_v2",
        importpath = "github.com/bsm/ginkgo/v2",
        sum = "h1:Ny8MWAHyOepLGlLKYmXG4IEkioBysk6GpaRTLC8zwWs=",
        version = "v2.12.0",
    )
    go_repository(
        name = "com_github_bsm_gomega",
        importpath = "github.com/bsm/gomega",
        sum = "h1:yeMWxP2pV2fG3FgAODIY8EiRE3dy0aeFYt4l7wh6yKA=",
        version = "v1.27.10",
    )
    go_repository(
        name = "com_github_bytedance_sonic",
        importpath = "github.com/bytedance/sonic",
        sum = "h1:9Csb3c9ZJhfUWeMtpCDCq6BUoH5ogfDFLUgQ/jG+R0k=",
        version = "v1.12.4",
    )
    go_repository(
        name = "com_github_bytedance_sonic_loader",
        importpath = "github.com/bytedance/sonic/loader",
        sum = "h1:1GgorWTqf12TA8mma4DDSbaQigE2wOgQo7iCjjJv3+E=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:DC2CZ1Ep5Y4k3ZQ899DldepgrayRUGE6BBZ/cd9Cj44=",
        version = "v2.2.0",
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
        sum = "h1:RrqgGjYQKalulkV8NGVIfkXQf6YYmOyiJKk8iXXhfZs=",
        version = "v22.5.0",
    )
    go_repository(
        name = "com_github_cpuguy83_go_md2man_v2",
        importpath = "github.com/cpuguy83/go-md2man/v2",
        sum = "h1:wfIWP927BUkWJb2NmU/kNDYIBTh/ziUX91+lVfRxZq4=",
        version = "v2.0.4",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:U9qPSI2PIWSS1VwoXQT9A3Wy9MM3WgvqSxFWenqJduM=",
        version = "v1.1.2-0.20180830191138-d8f796af33cc",
    )
    go_repository(
        name = "com_github_dgryski_go_rendezvous",
        importpath = "github.com/dgryski/go-rendezvous",
        sum = "h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=",
        version = "v0.0.0-20200823014737-9f7001d12a5f",
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
        sum = "h1:dAwr6QBTBZIkG8roQaJjGof0pp0EeF+tNV7YBP3F/8M=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_gabriel_vasile_mimetype",
        importpath = "github.com/gabriel-vasile/mimetype",
        sum = "h1:SKFKl7kD0RiPdbht0s7hFtjl489WcQ1VyPW8ZzUMYCA=",
        version = "v1.4.7",
    )
    go_repository(
        name = "com_github_gin_contrib_pprof",
        importpath = "github.com/gin-contrib/pprof",
        sum = "h1:Mzy+3HHtHbtwr4VewBTXZp/hR7pS6ZuZkueBIrQiLL4=",
        version = "v1.5.1",
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
        name = "com_github_go_logr_logr",
        importpath = "github.com/go-logr/logr",
        sum = "h1:6pFjapn8bFcIbiKo3XT4j/BhANplGihG6tvd+8rYgrY=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )
    go_repository(
        name = "com_github_go_mods_zerolog_gin",
        importpath = "github.com/go-mods/zerolog-gin",
        sum = "h1:QmOOU2pPkHuV4oPDaceelEouS6bwrOXNsIZdlpR3Ylg=",
        version = "v0.2.0",
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
        sum = "h1:/PwmTwZhS0dPkav3cdK9kV1FsAmrL8sThn8IHr/sO+o=",
        version = "v10.23.0",
    )
    go_repository(
        name = "com_github_go_redis_redismock_v9",
        importpath = "github.com/go-redis/redismock/v9",
        sum = "h1:ZrMYQeKPECZPjOj5u9eyOjg8Nnb0BS9lkVIZ6IpsKLw=",
        version = "v9.2.0",
    )
    go_repository(
        name = "com_github_go_task_slim_sprig_v3",
        importpath = "github.com/go-task/slim-sprig/v3",
        sum = "h1:sUs3vkvUymDpBKi3qH1YSqBQk9+9D/8M2mN1vB6EwHI=",
        version = "v3.0.0",
    )
    go_repository(
        name = "com_github_goccy_go_json",
        importpath = "github.com/goccy/go-json",
        sum = "h1:KZ5WoDbxAIgm2HNbYckL0se1fHD6rz5j4ywS6ebzDqA=",
        version = "v0.10.3",
    )
    go_repository(
        name = "com_github_godbus_dbus_v5",
        importpath = "github.com/godbus/dbus/v5",
        sum = "h1:9349emZab16e7zQvpmsbtjc18ykshndd8y2PG3sgJbA=",
        version = "v5.0.4",
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
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:ofyhxvXcZhMsU5ulbFiLKl/XBFqE1GSq7atu8tAmTRI=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_google_gofuzz",
        importpath = "github.com/google/gofuzz",
        sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_google_pprof",
        importpath = "github.com/google/pprof",
        sum = "h1:5iH8iuqE5apketRbSFBy+X1V0o+l+8NF1avt4HWl7cA=",
        version = "v0.0.0-20240827171923-fa2c70bbbfe5",
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
        sum = "h1:MtMxsa51/r9yyhkyLsVeVt0B+BGQZzpQiTQ4eHZ8bc4=",
        version = "v1.4.0",
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
        name = "com_github_inconshreveable_mousetrap",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:wN+x4NVGpMsO7ErUn/mUI3vEoE6Jt13X2s0bqwp9tc8=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_json_iterator_go",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_klauspost_compress",
        importpath = "github.com/klauspost/compress",
        sum = "h1:RlWWUY/Dr4fL8qk9YG7DTZ7PDgME2V4csBXA8L/ixi4=",
        version = "v1.17.2",
    )
    go_repository(
        name = "com_github_klauspost_cpuid_v2",
        importpath = "github.com/klauspost/cpuid/v2",
        sum = "h1:66ze0taIn2H33fBvCkXuv9BmCwDfafmiIVpKV9kKGuY=",
        version = "v2.2.9",
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
        name = "com_github_matttproud_golang_protobuf_extensions_v2",
        importpath = "github.com/matttproud/golang_protobuf_extensions/v2",
        sum = "h1:jWpvCLoY8Z/e3VKvlsiIGKtc+UG6U5vzxaoagmhXfyg=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_mcnijman_go_emailaddress",
        importpath = "github.com/mcnijman/go-emailaddress",
        sum = "h1:AGhgVDG3tCDaL0/Vc6erlPQjDuDN3dAT7rRdgFtetr0=",
        version = "v1.1.1",
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
        name = "com_github_mjl_adns",
        importpath = "github.com/mjl-/adns",
        sum = "h1:RxsPh1BrSAHeRDoEZuykc7lcnAQTC1akeVNOnVu37Ng=",
        version = "v0.0.0-20240624081245-1247f01fcf54",
    )
    go_repository(
        name = "com_github_mjl_autocert",
        importpath = "github.com/mjl-/autocert",
        sum = "h1:s6ay4bh4tmpPLdxjyeWG45mcwHfEluBMuGPkqxHWUJ4=",
        version = "v0.0.0-20231214125928-31b7400acb05",
    )
    go_repository(
        name = "com_github_mjl_bstore",
        importpath = "github.com/mjl-/bstore",
        sum = "h1:ntlu9MkfCkpm2XfBY4+Ws4KK9YzXzewr3+lCueFB+9c=",
        version = "v0.0.6",
    )
    go_repository(
        name = "com_github_mjl_mox",
        importpath = "github.com/mjl-/mox",
        sum = "h1:pa/ThG9GFl4dR2TyojjJe1v0OdCYmjwxCGCHy0B559o=",
        version = "v0.0.13",
    )
    go_repository(
        name = "com_github_mjl_sconf",
        importpath = "github.com/mjl-/sconf",
        sum = "h1:bdBcSFZCDFMm/UdBsgNCsjkYmKrSgYwp7rAOoufwHe4=",
        version = "v0.0.7",
    )
    go_repository(
        name = "com_github_mjl_sherpa",
        importpath = "github.com/mjl-/sherpa",
        sum = "h1:C5F8XQdV5nCuS4fvB+ye/ziUQrajEhOoj/t2w5T14BY=",
        version = "v0.6.7",
    )
    go_repository(
        name = "com_github_mjl_sherpadoc",
        importpath = "github.com/mjl-/sherpadoc",
        sum = "h1:BdlFNXfnTaA7qO54kof4xpNFJxYBTY0cIObRk7QAP6M=",
        version = "v0.0.16",
    )
    go_repository(
        name = "com_github_mjl_sherpaprom",
        importpath = "github.com/mjl-/sherpaprom",
        sum = "h1:1dlbkScsNafM5jURI44uiWrZMSwfZtcOFEEq7vx2C1Y=",
        version = "v0.0.2",
    )
    go_repository(
        name = "com_github_mjl_sherpats",
        importpath = "github.com/mjl-/sherpats",
        sum = "h1:2lSoJbb+jkjLOdlvoMxItq0QQrrnkH+rnm3PMRfpbmA=",
        version = "v0.0.6",
    )
    go_repository(
        name = "com_github_mjl_xfmt",
        importpath = "github.com/mjl-/xfmt",
        sum = "h1:6dLgd6U3bmDJKtTxsaSYYyMaORoO4hKBAJo4XKkPRko=",
        version = "v0.0.2",
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
        name = "com_github_nxadm_tail",
        importpath = "github.com/nxadm/tail",
        sum = "h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=",
        version = "v1.4.8",
    )
    go_repository(
        name = "com_github_onsi_ginkgo",
        importpath = "github.com/onsi/ginkgo",
        sum = "h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=",
        version = "v1.16.5",
    )
    go_repository(
        name = "com_github_onsi_ginkgo_v2",
        importpath = "github.com/onsi/ginkgo/v2",
        sum = "h1:YlVIbqct+ZmnEph770q9Q7NVAz4wwIiVNahee6JyUzo=",
        version = "v2.20.1",
    )
    go_repository(
        name = "com_github_onsi_gomega",
        importpath = "github.com/onsi/gomega",
        sum = "h1:Cwbd75ZBPxFSuZ6T+rN/WCb/gOc6YgFBXLlZLhC7Ds4=",
        version = "v1.35.1",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:YmeHyLY8mFWbdkNWwpr+qIL2bEqT0o95WSdkNHvL12M=",
        version = "v2.2.3",
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
        name = "com_github_prometheus_client_golang",
        importpath = "github.com/prometheus/client_golang",
        sum = "h1:HzFfmkOzH5Q8L8G+kSJKUx5dtG87sewO+FoDDqP5Tbk=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:VQw1hfvPvk3Uv6Qf29VrPF32JB6rtbgI6cYPYQjL0Qw=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_prometheus_common",
        importpath = "github.com/prometheus/common",
        sum = "h1:2BGz0eBc2hdMDLnO/8n0jeB3oPrt2D08CekT0lneoxM=",
        version = "v0.45.0",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        importpath = "github.com/prometheus/procfs",
        sum = "h1:jluTpSng7V9hY0O2R9DzzJHYb2xULk9VTR1V1R/k6Bo=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_redis_go_redis_v9",
        importpath = "github.com/redis/go-redis/v9",
        sum = "h1:HhLSs+B6O021gwzl+locl0zEDnyNkxMtf/Z3NNBMa9E=",
        version = "v9.7.0",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:KvO1DLK/DRN07sQ1LQKScxyZJuNnedQ5/wKSR38lUII=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_github_rs_xid",
        importpath = "github.com/rs/xid",
        sum = "h1:mKX4bl4iPYJtEIxp6CYiUuLQ/8DYMoz0PUdtGgMFRVc=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_rs_zerolog",
        importpath = "github.com/rs/zerolog",
        sum = "h1:1cU2KZkvPxNyfgEmhHAz/1A9Bz+llsdYzklWFzgp0r8=",
        version = "v1.33.0",
    )
    go_repository(
        name = "com_github_russross_blackfriday_v2",
        importpath = "github.com/russross/blackfriday/v2",
        sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
        version = "v2.1.0",
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
        sum = "h1:ON7AQg37yzcRPU69mt7gwhFEBwxI6P9T4Qu3N51bwOk=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_sagikazarmark_slog_shim",
        importpath = "github.com/sagikazarmark/slog-shim",
        sum = "h1:diDBnUNK9N/354PgrxMywXnAwEr1QZcOr6gto+ugjYE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_samber_lo",
        importpath = "github.com/samber/lo",
        sum = "h1:z7RynLwP5nbyRscyvcD043DWYoOcYRv3mV8lBeqOCLc=",
        version = "v1.47.0",
    )
    go_repository(
        name = "com_github_samber_slog_common",
        importpath = "github.com/samber/slog-common",
        sum = "h1:jTqqLBgoJshpoxlPSGiypyOanjH6tY+i9bwyYmIbjhI=",
        version = "v0.17.1",
    )
    go_repository(
        name = "com_github_samber_slog_zerolog_v2",
        importpath = "github.com/samber/slog-zerolog/v2",
        sum = "h1:M6Ui44B76xdaOC2ZG7MJA9fwjLIbyeuQpjah2uaItyU=",
        version = "v2.7.1",
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
        sum = "h1:ntdiHjuueXFgm5nzDRdOS4yfT43P5Fnud6DH50rz/7w=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_spf13_cobra",
        importpath = "github.com/spf13/cobra",
        sum = "h1:e5/vxKd/rZsfSJMUX1agtjeTDf+qv1/JdBF8gg5k9ZM=",
        version = "v1.8.1",
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
        sum = "h1:Xv5erBjTwe/5IxqUQTdXv5kgmIvbHo3QQyRwhJsOfJA=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:9NlTDc1FTs4qu0DDq7AEtTPNw6SVm7uBMsUCUjABIf8=",
        version = "v1.6.0",
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
        name = "com_github_yuin_goldmark",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
        version = "v1.4.13",
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
        name = "in_gopkg_check_v1",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_ini_v1",
        importpath = "gopkg.in/ini.v1",
        sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
        version = "v1.67.0",
    )
    go_repository(
        name = "in_gopkg_tomb_v1",
        importpath = "gopkg.in/tomb.v1",
        sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
        version = "v1.0.0-20141024135613-dd632973f1e7",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_etcd_go_bbolt",
        importpath = "go.etcd.io/bbolt",
        sum = "h1:yGEzV1wPz2yVCLsD8ZAiGHhHVlczyC9d1rP43/VCRJ0=",
        version = "v1.3.11",
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
        name = "io_rsc_qr",
        importpath = "rsc.io/qr",
        sum = "h1:6vBLea5/NRMVTz8V66gipeLycZMl/+UlFmk8DvqQ6WY=",
        version = "v0.2.0",
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
        sum = "h1:8Ar7bF+apOIoThw1EdZl0p1oWvMqTHmpA2fRTyZO8io=",
        version = "v1.35.2",
    )
    go_repository(
        name = "org_golang_x_arch",
        importpath = "golang.org/x/arch",
        sum = "h1:UsYJhbzPYGsT0HbEdmYcqtCv8UNGvnaL561NnIUvaKg=",
        version = "v0.12.0",
    )
    go_repository(
        name = "org_golang_x_crypto",
        importpath = "golang.org/x/crypto",
        sum = "h1:IOBPskki6Lysi0lo9qQvbxiQ+FvsCC/YWOecCHAixus=",
        version = "v0.33.0",
    )
    go_repository(
        name = "org_golang_x_exp",
        importpath = "golang.org/x/exp",
        sum = "h1:XdNn9LlyWAhLVp6P/i8QYBW+hlyhrhei9uErw2B5GJo=",
        version = "v0.0.0-20241108190413-2d47ceb2692f",
    )
    go_repository(
        name = "org_golang_x_mod",
        importpath = "golang.org/x/mod",
        sum = "h1:D4nJWe9zXqHOmWqj4VMOJhvzj7bEZg4wEYa759z1pH4=",
        version = "v0.22.0",
    )
    go_repository(
        name = "org_golang_x_net",
        importpath = "golang.org/x/net",
        sum = "h1:T5GQRQb2y08kTAByq9L4/bz8cipCdA8FbRTXewonqY8=",
        version = "v0.35.0",
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
        sum = "h1:GGz8+XQP4FvTTrjZPzNKTMFtSXH80RAzG+5ghFPgK9w=",
        version = "v0.11.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:QjkSwP/36a20jFYWkSue1YwXzLmsV5Gfq7Eiy72C1uc=",
        version = "v0.30.0",
    )
    go_repository(
        name = "org_golang_x_term",
        importpath = "golang.org/x/term",
        sum = "h1:L6pJp37ocefwRRtYPKSWOWzOtWSxVajvz2ldH/xi3iU=",
        version = "v0.29.0",
    )
    go_repository(
        name = "org_golang_x_text",
        importpath = "golang.org/x/text",
        sum = "h1:bofq7m3/HAFvbF51jz3Q9wLg3jkvSPuiZu/pD1XwgtM=",
        version = "v0.22.0",
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
        sum = "h1:qEKojBykQkQ4EynWy4S8Weg69NumxKdn40Fce3uc/8o=",
        version = "v0.27.0",
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
        name = "org_uber_go_mock",
        importpath = "go.uber.org/mock",
        sum = "h1:KAMbZvZPyBPWgD14IrIQ38QCyjwpvVVV6K/bHl1IwQU=",
        version = "v0.5.0",
    )
    go_repository(
        name = "org_uber_go_multierr",
        importpath = "go.uber.org/multierr",
        sum = "h1:blXXJkSxSSfBVBlC76pxqeO+LN3aDfLQo+309xJstO0=",
        version = "v1.11.0",
    )
    go_repository(
        name = "org_uber_go_zap",
        importpath = "go.uber.org/zap",
        sum = "h1:WefMeulhovoZ2sYXz7st6K0sLj7bBhpiFaud4r4zST8=",
        version = "v1.21.0",
    )
