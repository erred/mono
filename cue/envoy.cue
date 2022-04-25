overload_manager: {
  refresh_interval: "0.25s"
  resource_monitors: [{
    name: "envoy.resource_monitors.fixed_heap"
    typed_config: {
      "@type": "type.googleapis.com/envoy.extensions.resource_monitors.fixed_heap.v3.FixedHeapConfig"
      max_heap_size_bytes: 2147483648
    }
  }]
  actions: [{
    name: "envoy.overload_actions.shrink_heap"
    triggers: [{
      name: "envoy.resource_monitors.fixed_heap"
      threshold: value: 0.95
    }]
  },{
    name: "envoy.overload_actions.stop_accepting_requests"
    triggers: [{
      name: "envoy.resource_monitors.fixed_heap"
      threshold: value: 0.98
    }]
  }]
}

_http_filter: {
  name: "envoy.filters.network.http_connection_manager"
  typed_config: {
    "@type": "type.googleapis.com/envoy.extensions.filters.network.http_connection_manager.v3.HttpConnectionManager"
    stat_prefix: "ingress_http"
    codec_type: "AUTO"
    access_log: [_access_log]
    generate_request_id: true
    tracing: _tracing
    http_filters: [{
      name: "envoy.filters.http.router"
    }]
  }
}

_access_log: {
  name: "envoy.access_loggers.file"
  typed_config: {
    "@type": "type.googleapis.com/envoy.extensions.access_loggers.stream.v3.StdoutAccessLog"
    log_format: json_format: accesslog: {
      ts: "%START_TIME%"
      traceParent: "%REQ(traceparent)%"
      httpMethod: "%REQ(:method)%"
      httpUrl: "%REQ(x-envoy-original-path?:path)%"
      httpVersion: "%PROTOCOL%"
      httpHost: "%REQ(:authority)%"
      httpUseragent: "%REQ(user-agent)%"
      handleTimeMs: "%DURATION%"
      httpStatus: "%RESPONSE_CODE%"
      bytesWritten: "%BYTES_SENT%"
      bytesReceived: "%BYTES_RECEIVED%"
    }
  }
}
_tracing: provider: {
  name: "envoy.tracers.zipkin"
  typed_config: {
    "@type": "type.googleapis.com/envoy.config.trace.v3.ZipkinConfig"
    collector_cluster: "otelcol_zipkin"
    collector_endpoint: "/api/v2/spans"
    collector_endpoint_version: "HTTP_JSON"
  }
}

static_resources: {
  listeners: [{
    name: "http"
    address: socket_address: {
      address: "0.0.0.0"
      port_value: 80
    }
    filter_chains: [{
      filters: [ _http_filter & {
        typed_config: route_config: {
          name: "default_route"
          virtual_hosts: [{
            name: "redirect"
            domains: ["*"]
            routes: [{
              match: prefix: "/"
              redirect: https_redirect: true
            }]
          }]
        }
      }]
    }]
  },{
    name: "https"
    address: socket_address: {
      address: "0.0.0.0"
      port_value: 443
    }
    listener_filters: [{
      name: "envoy.filters.listener.tls_inspector"
    }]
    filter_chains: [{
      filter_chain_match: server_names: [
        "liao.dev",
        "*.liao.dev",
      ]
      transport_socket: {
        name: "envoy.transport_sockets.tls"
        typed_config: {
            "@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext"
            common_tls_context: tls_certificates: [{
                certificate_chain: filename: "/etc/tls/liao.dev/fullchain.cer"
                private_key: filename: "/etc/tls/liao.dev/liao.dev.key"
            }]
        }
      }
      filters: [ _http_filter & {
        typed_config: route_config: {
          name: "default_route"
          virtual_hosts: [for nme,domain in _virtual_host_dev {
            name: nme
            domains: [domain]
            routes: [{
              match: prefix: "/"
              route: cluster: nme
            }]
          }]
        }
      }]
    },{
      filter_chain_match: server_names: [
        "seankhliao.com",
        "*.seankhliao.com",
      ]
      transport_socket: {
        name: "envoy.transport_sockets.tls"
        typed_config: {
            "@type": "type.googleapis.com/envoy.extensions.transport_sockets.tls.v3.DownstreamTlsContext"
            common_tls_context: tls_certificates: [{
                certificate_chain: filename: "/etc/tls/seankhliao.com/fullchain.cer"
                private_key: filename: "/etc/tls/seankhliao.com/seankhliao.com.key"
            }]
        }
      }
      filters: [ _http_filter & {
        typed_config: route_config: {
          name: "default_route"
          virtual_hosts: [{
            name: "blog"
            domains: ["seankhliao.com"]
            routes: [{
              match: prefix: "/"
              route: cluster: "blog"
            }]
          },{
            name: "vanity"
            domains: ["go.seankhliao.com"]
            routes: [{
              match: prefix: "/"
              route: cluster: "vanity"
            }]
          }]
        }
      }]
    }]
  }]
  clusters: [for _name,_ports in _local_cluster_ports {
    name: _name
    connect_timeout: "0.25s"
    type: "STATIC"
    lb_policy: "ROUND_ROBIN"
    load_assignment: endpoints: [{
      lb_endpoints: [for _port in _ports {
          endpoint: address: socket_address: {
            address: "192.168.100.1"
            port_value: _port
          }
      }]
    }]
  }]
}

_virtual_host_dev: {
  liaodev: "liao.dev"
  earbug: "earbug.liao.dev"
  ghdefaults: "ghdefaults.liao.dev"
  medea: "medea.liao.dev"
  paste: "paste.liao.dev"
  stylesheet: "stylesheet.liao.dev"

}

_local_cluster_ports: {
  blog: [28001, 28101]
  vanity: [28002, 28102]
  ghdefaults: [28003, 28103]
  medea: [28004, 28104]
  stylesheet: [28005, 28105]
  liaodev: [28006, 28106]
  earbug: [28007]
  paste: [28008]
  otelcol_zipkin: [9411]
}
