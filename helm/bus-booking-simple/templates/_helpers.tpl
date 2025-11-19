{{/* Generate basic labels */}}
{{- define "bus-booking.labels" -}}
app.kubernetes.io/name: {{ .Chart.Name }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
generate a fullname for the release
*/}}
{{- define "bus-booking.fullname" -}}
{{- printf "%s" .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- end }}

{{/*
return the name of the app (release)
*/}}
{{- define "bus-booking.name" -}}
{{ .Release.Name }}
{{- end }}