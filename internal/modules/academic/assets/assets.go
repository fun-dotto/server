// Package assets は学年暦の JSON を academic モジュールに提供する。
// データ本体は internal/shared/academiccalendar にあり、app-api と共有する。
package assets

import "github.com/fun-dotto/server/internal/shared/academiccalendar"

// EventsJSON は行事予定の JSON。
var EventsJSON = academiccalendar.EventsJSON

// HolidaysJSON は休日の JSON。
var HolidaysJSON = academiccalendar.HolidaysJSON

// TermsJSON は開講時期ごとの期間の JSON。
var TermsJSON = academiccalendar.TermsJSON
