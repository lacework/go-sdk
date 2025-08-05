module github.com/lacework/go-sdk/v2

go 1.24.0

require (
	aead.dev/minisign v0.2.0
	cloud.google.com/go/compute v1.37.0
	cloud.google.com/go/compute/metadata v0.7.0
	cloud.google.com/go/oslogin v1.14.6
	cloud.google.com/go/resourcemanager v1.10.6
	github.com/AlecAivazis/survey/v2 v2.3.2
	github.com/BurntSushi/toml v1.3.2
	github.com/Masterminds/semver v1.5.0
	github.com/Netflix/go-expect v0.0.0-20200312175327-da48e75238e2
	github.com/abiosoft/colima v0.5.4
	github.com/aws/aws-sdk-go-v2 v1.37.2
	github.com/aws/aws-sdk-go-v2/config v1.30.3
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.222.0
	github.com/aws/aws-sdk-go-v2/service/ec2instanceconnect v1.28.2
	github.com/briandowns/spinner v1.17.0
	github.com/cenkalti/backoff/v4 v4.2.0
	github.com/fatih/color v1.13.0
	github.com/fatih/structs v1.1.0
	github.com/gammazero/workerpool v1.1.3
	github.com/hashicorp/go-version v1.6.0
	github.com/hashicorp/hc-install v0.4.0
	github.com/hashicorp/hcl/v2 v2.15.0
	github.com/hashicorp/terraform-exec v0.17.3
	github.com/hashicorp/terraform-json v0.14.0
	github.com/hinshun/vt10x v0.0.0-20180809195222-d55458df857c
	github.com/hokaccha/go-prettyjson v0.0.0-20190818114111-108c894c2c0e
	github.com/imdario/mergo v0.3.13
	github.com/kyokomi/emoji/v2 v2.2.12
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.5.0
	github.com/olekukonko/tablewriter v0.0.5
	github.com/peterbourgon/diskv/v3 v3.0.1
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.6
	github.com/spf13/viper v1.20.1
	github.com/stretchr/testify v1.10.0
	github.com/zclconf/go-cty v1.12.1
	go.uber.org/zap v1.24.0
	golang.org/x/crypto v0.38.0
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/text v0.25.0
	google.golang.org/grpc v1.72.1
	google.golang.org/protobuf v1.36.6
	gopkg.in/yaml.v2 v2.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.18.0
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.10.1
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization v1.0.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork v1.1.0
	github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions v1.3.0
	github.com/aws/aws-sdk-go-v2/credentials v1.18.3
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.48.4
	github.com/aws/aws-sdk-go-v2/service/eks v1.64.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.42.0
	github.com/aws/aws-sdk-go-v2/service/organizations v1.38.3
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.28.3
	github.com/aws/aws-sdk-go-v2/service/ssm v1.59.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.36.0
	github.com/aws/smithy-go v1.22.5
	github.com/gabriel-vasile/mimetype v1.4.8
	github.com/go-git/go-git/v5 v5.13.0
	github.com/go-resty/resty/v2 v2.11.0
	github.com/golang/protobuf v1.5.4
	github.com/google/uuid v1.6.0
	github.com/hashicorp/consul/sdk v0.13.1
	github.com/mattn/go-isatty v0.0.18
	github.com/mitchellh/hashstructure/v2 v2.0.2
	github.com/otiai10/copy v1.14.0
	github.com/pmezard/go-difflib v1.0.0
	github.com/spf13/cast v1.7.1
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56
	golang.org/x/oauth2 v0.30.0
	google.golang.org/api v0.235.0
)

require (
	cloud.google.com/go v0.120.0 // indirect
	cloud.google.com/go/auth v0.16.1 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/iam v1.5.2 // indirect
	cloud.google.com/go/longrunning v0.6.7 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.1 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.4.2 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/ProtonMail/go-crypto v1.1.3 // indirect
	github.com/agext/levenshtein v1.2.1 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.18.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.4.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.7.2 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.13.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.27.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.32.0 // indirect
	github.com/cloudflare/circl v1.3.7 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.5 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emirpasic/gods v1.18.1 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gammazero/deque v0.2.0 // indirect
	github.com/go-git/gcfg v1.5.1-0.20230307220236-3a3c6141e376 // indirect
	github.com/go-git/go-billy/v5 v5.6.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-viper/mapstructure/v2 v2.2.1 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.14.2 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51 // indirect
	github.com/kevinburke/ssh_config v1.2.0 // indirect
	github.com/kr/pty v1.1.8 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/mgutz/ansi v0.0.0-20170206155736-9520e82c474b // indirect
	github.com/mitchellh/go-wordwrap v0.0.0-20150314170334-ad45545899c7 // indirect
	github.com/pelletier/go-toml/v2 v2.2.3 // indirect
	github.com/pjbgf/sha1cd v0.3.0 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/sagikazarmark/locafero v0.7.0 // indirect
	github.com/sergi/go-diff v1.3.2-0.20230802210424-5b0b94c5c0d3 // indirect
	github.com/skeema/knownhosts v1.3.0 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/spf13/afero v1.12.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/xanzy/ssh-agent v0.3.3 // indirect
	go.opentelemetry.io/auto/sdk v1.1.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.60.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.60.0 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/metric v1.35.0 // indirect
	go.opentelemetry.io/otel/trace v1.35.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.9.0 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/term v0.32.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
	google.golang.org/genproto v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20250505200425-f936aa4a68b2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250512202823-5a2f75b736a9 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
)

replace github.com/kr/pty => github.com/creack/pty v1.1.7
