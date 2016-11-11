package action

const COMMON_TITLE = `
Action : {{.Action}}
BackupLinkFolder  : false
ClearBackupFolder : true
CreateTargetFolder: true
`

const list_template = `
---------------symbolic---------------
action : ignore
target : C:/useless/AAA
link   : bc@d:/|\d+$|/bin
match  :    d:/IntelliJ IDEA 15.0.2/bin
	     d:/apache-maven-3.3.3/bin
link   : v:/|cache$|
match  : v:/chrome_cache
	 v:/firefox_cache
	 v:/opera_cache
	 v:/safari_cache
`