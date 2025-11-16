{{/*
Expand the name of the chart.
*/}}
{{- define "bus-booking.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "bus-booking.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "bus-booking.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "bus-booking.labels" -}}
helm.sh/chart: {{ include "bus-booking.chart" . }}
{{ include "bus-booking.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "bus-booking.selectorLabels" -}}
app.kubernetes.io/name: {{ include "bus-booking.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Service-specific labels
*/}}
{{- define "bus-booking.serviceLabels" -}}
{{ include "bus-booking.selectorLabels" . }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "bus-booking.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "bus-booking.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Environment variables for database
*/}}
{{- define "bus-booking.databaseEnv" -}}
- name: DB_HOST
  value: {{ .Values.env.database.DB_HOST | quote }}
- name: DB_PORT
  value: {{ .Values.env.database.DB_PORT | quote }}
- name: DB_USER
  value: {{ .Values.env.database.DB_USER | quote }}
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "bus-booking.fullname" . }}-db-secret
      key: password
- name: DB_NAME
  value: {{ .Values.env.database.DB_NAME | quote }}
{{- end }}

{{/*
Environment variables for Redis
*/}}
{{- define "bus-booking.redisEnv" -}}
- name: REDIS_HOST
  value: {{ .Values.env.redis.REDIS_HOST | quote }}
- name: REDIS_PORT
  value: {{ .Values.env.redis.REDIS_PORT | quote }}
{{- end }}

{{/*
Environment variables for global config
*/}}
{{- define "bus-booking.globalEnv" -}}
- name: JWT_SECRET
  valueFrom:
    secretKeyRef:
      name: {{ include "bus-booking.fullname" . }}-secret
      key: jwt-secret
- name: LOG_LEVEL
  value: {{ .Values.env.global.LOG_LEVEL | quote }}
{{- end }}