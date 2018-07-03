#!{{.BusyboxShell}}
export PATH="{{.Path}}"
case "$1" in
  -h|--help|help|-V|--version) exec {{.SupBin}} "$@";;
  -*) exec {{sup_bin}} run {{.PrimarySvcIdent}} "$@";;
  *) exec {{.SupBin}} "$@";;
esac
