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
Common environment variables from ConfigMap
*/}}
{{- define "bus-booking.commonEnv" -}}
envFrom:
  - configMapRef:
      name: {{ include "bus-booking.fullname" . }}-config
{{- end }}

{{/*
Service-specific environment variables
*/}}
{{- define "bus-booking.serviceEnv" -}}
{{- $serviceName := .serviceName -}}
- name: SERVER_PORT
  value: {{ index .Values.env.ports (printf "%s_SERVICE_PORT" ($serviceName | upper)) | quote }}
{{- if ne $serviceName "gateway" }}
- name: DATABASE_NAME
  value: {{ index .Values.env.databases (printf "%s_DB_NAME" ($serviceName | upper)) | quote }}
- name: DATABASE_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "bus-booking.fullname" . }}-secret
      key: database-password
{{- end }}
- name: REDIS_DB
  value: {{ index .Values.env.redisDBs (printf "%s_REDIS_DB" ($serviceName | upper)) | quote }}
- name: JWT_SECRET_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "bus-booking.fullname" . }}-secret
      key: jwt-secret-key
- name: JWT_REFRESH_SECRET_KEY
  valueFrom:
    secretKeyRef:
      name: {{ include "bus-booking.fullname" . }}-secret
      key: jwt-refresh-secret-key
{{- if .Values.externalDatabase.redis.password }}
- name: REDIS_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ include "bus-booking.fullname" . }}-secret
      key: redis-password
{{- else }}
- name: REDIS_PASSWORD
  value: ""
{{- end }}
{{- end }}

