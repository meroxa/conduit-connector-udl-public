// Copyright © 2023 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package destination

import (
	"context"

	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/meroxa/conduit-connector-udl-public/udl"

	"fmt"
	"strings"
)

func OpenCDCPayload(rawPayload map[string]interface{}) string {
	p := rawPayload["payload"].(map[string]interface{})
	a := p["after"].(map[string]interface{})
	return a["opencdc.rawData"].(string)
}

func (d *Destination) writeEphemerisToUDL(ctx context.Context, records []sdk.Record) (int, error) {
	// var ephemerisData []udl.EphemerisIngest
	for _, r := range records {
		ephmerisRecord, err := ToUDLEphemeris(r.Payload.After.Bytes(), udl.EphemerisIngestDataMode(d.Config.DataMode), d.Config.ClassificationMarking)
		// sdk.Logger(ctx).Info().Msgf("ephmerisRecord: %+v", ephmerisRecord)
		// ephemerisData = append(ephemerisData, ephmerisRecord)
		if err != nil {
			sdk.Logger(ctx).Err(err).Msgf("ToUDLEphemeris failed")
			return 0, err
		}

		params := &udl.FiledropEphemPostIdParams{
			IdOnOrbit:       ephmerisRecord.ID,
			Classification:  d.Config.ClassificationMarking,
			DataMode:        "TEST",
			HasMnvr:         false,
			Type:            "ROUTINE",
			Category:        "EXTERNAL",
			EphemFormatType: "NASA",
			Source:          "Spire",
		}
		bodyReader := strings.NewReader(ephmerisRecord.String())
		response, err := d.client.FiledropEphemPostIdWithBody(ctx, params, "applications/json", bodyReader)
		if err != nil {
			return 0, err
		}

		sdk.Logger(context.Background()).Info().Msgf("Submitted Ephemeris Request Parameters - IdOnOrbit: %s, Classification: %s, DataMode: %s, HasMnvr: %t, Type: %s, Category: %s, EphemFormatType: %s, Source: %s", params.IdOnOrbit, params.Classification, params.DataMode, params.HasMnvr, params.Type, params.Category, params.EphemFormatType, params.Source)

		if response.StatusCode > 300 {
			return 0, fmt.Errorf(fmt.Sprintf("unsuccessful status code returned %d; response: %+v", response.StatusCode, response.Body))
		}

		sdk.Logger(context.Background()).Info().Msgf("Spire to Ephemeris UDL response: %+v:", response)
	}

	return 1, nil
}

func (d *Destination) writeAisToUDL(ctx context.Context, records []sdk.Record) (int, error) {
	var aisData []udl.AISIngest
	for _, r := range records {
		ais, err := ToUDLAis(r.Payload.After.Bytes(), udl.AISIngestDataMode(d.Config.DataMode), d.Config.ClassificationMarking)
		aisData = append(aisData, ais)
		if err != nil {
			sdk.Logger(ctx).Err(err).Msgf("ToUDLAis failed")
			return 0, err
		}
	}

	resp, err := d.client.FiledropUdlAisPostId(ctx, aisData)
	if err != nil || resp.StatusCode >= 300 {
		sdk.Logger(ctx).Err(err).Msgf("FiledropUdlAisPostId failed with status code: %v", resp.StatusCode)
		return 0, err
	}
	sdk.Logger(ctx).Info().Msgf("Spire to AIS UDL response: %+v", resp)

	return len(aisData), nil
}

func (d *Destination) writeElsetToUDL(ctx context.Context, records []sdk.Record) (int, error) {
	var elsets []udl.ElsetIngest
	for _, r := range records {
		elset, err := ToUDLElset(r.Payload.After.Bytes())
		if err != nil {
			sdk.Logger(ctx).Err(err).Msgf("ToUDLElset failed")
			return 0, err
		}
		elsets = append(elsets, elset)
	}

	resp, err := d.client.FiledropUdlElsetPostId(ctx, elsets)
	if err != nil || resp.StatusCode > 300 {
		sdk.Logger(ctx).Err(err).Msgf("FiledropUdlElsetPostId failed with status code: %v", resp.StatusCode)
		return 0, err
	}

	return len(elsets), nil
}
