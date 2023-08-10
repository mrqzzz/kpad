// Demo code for the TextArea primitive.
package main

import (
	"strings"
)

func main() {
	var corporate = `apiVersion: apps/v1
kind: Deployment
metadata:
  annotations:
    deployment.kubernetes.io/revision: "2"
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"apps/v1","kind":"Deployment","metadata":{"annotations":{"qlik.encryption.service":"{ \"serviceS2SName\": \"reporting\", \"serviceDisplayName\" : \"Reporting Service\" , \"serviceEncryptionContext\" : \"reporting\" }"},"labels":{"app":"reporting","app.kubernetes.io/component":"reporting","app.kubernetes.io/instance":"reporting","app.kubernetes.io/name":"reporting","app.kubernetes.io/version":"9.151.1","chart":"reporting-9.151.1","domain":"reporting","heritage":"Helm","qlik.encryption":"true","release":"reporting"},"name":"reporting","namespace":"default"},"spec":{"replicas":1,"selector":{"matchLabels":{"app":"reporting","release":"reporting"}},"template":{"metadata":{"annotations":{"checksum/configs":"8e03b70100b236c2ebb6c062eefeb6fd511357920123c86906f22388c22e2b24","checksum/secrets":"d7e48a4307a3af13c81636c9b463444f60fc875fabf471a1a05a5a61a9c3f7e7"},"labels":{"app":"reporting","app.kubernetes.io/component":"reporting","app.kubernetes.io/instance":"reporting","app.kubernetes.io/name":"reporting","app.kubernetes.io/version":"9.151.1","chart":"reporting-9.151.1","domain":"reporting","heritage":"Helm","messaging-nats-client":"true","release":"reporting"}},"spec":{"containers":[{"env":[{"name":"DT_RELEASE_VERSION","valueFrom":{"fieldRef":{"fieldPath":"metadata.labels['app.kubernetes.io/version']"}}},{"name":"NATS_CLIENT_ID","valueFrom":{"fieldRef":{"fieldPath":"metadata.name"}}},{"name":"PROXY_PRIVATE_KEY_FILE","value":"/run/secrets/qlik.com/reporting/proxyPrivateKey"},{"name":"PROXY_PUBLIC_KEY_FILE","value":"/run/secrets/qlik.com/reporting/proxyPublicKey"},{"name":"CUSTOMIZATION_URI","valueFrom":{"configMapKeyRef":{"key":"customizationUri","name":"reporting-configs"}}},{"name":"DATAPREP_URI","valueFrom":{"configMapKeyRef":{"key":"dataprepUri","name":"reporting-configs"}}},{"name":"EDGE_AUTH_URI","valueFrom":{"configMapKeyRef":{"key":"edgeAuthUri","name":"reporting-configs"}}},{"name":"ENGINE_REST_API_VERSION","valueFrom":{"configMapKeyRef":{"key":"engineRestApiVersion","name":"reporting-configs"}}},{"name":"ENGINE_URI","valueFrom":{"configMapKeyRef":{"key":"engineUri","name":"reporting-configs"}}},{"name":"FEATURE_FLAGS_URI","valueFrom":{"configMapKeyRef":{"key":"featureFlagsUri","name":"reporting-configs"}}},{"name":"LICENSES_URI","valueFrom":{"configMapKeyRef":{"key":"licensesUri","name":"reporting-configs"}}},{"name":"LOCALE_URI","valueFrom":{"configMapKeyRef":{"key":"localeUri","name":"reporting-configs"}}},{"name":"MAIN_WEB_CONTAINER_URI","valueFrom":{"configMapKeyRef":{"key":"mainWebContainerUri","name":"reporting-configs"}}},{"name":"NLBROKER_URI","valueFrom":{"configMapKeyRef":{"key":"nlbrokerUri","name":"reporting-configs"}}},{"name":"NOTES_URI","valueFrom":{"configMapKeyRef":{"key":"notesUri","name":"reporting-configs"}}},{"name":"PROXY_URI","valueFrom":{"configMapKeyRef":{"key":"proxyUri","name":"reporting-configs"}}},{"name":"RESOURCE_LIBRARY_URI","valueFrom":{"configMapKeyRef":{"key":"resourceLibraryUri","name":"reporting-configs"}}},{"name":"SENSE_CLIENT_URI","valueFrom":{"configMapKeyRef":{"key":"senseClientUri","name":"reporting-configs"}}},{"name":"SPACES_URI","valueFrom":{"configMapKeyRef":{"key":"spacesUri","name":"reporting-configs"}}},{"name":"TEMPORARY_CONTENTS_URI","valueFrom":{"configMapKeyRef":{"key":"temporaryContentsUri","name":"reporting-configs"}}},{"name":"TENANTS_URI","valueFrom":{"configMapKeyRef":{"key":"tenantsUri","name":"reporting-configs"}}},{"name":"USERS_URI","valueFrom":{"configMapKeyRef":{"key":"usersUri","name":"reporting-configs"}}},{"name":"PROXY_METRICS_PORT","value":"8484"},{"name":"PROXY_PORT","value":"8001"},{"name":"JAEGER_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}},{"name":"OTLP_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}}],"image":"ghcr.io/qlik-trial/reporting-proxy:3.10.0","imagePullPolicy":"IfNotPresent","livenessProbe":{"httpGet":{"path":"/health","port":"metric"},"initialDelaySeconds":10},"name":"rpr","ports":[{"containerPort":8484,"name":"metric","protocol":"TCP"},{"containerPort":8001}],"readinessProbe":{"httpGet":{"path":"/ready","port":"metric"},"initialDelaySeconds":10},"resources":{"limits":{"memory":"300Mi"},"requests":{"cpu":"200m","memory":"300Mi"}},"volumeMounts":[{"mountPath":"/run/secrets/qlik.com/reporting","name":"reporting-secrets","readOnly":true}]},{"args":["-l=$(LOG_LEVEL)","-http-server-enable=true","-http-server-host=0.0.0.0","-http-server-port=9288","-temp-dir=/home/wuser/rwr/tmp"],"env":[{"name":"DT_RELEASE_VERSION","valueFrom":{"fieldRef":{"fieldPath":"metadata.labels['app.kubernetes.io/version']"}}},{"name":"NATS_CLIENT_ID","valueFrom":{"fieldRef":{"fieldPath":"metadata.name"}}},{"name":"PROXY_PRIVATE_KEY_FILE","value":"/run/secrets/qlik.com/reporting/proxyPrivateKey"},{"name":"PROXY_PUBLIC_KEY_FILE","value":"/run/secrets/qlik.com/reporting/proxyPublicKey"},{"name":"HTTP_HOST","value":"0.0.0.0"},{"name":"HTTP_HOST_ENABLE","value":"true"},{"name":"HTTP_PORT","value":"9288"},{"name":"LOG_LEVEL","value":"INFO"},{"name":"JAEGER_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}},{"name":"OTLP_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}}],"image":"ghcr.io/qlik-trial/reporting-web-renderer:2.45.1","imagePullPolicy":"IfNotPresent","livenessProbe":{"httpGet":{"path":"/health","port":"http"}},"name":"rwr","ports":[{"containerPort":9288,"name":"http","protocol":"TCP"}],"readinessProbe":{"httpGet":{"path":"/ready","port":"http"}},"resources":{"limits":{"memory":"1.5Gi"},"requests":{"cpu":"1000m","memory":"1.5Gi"}},"securityContext":{"allowPrivilegeEscalation":false,"capabilities":{"drop":["ALL"]},"privileged":false,"readOnlyRootFilesystem":true,"runAsGroup":65200,"runAsNonRoot":true,"runAsUser":65200,"seccompProfile":{"type":"RuntimeDefault"}},"volumeMounts":[{"mountPath":"/home/wuser/rwr","name":"rwr-home"},{"mountPath":"/tmp","name":"rwr-tmp"},{"mountPath":"/home/wuser/.pki","name":"pki"},{"mountPath":"/run/secrets/qlik.com/reporting","name":"reporting-secrets","readOnly":true}]},{"args":["-l=\"$(LOG_LEVEL)\""],"env":[{"name":"DT_RELEASE_VERSION","valueFrom":{"fieldRef":{"fieldPath":"metadata.labels['app.kubernetes.io/version']"}}},{"name":"NATS_CLIENT_ID","valueFrom":{"fieldRef":{"fieldPath":"metadata.name"}}},{"name":"REDIS_ACL_PASSWORD_FILE","value":"/run/secrets/qlik.com/reporting/redisAclPassword"},{"name":"TOKEN_AUTH_PRIVATE_KEY_FILE","value":"/run/secrets/qlik.com/reporting/tokenAuthPrivateKey"},{"name":"TOKEN_AUTH_PRIVATE_KEY_ID_FILE","value":"/run/secrets/qlik.com/reporting/tokenAuthPrivateKeyId"},{"name":"AWSS3_ACCESSKEY","valueFrom":{"configMapKeyRef":{"key":"awss3Accesskey","name":"reporting-configs"}}},{"name":"AWSS3_BUCKET","valueFrom":{"configMapKeyRef":{"key":"awss3Bucket","name":"reporting-configs"}}},{"name":"AWSS3_REGION","valueFrom":{"configMapKeyRef":{"key":"awss3Region","name":"reporting-configs"}}},{"name":"AWSS3_SECRETKEY","valueFrom":{"configMapKeyRef":{"key":"awss3Secretkey","name":"reporting-configs"}}},{"name":"AWSS3_URL","valueFrom":{"configMapKeyRef":{"key":"awss3Url","name":"reporting-configs"}}},{"name":"EDGE_AUTH_URI","valueFrom":{"configMapKeyRef":{"key":"edgeAuthUri","name":"reporting-configs"}}},{"name":"ENCRYPTION_URI","valueFrom":{"configMapKeyRef":{"key":"encryptionUri","name":"reporting-configs"}}},{"name":"LOG_LEVEL","valueFrom":{"configMapKeyRef":{"key":"logLevel","name":"reporting-configs"}}},{"name":"REDIS_ACL_USER","valueFrom":{"configMapKeyRef":{"key":"redisAclUser","name":"reporting-configs"}}},{"name":"REDIS_URI","valueFrom":{"configMapKeyRef":{"key":"redisUri","name":"reporting-configs"}}},{"name":"HTTP_PORT","value":"8384"},{"name":"PORT","value":"52052"},{"name":"TOKEN_AUTH_ENABLED","value":"true"},{"name":"JAEGER_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}},{"name":"OTLP_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}}],"image":"ghcr.io/qlik-trial/reporting-composer:3.23.0","imagePullPolicy":"IfNotPresent","livenessProbe":{"httpGet":{"path":"/health","port":"http"},"initialDelaySeconds":10},"name":"cmp","ports":[{"containerPort":8384,"name":"http","protocol":"TCP"}],"readinessProbe":{"httpGet":{"path":"/ready","port":"http"},"initialDelaySeconds":10},"resources":{"limits":{"memory":"1.5Gi"},"requests":{"cpu":"200m","memory":"500Mi"}},"volumeMounts":[{"mountPath":"/run/secrets/qlik.com/reporting","name":"reporting-secrets","readOnly":true}]},{"env":[{"name":"DT_RELEASE_VERSION","valueFrom":{"fieldRef":{"fieldPath":"metadata.labels['app.kubernetes.io/version']"}}},{"name":"NATS_CLIENT_ID","valueFrom":{"fieldRef":{"fieldPath":"metadata.name"}}},{"name":"REDIS_ACL_PASSWORD_FILE","value":"/run/secrets/qlik.com/reporting/redisAclPassword"},{"name":"TOKEN_AUTH_PRIVATE_KEY_FILE","value":"/run/secrets/qlik.com/reporting/tokenAuthPrivateKey"},{"name":"TOKEN_AUTH_PRIVATE_KEY_ID_FILE","value":"/run/secrets/qlik.com/reporting/tokenAuthPrivateKeyId"},{"name":"ACCESS_CONTROLS_URI","valueFrom":{"configMapKeyRef":{"key":"accessControlsUri","name":"reporting-configs"}}},{"name":"AWSS3_ACCESSKEY","valueFrom":{"configMapKeyRef":{"key":"awss3Accesskey","name":"reporting-configs"}}},{"name":"AWSS3_BUCKET","valueFrom":{"configMapKeyRef":{"key":"awss3Bucket","name":"reporting-configs"}}},{"name":"AWSS3_REGION","valueFrom":{"configMapKeyRef":{"key":"awss3Region","name":"reporting-configs"}}},{"name":"AWSS3_SECRETKEY","valueFrom":{"configMapKeyRef":{"key":"awss3Secretkey","name":"reporting-configs"}}},{"name":"AWSS3_URL","valueFrom":{"configMapKeyRef":{"key":"awss3Url","name":"reporting-configs"}}},{"name":"CONTENT_SIZE_THRESHOLD","valueFrom":{"configMapKeyRef":{"key":"contentSizeThreshold","name":"reporting-configs"}}},{"name":"CONTENT_STORAGE","valueFrom":{"configMapKeyRef":{"key":"contentStorage","name":"reporting-configs"}}},{"name":"DEFAULT_REPORT_DEADLINE","valueFrom":{"configMapKeyRef":{"key":"defaultReportDeadline","name":"reporting-configs"}}},{"name":"EDGE_AUTH_URI","valueFrom":{"configMapKeyRef":{"key":"edgeAuthUri","name":"reporting-configs"}}},{"name":"ENABLE_DIRECT_QUERY","valueFrom":{"configMapKeyRef":{"key":"enableDirectQuery","name":"reporting-configs"}}},{"name":"ENABLE_FAIR_QUEUES","valueFrom":{"configMapKeyRef":{"key":"enableFairQueues","name":"reporting-configs"}}},{"name":"ENABLE_GUARD_RAIL","valueFrom":{"configMapKeyRef":{"key":"enableGuardRail","name":"reporting-configs"}}},{"name":"ENABLE_TRACING","valueFrom":{"configMapKeyRef":{"key":"enableTracing","name":"reporting-configs"}}},{"name":"ENCRYPTION_URI","valueFrom":{"configMapKeyRef":{"key":"encryptionUri","name":"reporting-configs"}}},{"name":"FEATURE_FLAGS_URI","valueFrom":{"configMapKeyRef":{"key":"featureFlagsUri","name":"reporting-configs"}}},{"name":"GUARD_RAIL_RSYN_PERIOD","valueFrom":{"configMapKeyRef":{"key":"guardRailRsynPeriod","name":"reporting-configs"}}},{"name":"KEYS_URI","valueFrom":{"configMapKeyRef":{"key":"keysUri","name":"reporting-configs"}}},{"name":"LOG_LEVEL","valueFrom":{"configMapKeyRef":{"key":"logLevel","name":"reporting-configs"}}},{"name":"MAX_REPORT_DEADLINE","valueFrom":{"configMapKeyRef":{"key":"maxReportDeadline","name":"reporting-configs"}}},{"name":"NATS_STREAMING_CLUSTER_ID","valueFrom":{"configMapKeyRef":{"key":"natsStreamingClusterId","name":"reporting-configs"}}},{"name":"NATS_URI","valueFrom":{"configMapKeyRef":{"key":"natsUri","name":"reporting-configs"}}},{"name":"PROXY_URI","valueFrom":{"configMapKeyRef":{"key":"proxyUri","name":"reporting-configs"}}},{"name":"QIX_SESSIONS_URI","valueFrom":{"configMapKeyRef":{"key":"qixSessionsUri","name":"reporting-configs"}}},{"name":"QMFE_IMPORT_MAP_OVERRIDE_URL_WHITELIST","valueFrom":{"configMapKeyRef":{"key":"qmfeImportMapOverrideUrlWhitelist","name":"reporting-configs"}}},{"name":"QUOTAS","valueFrom":{"configMapKeyRef":{"key":"quotas","name":"reporting-configs"}}},{"name":"REDIS_ACL_USER","valueFrom":{"configMapKeyRef":{"key":"redisAclUser","name":"reporting-configs"}}},{"name":"REDIS_URI","valueFrom":{"configMapKeyRef":{"key":"redisUri","name":"reporting-configs"}}},{"name":"RENDERER_TIMEOUTS","valueFrom":{"configMapKeyRef":{"key":"rendererTimeouts","name":"reporting-configs"}}},{"name":"REPORT_WORKERS","valueFrom":{"configMapKeyRef":{"key":"reportWorkers","name":"reporting-configs"}}},{"name":"REPORTING_TEMPLATES_URI","valueFrom":{"configMapKeyRef":{"key":"reportingTemplatesUri","name":"reporting-configs"}}},{"name":"REQUEST_LIMITS","valueFrom":{"configMapKeyRef":{"key":"requestLimits","name":"reporting-configs"}}},{"name":"SOLACE_MESSAGE_VPN","valueFrom":{"configMapKeyRef":{"key":"solaceMessageVpn","name":"reporting-configs"}}},{"name":"SOLACE_SKIP_CERT_VALIDATION","valueFrom":{"configMapKeyRef":{"key":"solaceSkipCertValidation","name":"reporting-configs"}}},{"name":"SOLACE_URI","valueFrom":{"configMapKeyRef":{"key":"solaceUri","name":"reporting-configs"}}},{"name":"TASK_REQUESTS_QUEUE_QUANTUM","valueFrom":{"configMapKeyRef":{"key":"taskRequestsQueueQuantum","name":"reporting-configs"}}},{"name":"TEMPORARY_CONTENTS_URI","valueFrom":{"configMapKeyRef":{"key":"temporaryContentsUri","name":"reporting-configs"}}},{"name":"TENANT_OVERWRITES","valueFrom":{"configMapKeyRef":{"key":"tenantOverwrites","name":"reporting-configs"}}},{"name":"USAGE_TRACKER_URI","valueFrom":{"configMapKeyRef":{"key":"usageTrackerUri","name":"reporting-configs"}}},{"name":"AUTH_ENABLED","value":"true"},{"name":"AUTH_JWD_AUD","value":"qlik.api.internal"},{"name":"AUTH_JWT_ISS","value":"qlik.api.internal"},{"name":"NO_JWT_VALIDATION","value":"false"},{"name":"REDIS_ENCRYPTION","value":"false"},{"name":"RENDERER_QUERY_STRING","value":""},{"name":"TOKEN_AUTH_ENABLED","value":"true"},{"name":"JAEGER_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}},{"name":"OTLP_AGENT_HOST","valueFrom":{"fieldRef":{"fieldPath":"status.hostIP"}}}],"image":"ghcr.io/qlik-trial/reporting-service:9.151.1","imagePullPolicy":"IfNotPresent","livenessProbe":{"failureThreshold":10,"httpGet":{"path":"/health","port":"http"},"initialDelaySeconds":15,"timeoutSeconds":5},"name":"rep","ports":[{"containerPort":8282,"name":"http","protocol":"TCP"}],"readinessProbe":{"failureThreshold":10,"httpGet":{"path":"/ready","port":"http"},"initialDelaySeconds":15,"timeoutSeconds":5},"resources":{"limits":{"memory":"1.5Gi"},"requests":{"cpu":"500m","memory":"1.5Gi"}},"volumeMounts":[{"mountPath":"/run/secrets/qlik.com/reporting","name":"reporting-secrets","readOnly":true}]}],"imagePullSecrets":[{"name":"artifactory-docker-secret"}],"volumes":[{"name":"reporting-secrets","secret":{"defaultMode":420,"secretName":"reporting-secrets"}},{"emptyDir":{},"name":"rwr-home"},{"emptyDir":{},"name":"rwr-tmp"},{"emptyDir":{},"name":"pki"}]}}}}
    qlik.encryption.service: '{ "serviceS2SName": "reporting", "serviceDisplayName"
      : "Reporting Service" , "serviceEncryptionContext" : "reporting" }'
  creationTimestamp: "2023-07-31T08:40:55Z"
  generation: 2
  labels:
    app: reporting
    app.kubernetes.io/component: reporting
    app.kubernetes.io/instance: reporting
    app.kubernetes.io/name: reporting
    app.kubernetes.io/version: 9.151.1
    chart: reporting-9.151.1
    domain: reporting
    heritage: Helm
    qlik.encryption: "true"
    release: reporting
  name: reporting
  namespace: default
  resourceVersion: "558897"
  uid: 5991be4e-a074-4821-afbc-cdc8baac81ec
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: reporting
      release: reporting
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      annotations:
        checksum/configs: 8e03b70100b236c2ebb6c062eefeb6fd511357920123c86906f22388c22e2b24
        checksum/secrets: d7e48a4307a3af13c81636c9b463444f60fc875fabf471a1a05a5a61a9c3f7e7
      creationTimestamp: null
      labels:
        app: reporting
        app.kubernetes.io/component: reporting
        app.kubernetes.io/instance: reporting
        app.kubernetes.io/name: reporting
        app.kubernetes.io/version: 9.151.1
        chart: reporting-9.151.1
        domain: reporting
        heritage: Helm
        messaging-nats-client: "true"
        release: reporting
    spec:
      containers:
      - env:
        - name: DT_RELEASE_VERSION
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels['app.kubernetes.io/version']
        - name: NATS_CLIENT_ID
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: PROXY_PRIVATE_KEY_FILE
          value: /run/secrets/qlik.com/reporting/proxyPrivateKey
        - name: PROXY_PUBLIC_KEY_FILE
          value: /run/secrets/qlik.com/reporting/proxyPublicKey
        - name: CUSTOMIZATION_URI
          valueFrom:
            configMapKeyRef:
              key: customizationUri
              name: reporting-configs
        - name: DATAPREP_URI
          valueFrom:
            configMapKeyRef:
              key: dataprepUri
              name: reporting-configs
        - name: EDGE_AUTH_URI
          valueFrom:
            configMapKeyRef:
              key: edgeAuthUri
              name: reporting-configs
        - name: ENGINE_REST_API_VERSION
          valueFrom:
            configMapKeyRef:
              key: engineRestApiVersion
              name: reporting-configs
        - name: ENGINE_URI
          valueFrom:
            configMapKeyRef:
              key: engineUri
              name: reporting-configs
        - name: FEATURE_FLAGS_URI
          valueFrom:
            configMapKeyRef:
              key: featureFlagsUri
              name: reporting-configs
        - name: LICENSES_URI
          valueFrom:
            configMapKeyRef:
              key: licensesUri
              name: reporting-configs
        - name: LOCALE_URI
          valueFrom:
            configMapKeyRef:
              key: localeUri
              name: reporting-configs
        - name: MAIN_WEB_CONTAINER_URI
          valueFrom:
            configMapKeyRef:
              key: mainWebContainerUri
              name: reporting-configs
        - name: NLBROKER_URI
          valueFrom:
            configMapKeyRef:
              key: nlbrokerUri
              name: reporting-configs
        - name: NOTES_URI
          valueFrom:
            configMapKeyRef:
              key: notesUri
              name: reporting-configs
        - name: PROXY_URI
          valueFrom:
            configMapKeyRef:
              key: proxyUri
              name: reporting-configs
        - name: RESOURCE_LIBRARY_URI
          valueFrom:
            configMapKeyRef:
              key: resourceLibraryUri
              name: reporting-configs
        - name: SENSE_CLIENT_URI
          valueFrom:
            configMapKeyRef:
              key: senseClientUri
              name: reporting-configs
        - name: SPACES_URI
          valueFrom:
            configMapKeyRef:
              key: spacesUri
              name: reporting-configs
        - name: TEMPORARY_CONTENTS_URI
          valueFrom:
            configMapKeyRef:
              key: temporaryContentsUri
              name: reporting-configs
        - name: TENANTS_URI
          valueFrom:
            configMapKeyRef:
              key: tenantsUri
              name: reporting-configs
        - name: USERS_URI
          valueFrom:
            configMapKeyRef:
              key: usersUri
              name: reporting-configs
        - name: PROXY_METRICS_PORT
          value: "8484"
        - name: PROXY_PORT
          value: "8001"
        - name: JAEGER_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        - name: OTLP_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        image: ghcr.io/qlik-trial/reporting-proxy:3.10.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: metric
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: rpr
        ports:
        - containerPort: 8484
          name: metric
          protocol: TCP
        - containerPort: 8001
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /ready
            port: metric
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          limits:
            memory: 300Mi
          requests:
            cpu: 200m
            memory: 300Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /run/secrets/qlik.com/reporting
          name: reporting-secrets
          readOnly: true
      - args:
        - -l=$(LOG_LEVEL)
        - -http-server-enable=true
        - -http-server-host=0.0.0.0
        - -http-server-port=9288
        - -temp-dir=/home/wuser/rwr/tmp
        env:
        - name: DT_RELEASE_VERSION
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels['app.kubernetes.io/version']
        - name: NATS_CLIENT_ID
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: PROXY_PRIVATE_KEY_FILE
          value: /run/secrets/qlik.com/reporting/proxyPrivateKey
        - name: PROXY_PUBLIC_KEY_FILE
          value: /run/secrets/qlik.com/reporting/proxyPublicKey
        - name: HTTP_HOST
          value: 0.0.0.0
        - name: HTTP_HOST_ENABLE
          value: "true"
        - name: HTTP_PORT
          value: "9288"
        - name: LOG_LEVEL
          value: INFO
        - name: JAEGER_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        - name: OTLP_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        image: ghcr.io/qlik-trial/reporting-web-renderer:2.45.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: rwr
        ports:
        - containerPort: 9288
          name: http
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /ready
            port: http
            scheme: HTTP
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          limits:
            memory: 1536Mi
          requests:
            cpu: "1"
            memory: 1536Mi
        securityContext:
          allowPrivilegeEscalation: false
          capabilities:
            drop:
            - ALL
          privileged: false
          readOnlyRootFilesystem: true
          runAsGroup: 65200
          runAsNonRoot: true
          runAsUser: 65200
          seccompProfile:
            type: RuntimeDefault
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /home/wuser/rwr
          name: rwr-home
        - mountPath: /tmp
          name: rwr-tmp
        - mountPath: /home/wuser/.pki
          name: pki
        - mountPath: /run/secrets/qlik.com/reporting
          name: reporting-secrets
          readOnly: true
      - args:
        - -l="$(LOG_LEVEL)"
        env:
        - name: DT_RELEASE_VERSION
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels['app.kubernetes.io/version']
        - name: NATS_CLIENT_ID
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: REDIS_ACL_PASSWORD_FILE
          value: /run/secrets/qlik.com/reporting/redisAclPassword
        - name: TOKEN_AUTH_PRIVATE_KEY_FILE
          value: /run/secrets/qlik.com/reporting/tokenAuthPrivateKey
        - name: TOKEN_AUTH_PRIVATE_KEY_ID_FILE
          value: /run/secrets/qlik.com/reporting/tokenAuthPrivateKeyId
        - name: AWSS3_ACCESSKEY
          valueFrom:
            configMapKeyRef:
              key: awss3Accesskey
              name: reporting-configs
        - name: AWSS3_BUCKET
          valueFrom:
            configMapKeyRef:
              key: awss3Bucket
              name: reporting-configs
        - name: AWSS3_REGION
          valueFrom:
            configMapKeyRef:
              key: awss3Region
              name: reporting-configs
        - name: AWSS3_SECRETKEY
          valueFrom:
            configMapKeyRef:
              key: awss3Secretkey
              name: reporting-configs
        - name: AWSS3_URL
          valueFrom:
            configMapKeyRef:
              key: awss3Url
              name: reporting-configs
        - name: EDGE_AUTH_URI
          valueFrom:
            configMapKeyRef:
              key: edgeAuthUri
              name: reporting-configs
        - name: ENCRYPTION_URI
          valueFrom:
            configMapKeyRef:
              key: encryptionUri
              name: reporting-configs
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              key: logLevel
              name: reporting-configs
        - name: REDIS_ACL_USER
          valueFrom:
            configMapKeyRef:
              key: redisAclUser
              name: reporting-configs
        - name: REDIS_URI
          valueFrom:
            configMapKeyRef:
              key: redisUri
              name: reporting-configs
        - name: HTTP_PORT
          value: "8384"
        - name: PORT
          value: "52052"
        - name: TOKEN_AUTH_ENABLED
          value: "true"
        - name: JAEGER_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        - name: OTLP_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        image: ghcr.io/qlik-trial/reporting-composer:3.23.0
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 3
          httpGet:
            path: /health
            port: http
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        name: cmp
        ports:
        - containerPort: 8384
          name: http
          protocol: TCP
        readinessProbe:
          failureThreshold: 3
          httpGet:
            path: /ready
            port: http
            scheme: HTTP
          initialDelaySeconds: 10
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 1
        resources:
          limits:
            memory: 1536Mi
          requests:
            cpu: 200m
            memory: 500Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /run/secrets/qlik.com/reporting
          name: reporting-secrets
          readOnly: true
      - env:
        - name: SKIP_QUEUE_PROCESSING
          value: "true"
        - name: DT_RELEASE_VERSION
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.labels['app.kubernetes.io/version']
        - name: NATS_CLIENT_ID
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: metadata.name
        - name: REDIS_ACL_PASSWORD_FILE
          value: /run/secrets/qlik.com/reporting/redisAclPassword
        - name: TOKEN_AUTH_PRIVATE_KEY_FILE
          value: /run/secrets/qlik.com/reporting/tokenAuthPrivateKey
        - name: TOKEN_AUTH_PRIVATE_KEY_ID_FILE
          value: /run/secrets/qlik.com/reporting/tokenAuthPrivateKeyId
        - name: ACCESS_CONTROLS_URI
          valueFrom:
            configMapKeyRef:
              key: accessControlsUri
              name: reporting-configs
        - name: AWSS3_ACCESSKEY
          valueFrom:
            configMapKeyRef:
              key: awss3Accesskey
              name: reporting-configs
        - name: AWSS3_BUCKET
          valueFrom:
            configMapKeyRef:
              key: awss3Bucket
              name: reporting-configs
        - name: AWSS3_REGION
          valueFrom:
            configMapKeyRef:
              key: awss3Region
              name: reporting-configs
        - name: AWSS3_SECRETKEY
          valueFrom:
            configMapKeyRef:
              key: awss3Secretkey
              name: reporting-configs
        - name: AWSS3_URL
          valueFrom:
            configMapKeyRef:
              key: awss3Url
              name: reporting-configs
        - name: CONTENT_SIZE_THRESHOLD
          valueFrom:
            configMapKeyRef:
              key: contentSizeThreshold
              name: reporting-configs
        - name: CONTENT_STORAGE
          valueFrom:
            configMapKeyRef:
              key: contentStorage
              name: reporting-configs
        - name: DEFAULT_REPORT_DEADLINE
          valueFrom:
            configMapKeyRef:
              key: defaultReportDeadline
              name: reporting-configs
        - name: EDGE_AUTH_URI
          valueFrom:
            configMapKeyRef:
              key: edgeAuthUri
              name: reporting-configs
        - name: ENABLE_DIRECT_QUERY
          valueFrom:
            configMapKeyRef:
              key: enableDirectQuery
              name: reporting-configs
        - name: ENABLE_FAIR_QUEUES
          valueFrom:
            configMapKeyRef:
              key: enableFairQueues
              name: reporting-configs
        - name: ENABLE_GUARD_RAIL
          valueFrom:
            configMapKeyRef:
              key: enableGuardRail
              name: reporting-configs
        - name: ENABLE_TRACING
          valueFrom:
            configMapKeyRef:
              key: enableTracing
              name: reporting-configs
        - name: ENCRYPTION_URI
          valueFrom:
            configMapKeyRef:
              key: encryptionUri
              name: reporting-configs
        - name: FEATURE_FLAGS_URI
          valueFrom:
            configMapKeyRef:
              key: featureFlagsUri
              name: reporting-configs
        - name: GUARD_RAIL_RSYN_PERIOD
          valueFrom:
            configMapKeyRef:
              key: guardRailRsynPeriod
              name: reporting-configs
        - name: KEYS_URI
          valueFrom:
            configMapKeyRef:
              key: keysUri
              name: reporting-configs
        - name: LOG_LEVEL
          valueFrom:
            configMapKeyRef:
              key: logLevel
              name: reporting-configs
        - name: MAX_REPORT_DEADLINE
          valueFrom:
            configMapKeyRef:
              key: maxReportDeadline
              name: reporting-configs
        - name: NATS_STREAMING_CLUSTER_ID
          valueFrom:
            configMapKeyRef:
              key: natsStreamingClusterId
              name: reporting-configs
        - name: NATS_URI
          valueFrom:
            configMapKeyRef:
              key: natsUri
              name: reporting-configs
        - name: PROXY_URI
          valueFrom:
            configMapKeyRef:
              key: proxyUri
              name: reporting-configs
        - name: QIX_SESSIONS_URI
          valueFrom:
            configMapKeyRef:
              key: qixSessionsUri
              name: reporting-configs
        - name: QMFE_IMPORT_MAP_OVERRIDE_URL_WHITELIST
          valueFrom:
            configMapKeyRef:
              key: qmfeImportMapOverrideUrlWhitelist
              name: reporting-configs
        - name: QUOTAS
          valueFrom:
            configMapKeyRef:
              key: quotas
              name: reporting-configs
        - name: REDIS_ACL_USER
          valueFrom:
            configMapKeyRef:
              key: redisAclUser
              name: reporting-configs
        - name: REDIS_URI
          valueFrom:
            configMapKeyRef:
              key: redisUri
              name: reporting-configs
        - name: RENDERER_TIMEOUTS
          valueFrom:
            configMapKeyRef:
              key: rendererTimeouts
              name: reporting-configs
        - name: REPORT_WORKERS
          valueFrom:
            configMapKeyRef:
              key: reportWorkers
              name: reporting-configs
        - name: REPORTING_TEMPLATES_URI
          valueFrom:
            configMapKeyRef:
              key: reportingTemplatesUri
              name: reporting-configs
        - name: REQUEST_LIMITS
          valueFrom:
            configMapKeyRef:
              key: requestLimits
              name: reporting-configs
        - name: SOLACE_MESSAGE_VPN
          valueFrom:
            configMapKeyRef:
              key: solaceMessageVpn
              name: reporting-configs
        - name: SOLACE_SKIP_CERT_VALIDATION
          valueFrom:
            configMapKeyRef:
              key: solaceSkipCertValidation
              name: reporting-configs
        - name: SOLACE_URI
          valueFrom:
            configMapKeyRef:
              key: solaceUri
              name: reporting-configs
        - name: TASK_REQUESTS_QUEUE_QUANTUM
          valueFrom:
            configMapKeyRef:
              key: taskRequestsQueueQuantum
              name: reporting-configs
        - name: TEMPORARY_CONTENTS_URI
          valueFrom:
            configMapKeyRef:
              key: temporaryContentsUri
              name: reporting-configs
        - name: TENANT_OVERWRITES
          valueFrom:
            configMapKeyRef:
              key: tenantOverwrites
              name: reporting-configs
        - name: USAGE_TRACKER_URI
          valueFrom:
            configMapKeyRef:
              key: usageTrackerUri
              name: reporting-configs
        - name: AUTH_ENABLED
          value: "true"
        - name: AUTH_JWD_AUD
          value: qlik.api.internal
        - name: AUTH_JWT_ISS
          value: qlik.api.internal
        - name: NO_JWT_VALIDATION
          value: "false"
        - name: REDIS_ENCRYPTION
          value: "false"
        - name: RENDERER_QUERY_STRING
        - name: TOKEN_AUTH_ENABLED
          value: "true"
        - name: JAEGER_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        - name: OTLP_AGENT_HOST
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.hostIP
        image: ghcr.io/qlik-trial/reporting-service:9.151.1
        imagePullPolicy: IfNotPresent
        livenessProbe:
          failureThreshold: 10
          httpGet:
            path: /health
            port: http
            scheme: HTTP
          initialDelaySeconds: 15
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5
        name: rep
        ports:
        - containerPort: 8282
          name: http
          protocol: TCP
        readinessProbe:
          failureThreshold: 10
          httpGet:
            path: /ready
            port: http
            scheme: HTTP
          initialDelaySeconds: 15
          periodSeconds: 10
          successThreshold: 1
          timeoutSeconds: 5
        resources:
          limits:
            memory: 1536Mi
          requests:
            cpu: 500m
            memory: 1536Mi
        terminationMessagePath: /dev/termination-log
        terminationMessagePolicy: File
        volumeMounts:
        - mountPath: /run/secrets/qlik.com/reporting
          name: reporting-secrets
          readOnly: true
      dnsPolicy: ClusterFirst
      imagePullSecrets:
      - name: artifactory-docker-secret
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: {}
      terminationGracePeriodSeconds: 30
      volumes:
      - name: reporting-secrets
        secret:
          defaultMode: 420
          secretName: reporting-secrets
      - emptyDir: {}
        name: rwr-home
      - emptyDir: {}
        name: rwr-tmp
      - emptyDir: {}
        name: pki
status:
  availableReplicas: 1
  conditions:
  - lastTransitionTime: "2023-07-31T08:42:40Z"
    lastUpdateTime: "2023-07-31T08:42:40Z"
    message: Deployment has minimum availability.
    reason: MinimumReplicasAvailable
    status: "True"
    type: Available
  - lastTransitionTime: "2023-07-31T08:40:55Z"
    lastUpdateTime: "2023-08-01T10:25:47Z"
    message: ReplicaSet "reporting-67477cbd96" has successfully progressed.
    reason: NewReplicaSetAvailable
    status: "True"
    type: Progressing
  observedGeneration: 2
  readyReplicas: 1
  replicas: 1
  updatedReplicas: 1`

	corporate = `－
`

	e := &Editor{}
	e.Edit(corporate)

}

