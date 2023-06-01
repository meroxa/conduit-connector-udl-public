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

	"github.com/meroxa/udl-go"
)

func (d *Destination) writeAisToUDL(ctx context.Context, records []sdk.Record) (int, error) {
	var aisData []udl.AISIngest
	for _, r := range records {
		ais, err := ToUDLAis(r.Payload.After.Bytes(), udl.AISIngestDataMode(d.Config.DataMode))
		sdk.Logger(ctx).Debug().Msgf("ais output: %+v", ais)
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
	sdk.Logger(ctx).Debug().Msgf("resp: %+v", resp)

	return len(aisData), nil
}

func (d *Destination) writeElsetToUDL(ctx context.Context, records []sdk.Record) (int, error) {
	var elsets []udl.ElsetIngest
	for _, r := range records {
		sdk.Logger(ctx).Debug().Msgf("record: %+v", r)
		elset, err := ToUDLElset(r.Payload.After.Bytes())
		sdk.Logger(ctx).Debug().Msgf("elset: %+v", elset)
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