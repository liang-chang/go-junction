package action

const COMMON_TITLE = `
Action  : {{.Action}}
BackupLinkFolder  : {{.Config.BackupLinkFolder}}
ClearBackupFolder : {{.Config.ClearBackupFolder}}
CreateTargetFolder: {{.Config.CreateTargetFolder}}
PathAlias :{{range $name,$path := .PathAlias }}
           {{$name}} => {{$path -}}{{end -}}
`

const list_template = `
{{- range $si, $symbo := .Symbolic }}
---------------symbolic---------------
skip   : {{$symbo.Skip}}
target : {{$symbo.Target}}
{{range $ii,$linkConf := $symbo.LinkConfig}}
link   : {{index $symbo.Link $ii -}}
	{{if gt (len $linkConf.MatchFolder) 0}}
	{{- range $mi,$mfolder := $linkConf.MatchFolder -}}
{{if eq $mi 0}}
match  : {{$mfolder}}{{else}}	     {{$mfolder}}{{end }}
{{end}}{{else}}
match  : No directory match !
{{end}}{{end}}{{end}}
`

const check_template = `
{{- range $si, $symbo := .Symbolic }}
---------------symbolic---------------
skip   : {{$symbo.Skip}}
target : {{$symbo.Target}}
{{range $ii,$linkConf := $symbo.LinkConfig }}
link   : {{index $symbo.Link $ii -}}
{{- range $mi,$mfolder := $linkConf.MatchFolder -}}
{{if eq $mi 0}}
match  : {{$mfolder -}}{{else}}         {{$mfolder -}}{{end}}
{{end}}{{end}}{{end}}
-------------------------------
Warning : {{.WarnCount}}
Error   : {{.ErrorCount}}
`