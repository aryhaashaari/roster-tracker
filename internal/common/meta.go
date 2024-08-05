package common

import (
	"fmt"
	"github.com/pkg/errors"
	"gitlab.privy.id/privypass/privypass-boilerplate/internal/consts"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Meta struct {
	Pagination
	Ordering
	Filters       []Filter
	filters       map[string]string
	DateTimeRange *DateTimeRange

	opts MetaOption
}

type Ordering struct {
	OrderBy        string
	OrderDirection string
}

type Pagination struct {
	PerPage int `json:"per_page"`
	Page    int `json:"page"`
	Total   int `json:"total"`
}

type Filter struct {
	Key   string
	Value string
}

type DateTimeRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func NewMeta(v url.Values, opts ...MetaOptions) (*Meta, error) {
	metaOpts := defaultMetaOption
	for _, optFn := range opts {
		optFn(&metaOpts)
	}

	pagination := PaginationFromURL(v)
	filtering := FilteringFromURL(v, metaOpts.AllFilterKeys, metaOpts.CustomFilterKeys...)

	dateTime, err := DateTimeRangeFromURL(v, metaOpts.FormatDateTimeRange, metaOpts.StartKeyDateTimeRange, metaOpts.EndKeyDateTimeRange)
	if err != nil {
		return nil, err
	}

	return &Meta{
		Pagination:    pagination,
		Filters:       filtering.filtersSlice,
		filters:       filtering.filtersMap,
		DateTimeRange: dateTime,
		Ordering:      OrderingFromURL(v, metaOpts.OrderByKey, metaOpts.OrderDirectionKey),
		opts:          metaOpts,
	}, nil
}

// PaginationFromURL gets pagination meta from request URL.
func PaginationFromURL(u url.Values) Pagination {
	const DefaultPerPage = 10

	p := Pagination{
		PerPage: DefaultPerPage,
		Page:    1,
	}

	pps := u.Get("per_page")
	if v, err := strconv.Atoi(pps); err == nil {
		if v <= 0 {
			v = DefaultPerPage
		}

		p.PerPage = v
	}

	ps := u.Get("page")
	if v, err := strconv.Atoi(ps); err == nil {
		if v < 1 {
			v = 1
		}

		p.Page = v
	}

	return p
}

type filteringFromURL struct {
	filtersSlice []Filter
	filtersMap   map[string]string
}

func FilteringFromURL(v url.Values, allowAllKeys bool, customKeys ...string) filteringFromURL {
	filters := make([]Filter, 0)
	fMap := make(map[string]string)

	//	find the default
	searchBy := v.Get("search_by")
	search := v.Get("search")

	if len(searchBy) > 0 && len(search) > 0 {
		filters = append(filters, Filter{
			Key:   searchBy,
			Value: search,
		})
		fMap[searchBy] = search
	}

	if allowAllKeys {
		for k, val := range v {
			k = strings.ToLower(k)

			if len(val) > 0 && len(val[0]) > 0 {
				filters = append(filters, Filter{
					Key:   k,
					Value: val[0],
				})

				fMap[k] = val[0]
			}
		}
	}

	for _, key := range customKeys {
		value := v.Get(key)

		if len(value) > 0 {
			filters = append(filters, Filter{
				Key:   strings.ToLower(key),
				Value: value,
			})
			fMap[strings.ToLower(key)] = value
		}
	}

	return filteringFromURL{
		filtersSlice: filters,
		filtersMap:   fMap,
	}
}

func OrderingFromURL(v url.Values, orderByKey string, orderDirectionKey string) Ordering {
	o := Ordering{
		OrderBy:        "created_at",
		OrderDirection: "DESC",
	}

	orderBy := v.Get(orderByKey)
	if len(orderBy) > 0 {
		o.OrderBy = strings.ToLower(orderBy)
	}

	orderDirection := strings.ToUpper(v.Get(orderDirectionKey))
	if len(orderDirection) > 0 {
		if orderDirection == "ASC" || orderDirection == "DESC" {
			o.OrderDirection = orderDirection
		}
	}

	return o
}

func DateTimeRangeFromURL(v url.Values, format string, startQuery, endQuery string) (*DateTimeRange, error) {
	ts := v.Get(startQuery)
	te := v.Get(endQuery)
	if len(ts) == 0 || len(te) == 0 {
		return nil, nil
	}

	dts, err := time.Parse(format, ts)
	if err != nil {
		return nil, errors.Wrap(consts.ErrInvalidMetaData(startQuery), fmt.Sprintf("format must be %s", format))
	}

	dte, err := time.Parse(format, te)
	if err != nil {
		return nil, errors.Wrap(consts.ErrInvalidMetaData(endQuery), fmt.Sprintf("format must be %s", format))
	}

	if dts.After(dte) {
		return nil, consts.ErrInvalidMetaData(fmt.Sprintf("%s greater than %s", startQuery, endQuery))
	}

	return &DateTimeRange{
		StartDate: dts,
		EndDate:   dte,
	}, nil
}

type MetaOption struct {
	CustomFilterKeys      []string
	FormatDateTimeRange   string
	StartKeyDateTimeRange string
	EndKeyDateTimeRange   string
	OrderByKey            string
	OrderDirectionKey     string
	AllFilterKeys         bool
}

var defaultMetaOption = MetaOption{
	CustomFilterKeys:      make([]string, 0),
	FormatDateTimeRange:   "2006-01-02 15:04",
	StartKeyDateTimeRange: "start_date",
	EndKeyDateTimeRange:   "end_date",
	OrderByKey:            "order_by",
	OrderDirectionKey:     "order_type",
}

type MetaOptions func(meta *MetaOption)

func WithCustomFilterKeys(keys ...string) MetaOptions {
	return func(meta *MetaOption) {
		meta.CustomFilterKeys = append(meta.CustomFilterKeys, keys...)
	}
}

func WithCustomDateTimeRange(format string, keyStart string, keyEnd string) MetaOptions {
	return func(meta *MetaOption) {
		meta.FormatDateTimeRange = format
		meta.StartKeyDateTimeRange = keyStart
		meta.EndKeyDateTimeRange = keyEnd
	}
}

func WithCustomOrderByKey(orderByKey string, orderDirectionKey string) MetaOptions {
	return func(meta *MetaOption) {
		meta.OrderByKey = orderByKey
		meta.OrderDirectionKey = orderDirectionKey
	}
}

func WithAllFilterKeys() MetaOptions {
	return func(meta *MetaOption) {
		meta.AllFilterKeys = true
	}
}

type Query struct {
	Limit          int               `json:"limit"`
	Offset         int               `json:"offset"`
	Filters        map[string]string `json:"filters"`
	DateTimeRange  *DateTimeRange    `json:"date_time_range"`
	OrderBy        *string           `json:"order_by"`
	OrderDirection string            `json:"order_direction"`
}

type ToQuerier interface {
	Sortable(field string) bool
	Searchable(field string) bool
}

func (m *Meta) ToQuery(fn ToQuerier) Query {
	limit := m.PerPage
	offset := (m.Page - 1) * limit
	filters := make(map[string]string)

	var orderBy *string = nil

	if fn.Sortable(m.OrderBy) {
		orderBy = &m.OrderBy
	}

	for _, f := range m.Filters {
		f.Key = strings.ToLower(f.Key)

		if fn.Searchable(f.Key) {
			filters[f.Key] = f.Value
		}
	}

	return Query{
		Limit:          limit,
		Offset:         offset,
		Filters:        filters,
		DateTimeRange:  m.DateTimeRange,
		OrderBy:        orderBy,
		OrderDirection: m.OrderDirection,
	}
}

func (m *Meta) ToResponse() map[string]any {
	res := make(map[string]any)

	res["per_page"] = m.PerPage
	res["page"] = m.Page
	res["total"] = m.Total

	res[m.opts.OrderByKey] = m.OrderBy
	res[m.opts.OrderDirectionKey] = m.OrderDirection

	if m.DateTimeRange != nil {
		res[m.opts.StartKeyDateTimeRange] = m.DateTimeRange.StartDate.Format(m.opts.FormatDateTimeRange)
		res[m.opts.EndKeyDateTimeRange] = m.DateTimeRange.EndDate.Format(m.opts.FormatDateTimeRange)
	}

	for _, v := range m.Filters {
		res[v.Key] = v.Value
	}

	return res
}

type ToURLValuesOption struct {
	DateTimeFormat    string
	StartKeyDateTime  string
	EndKeyDateTime    string
	OrderByKey        string
	OrderDirectionKey string
	SearchByKey       string
	SearchKey         string
	SearchValue       string
	SearchByValue     string
}

func (m *Meta) AppendToURLValues(v *url.Values) {
	toURLValuesOpt := ToURLValuesOption{
		DateTimeFormat:    m.opts.FormatDateTimeRange,
		StartKeyDateTime:  m.opts.StartKeyDateTimeRange,
		EndKeyDateTime:    m.opts.EndKeyDateTimeRange,
		OrderByKey:        m.opts.OrderByKey,
		OrderDirectionKey: m.opts.OrderDirectionKey,
		SearchByKey:       "",
		SearchKey:         "",
		SearchByValue:     "",
		SearchValue:       "",
	}

	v.Set("per_page", fmt.Sprintf("%d", m.PerPage))
	v.Set("page", fmt.Sprintf("%d", m.Page))

	if m.DateTimeRange != nil {
		v.Set(toURLValuesOpt.StartKeyDateTime, m.DateTimeRange.StartDate.Format(toURLValuesOpt.DateTimeFormat))
		v.Set(toURLValuesOpt.EndKeyDateTime, m.DateTimeRange.EndDate.Format(toURLValuesOpt.DateTimeFormat))
	}

	v.Set(toURLValuesOpt.OrderByKey, m.OrderBy)
	v.Set(toURLValuesOpt.OrderDirectionKey, m.OrderDirection)

	if len(toURLValuesOpt.SearchByKey) > 0 && len(toURLValuesOpt.SearchKey) > 0 && len(toURLValuesOpt.SearchValue) > 0 {
		v.Set(toURLValuesOpt.SearchByKey, toURLValuesOpt.SearchByValue)
		v.Set(toURLValuesOpt.SearchKey, toURLValuesOpt.SearchValue)
	}

	for _, f := range m.Filters {
		v.Set(f.Key, f.Value)
	}
}

func (m *Meta) GetFilter(key string) string {
	if d, ok := m.filters[key]; ok {
		return d
	}

	return ""
}

func (m *Meta) ReplaceFilter(key string, value string) {
	for i, f := range m.Filters {
		if strings.Compare(f.Key, key) == 0 {
			m.Filters[i] = Filter{
				Key:   f.Key,
				Value: value,
			}
			break
		}
	}

	m.filters[key] = value
}

func (m *Meta) AddFilter(key string, value string) {
	m.Filters = append(m.Filters, Filter{
		Key:   key,
		Value: value,
	})

	m.filters[key] = value
}
