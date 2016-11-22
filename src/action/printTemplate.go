package action

const COMMON_TITLE = `
Action  : {{.Action}}
BackupLinkFolder  : {{.Config.BackupLinkFolder}}
ClearBackupFolder : {{.Config.ClearBackupFolder}}
CreateTargetFolder: {{.Config.CreateTargetFolder}}+
{{- $index:=0}}
PathAlias :{{range $name,$path := .PathAlias }}{{if eq $index 0}}{{$index := 1}}{{$index}}{{ $name}} => {{$path }}{{else}}
{{$name}} => {{$path }}
{{end}}{{end}}
`

const list_template = `
{{- range $si, $symbo := .Symbolic -}}
---------------symbolic---------------
action : {{$symbo.Action}}
target : {{$symbo.Target}}
    {{range $ii,$linkConf := $symbo.LinkConfig}}
link   : {{index $symbo.Link $ii -}}
	 {{- range $mi,$mfolder := $linkConf.MatchFolder -}}
{{if eq $mi 0}}
match  : {{$mfolder -}}{{else}}         {{$mfolder -}}
{{end}}
{{end}}{{end}}{{end}}
`

const check_template = `
{{- range $si, $symbo := .Symbolic }}
---------------symbolic---------------
action : {{$symbo.Action}}
target : {{$symbo.Target}}
    {{range $ii,$linkConf := $symbo.LinkConfig }}
link   : {{index $symbo.Link $ii }}
	 {{- range $mi,$mfolder := $linkConf.MatchFolder -}}
{{if eq $mi 0}}
match  : {{$mfolder -}}{{else}}         {{$mfolder -}}
{{end}}
{{end}}{{end}}{{end}}
`