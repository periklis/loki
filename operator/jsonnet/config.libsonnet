local utils = (import 'github.com/grafana/jsonnet-libs/mixin-utils/utils.libsonnet');

{
  loki: {
    local withDatasource = function(ds) ds + (
      if ds.name == 'datasource' then {
        query: 'prometheus',
        current: {
          selected: true,
          text: 'default',
          value: 'default',
        },
      } else {}
    ),

    local withNamespace = function(ns) ns + (
      if ns.label == 'namespace' then {
        datasource: '${datasource}',
        current: {
          selected: false,
          text: 'openshift-logging',
          value: 'openshift-logging',
        },
        definition: 'label_values(loki_build_info, namespace)',
        query: 'label_values(loki_build_info, namespace)',
      } else {}
    ),

    local defaultLokiTags = function(t)
      std.uniq(t + ['logging', 'loki-mixin']),

    local replaceMatchers = function(replacements)
      function(p) p {
        targets: [
          t {
            expr: std.foldl(function(x, rp) std.strReplace(x, rp.from, rp.to), replacements, t.expr),
          }
          for t in p.targets
          if std.objectHas(p, 'targets')
        ],
      },

    // dropPanels removes unnecessary panels from the loki dashboards
    // that are of obsolete usage on our AWS-based deployment environment.
    local dropPanels = function(panels, dropList, fn)
      [
        p
        for p in panels
        if !std.member(dropList, p.title) && fn(p)
      ],

    // mapPanels applies recursively a set of functions over all panels.
    // Note: A Grafana dashboard panel can include other panels.
    // Example: Replace job label in expression and add axis units for all panels.
    local mapPanels = function(funcs, panels)
      [
        // Transform the current panel by applying all transformer funcs.
        // Keep the last version after foldl ends.
        std.foldl(function(agg, fn) fn(agg), funcs, p) + (
          // Recursively apply all transformer functions to any
          // children panels.
          if std.objectHas(p, 'panels') then {
            panels: mapPanels(funcs, p.panels),
          } else {}
        )
        for p in panels
      ],

    // mapTemplateParameters applies a static list of transformer functions to
    // all dashboard template parameters. The static list includes:
    // - cluster-prometheus as datasource based on environment.
    local mapTemplateParameters = function(ls)
      [
        std.foldl(function(x, fn) fn(x), [withDatasource, withNamespace], item)
        for item in ls
        if item.name != 'cluster'
      ],

    prometheusRules+: {
      groups: [
        g {
          rules: [
            r {
              expr: std.strReplace(r.expr, 'cluster, ', ''),
              record: std.strReplace(r.record, 'cluster_', ''),
            }
            for r in g.rules
          ],
        }
        for g in super.groups
      ],
    },

    local dropHistograms = function(p)
      local elems = std.filter(function(p) p.type == 'heatmap', p.panels);
      std.length(elems) == 0,

    grafanaDashboards+: {
      'loki-retention.json'+: {
        local dropList = ['Logs', 'Per Table Marker', 'Sweeper', ''],
        local replacements = [
          { from: 'cluster=~"$cluster",', to: '' },
          { from: 'container="compactor"', to: 'container=~".+-compactor"' },
          { from: 'job=~"($namespace)/compactor"', to: 'namespace="$namespace", job=~".+-compactor-http"' },
        ],
        uid: 'RetCujSHzC8gd9i5fck9a3v9n2EvTzA',
        title: 'OpenShift Logging / LokiStack / Retention',
        tags: defaultLokiTags(super.tags),
        rows: [
          r {
            panels: mapPanels([replaceMatchers(replacements)], r.panels),
          }
          for r in dropPanels(super.rows, dropList, function(p) true)
        ],
        templating+: {
          list: mapTemplateParameters(super.list),
        },
      },
      'loki-chunks.json'+: {
        local dropList = ['Utilization'],
        uid: 'GtCujSHzC8gd9i5fck9a3v9n2EvTzA',
        title: 'OpenShift Logging / LokiStack / Chunks',
        tags: defaultLokiTags(super.tags),
        showMultiCluster:: false,
        namespaceQuery:: 'label_values(loki_build_info, namespace)',
        namespaceType:: 'query',
        labelsSelector:: 'namespace="$namespace", job=~".+-ingester-http"',
        rows: [
          r
          for r in dropPanels(super.rows, dropList, dropHistograms)
        ],
        templating+: {
          list: mapTemplateParameters(super.list),
        },
      },
      'loki-reads.json'+: {
        local dropList = ['BigTable', 'Ingester - Zone Aware'],

        uid: '62q5jjYwhVSaz4Mcrm8tV3My3gcKED',
        title: 'OpenShift Logging / LokiStack / Reads',
        tags: defaultLokiTags(super.tags),
        showMultiCluster:: false,
        namespaceQuery:: 'label_values(loki_build_info, namespace)',
        namespaceType:: 'query',
        matchers:: {
          cortexgateway:: [],
          queryFrontend:: [
            utils.selector.eq('namespace', '$namespace'),
            utils.selector.re('job', '.+-query-frontend-http'),
          ],
          querier:: [
            utils.selector.eq('namespace', '$namespace'),
            utils.selector.re('job', '.+-querier-http'),
          ],
          ingester:: [
            utils.selector.eq('namespace', '$namespace'),
            utils.selector.re('job', '.+-ingester-http'),
          ],
          ingesterZoneAware:: [],
          querierOrIndexGateway:: [
            utils.selector.eq('namespace', '$namespace'),
            utils.selector.re('job', '.+-index-gateway-http'),
          ],
        },
        rows: dropPanels(super.rows, dropList, function(p) true),
        templating+: {
          list: mapTemplateParameters(super.list),
        },
      },
      'loki-writes.json'+: {
        local dropList = ['Ingester - Zone Aware'],
        uid: 'F6nRYKuXmFVpVSFQmXr7cgXy5j7UNr',
        title: 'OpenShift Logging / LokiStack / Writes',
        tags: defaultLokiTags(super.tags),
        showMultiCluster:: false,
        namespaceQuery:: 'label_values(loki_build_info, namespace)',
        namespaceType:: 'query',
        matchers:: {
          cortexgateway:: [],
          distributor:: [
            utils.selector.eq('namespace', '$namespace'),
            utils.selector.re('job', '.+-distributor-http'),
            utils.selector.eq('route', 'loki_api_v1_push'),
          ],
          ingester:: [
            utils.selector.eq('namespace', '$namespace'),
            utils.selector.re('job', '.+-ingester-http'),
          ],
          ingester_zone:: [],
        },
        rows: dropPanels(super.rows, dropList, function(p) true),
        templating+: {
          list: mapTemplateParameters(super.list),
        },
      },
    },
  },
}