func populateDropdown(root *Node) {
	items := []string{}
	if root != nil {
		for _, child := range root.Children {
			items = append(items, child.FieldName+"   "+child.FieldType)
		}

	}
	//drop.SetOptions(items, nil)
}

func buildCurrentPath(txt string, x int, y int) string {
	var path []string
	for i := y - 1; i >= 0; i-- {
		st, x1, _ := getLeftmostWordAtLine(txt, i)
		if x1 < x && st != "" {
			if st[len(st)-1:] == ":" {
				st = st[:len(st)-1]
			}
			path = append([]string{st}, path...)
			x = x1
		}
	}
	result := strings.Join(path, ".")
	return result
}

func getCurrentWordSelection(txt string, selPos int) (selStart int, selEnd int) {
	//if selPos < len(txt) && selPos > 0 && txt[selPos:selPos+1] == "\n" {
	//	selPos--
	//}
	selStart = len(txt) - 1
	selEnd = selStart + 1

	// expand left
	for i := selPos; i > -1; i-- {
		if i < len(txt) {
			s := txt[i : i+1]
			if !isLetter(s) {
				selStart = i
				if i != selPos {
					selStart++
				}
				break
			}
			if i == 0 {
				selStart--
				break
			}
		}
	}
	// expand right
	for i := selPos; i < len(txt); i++ {
		s := txt[i : i+1]
		if !isLetter(s) {
			selEnd = i
			break
		}
	}
	if selEnd < selStart {
		selEnd = selStart
	}
	return
}

func getLeftmostWordAtLine(txt string, y int) (word string, x1 int, x2 int) {
	n := 0
	start := -1
	end := 0
	if y < 0 {
		return "", 0, 0
	}
	//  get the text at line y
	for i := 0; i < len(txt); i++ {
		if txt[i:i+1] == "\n" {
			n++
			if n == y {
				start = i
			}
			if n == y+1 {
				break
			}
		}
		end = i
	}
	st := txt[start+1 : end+1]

	// find the start of the word
	for i := 0; i < len(st); i++ {
		s := st[i : i+1]
		if isLetter(s) {

			if s == "-" {
				continue
			}

			x1 = i
			break
		}
	}
	// find the end of the word
	x2 = len(st)
	for i := x1; i < len(st); i++ {
		s := st[i : i+1]
		if !isLetter(s) {
			x2 = i
			break
		}
	}
	word = st[x1:x2]
	return
}

func isLetter(s string) bool {
	if s == " " || s == "\t" || s == "\n" {
		return false
	}
	return true
}
