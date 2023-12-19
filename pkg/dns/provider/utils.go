/*
 * Copyright 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 *
 */

package provider

import (
	"fmt"
	"reflect"

	"github.com/gardener/controller-manager-library/pkg/resources"
	api "github.com/gardener/external-dns-management/pkg/apis/dns/v1alpha1"
	"github.com/gardener/external-dns-management/pkg/dns"
)

type NullMetrics struct{}

var _ Metrics = &NullMetrics{}

func (m *NullMetrics) AddGenericRequests(_ string, _ int) {
}

func (m *NullMetrics) AddZoneRequests(_, _ string, _ int) {
}

func copyZones(src map[dns.ZoneID]*dnsHostedZone) dnsHostedZones {
	dst := dnsHostedZones{}
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func errorValue(format string, err error) string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf(format, err.Error())
}

func filterZoneByProvider(zones []*dnsHostedZone, provider DNSProvider) *dnsHostedZone {
	if provider != nil {
		for _, zone := range zones {
			if provider.IncludesZone(zone.Id()) {
				return zone
			}
		}
	}
	return nil
}

func assureRateLimit(mod *resources.ModificationState, t **api.RateLimit, s *api.RateLimit) {
	if s == nil && *t != nil {
		*t = nil
		mod.Modify(true)
	} else if s != nil {
		if *t == nil || !reflect.DeepEqual(**t, *s) {
			*t = s
			mod.Modify(true)
		}
	}
}
