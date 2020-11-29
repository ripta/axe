module github.com/ripta/axe

go 1.15

require (
	github.com/gdamore/tcell/v2 v2.0.1-0.20201017141208-acf90d56d591
	github.com/jonboulle/clockwork v0.1.0
	github.com/mattn/go-runewidth v0.0.9
	github.com/spf13/cobra v1.1.1
	golang.org/x/sys v0.0.0-20201126233918-771906719818 // indirect
	golang.org/x/text v0.3.4 // indirect
	k8s.io/api v0.19.3
	k8s.io/apimachinery v0.19.3
	k8s.io/cli-runtime v0.19.3
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/klog v1.0.0
	k8s.io/klog/v2 v2.4.0 // indirect
	k8s.io/kube-openapi v0.0.0-20201106092651-fd18780d1fff // indirect
	k8s.io/kubectl v0.19.3
	k8s.io/utils v0.0.0-20201104234853-8146046b121e // indirect
)

replace k8s.io/api => k8s.io/api v0.16.7

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.7

replace k8s.io/apimachinery => k8s.io/apimachinery v0.16.7

replace k8s.io/apiserver => k8s.io/apiserver v0.16.7

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.7

replace k8s.io/client-go => k8s.io/client-go v0.16.7

replace k8s.io/code-generator => k8s.io/code-generator v0.16.7

replace k8s.io/component-base => k8s.io/component-base v0.16.7

replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20200204173128-addea2498afe

replace k8s.io/kubectl => k8s.io/kubectl v0.16.7
