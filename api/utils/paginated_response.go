package utils

type PaginatedResponse[T any] struct {
	Page        int  `json:"page"`
	HasNextPage bool `json:"hasNextPage"`
	Items       []T  `json:"items"`
}

type CursorPaginatedResponse[T any] struct {
	NextCursor *int64 `json:"nextCursor"`
	Items      []T    `json:"items"`
}

func PaginatedResponseFrom[T any](
	items []T,
	page int,
	limit int,
) PaginatedResponse[T] {
	hasNextPage := len(items) > limit

	if hasNextPage {
		items = items[:limit]
	}

	return PaginatedResponse[T]{
		Page:        page,
		HasNextPage: hasNextPage,
		Items:       items,
	}
}

func CursorPaginatedResponseFrom[T any](
	items []T,
	limit int,
	getCursor func(T) int64,
) CursorPaginatedResponse[T] {
	var nextCursor *int64

	if len(items) > limit {
		items = items[:limit]
		cursor := getCursor(items[len(items)-1])
		nextCursor = &cursor
	}

	return CursorPaginatedResponse[T]{
		NextCursor: nextCursor,
		Items:      items,
	}
}
