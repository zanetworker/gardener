{{- define "kube-controller-manager.featureGates" -}}
{{- if .Values.featureGates }}
- --feature-gates={{ range $feature, $enabled := .Values.featureGates }}{{ $feature }}={{ $enabled }},{{ end }}
{{- end }}
{{- if semverCompare "1.11.x" .Values.kubernetesVersion }}
- --feature-gates=ScheduleDaemonSetPods=true
{{- end }}
{{- end -}}

{{- define "kube-controller-manager.port" -}}
{{- if semverCompare ">= 1.13" .Values.kubernetesVersion -}}
10257
{{- else -}}
10252
{{- end -}}
{{- end -}}
