// Copyright 2023 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

import moment from "moment-timezone";
import React, { useContext } from "react";

import { Timezone } from "src/timestamp";

import { TimezoneContext } from "../contexts";

import { dateFormat, timeFormat } from "./timeScaleDropdown";
import { TimeScale } from "./timeScaleTypes";
import { toRoundedDateRange } from "./utils";

export const FormattedTimescale = (props: {
  ts: TimeScale;
  requestTime?: moment.Moment;
}) => {
  const timezone = useContext(TimezoneContext);

  const [start, end] = toRoundedDateRange(props.ts);
  const startTz = start.tz(timezone);
  const endTz = end.tz(timezone);
  const endDayIsToday = endTz.isSame(moment.tz(timezone), "day");
  const startEndOnSameDay = endTz.isSame(startTz, "day");
  const omitDayFormat = endDayIsToday && startEndOnSameDay;
  const dateStart = omitDayFormat ? "" : startTz.format(dateFormat);
  const dateEnd =
    omitDayFormat || startEndOnSameDay ? "" : endTz.format(dateFormat);
  const timeStart = startTz.format(timeFormat);
  const timeEnd =
    props.ts.key !== "Custom" && props.requestTime?.isValid()
      ? props.requestTime.tz(timezone).format(timeFormat)
      : endTz.format(timeFormat);

  return (
    <>
      {dateStart} {timeStart} to {dateEnd} {timeEnd} <Timezone />
    </>
  );
};
