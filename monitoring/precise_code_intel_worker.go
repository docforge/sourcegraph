package main

func PreciseCodeIntelWorker() *Container {
	return &Container{
		Name:        "precise-code-intel-worker",
		Title:       "Precise Code Intel Worker",
		Description: "Handles conversion of uploaded precise code intelligence bundles.",
		Groups: []Group{
			{
				Title: "General",
				Rows: []Row{
					{
						{
							Name:              "upload_queue_size",
							Description:       "upload queue size",
							Query:             `max(src_upload_queue_uploads_total)`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 100},
							PanelOptions:      PanelOptions().LegendFormat("uploads queued for processing"),
							PossibleSolutions: "none",
						},
						{
							Name:              "upload_queue_growth_rate",
							Description:       "upload queue growth rate every 5m",
							Query:             `sum(increase(src_upload_queue_uploads_total[30m])) / sum(increase(src_upload_queue_processor_total[30m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 5},
							PanelOptions:      PanelOptions().LegendFormat("upload queue growth rate"),
							PossibleSolutions: "none",
						},
						{
							Name:        "upload_process_errors",
							Description: "upload process errors every 5m",
							// TODO(efritz) - ensure these differentiate malformed dumps and system errors
							Query:             `sum(increase(src_upload_queue_processor_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
					},
					{
						{
							Name:        "99th_percentile_store_duration",
							Description: "99th percentile successful database query duration over 5m",
							// TODO(efritz) - ensure these exclude error durations
							Query:             `histogram_quantile(0.99, sum by (le)(rate(src_code_intel_store_duration_seconds_bucket{job="precise-code-intel-worker"}[5m])))`,
							DataMayNotExist:   true,
							DataMayBeNaN:      true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("store operation").Unit(Seconds),
							PossibleSolutions: "none",
						},
						{
							Name:              "store_errors",
							Description:       "database errors every 5m",
							Query:             `increase(src_code_intel_store_errors_total{job="precise-code-intel-worker"}[5m])`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("store operation"),
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Upload resetter - re-queues uploads that did not complete processing",
				Hidden: true,
				Rows: []Row{
					{
						{
							Name:              "processing_uploads_reset",
							Description:       "uploads reset to queued state every 5m",
							Query:             `sum(increase(src_upload_queue_resets_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("uploads"),
							PossibleSolutions: "none",
						},
						{
							Name:              "processing_uploads_reset_failures",
							Description:       "uploads errored after repeated resets every 5m",
							Query:             `sum(increase(src_upload_queue_max_resets_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("uploads"),
							PossibleSolutions: "none",
						},
						{
							Name:              "upload_resetter_errors",
							Description:       "upload resetter errors every 5m",
							Query:             `sum(increase(src_upload_queue_reset_errors_total[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("errors"),
							PossibleSolutions: "none",
						},
					},
				},
			},
			{
				Title:  "Internal service requests",
				Hidden: true,
				Rows: []Row{
					{
						{
							Name:              "99th_percentile_bundle_manager_transfer_duration",
							Description:       "99th percentile successful bundle manager data transfer duration over 5m",
							Query:             `histogram_quantile(0.99, sum by (le,category)(rate(src_precise_code_intel_bundle_manager_request_duration_seconds_bucket{job="precise-code-intel-worker",category="transfer"}[5m])))`,
							DataMayNotExist:   true,
							DataMayBeNaN:      true,
							Warning:           Alert{GreaterOrEqual: 300},
							PanelOptions:      PanelOptions().LegendFormat("{{category}}").Unit(Seconds),
							PossibleSolutions: "none",
						},
						{
							Name:              "bundle_manager_error_responses",
							Description:       "bundle manager error responses every 5m",
							Query:             `sum by (category)(increase(src_precise_code_intel_bundle_manager_request_duration_seconds_count{job="precise-code-intel-worker",code!~"2.."}[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 5},
							PanelOptions:      PanelOptions().LegendFormat("{{category}}"),
							PossibleSolutions: "none",
						},
					},
					{
						{
							Name:              "99th_percentile_gitserver_duration",
							Description:       "99th percentile successful gitserver query duration over 5m",
							Query:             `histogram_quantile(0.99, sum by (le,category)(rate(src_gitserver_request_duration_seconds_bucket{job="precise-code-intel-worker"}[5m])))`,
							DataMayNotExist:   true,
							DataMayBeNaN:      true,
							Warning:           Alert{GreaterOrEqual: 20},
							PanelOptions:      PanelOptions().LegendFormat("{{category}}").Unit(Seconds),
							PossibleSolutions: "none",
						},
						{
							Name:              "gitserver_error_responses",
							Description:       "gitserver error responses every 5m",
							Query:             `sum by (category)(increase(src_gitserver_request_duration_seconds_count{job="precise-code-intel-worker",code!~"2.."}[5m]))`,
							DataMayNotExist:   true,
							Warning:           Alert{GreaterOrEqual: 5},
							PanelOptions:      PanelOptions().LegendFormat("{{category}}"),
							PossibleSolutions: "none",
						},
					},
					{
						sharedFrontendInternalAPIErrorResponses("precise-code-intel-worker"),
					},
				},
			},
			{
				Title:  "Container monitoring (not available on server)",
				Hidden: true,
				Rows: []Row{
					{
						sharedContainerRestarts("precise-code-intel-worker"),
						sharedContainerMemoryUsage("precise-code-intel-worker"),
						sharedContainerCPUUsage("precise-code-intel-worker"),
					},
				},
			},
			{
				Title:  "Provisioning indicators (not available on server)",
				Hidden: true,
				Rows: []Row{
					{
						sharedProvisioningCPUUsage7d("precise-code-intel-worker"),
						sharedProvisioningMemoryUsage7d("precise-code-intel-worker"),
					},
					{
						sharedProvisioningCPUUsage5m("precise-code-intel-worker"),
						sharedProvisioningMemoryUsage5m("precise-code-intel-worker"),
					},
				},
			},
		},
	}
}
